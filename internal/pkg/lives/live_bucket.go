package lives

import (
	"context"
	"sync"
	"time"
)

type PlatformBucket struct {
	mu       sync.Mutex
	platform string

	capacity     int
	tokens       int
	currentUsed  int
	lastRefillAt time.Time
}

func NewPlatformBucket(platform string, capacity int) *PlatformBucket {
	now := time.Now()
	return &PlatformBucket{
		platform:     platform,
		capacity:     capacity,
		tokens:       capacity,
		lastRefillAt: now,
		currentUsed:  0,
	}
}

func (b *PlatformBucket) TryAcquire(n int) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if n <= 0 {
		return true
	}
	if b.tokens >= n {
		b.tokens -= n
		b.currentUsed += n
		return true
	}
	return false
}

func (b *PlatformBucket) Acquire(ctx context.Context, n int) error {
	if b.TryAcquire(n) {
		return nil
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 7*time.Second)
	defer cancel()

	ticker := time.NewTicker(321 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeoutCtx.Done():
			return timeoutCtx.Err()
		case <-ticker.C:
			if b.TryAcquire(n) {
				return nil
			}
		}
	}
}

func (b *PlatformBucket) refill(now time.Time) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.tokens = b.capacity
	b.currentUsed = 0
	b.lastRefillAt = now
}

type BucketSnapshot struct {
	Platform        string    `json:"platform"`
	Capacity        int       `json:"capacity"`
	TokensAvailable int       `json:"tokensAvailable"`
	LastRefillAt    time.Time `json:"lastRefillAt"`
	CurrentUsed     int       `json:"currentUsed"`
}

func (b *PlatformBucket) snapshot() BucketSnapshot {
	b.mu.Lock()
	defer b.mu.Unlock()
	return BucketSnapshot{
		Capacity:        b.capacity,
		CurrentUsed:     b.currentUsed,
		LastRefillAt:    b.lastRefillAt,
		Platform:        b.platform,
		TokensAvailable: b.tokens,
	}
}
