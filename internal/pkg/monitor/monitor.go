package monitor

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/monitor/entity"
	"github.com/shichen437/gowlive/internal/pkg/sse"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type MonitorHandler struct{}

func (h *MonitorHandler) OnConnect(channel string) {
	if channel == consts.SSE_CHANNEL_MONITOR {
		HardwareInfo()
		sse.StartAutoSend(channel, 3*time.Second, HardwareInfo)
	}
}

func (h *MonitorHandler) OnDisconnect(channel string) {
	if channel == consts.SSE_CHANNEL_MONITOR {
		sse.StopAutoSend(channel)
	}
}

func HardwareInfo() {
	cpuInfo, _ := cpu.Info()
	var cpuData *entity.CpuInfo
	if len(cpuInfo) > 0 {
		gconv.Struct(cpuInfo[0], &cpuData)
		value, err := gcache.Get(gctx.GetInitCtx(), consts.CpuPercentCacheKey)
		if err == nil {
			cpuData.Percent = value.Float64()
		}
	}
	var memData *entity.MemoryInfo
	memInfo, _ := mem.VirtualMemory()
	if memInfo != nil {
		gconv.Struct(memInfo, &memData)
	}
	var diskData *entity.DiskInfo
	diskInfo, _ := disk.Usage("/")
	if diskInfo != nil {
		gconv.Struct(diskInfo, &diskData)
	}
	msg := sse.GetSseMsgStr(consts.SSE_EVENT_TYPE_MONITOR, g.Map{
		"cpu":  cpuData,
		"mem":  memData,
		"disk": diskData,
	})
	sse.BroadcastMessage(consts.SSE_CHANNEL_MONITOR, msg)
}
