package gotify

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/gowlive/internal/app/system/model"
	mp "github.com/shichen437/gowlive/internal/pkg/message_push"
)

const (
	channelType = "gotify"
)

func init() {
	mp.Register(channelType, &builder{})
}

type builder struct{}

func (b *builder) Build(channelType string, anchor string) (mp.MessagePush, error) {
	return &MessagePush{
		Anchor: anchor,
	}, nil
}

type MessagePush struct {
	Anchor string
}

func (p *MessagePush) Push(ctx context.Context, channel *model.PushChannel) (err error) {
	return gotify(channel.Url, "开播通知", "你关注的主播 【 "+p.Anchor+" 】开播了！")
}

func (p *MessagePush) PushMessage(ctx context.Context, channel *model.PushChannel, model *mp.MessageModel) (err error) {
	return gotify(channel.Url, model.Title, model.Content)
}

func gotify(url string, title, message string) (err error) {
	c := g.Client()
	data := g.Map{
		"title":   title,
		"message": message,
	}
	_, err = c.Post(context.Background(), url, data)
	return err
}
