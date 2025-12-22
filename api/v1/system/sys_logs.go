package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/gowlive/api/v1/common"
	"github.com/shichen437/gowlive/internal/app/system/model/entity"
)

type GetLogsListReq struct {
	g.Meta `path:"/system/logs/list" method:"get" tags:"日志管理" summary:"获取日志列表"`
	common.PageReq
	Type   *int   `json:"type" v:"in:1,2#system.logs.valid.TypeInvalid"`
	Status *int   `json:"status" v:"in:0,1#system.logs.valid.StatusInvalid"`
	Sort   string `json:"sort"`
}

type GetLogsListRes struct {
	g.Meta `mime:"application/json"`
	Rows   []*entity.SysLogs `json:"rows"`
	Total  int               `json:"total"`
}

type DeleteLogsReq struct {
	g.Meta `path:"/system/logs/{id}" method:"delete" tags:"日志管理" summary:"删除日志"`
	Id     *int `json:"" v:"required#system.logs.valid.IDBlank"`
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

type GetTerminalLogsListReq struct {
	g.Meta `path:"/system/logs/terminal" method:"get" tags:"日志管理" summary:"获取终端日志"`
	Since  int64 `json:"since"`
	Limit  int   `json:"limit" v:"required|min:1|max:100"`
}

type GetTerminalLogsListRes struct {
	g.Meta `mime:"application/json"`
	Rows   []*LogItem `json:"rows"`
	Next   int64      `json:"next"`
}

type LogItem struct {
	Time  int64  `json:"time"`
	Level string `json:"level"`
	Msg   string `json:"msg"`
}
