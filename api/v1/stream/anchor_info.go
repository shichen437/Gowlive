package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/gowlive/api/v1/common"
	"github.com/shichen437/gowlive/internal/app/stream/model/entity"
)

type GetAnchorListReq struct {
	g.Meta `path:"/live/anchor/list" method:"get" tags:"主播数据" summary:"获取主播数据列表"`
	common.PageReq
	Platform string `json:"platform" dc:"平台"`
	Nickname string `json:"nickname" dc:"主播昵称"`
}
type GetAnchorListRes struct {
	g.Meta `mime:"application/json"`
	Rows   []*entity.AnchorInfo `json:"rows"`
	Total  int                  `json:"total"`
}

type PostAnchorReq struct {
	g.Meta `path:"/live/anchor" method:"post" tags:"主播数据" summary:"新增主播"`
	Url    string `json:"url" dc:"主播URL" v:"required#主播URL不能为空"`
}
type PostAnchorRes struct {
	g.Meta `mime:"application/json"`
}

type DeleteAnchorReq struct {
	g.Meta `path:"/live/anchor/{id}" method:"delete" tags:"主播数据" summary:"删除主播"`
	Id     *int `json:"id" dc:"主播ID" v:"required#主播ID不能为空"`
}
type DeleteAnchorRes struct {
	g.Meta `mime:"application/json"`
}
