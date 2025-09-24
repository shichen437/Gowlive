package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/gowlive/api/v1/common"
	"github.com/shichen437/gowlive/internal/app/system/model/entity"
)

type GetLogsListReq struct {
	g.Meta `path:"/system/logs/list" method:"get" tags:"日志管理" summary:"获取日志列表"`
	common.PageReq
	Type   *int   `json:"type" v:"in:1,2#日志类型错误"`
	Status *int   `json:"status" v:"in:0,1#日志状态错误"`
	Sort   string `json:"sort"`
}

type GetLogsListRes struct {
	g.Meta `mime:"application/json"`
	Rows   []*entity.SysLogs `json:"rows"`
	Total  int               `json:"total"`
}

type DeleteLogsReq struct {
	g.Meta `path:"/system/logs/{id}" method:"delete" tags:"日志管理" summary:"删除日志"`
	Id     *int `json:"" v:"required#日志ID不能为空"`
}
type DeleteLogsRes struct {
	g.Meta `mime:"application/json"`
}

type DeleteAllLogsReq struct {
	g.Meta `path:"/system/logs/all" method:"delete" tags:"日志管理" summary:"删除全部日志"`
}

type DeleteAllLogsRes struct {
	g.Meta `mime:"application/json"`
}
