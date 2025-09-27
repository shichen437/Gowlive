package message_push

import (
	"context"
	"sync"

	"github.com/shichen437/gowlive/internal/app/system/model"
	"github.com/shichen437/gowlive/internal/pkg/manager"
	"github.com/shichen437/gowlive/internal/pkg/service"

	"github.com/gogf/gf/v2/errors/gerror"
)

var (
	builders sync.Map
)

type MessageModel struct {
	Title   string
	Content string
}

type MessagePush interface {
	Push(ctx context.Context, channel *model.PushChannel) (err error)
	PushMessage(ctx context.Context, channel *model.PushChannel, mp *MessageModel) (err error)
}

func Register(channelType string, b Builder) {
	builders.Store(channelType, b)
}

type Builder interface {
	Build(string, string) (MessagePush, error)
}

func LivePush(ctx context.Context, anchor string) {
	channels := service.GetAllPushChannel(ctx)
	manager.GetNotfiyManager().AddInfoNotify("开播通知", "您关注的主播 【"+anchor+"】 开播了")
	for _, v := range channels {
		channelPush(ctx, v, anchor)
	}
}

func PushMessage(ctx context.Context, mp *MessageModel) (err error) {
	channels := service.GetAllPushChannel(ctx)
	for _, v := range channels {
		channelPushMessage(ctx, v, mp)
	}
	return
}

func channelPush(ctx context.Context, v *model.PushChannel, anchor string) error {
	b, err := getBuilder(v.Type)
	if err != nil {
		return err
	}
	builder, err := b.Build(v.Type, anchor)
	if err != nil {
		return gerror.New("不支持的渠道类型！")
	}
	err = builder.Push(ctx, v)
	if err != nil {
		return gerror.New(v.Type + "消息推送失败")
	}
	return nil
}

func channelPushMessage(ctx context.Context, v *model.PushChannel, mp *MessageModel) error {
	b, err := getBuilder(v.Type)
	if err != nil {
		return err
	}
	builder, err := b.Build(v.Type, "")
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
