package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/gowlive/internal/app/system/model"
)

type GetOverviewReq struct {
	g.Meta `path:"/system/overview" method:"get" tags:"系统" summary:"获取系统概览"`
}
type GetOverviewRes struct {
	g.Meta `mime:"application/json"`
	Data   *model.Overview `json:"data"`
}
