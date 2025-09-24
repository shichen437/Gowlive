package listeners

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/gowlive/internal/pkg/interfaces"
	"github.com/shichen437/gowlive/internal/pkg/lives"
)

const (
	begin uint32 = iota
	pending
	running
	stopped
)

type manager struct {
	lock     sync.RWMutex
	listener Listener
	session  *lives.LiveSession
}

type Manager interface {
	interfaces.Module
	AddListener(ctx context.Context) error
	RemoveListener(ctx context.Context) error
	GetListener(ctx context.Context) (Listener, error)
	HasListener(ctx context.Context) bool
}

func NewManager(ctx context.Context, session *lives.LiveSession) Manager {
	return &manager{
		session: session,
	}
}

func (m *manager) Start(ctx context.Context) error {
	g.Log().Infof(ctx, "ListenerManager Started for session %d!", m.session.Id)
	return m.AddListener(ctx)
}

func (m *manager) Close(ctx context.Context) {
	m.RemoveListener(ctx)
	g.Log().Infof(ctx, "ListenerManager Closed for session %d!", m.session.Id)
}

func (m *manager) AddListener(ctx context.Context) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.listener != nil {
		g.Log().Warning(ctx, "this live has a listener")
		return nil
	}
	listener := NewListener(m.session)
	m.listener = listener
	return listener.Start()
}

func (m *manager) RemoveListener(ctx context.Context) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	listener := m.listener
	if listener == nil {
		g.Log().Warning(ctx, "this live has not a listener")
		return nil
	}
	listener.Close()
	m.listener = nil
	return nil
}

func (m *manager) GetListener(ctx context.Context) (Listener, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	listener := m.listener
	if listener == nil {
		return nil, gerror.Newf("this live has not a listener")
	}
	return listener, nil
}

func (m *manager) HasListener(ctx context.Context) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.listener != nil
}
