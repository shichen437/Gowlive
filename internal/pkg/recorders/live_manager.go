package recorders

import (
	"context"
	"sync"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/gowlive/internal/pkg/events"
	"github.com/shichen437/gowlive/internal/pkg/interfaces"
	"github.com/shichen437/gowlive/internal/pkg/lives"
)

func NewManager(ctx context.Context, session *lives.LiveSession) Manager {
	return &manager{
		session: session,
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

var (
	newRecorder = NewRecorder
)

type manager struct {
	lock     sync.RWMutex
	recorder Recorder
	session  *lives.LiveSession
}

func (m *manager) registryListener(ctx context.Context, ed events.Dispatcher) {
	ed.AddEventListener("LiveStart", events.NewEventListener(func(event *events.Event) {
		session := event.Object.(*lives.LiveSession)
		if session.Id != m.session.Id {
			return
		}
		err := m.AddRecorder(ctx)
		if err != nil {
			g.Log().Error(ctx, "failed to add recorder")
			return
		}
		liveStartBiz(ctx, session.Id, session.State.Anchor)
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
			liveEndBiz(ctx, session.Id, session.State.Anchor)
		}
	})
	ed.AddEventListener("LiveEnd", removeEvtListener)
	ed.AddEventListener("ListenStop", removeEvtListener)
}

func (m *manager) Start(ctx context.Context) error {
	ed := m.session.EventDispatcher.(events.Dispatcher)
	m.registryListener(ctx, ed)
	g.Log().Infof(ctx, "RecorderManager Started for session %d!", m.session.Id)
	return nil
}

func (m *manager) Close(ctx context.Context) {
	if err := m.RemoveRecorder(ctx); err == nil {
		liveEndBiz(ctx, m.session.Id, m.session.State.Anchor)
	}
	g.Log().Infof(ctx, "RecorderManager Closed for session %d!", m.session.Id)
}

func (m *manager) AddRecorder(ctx context.Context) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.recorder != nil {
		return gerror.New("this live has a recorder")
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
