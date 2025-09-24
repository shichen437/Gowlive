package stream_parser

import (
	"context"
	"errors"

	"github.com/shichen437/gowlive/internal/pkg/lives"
)

type Builder interface {
	Build(pm map[string]string) (Parser, error)
}

type Parser interface {
	ParseLiveStream(ctx context.Context, streamInfo *lives.StreamUrlInfo, file, refer string) error
	Stop() error
}

type StatusParser interface {
	Parser
	Status() (map[string]string, error)
}

var m = make(map[string]Builder)

func Register(name string, b Builder) {
	m[name] = b
}

func New(name string, pm map[string]string) (Parser, error) {
	builder, ok := m[name]
	if !ok {
		return nil, errors.New("unknown parser")
	}
	return builder.Build(pm)
}
