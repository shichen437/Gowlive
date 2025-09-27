package recorders

import (
	"bytes"
	"context"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shichen437/gowlive/internal/pkg/events"
	"github.com/shichen437/gowlive/internal/pkg/lives"
	parser "github.com/shichen437/gowlive/internal/pkg/stream_parser"
	ffmpeg_parser "github.com/shichen437/gowlive/internal/pkg/stream_parser/ffmpeg"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

const (
	begin uint32 = iota
	pending
	running
	stopped
)

type Recorder interface {
	Start(ctx context.Context) error
	Close()
	StartTime() time.Time
	GetStatus() (map[string]string, error)
}

type recorder struct {
	session *lives.LiveSession

	ed         events.Dispatcher
	startTime  time.Time
	parser     parser.Parser
	parserLock *sync.RWMutex

	stop  chan struct{}
	state uint32
}

func NewRecorder(session *lives.LiveSession) (Recorder, error) {
	return &recorder{
		session:    session,
		startTime:  time.Now(),
		ed:         session.EventDispatcher.(events.Dispatcher),
		state:      begin,
		stop:       make(chan struct{}),
		parserLock: new(sync.RWMutex),
	}, nil
}

func (r *recorder) Start(ctx context.Context) error {
	if !atomic.CompareAndSwapUint32(&r.state, begin, pending) {
		return nil
	}
	go r.run(ctx)
	go r.monitorDiskSpace(ctx)
	r.ed.DispatchEvent(events.NewEvent("RecorderStart", r.session))
	atomic.CompareAndSwapUint32(&r.state, pending, running)
	return nil
}

func (r *recorder) Close() {
	if !atomic.CompareAndSwapUint32(&r.state, running, stopped) {
		return
	}
	close(r.stop)
	if p := r.getParser(); p != nil {
		if err := p.Stop(); err != nil {
			g.Log().Error(gctx.New(), "failed to end recorder")
		}
	}
	r.ed.DispatchEvent(events.NewEvent("RecorderStop", r.session))
}

func (r *recorder) StartTime() time.Time {
	return r.startTime
}

func (r *recorder) GetStatus() (map[string]string, error) {
	statusP, ok := r.getParser().(parser.StatusParser)
	if !ok {
		return nil, gerror.New("parser does not implement StatusParser")
	}
	return statusP.Status()
}

func (r *recorder) tryRecord(ctx context.Context) error {
	state := r.session.GetState()
	if !state.IsLive {
		return gerror.New("stream is not live")
	}

	streamInfos := state.StreamInfos
	if len(streamInfos) == 0 {
		g.Log().Warningf(ctx, "Stream info not in session state for liveId %d, fetching manually.", r.session.Id)
		info, err := r.session.LiveApi.GetInfo()
		if err != nil {
			return gerror.Wrap(err, "failed to get stream info in tryRecord fallback")
		}
		if !info.IsLive {
			return gerror.New("stream is not live (checked in fallback)")
		}
		streamInfos = info.StreamInfos
		if len(streamInfos) == 0 {
			return gerror.New("no stream info found in fallback")
		}
	}
	streamInfo := streamInfos[0]

	fileName, outputPath, err := r.getOutPathAndFilename(&state)
	if err != nil {
		return gerror.Wrap(err, "failed to get output path and filename")
	}
	if err = os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return gerror.Wrap(err, "failed to create output path")
	}
	parserCfg := map[string]string{
		"timeout_in_us": strconv.Itoa(100000000),
	}
	p, err := newParser(parserCfg)
	if err != nil {
		return gerror.Wrap(err, "failed to init parser")
	}
	r.setAndCloseParser(p)
	r.startTime = time.Now()
	err = r.parser.ParseLiveStream(ctx, streamInfo, fileName, r.session.Config.RoomUrl)
	removeEmptyFile(fileName)
	return err
}

func newParser(cfg map[string]string) (parser.Parser, error) {
	return parser.New(ffmpeg_parser.Name, cfg)
}

func removeEmptyFile(file string) {
	if stat, err := os.Stat(file); err == nil && stat.Size() == 0 {
		os.Remove(file)
	}
}

func (r *recorder) run(ctx context.Context) {
	const maxBackoff = 5 * time.Minute
	backoff := 5 * time.Second

	for {
		select {
		case <-r.stop:
			return
		default:
			err := r.tryRecord(ctx)
			if err != nil {
				g.Log().Warningf(ctx, "recording process for session %d stopped with error: %v. Retrying in %s", r.session.Id, err, backoff)
				select {
				case <-time.After(backoff):
					backoff *= 2
					if backoff > maxBackoff {
						backoff = maxBackoff
					}
				case <-r.stop:
					return
				}
			} else {
				g.Log().Infof(ctx, "recording process for session %d stopped cleanly. Checking again in %s", r.session.Id, backoff)
				select {
				case <-time.After(backoff):
					backoff = 5 * time.Second
				case <-r.stop:
					return
				}
			}
		}
	}
}

func (r *recorder) getParser() parser.Parser {
	r.parserLock.RLock()
	defer r.parserLock.RUnlock()
	return r.parser
}

func (r *recorder) setAndCloseParser(p parser.Parser) {
	r.parserLock.Lock()
	defer r.parserLock.Unlock()
	if r.parser != nil {
		if err := r.parser.Stop(); err != nil {
			g.Log().Error(gctx.New(), "failed to end recorder", err)
		}
	}
	r.parser = p
}

func (r *recorder) getOutPathAndFilename(info *lives.LiveState) (string, string, error) {
	format := r.session.Config.Format
	if format == "" {
		format = "flv"
	}
	buf := new(bytes.Buffer)
	outTmpl := utils.GetOutputPathTemplate()
	err := outTmpl.Execute(buf, info)
	if err != nil {
		return "", "", gerror.New("failed to get outputPath template")
	}
	outputPath := buf.String()
	filenameTmpl := utils.GetFilenameTemplate(outputPath, format)
	buf.Reset()
	err = filenameTmpl.Execute(buf, info)
	if err != nil {
		return "", "", gerror.New("failed to get filename template")
	}
	filename := buf.String()
	return filename, outputPath, nil
}

func (r *recorder) monitorDiskSpace(ctx context.Context) {
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if atomic.LoadUint32(&r.state) != running {
				return
			}
			if utils.GetDiskUsage() > 95 {
				g.Log().Warningf(ctx, "Disk usage is over 95%%. Stopping recording for session %d.", r.session.Id)
				r.ed.DispatchEvent(events.NewEvent("RecordingStoppedDueToDiskSpace", r.session))
				return
			}
		case <-r.stop:
			return
		}
	}
}
