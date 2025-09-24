package registry

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/events"
	"github.com/shichen437/gowlive/internal/pkg/listeners"
	"github.com/shichen437/gowlive/internal/pkg/lives"
	"github.com/shichen437/gowlive/internal/pkg/manager"
	"github.com/shichen437/gowlive/internal/pkg/recorders"
	"github.com/shichen437/gowlive/internal/pkg/service"
)

var (
	sessionRegistry = NewSessionRegistry()
)

type SessionRegistry struct {
	mu       sync.RWMutex
	sessions map[int]*lives.LiveSession
}

func NewSessionRegistry() *SessionRegistry {
	return &SessionRegistry{
		sessions: make(map[int]*lives.LiveSession),
	}
}

func Get() *SessionRegistry {
	return sessionRegistry
}

func (r *SessionRegistry) Add(ctx context.Context, liveId int) error {
	session := service.GenLiveSessionById(ctx, liveId)
	if session == nil {
		manager.GetLogManager().AddErrorLog(consts.LogTypeLive, "获取直播会话为空")
		return gerror.New("获取直播会话为空")
	}
	return r.add(ctx, session)
}

func (r *SessionRegistry) AddAll(ctx context.Context, liveIds []int) {
	sessions := service.GenLiveSessionsByIds(ctx, liveIds)
	var wg sync.WaitGroup
	for _, session := range sessions {
		wg.Add(1)
		go func(s *lives.LiveSession) {
			defer wg.Done()
			if err := r.add(ctx, s); err != nil {
				g.Log().Errorf(ctx, "添加直播会话失败, liveId: %d, err: %v", s.Id, err)

			}
		}(session)
	}
	wg.Wait()
}

func (r *SessionRegistry) add(ctx context.Context, session *lives.LiveSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.sessions[session.Id]; ok {
		return gerror.Newf("session with ID %d already exists", session.Id)
	}

	liveApi, err := lives.New(session.Config.RoomUrl)
	if err != nil {
		manager.GetLogManager().AddErrorLog(consts.LogTypeLive, "解析 Api 失败: "+session.Config.RoomUrl)
		return gerror.Newf("failed to create live session : %s", err)
	}
	session.LiveApi = liveApi

	session.EventDispatcher = events.NewDispatcher(ctx)
	session.ListenerManager = listeners.NewManager(ctx, session)
	session.RecorderManager = recorders.NewManager(ctx, session)

	if err := session.EventDispatcher.Start(ctx); err != nil {
		return gerror.Newf("failed to start EventDispatcher for session : %s", err)
	}
	if err := session.RecorderManager.Start(ctx); err != nil {
		session.EventDispatcher.Close(ctx)
		return gerror.Newf("failed to start RecorderManager for session : %s", err)
	}
	if err := session.ListenerManager.Start(ctx); err != nil {
		session.RecorderManager.Close(ctx)
		session.EventDispatcher.Close(ctx)
		return gerror.Newf("failed to start ListenerManager for session : %s", err)
	}

	r.sessions[session.Id] = session
	manager.GetLogManager().AddSuccessLog(consts.LogTypeLive, "创建直播会话成功："+session.Config.RoomUrl)
	g.Log().Infof(ctx, "Successfully added and started session %d (%s)", session.Id, session.Config.RoomUrl)
	return nil
}

func (r *SessionRegistry) Remove(ctx context.Context, sessionId int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	session, ok := r.sessions[sessionId]
	if !ok {
		return gerror.Newf("session with ID %d not found", sessionId)
	}

	if session.ListenerManager != nil {
		session.ListenerManager.Close(ctx)
	}
	if session.RecorderManager != nil {
		session.RecorderManager.Close(ctx)
	}
	if session.EventDispatcher != nil {
		session.EventDispatcher.Close(ctx)
	}

	delete(r.sessions, sessionId)
	g.Log().Infof(ctx, "Successfully stopped and removed session %d", sessionId)
	manager.GetLogManager().AddSuccessLog(consts.LogTypeLive, "删除直播会话成功")
	return nil
}

func (r *SessionRegistry) Get(sessionId int) (*lives.LiveSession, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	session, ok := r.sessions[sessionId]
	return session, ok
}

func (r *SessionRegistry) Exists(sessionId int) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.sessions[sessionId]
	return ok
}

func (r *SessionRegistry) IsRecording(sessionId int) bool {
	r.mu.RLock()
	session, ok := r.sessions[sessionId]
	r.mu.RUnlock()

	if !ok || session.RecorderManager == nil {
		return false
	}

	return session.RecorderManager.(recorders.Manager).HasRecorder(context.Background())
}

func (r *SessionRegistry) RecordingCount() int {
	count := 0
	for _, session := range r.sessions {
		if r.IsRecording(session.Id) {
			count++
		}
	}
	return count
}
