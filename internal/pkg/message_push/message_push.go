package message_push

import (
	"context"
	"fmt"
	"sync"

	"github.com/shichen437/gowlive/internal/app/system/model"
	"github.com/shichen437/gowlive/internal/pkg/manager"
	"github.com/shichen437/gowlive/internal/pkg/service"

	"github.com/gogf/gf/v2/errors/gerror"
)

const (
	LiveStartTitle          = "开播通知"
	LiveEndTitle            = "下播通知"
	LiveStartNotifyTemplate = "你关注的主播 【 %s 】 开播啦！"
	LiveEndNotifyTemplate   = "你关注的主播 【 %s 】 已下播。"
)

var (
	builders sync.Map
)

type MessageModel struct {
	Title   string
	Content string
}

type MessagePush interface {
	PushMessage(ctx context.Context, channel *model.PushChannel, mp *MessageModel) (err error)
}

func Register(channelType string, b Builder) {
	builders.Store(channelType, b)
}

type Builder interface {
	Build(string) (MessagePush, error)
}

func LivePush(ctx context.Context, anchor string, isLiving bool) {
	var title, content string
	if isLiving {
		title = LiveStartTitle
		content = fmt.Sprintf(LiveStartNotifyTemplate, anchor)
	} else {
		title = LiveEndTitle
		content = fmt.Sprintf(LiveEndNotifyTemplate, anchor)
	}
	manager.GetNotifyManager().AddInfoNotify(title, content)
	PushMessage(ctx, &MessageModel{Title: title, Content: content})
}

func PushMessage(ctx context.Context, mp *MessageModel) (err error) {
	channels := service.GetAllPushChannel(ctx)
	for _, v := range channels {
		channelPushMessage(ctx, v, mp)
	}
	return
}

func channelPushMessage(ctx context.Context, v *model.PushChannel, mp *MessageModel) error {
	b, err := getBuilder(v.Type)
	if err != nil {
		return err
	}
	builder, err := b.Build(v.Type)
	if err != nil {
		return gerror.New("不支持的渠道类型！")
	}
	err = builder.PushMessage(ctx, v, mp)
	if err != nil {
		return gerror.New(v.Type + "消息推送失败")
	}
	return nil
}

func getBuilder(channelType string) (Builder, error) {
	builder, ok := builders.Load(channelType)
	if !ok {
		return nil, gerror.New("不支持的渠道类型！")
	}
	return builder.(Builder), nil
}
