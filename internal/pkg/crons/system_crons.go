package crons

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/shichen437/gowlive/internal/pkg/crons/system"
)

var (
	checkLatestVersion = "checkLatestVersion"
	storageWarning     = "storageWarning"
	anchorInfo         = "anchorInfo"
)

func SystemCron(ctx context.Context) {
	gcron.Add(ctx, "@hourly", func(ctx context.Context) {
		g.Log().Info(ctx, "Add job - "+checkLatestVersion)
		system.CheckVersion(ctx)
	}, checkLatestVersion)
	gcron.Add(ctx, "@hourly", func(ctx context.Context) {
		g.Log().Info(ctx, "Add job - "+storageWarning)
		system.StorageWarning(ctx)
	}, storageWarning)
	// 每天陵城 5 点执行
	gcron.Add(ctx, "# 0 5 * * *", func(ctx context.Context) {
		g.Log().Info(ctx, "Add job - "+anchorInfo)
		AnchorInfoCron(ctx)
	}, anchorInfo)
}
