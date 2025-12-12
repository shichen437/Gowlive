package anchor

import (
	"context"
	"net/url"
	"sync"

	"github.com/shichen437/gowlive/internal/pkg/utils"
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
		return nil, utils.TfError(context.TODO(), "ext.domain.not.supported", domain)
	}
	return builder.(Builder), nil
}

func New(sUrl string) (anchorApi AnchorApi, err error) {
	u, err := url.Parse(sUrl)
	if err != nil {
		return nil, utils.TfError(context.TODO(), "ext.anchor.domain.invalid", sUrl)
	}
	builder, err := getBuilder(u.Hostname())
	if err != nil {
		return nil, err
	}
	anchorApi, err = builder.Build(u)
	if err != nil {
		return nil, utils.TfError(context.TODO(), "ext.anchor.api.build.failed", sUrl)
	}
	return
}
