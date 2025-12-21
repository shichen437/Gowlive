package bigo

import "net/url"

const (
	domain    = "www.bigo.tv"
	platform  = "bigo"
	anchorApi = "https://ta.bigo.tv/official_website_tiebar/anchor/info"
	regexpUid = `^/?(?:[a-z]{2}(?:-[A-Za-z]{2,})?/)?user/([^/?#]+)`
)

type Bigo struct {
	Url      *url.URL
	Platform string
	Proxy    string
}

type BigoUserProfileResp struct {
	Code int                  `json:"code"`
	Data *BigoUserProfileData `json:"data"`
}

type BigoUserProfileData struct {
	Nickname       string `json:"nickname"`
	Sign           string `json:"introduce"`
	SecId          string `json:"uid"`
	BigoId         string `json:"bigoId"`
	VideoCount     string `json:"postCount"`
	FollowingCount int    `json:"followCnt"`
	FollowerCount  int    `json:"fansCnt"`
}
