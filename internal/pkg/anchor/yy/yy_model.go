package yy

import "net/url"

type YY struct {
	Url      *url.URL
	Platform string
}

const (
	domain         = "www.yy.com"
	platform       = "douyin"
	video_data_url = "https://webuser.yy.com/u/videos/data/"
	regexpProps    = `window\.pageProps\s*=\s*([\s\S]*?)</script>`
)
