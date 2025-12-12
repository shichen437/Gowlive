package bigo

import "net/url"

const (
	domain        = "www.bigo.tv"
	platform      = "bigo"
	referer       = "https://www.bigo.tv/"
	studioInfoApi = "https://ta.bigo.tv/official_website/studio/getInternalStudioInfo"
	regexpUid     = `(?i)^https?://(?:www\.)?bigo\.tv/(?:[a-z]{2,}/)?([^/]+)$`
)

type Bigo struct {
	Url         *url.URL
	Platform    string
	RespCookies map[string]string
}

type StudioInfo struct {
	Alive     int    `json:"alive"`
	HlsSrc    string `json:"hls_src"`
	Nickname  string `json:"nick_name"`
	RoomId    string `json:"roomId"`
	RoomTopic string `json:"roomTopic"`
	SiteId    string `json:"siteId"`
	Uid       int    `json:"uid"`
}

type StudioInfoResp struct {
	Code int         `json:"code"`
	Data *StudioInfo `json:"data"`
}
