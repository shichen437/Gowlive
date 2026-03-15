package webhook

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/shichen437/gowlive/internal/app/system/model"
	mp "github.com/shichen437/gowlive/internal/pkg/message_push"
)

const (
	channelType = "custom"
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
	custom := channel.Custom
	if custom == nil {
		return gerror.New("未找到自定义配置信息")
	}
	c := g.Client()
	if custom.RequestHeaders != "" {
		c.SetHeaderMap(gconv.MapStrStr(custom.RequestHeaders))
	}
	url := custom.WebhookUrl
	url = strings.ReplaceAll(url, "{{msgTitle}}", model.Title)
	url = strings.ReplaceAll(url, "{{msgContent}}", model.Content)
	if custom.RequestMethod == 0 {
		_, err = c.Get(ctx, url)
	}
	if custom.RequestMethod == 1 {
		c.SetContentType("application/json;charset=utf-8")
		body := custom.RequestBody
		if body == "" {
			body = "{}"
		} else {
			body = strings.ReplaceAll(body, "{{msgTitle}}", model.Title)
			body = strings.ReplaceAll(body, "{{msgContent}}", model.Content)
		}
		_, err = c.Post(ctx, url, body)
	}
	if err != nil {
		g.Log().Error(ctx, err)
	}
	return err
}
