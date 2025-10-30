package cmd

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shichen437/gowlive/internal/pkg/crons"
)

func JobInit() {
	ctx := gctx.GetInitCtx()
	g.Log().Info(ctx, "Job server initing...")
	crons.SystemCron(ctx)
	g.Log().Info(ctx, "Job server init done")
}
