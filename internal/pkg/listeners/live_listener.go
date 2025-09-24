package listeners

import (
	"sync/atomic"
	"time"

	"github.com/shichen437/gowlive/internal/pkg/events"
	"github.com/shichen437/gowlive/internal/pkg/lives"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/lthibault/jitterbug"
)

type Listener interface {
	Start() error
	Close()
}

type listener struct {
	session *lives.LiveSession
	ed      events.Dispatcher
	state   uint32
	stop    chan struct{}
}

func NewListener(session *lives.LiveSession) Listener {
	return &listener{
		session: session,
		ed:      session.EventDispatcher.(events.Dispatcher),
		state:   begin,
		stop:    make(chan struct{}),
	}
}

func (l *listener) Start() error {
	if !atomic.CompareAndSwapUint32(&l.state, begin, pending) {
		return nil
	}
	defer atomic.CompareAndSwapUint32(&l.state, pending, running)

	l.ed.DispatchEvent(events.NewEvent("ListenStart", l.session))
	l.refresh()
	go l.run()
	return nil
}

func (l *listener) Close() {
	if !atomic.CompareAndSwapUint32(&l.state, running, stopped) {
		return
	}
	l.ed.DispatchEvent(events.NewEvent("ListenStop", l.session))
	close(l.stop)
}

func (l *listener) refresh() {
	info, err := l.session.LiveApi.GetInfo()
	if err != nil {
		g.Log().Error(gctx.GetInitCtx(), "failed to get live info")
		return
	}

	oldState := l.session.GetState()
	oldStatus := status{
		roomName:   oldState.RoomName,
		anchor:     oldState.Anchor,
		roomStatus: oldState.IsLive,
	}

	var evtTyp events.EventType
	latestStatus := getLatestStatus(info)

	isStatusChanged := true
	switch oldStatus.Diff(latestStatus) {
	case 0:
		isStatusChanged = false
	case statusToTrueEvt:
		evtTyp = "LiveStart"
	case statusToFalseEvt:
		evtTyp = "LiveEnd"
	}

	if isStatusChanged {
		l.ed.DispatchEvent(events.NewEvent(evtTyp, l.session))
	}

	l.session.UpdateState(*info)

	if !info.IsLive && (oldState.RoomName != info.RoomName || oldState.Anchor != info.Anchor) {
		l.ed.DispatchEvent(events.NewEvent("NameChanged", l.session))
	}
}

func getLatestStatus(info *lives.LiveState) status {
	return status{
		roomName:   info.RoomName,
		anchor:     info.Anchor,
		roomStatus: info.IsLive,
	}
}

func (l *listener) run() {
	interval := max(l.session.Config.Interval, 30)
	ticker := jitterbug.New(
		time.Duration(interval)*time.Second,
		jitterbug.Norm{
			Stdev: time.Second * 5,
		},
	)
	defer ticker.Stop()

	for {
		select {
		case <-l.stop:
			return
		case <-ticker.C:
			l.refresh()
		}
	}
}
