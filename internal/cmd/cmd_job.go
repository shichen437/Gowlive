package cmd

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shichen437/gowlive/internal/pkg/crons"
)

func JobInit() {
	g.Log().Info(gctx.GetInitCtx(), "Job server initing...")
	crons.SystemCron(gctx.GetInitCtx())
	g.Log().Info(gctx.GetInitCtx(), "Job server init done")
}
