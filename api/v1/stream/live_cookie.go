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
	Platform string `json:"platform" dc:"直播平台" v:"required#stream.cookie.valid.PlatformRequired"`
	Cookie   string `json:"cookie" dc:"直播cookie" v:"required#stream.cookie.valid.CookieRequired"`
	Remark   string `json:"remark" dc:"备注" v:"max-length:45#stream.cookie.valid.RemarkMaxLength"`
}
type PostLiveCookieRes struct {
	g.Meta `mime:"application/json"`
}

type PutLiveCookieReq struct {
	g.Meta `path:"/live/cookie" method:"put" tags:"直播管理" summary:"修改直播cookie"`
	Id     int    `json:"id" dc:"直播cookie id" v:"required#stream.cookie.valid.IDRequired"`
	Cookie string `json:"cookie" dc:"直播cookie" v:"required#stream.cookie.valid.CookieRequired"`
	Remark string `json:"remark" dc:"备注" v:"max-length:45#stream.cookie.valid.RemarkMaxLength"`
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
