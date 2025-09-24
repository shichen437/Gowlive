package v1

import "github.com/gogf/gf/v2/frame/g"

type PostMediaParseReq struct {
	g.Meta `path:"/media/parse" method:"post" tags:"媒体解析" summary:"媒体解析"`
	Url    string `json:"url" v:"required#媒体链接不能为空"`
}

type PostMediaParseRes struct {
	g.Meta `mime:"application/json"`
}
