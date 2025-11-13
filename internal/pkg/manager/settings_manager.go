package manager

import (
	"sync"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shichen437/gowlive/internal/pkg/service"
)

var (
	settingsManager *SettingsManager
	sOnce           sync.Once
)

type SettingsManager struct {
	mu       sync.RWMutex
	settings map[string]int
}

func GetSettingsManager() *SettingsManager {
	sOnce.Do(func() {
		settingsManager = &SettingsManager{
			mu:       sync.RWMutex{},
			settings: make(map[string]int),
		}
	})
	return settingsManager
}

func (s *SettingsManager) SaveSetting(key string, value int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.settings[key] = value
}

func (s *SettingsManager) GetSetting(key string) int {
	s.mu.RLock()

	sVal, ok := s.settings[key]
	if ok {
		s.mu.RUnlock()
		return sVal
	}
	s.mu.RUnlock()
	sVal = service.GetSettings(gctx.GetInitCtx(), key)
	s.SaveSetting(key, sVal)
	return sVal
}
