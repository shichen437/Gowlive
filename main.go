package main

import (
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"

	"github.com/shichen437/gowlive/internal/cmd"
)

func main() {
	if err := cmd.Migrate(); err != nil {
		g.Log().Fatalf(gctx.GetInitCtx(), "Start Failed: %+v", err)
	}
	cmd.Main.Run(gctx.New())
}
