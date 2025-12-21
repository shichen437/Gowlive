package manager

import (
	"context"
	"math/rand"
	"strings"
	"sync"
)

var (
	proxyManager *ProxyManager
	proxyOnce    sync.Once
)

type ProxyManager struct {
	mu      sync.RWMutex
	proxies map[string][]string
}

func GetProxyManager() *ProxyManager {
	proxyOnce.Do(func() {
		proxyManager = &ProxyManager{
			proxies: make(map[string][]string),
		}
	})
	return proxyManager
}

func (r *ProxyManager) SaveProxy(ctx context.Context, platform, proxy string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	proxies := strings.Split(proxy, ",")
	r.proxies[platform] = proxies
}

func (r *ProxyManager) GetProxies(platform string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	proxies, ok := r.proxies[platform]
	if !ok {
		return nil
	}
	return proxies
}

func (r *ProxyManager) GetRandomProxy(platform string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	proxies, ok := r.proxies[platform]
	if !ok || len(proxies) == 0 {
		return ""
	}
	if len(proxies) == 1 {
		return proxies[0]
	}
	return proxies[rand.Intn(len(proxies))]
}

func (r *ProxyManager) RemoveProxy(ctx context.Context, platform string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.proxies, platform)
}
