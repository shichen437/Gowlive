package douyin

import "net/url"

const (
	domain                     = "live.douyin.com"
	platform                   = "douyin"
	randomCookieChars          = "1234567890abcdef"
	mainInfoLineCatcherRegex   = `self.__pace_f.push\(\[1,\s*"[^:]*:([^<]*,null,\{\\"state\\"[^<]*\])\\n"\]\)`
	commonInfoLineCatcherRegex = `self.__pace_f.push\(\[1,\s*\"(\{.*\})\"\]\)`
)

var (
	cameraQualities = []qualityConf{
		{key: "FULL_HD1", name: "uhd"},
		{key: "HD1", name: "hd"},
		{key: "SD1", name: "sd"},
		{key: "SD2", name: "md"},
	}
)

type Douyin struct {
	Url         *url.URL
	Platform    string
	UserAgent   string
	RespCookies map[string]string
}

type qualityConf struct {
	key  string
	name string
}
