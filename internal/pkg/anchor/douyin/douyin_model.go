package douyin

import "net/url"

const (
	domain            = "www.douyin.com"
	platform          = "douyin"
	randomCookieChars = "1234567890abcdef"
	user_info_url     = "https://www.douyin.com/web/api/v2/user/info/?sec_uid="
)

type Douyin struct {
	Url      *url.URL
	Platform string
}

type AnchorInfoResp struct {
	UserInfo struct {
		UniqueID       string `json:"unique_id"`
		Nickname       string `json:"nickname"`
		Signature      string `json:"signature"`
		FollowingCount int    `json:"following_count"`
		FollowerCount  int    `json:"mplatform_followers_count"`
		VideoCount     int    `json:"aweme_count"`
		LikeCount      string `json:"total_favorited"`
	} `json:"user_info"`
}
