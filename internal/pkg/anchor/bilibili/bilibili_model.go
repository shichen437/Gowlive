package bilibili

import (
	"net/url"

	"github.com/gogf/gf/v2/frame/g"
)

const (
	domain   = "space.bilibili.com"
	platform = "bilibili"

	userProfileInfoUrl = "https://api.bilibili.com/x/web-interface/card?mid="
)

var (
	bilibiliHeaders = g.MapStrStr{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
		"Referer":         "https://space.bilibili.com/",
		"Origin":          "https://www.bilibili.com",
		"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
	}
)

type Bilibili struct {
	Url      *url.URL
	Platform string
}
