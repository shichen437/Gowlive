package bilibili

import "net/url"

const (
	domain       = "live.bilibili.com"
	platform     = "bilibili"
	roomInitUrl  = "https://api.live.bilibili.com/room/v1/Room/room_init"
	roomApiUrl   = "https://api.live.bilibili.com/room/v1/Room/get_info"
	userApiUrl   = "https://api.live.bilibili.com/live_user/v1/UserInfo/get_anchor_in_room"
	liveApiUrlv2 = "https://api.live.bilibili.com/xlive/web-room/v2/index/getRoomPlayInfo"
	parseAddr    = "data.playurl_info.playurl.stream.0.format.0.codec.0"
)

type Bilibili struct {
	Url         *url.URL
	Platform    string
	RoomID      string
	UserAgent   string
	RespCookies map[string]string
}
