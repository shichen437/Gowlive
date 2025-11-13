package manager

import (
	"context"
	"sync"
)

var (
	cookieManager *CookieManager
	cookieOnce    sync.Once
)

type CookieManager struct {
	mu      sync.RWMutex
	cookies map[string]string
}

func GetCookieManager() *CookieManager {
	cookieOnce.Do(func() {
		cookieManager = &CookieManager{
			cookies: make(map[string]string),
		}
	})
	return cookieManager
}

func (r *CookieManager) Save(ctx context.Context, platform, cookie string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cookies[platform] = cookie
}

func (r *CookieManager) Get(ctx context.Context, platform string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cookie, ok := r.cookies[platform]
	if !ok {
		return ""
	}
	return cookie
}

func (r *CookieManager) Remove(ctx context.Context, platform string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.cookies, platform)
}
