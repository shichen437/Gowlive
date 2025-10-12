package anchor

import (
	"context"
	"net/url"
	"sync"

	"github.com/gogf/gf/v2/errors/gerror"
)

var (
	builders sync.Map
)

type AnchorApi interface {
	ParseAnchorInfo(ctx context.Context) (*AnchorInfo, error)
}

func Register(domain string, b Builder) {
	builders.Store(domain, b)
}

type Builder interface {
	Build(*url.URL) (AnchorApi, error)
}

func getBuilder(domain string) (Builder, error) {
	builder, ok := builders.Load(domain)
	if !ok {
		return nil, gerror.Newf("不支持的域名或平台: %s", domain)
	}
	return builder.(Builder), nil
}

func New(sUrl string) (anchorApi AnchorApi, err error) {
	u, err := url.Parse(sUrl)
	if err != nil {
		return nil, gerror.Wrapf(err, "URL格式不正确: %s", sUrl)
	}
	builder, err := getBuilder(u.Hostname())
	if err != nil {
		return nil, err
	}
	anchorApi, err = builder.Build(u)
	if err != nil {
		return nil, gerror.Wrapf(err, "构建主播实例失败: %s", sUrl)
	}
	return
}
