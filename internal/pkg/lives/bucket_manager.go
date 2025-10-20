package lives

import (
	"context"
	"sync"
	"time"
)

var (
	instance      *TokenBucketManager
	once          sync.Once
	defaultLevels = []Level{
		{LimitPerSec: 5},
		{LimitPerSec: 10},
		{LimitPerSec: 20},
	}
)

func GetBucketManager() *TokenBucketManager {
	once.Do(func() {
		instance = NewTokenBucketManager(defaultLevels)
		instance.Start()
	})
	return instance
}

type TokenBucketManager struct {
	mu        sync.RWMutex
	buckets   map[string]*PlatformBucket
	levels    []Level
	ticker    *time.Ticker
	stopChan  chan struct{}
	startOnce sync.Once
	stopOnce  sync.Once
}

func NewTokenBucketManager(levels []Level) *TokenBucketManager {
	if len(levels) == 0 {
		levels = defaultLevels
	}
	return &TokenBucketManager{
		buckets:  make(map[string]*PlatformBucket),
		levels:   levels,
		stopChan: make(chan struct{}),
	}
}

func (m *TokenBucketManager) Start() {
	m.startOnce.Do(func() {
		m.ticker = time.NewTicker(1 * time.Second)
		go func() {
			for {
				select {
				case now := <-m.ticker.C:
					m.refillAll(now)
				case <-m.stopChan:
					return
				}
			}
		}()
	})
}

func (m *TokenBucketManager) Stop() {
	m.stopOnce.Do(func() {
		close(m.stopChan)
		if m.ticker != nil {
			m.ticker.Stop()
		}
	})
}

func (m *TokenBucketManager) EnsureBucket(platform string) *PlatformBucket {
	m.mu.RLock()
	b := m.buckets[platform]
	m.mu.RUnlock()
	if b != nil {
		return b
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if bb, ok := m.buckets[platform]; ok {
		return bb
	}
	nb := NewPlatformBucket(platform, m.levels)
	m.buckets[platform] = nb
	return nb
}

func (m *TokenBucketManager) TryAcquire(ctx context.Context, platform string) bool {
	return m.TryAcquireBatch(ctx, platform, 1)
}

func (m *TokenBucketManager) TryAcquireBatch(ctx context.Context, platform string, n int) bool {
	b := m.EnsureBucket(platform)
	return b.TryAcquire(n)
}

func (m *TokenBucketManager) Acquire(ctx context.Context, platform string) error {
	return m.AcquireBatch(ctx, platform, 1)
}

func (m *TokenBucketManager) AcquireBatch(ctx context.Context, platform string, n int) error {
	b := m.EnsureBucket(platform)
	return b.Acquire(ctx, n)
}

func (m *TokenBucketManager) Status() map[string]BucketSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make(map[string]BucketSnapshot, len(m.buckets))
	for p, b := range m.buckets {
		out[p] = b.snapshot()
	}
	return out
}

func (m *TokenBucketManager) refillAll(now time.Time) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, b := range m.buckets {
		b.refill(now)
	}
}
