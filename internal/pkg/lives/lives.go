package lives

import (
	"net/url"
	"sync"

	"github.com/gogf/gf/v2/errors/gerror"
)

var (
	builders sync.Map
)

type LiveApi interface {
	GetInfo() (*LiveState, error)
}

func Register(domain string, b Builder) {
	builders.Store(domain, b)
}

type Builder interface {
	Build(*url.URL) (LiveApi, error)
}

func getBuilder(domain string) (Builder, error) {
	builder, ok := builders.Load(domain)
	if !ok {
		return nil, gerror.New("unknown domain")
	}
	return builder.(Builder), nil
}

func New(sUrl string) (liveApi LiveApi, err error) {
	url, err := url.Parse(sUrl)
	if err != nil {
		return nil, gerror.New("不支持的域名或平台: " + sUrl)
	}
	builder, err := getBuilder(url.Hostname())
	if err != nil {
		return nil, gerror.New("不支持的域名或平台: " + url.Hostname())
	}
	liveApi, err = builder.Build(url)
	if err != nil {
		return nil, gerror.New("不支持的域名或平台: " + url.Hostname())
	}
	return
}
