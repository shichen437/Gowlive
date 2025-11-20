package recorders

import (
	"context"
	"sync"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/lthibault/jitterbug"
	"github.com/shichen437/gowlive/internal/pkg/events"
	"github.com/shichen437/gowlive/internal/pkg/interfaces"
	"github.com/shichen437/gowlive/internal/pkg/lives"
)

var (
	newRecorder = NewRecorder
)

type manager struct {
	lock     sync.RWMutex
	recorder Recorder
	session  *lives.LiveSession
	stop     chan struct{}
	stopOnce sync.Once
}

func NewManager(ctx context.Context, session *lives.LiveSession) Manager {
	return &manager{
		session: session,
		stop:    make(chan struct{}),
	}
}

type Manager interface {
	interfaces.Module
	AddRecorder(ctx context.Context) error
	RemoveRecorder(ctx context.Context) error
	RestartRecorder(ctx context.Context) error
	GetRecorder(ctx context.Context) (Recorder, error)
	HasRecorder(ctx context.Context) bool
}

func (m *manager) registryListener(ctx context.Context, ed events.Dispatcher) {
	ed.AddEventListener("LiveStart", events.NewEventListener(func(event *events.Event) {
		session := event.Object.(*lives.LiveSession)
		if session.Id != m.session.Id {
			return
		}
		liveStartBiz(ctx, session)
		err := m.AddRecorder(ctx)
		if err != nil {
			g.Log().Errorf(ctx, "failed to add recorder for session %d: %v", m.session.Id, err)
			return
		}
	}))

	ed.AddEventListener("NameChanged", events.NewEventListener(func(event *events.Event) {
		session := event.Object.(*lives.LiveSession)
		if session.Id != m.session.Id {
			return
		}
		m.updateName(ctx, session)
	}))

	removeEvtListener := events.NewEventListener(func(event *events.Event) {
		session := event.Object.(*lives.LiveSession)
		if session.Id != m.session.Id {
			return
		}
		if err := m.RemoveRecorder(ctx); err == nil {
			liveEndBiz(ctx, session)
		}
	})
	ed.AddEventListener("LiveEnd", removeEvtListener)
	ed.AddEventListener("ListenStop", removeEvtListener)

	ed.AddEventListener("RecordingStoppedDueToDiskSpace", events.NewEventListener(func(event *events.Event) {
		session := event.Object.(*lives.LiveSession)
		if session.Id != m.session.Id {
			return
		}
		g.Log().Warningf(ctx, "Received RecordingStoppedDueToDiskSpace event for session %d. Removing recorder.", m.session.Id)
		if err := m.RemoveRecorder(ctx); err == nil {
			liveEndBiz(ctx, session)
		} else {
			g.Log().Errorf(ctx, "Failed to remove recorder on RecordingStoppedDueToDiskSpace event for session %d: %v", m.session.Id, err)
		}
	}))
}

func (m *manager) Start(ctx context.Context) error {
	ed := m.session.EventDispatcher.(events.Dispatcher)
	m.registryListener(ctx, ed)
	go m.reconciliationLoop(ctx)
	g.Log().Infof(ctx, "RecorderManager Started for session %d!", m.session.Id)
	return nil
}

func (m *manager) Close(ctx context.Context) {
	m.stopOnce.Do(func() {
		close(m.stop)
	})
	if err := m.RemoveRecorder(ctx); err == nil {
		liveEndBiz(ctx, m.session)
	}
	g.Log().Infof(ctx, "RecorderManager Closed for session %d!", m.session.Id)
}

func (m *manager) AddRecorder(ctx context.Context) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.recorder != nil {
		return gerror.New("this live has a recorder")
	}
	if diskOverLimit() {
		g.Log().Warningf(ctx, "Disk usage already reached limit, cannot start recording for session %d", m.session.Id)
		return gerror.New("disk usage reached limit")
	}
	recorder, err := newRecorder(m.session)
	if err != nil {
		return err
	}
	m.recorder = recorder

	return recorder.Start(ctx)
}

func (m *manager) RestartRecorder(ctx context.Context) error {
	if err := m.RemoveRecorder(ctx); err != nil {
		return err
	}
	time.Sleep(20 * time.Second)
	if err := m.AddRecorder(ctx); err != nil {
		return err
	}
	return nil
}

func (m *manager) RemoveRecorder(ctx context.Context) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	recorder := m.recorder
	if recorder == nil {
		return gerror.New("this live has not a recorder")
	}
	recorder.Close()
	m.recorder = nil
	return nil
}

func (m *manager) GetRecorder(ctx context.Context) (Recorder, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	r := m.recorder
	if r == nil {
		return nil, gerror.New("this live has not a recorder")
	}
	return r, nil
}

func (m *manager) HasRecorder(ctx context.Context) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.recorder != nil
}

// 周期性状态协调
func (m *manager) reconciliationLoop(ctx context.Context) {
	ticker := jitterbug.New(
		2*time.Minute,
		&jitterbug.Norm{Stdev: 30 * time.Second},
	)
	defer ticker.Stop()
	for {
		select {
		case <-m.stop:
			g.Log().Infof(ctx, "RecorderManager Closed for session %d!", m.session.Id)
			return
		case <-ticker.C:
			m.reconcile(ctx)
		}
	}
}

func (m *manager) reconcile(ctx context.Context) {
	if !m.session.GetState().IsLive || m.HasRecorder(ctx) {
		return
	}
	if diskOverLimit() {
		return
	}
	err := m.AddRecorder(ctx)
	if err != nil {
		g.Log().Errorf(ctx, "Reconciliation failed to add recorder for session %d: %v", m.session.Id, err)
	}
	liveStartBiz(ctx, m.session)
}
