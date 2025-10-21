package lives

import (
	"context"
	"sync"
	"time"
)

type Level struct {
	LimitPerSec int
}

type PlatformBucket struct {
	mu sync.Mutex

	platform string

	levels       []Level
	currentIdx   int
	capacity     int
	tokens       int
	lastRefillAt time.Time

	consecutiveAtLimit int

	lastSecUsed int
}

func NewPlatformBucket(platform string, levels []Level) *PlatformBucket {
	if len(levels) == 0 {
		levels = defaultLevels
	}
	now := time.Now()
	return &PlatformBucket{
		platform:     platform,
		levels:       levels,
		currentIdx:   0,
		capacity:     levels[0].LimitPerSec,
		tokens:       levels[0].LimitPerSec,
		lastRefillAt: now,
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
		b.lastSecUsed += n
		return true
	}
	return false
}

func (b *PlatformBucket) Acquire(ctx context.Context, n int) error {
	if b.TryAcquire(n) {
		return nil
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ticker := time.NewTicker(200 * time.Millisecond)
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

	atLimit := b.lastSecUsed >= b.capacity
	if atLimit {
		b.consecutiveAtLimit++
	} else {
		b.consecutiveAtLimit = 0
	}
	for b.consecutiveAtLimit >= 5 && b.currentIdx < len(b.levels)-1 {
		b.currentIdx++
		b.capacity = b.levels[b.currentIdx].LimitPerSec
		b.consecutiveAtLimit = 0
		b.tokens = b.capacity
	}

	b.lastSecUsed = 0
	b.lastRefillAt = now
}

type BucketSnapshot struct {
	Platform           string    `json:"platform"`
	CapacityPerSec     int       `json:"capacityPerSec"`
	TokensAvailable    int       `json:"tokensAvailable"`
	CurrentLevelIndex  int       `json:"currentLevelIndex"`
	ConsecutiveAtLimit int       `json:"consecutiveAtLimit"`
	LastRefillAt       time.Time `json:"lastRefillAt"`
	LastSecondUsed     int       `json:"lastSecondUsed"`
	Levels             []int     `json:"levels"`
}

func (b *PlatformBucket) snapshot() BucketSnapshot {
	b.mu.Lock()
	defer b.mu.Unlock()
	levels := make([]int, len(b.levels))
	for i, lv := range b.levels {
		levels[i] = lv.LimitPerSec
	}
	return BucketSnapshot{
		Platform:           b.platform,
		CapacityPerSec:     b.capacity,
		TokensAvailable:    b.tokens,
		CurrentLevelIndex:  b.currentIdx,
		ConsecutiveAtLimit: b.consecutiveAtLimit,
		LastRefillAt:       b.lastRefillAt,
		LastSecondUsed:     b.lastSecUsed,
		Levels:             levels,
	}
}
