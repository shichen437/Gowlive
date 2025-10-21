package monitor

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/metrics"
	"github.com/shichen437/gowlive/internal/pkg/sse"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

type HealthHandler struct{}

func (m *HealthHandler) OnConnect(channel string) {
	if channel == consts.SSE_CHANNEL_HEALTH {
		HealthInfo()
		sse.StartAutoSend(channel, 20*time.Second, HealthInfo)
	}
}

func (m *HealthHandler) OnDisconnect(channel string) {
	if channel == consts.SSE_CHANNEL_HEALTH {
		sse.StopAutoSend(channel)
	}
}

func HealthInfo() {
	diskUsage := utils.GetDiskUsage()
	result := metrics.GetIndicatorManager().SummaryAll(gctx.GetInitCtx())
	percent := 0
	num := 0
	for _, v := range result {
		if v.TotalPercent > 0 {
			percent += int(v.TotalPercent)
			num++
		}
	}
	if num > 0 {
		percent /= num
	}
	msg := sse.GetSseMsgStr(consts.SSE_EVENT_TYPE_HEALTH, g.Map{
		"errorPercent": percent,
		"diskUsage":    diskUsage,
	})
	sse.BroadcastMessage(consts.SSE_CHANNEL_HEALTH, msg)
}
