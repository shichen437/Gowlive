// filename: listeners/listeners.go
package listeners

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/lthibault/jitterbug"

	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/events"
	"github.com/shichen437/gowlive/internal/pkg/lives"
	"github.com/shichen437/gowlive/internal/pkg/providers"
)

type Listener interface {
	Start() error
	Close()
}

type listener struct {
	session *lives.LiveSession
	ed      events.Dispatcher

	state uint32

	stop     chan struct{}
	stopOnce sync.Once

	mu sync.RWMutex
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

	defer func() {
		atomic.StoreUint32(&l.state, running)
	}()

	l.ed.DispatchEvent(events.NewEvent("ListenStart", l.session))

	isLive := l.refresh()

	if l.session.Config.MonitorType == consts.MonitorTypeIntelligent {
		go l.runForIntelligent(isLive)
	} else {
		go l.run()
	}
	return nil
}

func (l *listener) Close() {
	for {
		cur := atomic.LoadUint32(&l.state)
		if cur == stopped {
			return
		}
		if cur != running && cur != pending {
			if atomic.CompareAndSwapUint32(&l.state, cur, stopped) {
				break
			}
			continue
		}
		if atomic.CompareAndSwapUint32(&l.state, cur, stopped) {
			break
		}
	}

	l.ed.DispatchEvent(events.NewEvent("ListenStop", l.session))

	l.stopOnce.Do(func() {
		close(l.stop)
	})
}

func (l *listener) refresh() bool {
	info, err := l.session.LiveApi.GetInfo()
	if err != nil {
		g.Log().Errorf(gctx.GetInitCtx(), "获取直播信息失败: %v", err)
		return false
	}

	l.mu.RLock()
	oldState := l.session.GetState()
	l.mu.RUnlock()

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

	if info.Anchor != "" {
		l.mu.Lock()
		l.session.UpdateState(*info)
		l.mu.Unlock()
	}

	if !info.IsLive && info.Anchor != "" && (oldState.RoomName != info.RoomName || oldState.Anchor != info.Anchor) {
		l.ed.DispatchEvent(events.NewEvent("NameChanged", l.session))
	}
	return info.IsLive
}

func getLatestStatus(info *lives.LiveState) status {
	return status{
		roomName:   info.RoomName,
		anchor:     info.Anchor,
		roomStatus: info.IsLive,
	}
}

func (l *listener) calcInterval(isLive bool) int {
	l.mu.RLock()
	minInterval := l.session.Config.Interval
	id := l.session.Id
	platform := l.session.State.Platform
	l.mu.RUnlock()

	var interval int
	if isLive {
		interval = providers.GetRegularInterval(gctx.GetInitCtx(), id, platform)
	} else {
		interval = providers.GetIntelligentInterval(gctx.GetInitCtx(), id, platform)
	}
	if minInterval > 0 {
		interval = max(minInterval, interval)
	}
	return max(interval, consts.DefaultInterval)
}

func (l *listener) run() {
	interval := max(l.session.Config.Interval, consts.DefaultInterval)
	ticker := jitterbug.New(
		time.Duration(interval)*time.Second,
		jitterbug.Norm{
			Stdev: time.Second * 5,
		},
	)

	for {
		select {
		case <-l.stop:
			ticker.Stop()
			return
		case <-ticker.C:
			l.refresh()
		}
	}
}

func (l *listener) runForIntelligent(isLive bool) {
	interval := l.calcInterval(isLive)

	ticker := jitterbug.New(
		time.Duration(interval)*time.Second,
		jitterbug.Norm{
			Stdev: time.Second * 5,
		},
	)

	for {
		select {
		case <-l.stop:
			ticker.Stop()
			return
		case <-ticker.C:
			currentIsLive := l.refresh()

			newInterval := l.calcInterval(currentIsLive)

			if newInterval != interval {
				old := ticker
				interval = newInterval
				ticker = jitterbug.New(
					time.Duration(interval)*time.Second,
					jitterbug.Norm{
						Stdev: time.Second * 5,
					},
				)
				old.Stop()
			}
		}
	}
}
