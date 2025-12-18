package yy

import "net/url"

const (
	domain       = "www.yy.com"
	platform     = "yy"
	referer      = "https://www.yy.com/"
	streamApi    = "https://stream-manager.yy.com/v3/channel/streams?"
	detailApi    = "https://www.yy.com/live/detail?"
	regexpAnchor = `nick:\s*"(.*?)",\s*logo`
	regexpCid    = `sid\s*:\s*"(.*?)",\s*ssid`
)

type YY struct {
	Url         *url.URL
	Platform    string
	UserAgent   string
	RespCookies map[string]string
}
