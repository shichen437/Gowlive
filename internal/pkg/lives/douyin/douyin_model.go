package douyin

import "net/url"

const (
	domain                     = "live.douyin.com"
	platform                   = "douyin"
	randomCookieChars          = "1234567890abcdef"
	mainInfoLineCatcherRegex   = `self.__pace_f.push\(\[1,\s*"[^:]*:([^<]*,null,\{\\"state\\"[^<]*\])\\n"\]\)`
	commonInfoLineCatcherRegex = `self.__pace_f.push\(\[1,\s*\"(\{.*\})\"\]\)`
)

type Douyin struct {
	Url         *url.URL
	Platform    string
	RespCookies map[string]string
}
