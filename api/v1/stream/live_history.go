package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/gowlive/api/v1/common"
	"github.com/shichen437/gowlive/internal/app/stream/model/entity"
)

type GetLiveHistoryListReq struct {
	g.Meta `path:"/live/history/list" method:"get" tags:"直播管理" summary:"获取直播历史列表"`
	common.PageReq
	LiveId *int `json:"liveId" dc:"直播id"`
}
type GetLiveHistoryListRes struct {
	g.Meta `mime:"application/json"`
	Rows   []*entity.LiveHistory `json:"rows"`
	Total  int                   `json:"total"`
}

type DeleteLiveHistoryReq struct {
	g.Meta `path:"/live/history/{id}" method:"delete" tags:"直播管理" summary:"删除直播历史"`
	Id     int `json:"id" dc:"直播历史id"`
}
type DeleteLiveHistoryRes struct {
	g.Meta `mime:"application/json"`
}

type DeleteAllHistoryReq struct {
	g.Meta `path:"/live/history/all" method:"delete" tags:"直播管理" summary:"清空直播历史"`
}
type DeleteAllHistoryRes struct {
	g.Meta `mime:"application/json"`
}
