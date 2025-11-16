package lives

import (
	"net/url"
	"sync"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shichen437/gowlive/internal/pkg/interfaces"
)

type LiveSession struct {
	Id              int
	Config          LiveConfig
	State           LiveState
	StartedAt       *gtime.Time
	LiveApi         LiveApi
	EventDispatcher interfaces.Module
	ListenerManager interfaces.Module
	RecorderManager interfaces.Module
	mu              sync.RWMutex
}

func NewLiveSession(id int, config LiveConfig, liveAPI LiveApi) *LiveSession {
	return &LiveSession{
		Id:      id,
		Config:  config,
		LiveApi: liveAPI,
	}
}

func (s *LiveSession) GetState() LiveState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.State
}

func (s *LiveSession) UpdateState(newState LiveState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.State.Anchor = newState.Anchor
	s.State.RoomName = newState.RoomName
	s.State.IsLive = newState.IsLive
	s.State.StreamInfos = newState.StreamInfos
}

type LiveConfig struct {
	Id             int
	RoomUrl        string
	Interval       int
	Format         string
	MonitorType    int
	MonitorStartAt string
	MonitorStopAt  string
	Quality        int
	SegmentTime    int
}

type LiveState struct {
	LiveId      int
	RoomName    string
	Anchor      string
	Platform    string
	IsLive      bool
	StreamInfos []*StreamUrlInfo
}

type StreamUrlInfo struct {
	Url                  *url.URL
	Name                 string
	Description          string
	Resolution           int
	Vbitrate             int
	HeadersForDownloader map[string]string
}
