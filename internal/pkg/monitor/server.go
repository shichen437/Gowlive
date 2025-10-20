package monitor

import (
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/sse"
)

func init() {
	sse.RegisterChannelHandler(consts.SSE_CHANNEL_MONITOR, &MonitorHandler{})
	sse.RegisterChannelHandler(consts.SSE_CHANNEL_METRIC, &MetricHandler{})
}
