package v1

import "github.com/gogf/gf/v2/frame/g"

type GetOpenlistStatusReq struct {
	g.Meta `path:"/third/openlist/status" method:"get" tags:"Openlist" summary:"获取openlist状态"`
}
type GetOpenlistStatusRes struct {
	g.Meta `mime:"application/json"`
	Status int `json:"status"`
}
