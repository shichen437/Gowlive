package gotify

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/gowlive/internal/app/system/model"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/manager"
	mp "github.com/shichen437/gowlive/internal/pkg/message_push"
)

const (
	channelType = "gotify"
)

func init() {
	mp.Register(channelType, &builder{})
}

type builder struct{}
type MessagePush struct{}

func (b *builder) Build(channelType string) (mp.MessagePush, error) {
	return &MessagePush{}, nil
}

func (p *MessagePush) PushMessage(ctx context.Context, channel *model.PushChannel, model *mp.MessageModel) (err error) {
	return gotify(ctx, channel.Url, model.Title, model.Content)
}

func gotify(ctx context.Context, url string, title, message string) (err error) {
	c := g.Client()
	data := g.Map{
		"title":   title,
		"message": message,
	}
	_, err = c.Post(context.Background(), url, data)
	if err != nil {
		g.Log().Errorf(ctx, "Gotify推送失败, 错误信息: %v", err)
		manager.GetLogManager().AddErrorLog(consts.LogTypePush, "Gotify推送失败")
	}
	return err
}
