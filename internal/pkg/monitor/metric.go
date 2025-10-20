package monitor

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/metrics"
	"github.com/shichen437/gowlive/internal/pkg/sse"
)

type MetricHandler struct{}

func (m *MetricHandler) OnConnect(channel string) {
	if channel == consts.SSE_CHANNEL_METRIC {
		MetricInfo()
		sse.StartAutoSend(channel, 5*time.Second, MetricInfo)
	}
}

func (m *MetricHandler) OnDisconnect(channel string) {
	if channel == consts.SSE_CHANNEL_METRIC {
		sse.StopAutoSend(channel)
	}
}

func MetricInfo() {
	result := metrics.GetIndicatorManager().SummaryAll(gctx.GetInitCtx())
	if len(result) == 0 {
		return
	}
	msg := sse.GetSseMsgStr(consts.SSE_EVENT_TYPE_METRIC, g.Map{
		"data": result,
	})
	sse.BroadcastMessage(consts.SSE_CHANNEL_METRIC, msg)
}
