package system

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/os/gcache"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shirou/gopsutil/v3/cpu"
)

func CpuPercent(ctx context.Context) {
	cpuInfo, _ := cpu.Info()
	percents, _ := cpu.Percent(5*time.Second, false)
	if len(cpuInfo) > 0 {
		if len(percents) > 0 {
			gcache.Set(ctx, consts.CpuPercentCacheKey, percents[0], 10*time.Second)
		}
	}
}
