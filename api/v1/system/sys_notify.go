package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/gowlive/api/v1/common"
	"github.com/shichen437/gowlive/internal/app/system/model/entity"
)

type GetNotifyListReq struct {
	g.Meta `path:"/system/notify/list" method:"get" tags:"通知管理" summary:"获取通知列表"`
	common.PageReq
	Status *int `json:"status"`
}

type GetNotifyListRes struct {
	g.Meta `mime:"application/json"`
	Rows   []*entity.SysNotify `json:"rows"`
	Total  int                 `json:"total"`
}

type PutMarkNotifyReadReq struct {
	g.Meta `path:"/system/notify/{id}" method:"put" tags:"通知管理" summary:"标记已读"`
	Id     *int `json:"id" v:"required#通知ID不能为空"`
}
type PutMarkNotifyReadRes struct {
	g.Meta `mime:"application/json"`
}

type PutMarkNotifyAllReadReq struct {
	g.Meta `path:"/system/notify/all" method:"put" tags:"通知管理" summary:"全部标记已读"`
}
type PutMarkNotifyAllReadRes struct {
	g.Meta `mime:"application/json"`
}

type DeleteNotifyReq struct {
	g.Meta `path:"/system/notify/{id}" method:"delete" tags:"通知管理" summary:"删除通知"`
	Id     *int `json:"id" v:"required#通知ID不能为空"`
}
type DeleteNotifyRes struct {
	g.Meta `mime:"application/json"`
}

type DeleteAllNotifyReq struct {
	g.Meta `path:"/system/notify/all" method:"delete" tags:"通知管理" summary:"删除通知"`
}
type DeleteAllNotifyRes struct {
	g.Meta `mime:"application/json"`
}
