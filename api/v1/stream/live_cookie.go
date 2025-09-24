package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/gowlive/internal/app/stream/model/entity"
)

type GetAllCookieReq struct {
	g.Meta `path:"/live/cookie/list" method:"get" tags:"直播管理" summary:"获取所有直播cookie"`
}
type GetAllCookieRes struct {
	g.Meta `mime:"application/json"`
	Rows   []*entity.LiveCookie `json:"rows"`
}

type PostLiveCookieReq struct {
	g.Meta   `path:"/live/cookie" method:"post" tags:"直播管理" summary:"添加直播cookie"`
	Platform string `json:"platform" dc:"直播平台" v:"required#直播平台不能为空"`
	Cookie   string `json:"cookie" dc:"直播cookie" v:"required#直播cookie不能为空"`
	Remark   string `json:"remark" dc:"备注" v:"max-length:45#备注最大长度为 45"`
}
type PostLiveCookieRes struct {
	g.Meta `mime:"application/json"`
}

type PutLiveCookieReq struct {
	g.Meta `path:"/live/cookie" method:"put" tags:"直播管理" summary:"修改直播cookie"`
	Id     int    `json:"id" dc:"直播cookie id" v:"required#直播cookie id不能为空"`
	Cookie string `json:"cookie" dc:"直播cookie" v:"required#直播cookie不能为空"`
	Remark string `json:"remark" dc:"备注" v:"max-length:45#备注最大长度为 45"`
}
type PutLiveCookieRes struct {
	g.Meta `mime:"application/json"`
}

type DeleteLiveCookieReq struct {
	g.Meta `path:"/live/cookie/{id}" method:"delete" tags:"直播管理" summary:"删除直播cookie"`
	Id     int `json:"id" dc:"直播cookie id"`
}
type DeleteLiveCookieRes struct {
	g.Meta `mime:"application/json"`
}
