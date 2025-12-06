package ffmpeg_parser

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/shichen437/gowlive/internal/pkg/lives"
	parser "github.com/shichen437/gowlive/internal/pkg/stream_parser"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

const (
	Name      = "ffmpeg"
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36"
)

func init() {
	parser.Register(Name, new(builder))
}

type builder struct{}

type Parser struct {
	cmd         *exec.Cmd
	cmdStdIn    io.WriteCloser
	cmdStdout   io.ReadCloser
	closeOnce   *sync.Once
	waitOnce    *sync.Once
	debug       bool
	timeoutInUs string
	referer     string
	format      string
	st          string // 切片时间
	fr          bool   // 固定分辨率
	statusReq   chan struct{}
	statusResp  chan map[string]string
	cmdLock     sync.Mutex
}

func (b *builder) Build(cfg map[string]string) (parser.Parser, error) {
	debug, fr := false, false
	if debugFlag, ok := cfg["debug"]; ok && debugFlag != "" {
		debug = true
	}
	if frFlag, ok := cfg["fr"]; ok && frFlag == "true" {
		fr = true
	}
	return &Parser{
		debug:       debug,
		closeOnce:   new(sync.Once),
		waitOnce:    new(sync.Once),
		statusReq:   make(chan struct{}, 1),
		statusResp:  make(chan map[string]string, 1),
		timeoutInUs: cfg["timeout_in_us"],
		referer:     cfg["referer"],
		format:      cfg["format"],
		st:          cfg["st"],
		fr:          fr,
	}, nil
}

func (p *Parser) Status() (map[string]string, error) {
	select {
	case p.statusReq <- struct{}{}:
	default:
	}
	resp := <-p.statusResp
	if resp == nil {
		return nil, io.EOF
	}
	return resp, nil
}

func (p *Parser) ParseLiveStream(ctx context.Context, streamInfo *lives.StreamUrlInfo, file string) (err error) {
	url := streamInfo.Url
	ffmpegPath, err := utils.GetDefaultFFmpegPath()
	if err != nil {
		return err
	}
	headers := streamInfo.HeadersForDownloader
	ffUserAgent, exists := headers["User-Agent"]
	if !exists {
		ffUserAgent = userAgent
	}

	args := p.buildArgs(ffUserAgent, file, url, headers)

	func() {
		p.cmdLock.Lock()
		defer p.cmdLock.Unlock()
		p.cmd = exec.Command(ffmpegPath, args...)
		p.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		if p.cmdStdIn, err = p.cmd.StdinPipe(); err != nil {
			return
		}
		if p.cmdStdout, err = p.cmd.StdoutPipe(); err != nil {
			return
		}
		if p.debug {
			p.cmd.Stderr = os.Stderr
		}
		if err = p.cmd.Start(); err != nil {
			if p.cmd.Process != nil {
				p.cmd.Process.Kill()
			}
			return
		}
	}()
	if err != nil {
		return err
	}

	go p.scheduler()
	p.waitOnce.Do(func() {
		err = p.cmd.Wait()
	})
	if err != nil {
		return err
	}
	defer func() {
		if p.cmdStdIn != nil {
			p.cmdStdIn.Close()
		}
		if p.cmdStdout != nil {
			p.cmdStdout.Close()
		}
	}()
	return nil
}

func (p *Parser) Stop() (err error) {
	p.closeOnce.Do(func() {
		p.cmdLock.Lock()
		defer p.cmdLock.Unlock()

		if p.cmdStdIn != nil {
			p.cmdStdIn.Close()
		}
		if p.cmdStdout != nil {
			p.cmdStdout.Close()
		}

		if p.cmd != nil && p.cmd.Process != nil && p.cmd.ProcessState == nil {
			_ = syscall.Kill(-p.cmd.Process.Pid, syscall.SIGTERM)

			done := make(chan error, 1)
			go func() {
				p.waitOnce.Do(func() {
					done <- p.cmd.Wait()
				})
				if p.cmd.ProcessState != nil {
					done <- nil
				}
			}()

			select {
			case <-time.After(3 * time.Second):
				_ = syscall.Kill(-p.cmd.Process.Pid, syscall.SIGKILL)
				<-done
			case err = <-done:
			}
		}
	})
	return err
}

func (p *Parser) scanFFmpegStatus() <-chan []byte {
	ch := make(chan []byte)
	br := bufio.NewScanner(p.cmdStdout)
	br.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	var buf bytes.Buffer
	go func() {
		defer close(ch)
		for br.Scan() {
			line := br.Bytes()
			buf.Write(line)
			buf.WriteByte('\n')

			if bytes.Equal(line, []byte("progress=continue")) || bytes.Equal(line, []byte("progress=end")) {
				ch <- append([]byte(nil), buf.Bytes()...)
				buf.Reset()
			}
		}
	}()
	return ch
}

func (p *Parser) decodeFFmpegStatus(b []byte) (status map[string]string) {
	status = map[string]string{
		"parser": Name,
	}
	s := bufio.NewScanner(bytes.NewReader(b))
	s.Split(bufio.ScanLines)
	for s.Scan() {
		split := bytes.SplitN(s.Bytes(), []byte("="), 2)
		if len(split) != 2 {
			continue
		}
		status[string(bytes.TrimSpace(split[0]))] = string(bytes.TrimSpace(split[1]))
	}
	return
}

func (p *Parser) scheduler() {
	defer close(p.statusResp)
	statusCh := p.scanFFmpegStatus()
	for {
		select {
		case <-p.statusReq:
			select {
			case b, ok := <-statusCh:
				if !ok {
					p.statusResp <- nil
					return
				}
				p.statusResp <- p.decodeFFmpegStatus(b)
			case <-time.After(time.Second * 3):
				p.statusResp <- nil
			}
		default:
			if _, ok := <-statusCh; !ok {
				return
			}
		}
	}
}
