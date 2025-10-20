package metrics

import (
	"context"
	"maps"
	"math"
	"sync"
	"time"

	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/manager"
	"github.com/shichen437/gowlive/internal/pkg/message_push"
)

var (
	instance *IndicatorManager
	once     sync.Once
)

func GetIndicatorManager() *IndicatorManager {
	once.Do(func() {
		instance = NewIndicatorManager()
	})
	return instance
}

type bucket struct {
	startSec int64
	req      int64
	err      int64
	mreq     int64
	merr     int64
}

type platformMetrics struct {
	mu sync.Mutex

	bucketDur   time.Duration
	windowDur   time.Duration
	bucketCount int

	// 环形缓冲
	buckets []bucket
	headIdx int
	headSec int64

	// 滚动总和（窗口内所有桶的和）
	totalReq int64
	totalErr int64
	mainReq  int64
	mainErr  int64

	lastGen time.Time
}

type IndicatorManager struct {
	bucketDur time.Duration
	windowDur time.Duration
	bCount    int

	pmMu      sync.RWMutex
	platforms map[string]*platformMetrics
}

type Summary struct {
	Platform      string    `json:"platform"`
	TotalRequests int64     `json:"totalRequests"`
	TotalErrors   int64     `json:"totalErrors"`
	TotalPercent  float64   `json:"totalPercent"`
	MainRequests  int64     `json:"mainRequests"`
	MainErrors    int64     `json:"mainErrors"`
	MainPercent   float64   `json:"mainPercent"`
	GeneratedAt   time.Time `json:"generatedAt"`
}

func NewIndicatorManager() *IndicatorManager {
	bd := 5 * time.Second
	wd := 5 * time.Minute
	bc := int(wd / bd)
	return &IndicatorManager{
		bucketDur: bd,
		windowDur: wd,
		bCount:    bc,
		platforms: make(map[string]*platformMetrics),
	}
}

func (m *IndicatorManager) Record(ctx context.Context, platform string, result bool, isMain bool) {
	pm := m.ensurePlatform(platform)
	now := time.Now()
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.advance(now)

	b := &pm.buckets[pm.headIdx]
	b.req++
	pm.totalReq++
	if !result {
		b.err++
		pm.totalErr++
		if pm.totalErr > 10 && calcPercent(pm.totalReq, pm.totalErr) >= consts.RequestFailedThreshold {
			go func() {
				nCtx := gctx.GetInitCtx()
				val, err := gcache.Get(nCtx, consts.RequestsFailedWarningCKP+platform)
				if err != nil || val == nil {
					return
				}
				valInt := val.Int()
				if valInt == 0 {
					gcache.Set(nCtx, consts.RequestsFailedWarningCKP+platform, 1, consts.RequestsFailedWarningTtl)
					content := "直播平台【" + platform + "】近五分钟内，已超过 5% 的请求失败，请及时处理！"
					message_push.PushMessage(nCtx, &message_push.MessageModel{
						Title:   "直播监控告警",
						Content: content,
					})
					manager.GetNotfiyManager().AddWarningNotify("直播监控告警", content)
				}
			}()
		}
	}
	if isMain {
		b.mreq++
		pm.mainReq++
		if !result {
			b.merr++
			pm.mainErr++
		}
	}
}

func (m *IndicatorManager) Summary(ctx context.Context, platform string) Summary {
	pm := m.getPlatform(platform)
	if pm == nil {
		return Summary{}
	}
	now := time.Now()

	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.advance(now)

	totalReq := pm.totalReq
	totalErr := pm.totalErr
	mainReq := pm.mainReq
	mainErr := pm.mainErr

	return Summary{
		Platform:      platform,
		TotalRequests: totalReq,
		TotalErrors:   totalErr,
		TotalPercent:  calcPercent(totalReq, totalErr),
		MainRequests:  mainReq,
		MainErrors:    mainErr,
		MainPercent:   calcPercent(mainReq, mainErr),
		GeneratedAt:   now,
	}
}

func (m *IndicatorManager) SummaryAll(ctx context.Context) map[string]Summary {
	now := time.Now()

	m.pmMu.RLock()
	plats := make(map[string]*platformMetrics, len(m.platforms))
	maps.Copy(plats, m.platforms)
	m.pmMu.RUnlock()

	out := make(map[string]Summary, len(plats))
	for platform, pm := range plats {
		pm.mu.Lock()
		pm.advance(now)
		s := Summary{
			Platform:      platform,
			TotalRequests: pm.totalReq,
			TotalErrors:   pm.totalErr,
			TotalPercent:  calcPercent(pm.totalReq, pm.totalErr),
			MainRequests:  pm.mainReq,
			MainErrors:    pm.mainErr,
			MainPercent:   calcPercent(pm.mainReq, pm.mainErr),
			GeneratedAt:   now,
		}
		pm.mu.Unlock()

		if s.TotalRequests == 0 {
			continue
		}
		out[platform] = s
	}
	return out
}

func (m *IndicatorManager) WithWindow(window, bucket time.Duration) *IndicatorManager {
	if window > 0 {
		m.windowDur = window
	}
	if bucket > 0 {
		m.bucketDur = bucket
	}
	m.bCount = int(m.windowDur / m.bucketDur)
	return m
}

func (m *IndicatorManager) ensurePlatform(platform string) *platformMetrics {
	m.pmMu.RLock()
	if pm, ok := m.platforms[platform]; ok {
		m.pmMu.RUnlock()
		return pm
	}
	m.pmMu.RUnlock()

	m.pmMu.Lock()
	defer m.pmMu.Unlock()
	if pm, ok := m.platforms[platform]; ok {
		return pm
	}

	pm := &platformMetrics{
		bucketDur:   m.bucketDur,
		windowDur:   m.windowDur,
		bucketCount: m.bCount,
		buckets:     make([]bucket, m.bCount),
		headIdx:     0,
		headSec:     alignSec(time.Now(), m.bucketDur),
	}
	for i := 0; i < pm.bucketCount; i++ {
		pm.buckets[i].startSec = pm.headSec + int64(i*int(pm.bucketDur/time.Second))
	}
	m.platforms[platform] = pm
	return pm
}

func (m *IndicatorManager) getPlatform(platform string) *platformMetrics {
	m.pmMu.RLock()
	defer m.pmMu.RUnlock()
	return m.platforms[platform]
}

func (pm *platformMetrics) advance(now time.Time) {
	curSec := alignSec(now, pm.bucketDur)

	if curSec == pm.headSec {
		return
	}
	shift := int((curSec - pm.headSec) / int64(pm.bucketDur/time.Second))
	if shift <= 0 {
		return
	}
	if shift >= pm.bucketCount {
		for i := 0; i < pm.bucketCount; i++ {
			b := &pm.buckets[i]
			pm.totalReq -= b.req
			pm.totalErr -= b.err
			pm.mainReq -= b.mreq
			pm.mainErr -= b.merr

			b.req, b.err, b.mreq, b.merr = 0, 0, 0, 0
			b.startSec = curSec + int64(i*int(pm.bucketDur/time.Second))
		}
		pm.headIdx = 0
		pm.headSec = curSec
		return
	}

	for range shift {
		nextIdx := (pm.headIdx + 1) % pm.bucketCount
		nb := &pm.buckets[nextIdx]
		pm.totalReq -= nb.req
		pm.totalErr -= nb.err
		pm.mainReq -= nb.mreq
		pm.mainErr -= nb.merr

		nb.req, nb.err, nb.mreq, nb.merr = 0, 0, 0, 0
		nb.startSec = pm.headSec + int64(int(pm.bucketDur/time.Second))

		pm.headIdx = nextIdx
		pm.headSec = nb.startSec
	}
	pm.headSec = curSec
}

func alignSec(t time.Time, bucketDur time.Duration) int64 {
	sec := t.Unix()
	bs := int64(bucketDur / time.Second)
	if bs <= 0 {
		bs = 1
	}
	return (sec / bs) * bs
}

func calcPercent(total, err int64) float64 {
	if total <= 0 {
		return 0
	}
	p := float64(err) / float64(total) * 100.0
	return math.Round(p*100) / 100.0
}
