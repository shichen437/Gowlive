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
	"github.com/shichen437/gowlive/internal/pkg/utils"
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

	sessionsByPlatform := make(map[string][]*lives.LiveSession)
	for _, session := range sessions {
		if session == nil || session.State.Platform == "" {
			g.Log().Warningf(ctx, "Skipping session %d due to nil or missing platform", session.Id)
			continue
		}
		sessionsByPlatform[session.State.Platform] = append(sessionsByPlatform[session.State.Platform], session)
	}

	var wg sync.WaitGroup
	for platform, platformSessions := range sessionsByPlatform {
		wg.Add(1)
		go func(p string, ss []*lives.LiveSession) {
			defer wg.Done()
			g.Log().Infof(ctx, "Starting to add sessions for platform: %s, count: %d", p, len(ss))
			for _, session := range ss {
				if err := r.add(ctx, session); err != nil {
					g.Log().Errorf(ctx, "添加直播会话失败, liveId: %d, platform: %s, err: %v", session.Id, p, err)
				}
			}
			g.Log().Infof(ctx, "Finished adding sessions for platform: %s", p)
		}(platform, platformSessions)
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
		manager.GetLogManager().AddErrorLog(consts.LogTypeLive, utils.T(ctx, "ext.live.parse.api.error")+session.Config.RoomUrl)
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
	manager.GetLogManager().AddSuccessLog(consts.LogTypeLive, utils.T(ctx, "ext.live.add.session.success")+session.Config.RoomUrl)
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
	manager.GetLogManager().AddSuccessLog(consts.LogTypeLive, utils.T(ctx, "ext.live.remove.session.success"))
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

func (r *SessionRegistry) StopAll(ctx context.Context) {
	r.mu.RLock()
	var ids []int
	for id := range r.sessions {
		ids = append(ids, id)
	}
	r.mu.RUnlock()

	if len(ids) == 0 {
		return
	}

	g.Log().Info(ctx, "Stopping all live sessions...")
	var wg sync.WaitGroup
	for _, id := range ids {
		wg.Add(1)
		go func(sessionId int) {
			defer wg.Done()
			err := r.Remove(ctx, sessionId)
			if err != nil {
				g.Log().Errorf(ctx, "Failed to stop session %d: %v", sessionId, err)
			}
		}(id)
	}
	wg.Wait()
	g.Log().Info(ctx, "All live sessions stopped.")
}
