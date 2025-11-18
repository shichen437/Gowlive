package stream_parser

import (
	"context"
	"slices"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/shichen437/gowlive/internal/pkg/lives"
)

var AssistFormatArray = []string{"flv", "mp4", "mp3", "mkv", "ts"}

type Builder interface {
	Build(pm map[string]string) (Parser, error)
}

type Parser interface {
	ParseLiveStream(ctx context.Context, streamInfo *lives.StreamUrlInfo, file string) error
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
	if _, ok := pm["format"]; !ok || !slices.Contains(AssistFormatArray, pm["format"]) {
		return nil, gerror.New("不支持的录制格式")
	}
	builder, ok := m[name]
	if !ok {
		return nil, gerror.New("未知解析器")
	}
	return builder.Build(pm)
}
