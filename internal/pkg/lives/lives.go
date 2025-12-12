package lives

import (
	"context"
	"net/url"
	"sync"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/shichen437/gowlive/internal/pkg/utils"
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
		return nil, utils.TfError(context.TODO(), "ext.domain.not.supported", sUrl)
	}
	builder, err := getBuilder(url.Hostname())
	if err != nil {
		return nil, utils.TfError(context.TODO(), "ext.domain.not.supported", url.Hostname())
	}
	liveApi, err = builder.Build(url)
	if err != nil {
		return nil, utils.TfError(context.TODO(), "ext.domain.not.supported", url.Hostname())
	}
	return
}
