package bilibili

const (
	domain       = "live.bilibili.com"
	platform     = "bilibili"
	userAgent    = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"
	roomInitUrl  = "https://api.live.bilibili.com/room/v1/Room/room_init"
	roomApiUrl   = "https://api.live.bilibili.com/room/v1/Room/get_info"
	userApiUrl   = "https://api.live.bilibili.com/live_user/v1/UserInfo/get_anchor_in_room"
	liveApiUrlv2 = "https://api.live.bilibili.com/xlive/web-room/v2/index/getRoomPlayInfo"
	parseAddr    = "data.playurl_info.playurl.stream.0.format.0.codec.0"
)
