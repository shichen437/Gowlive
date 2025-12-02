package ffmpeg_parser

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gogf/gf/v2/util/gconv"
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
		statusReq:   make(chan struct{}, 1),
		statusResp:  make(chan map[string]string, 1),
		timeoutInUs: cfg["timeout_in_us"],
		referer:     cfg["referer"],
		format:      cfg["format"],
		st:          cfg["st"],
		fr:          fr,
	}, nil
}

func (p *Parser) scanFFmpegStatus() <-chan []byte {
	ch := make(chan []byte)
	br := bufio.NewScanner(p.cmdStdout)
	br.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		if idx := bytes.Index(data, []byte("progress=continue\n")); idx >= 0 {
			return idx + 1, data[0:idx], nil
		}

		return 0, nil, nil
	})
	go func() {
		defer close(ch)
		for br.Scan() {
			ch <- br.Bytes()
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

func (p *Parser) Status() (map[string]string, error) {
	// TODO: check parser is running
	p.statusReq <- struct{}{}
	return <-p.statusResp, nil
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
	err = p.cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (p *Parser) buildArgs(ffUserAgent, file string, sUrl *url.URL, headers map[string]string) []string {
	referer, exists := headers["Referer"]
	if !exists {
		referer = p.referer
	}
	args := p.basicArgs(ffUserAgent, referer, sUrl)
	for k, v := range headers {
		if k == "User-Agent" || k == "Referer" {
			continue
		}
		args = append(args, "-headers", k+": "+v)
	}
	switch strings.ToLower(p.format) {
	case "flv":
		args = p.flvArgs(file, args)
	case "mp4":
		args = p.mp4Args(file, args)
	case "mkv":
		args = p.mkvArgs(file, args)
	case "ts":
		args = p.tsArgs(file, args)
	default:
		args = p.mp3Args(file, args)
	}
	return args
}

func (p *Parser) Stop() (err error) {
	p.closeOnce.Do(func() {
		p.cmdLock.Lock()
		defer p.cmdLock.Unlock()
		if p.cmd != nil && p.cmd.Process != nil && p.cmd.ProcessState == nil {
			err = syscall.Kill(-p.cmd.Process.Pid, syscall.SIGTERM)
		}
	})
	return err
}

func (p *Parser) basicArgs(ffUserAgent, referer string, sUrl *url.URL) []string {
	return []string{
		"-nostats",
		"-progress", "-",
		"-y", "-re",
		"-fflags", "+genpts+igndts+discardcorrupt",
		"-err_detect", "ignore_err",
		"-reconnect", "1",
		"-reconnect_streamed", "1",
		"-reconnect_delay_max", "5",
		"-user_agent", ffUserAgent,
		"-referer", referer,
		"-rw_timeout", p.timeoutInUs,
		"-i", sUrl.String(),
	}
}

func (p *Parser) flvArgs(file string, args []string) []string {
	args = append(args, "-c", "copy")
	if gconv.Int(p.st) > 0 {
		template := utils.BuildSegmentTemplate(file, ".flv")
		args = append(args,
			"-f", "segment",
			"-segment_time", p.st,
			"-reset_timestamps", "1",
			template,
		)
	} else {
		args = append(args, "-f", "flv", utils.EnsureSuffix(file, ".flv"))
	}
	return args
}

func (p *Parser) mp4Args(file string, args []string) []string {
	args = append(args, "-c", "copy", "-bsf:a", "aac_adtstoasc")
	if gconv.Int(p.st) > 0 {
		template := utils.BuildSegmentTemplate(file, ".mp4")
		args = append(args,
			"-f", "segment",
			"-segment_time", p.st,
			"-reset_timestamps", "1",
			"-segment_format_options", "movflags=+faststart",
			template,
		)
	} else {
		args = append(args, "-movflags", "+faststart", "-f", "mp4", utils.EnsureSuffix(file, ".mp4"))
	}
	return args
}

func (p *Parser) mkvArgs(file string, args []string) []string {
	if p.fr {
		args = append(args, "-vf", "scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2:black")
		args = append(args, "-c:v", "libx264", "-preset", "fast", "-crf", "23")
		args = append(args, "-c:a", "aac", "-b:a", "128k")
	} else {
		args = append(args, "-c", "copy")
	}

	if gconv.Int(p.st) > 0 {
		template := utils.BuildSegmentTemplate(file, ".mkv")
		args = append(args,
			"-f", "segment",
			"-segment_time", p.st,
			"-reset_timestamps", "1",
			template,
		)
	} else {
		args = append(args, "-f", "matroska", utils.EnsureSuffix(file, ".mkv"))
	}
	return args
}

func (p *Parser) tsArgs(file string, args []string) []string {
	args = append(args, "-c", "copy")

	if gconv.Int(p.st) > 0 {
		template := utils.BuildSegmentTemplate(file, ".ts")
		args = append(args,
			"-f", "segment",
			"-segment_time", p.st,
			"-reset_timestamps", "1",
			template,
		)
	} else {
		args = append(args, "-f", "mpegts", utils.EnsureSuffix(file, ".ts"))
	}
	return args
}

func (p *Parser) mp3Args(file string, args []string) []string {
	args = append(args, "-vn", "-c:a", "libmp3lame", "-b:a", "192k")
	if gconv.Int(p.st) > 0 {
		template := utils.BuildSegmentTemplate(file, ".mp3")
		args = append(args,
			"-f", "segment",
			"-segment_time", p.st,
			"-reset_timestamps", "1",
			template,
		)
	} else {
		args = append(args, "-f", "mp3", utils.EnsureSuffix(file, ".mp3"))
	}
	return args
}
