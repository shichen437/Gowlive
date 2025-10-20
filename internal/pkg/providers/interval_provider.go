package providers

import (
	"context"

	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/metrics"
	"github.com/shichen437/gowlive/internal/pkg/service"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

func GetIntelligentInterval(ctx context.Context, liveId int, platform string) int {
	flag := service.IsNonLiveHours(ctx, liveId)
	if flag {
		return getNonLiveInterval()
	}
	return GetRegularInterval(ctx, liveId, platform)
}

// 常规间隔，开播中直接获取, 30s-110s
func GetRegularInterval(ctx context.Context, liveId int, platform string) int {
	summary := metrics.GetIndicatorManager().Summary(ctx, platform)
	// 5 分钟指标大于 10%，触发探测间隔
	if summary != (metrics.Summary{}) && summary.TotalErrors > 10 && summary.MainPercent > consts.RequestFailedThreshold {
		return getDetectInterval()
	}
	return utils.RandomSecondsBatesInt(30, 110, 3)
}

// 探测间隔，介于常规间隔和非直播时段间隔之间, 6min-8min
func getDetectInterval() int {
	return utils.RandomSecondsBatesInt(360, 480, 4)
}

// 非直播时段, 14-16min
func getNonLiveInterval() int {
	return utils.RandomSecondsBatesInt(60*14, 60*16, 5)
}
