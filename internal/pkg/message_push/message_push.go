package message_push

import (
	"context"
	"sync"
	"time"

	"github.com/shichen437/gowlive/internal/app/system/model"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/manager"
	"github.com/shichen437/gowlive/internal/pkg/service"
	"github.com/shichen437/gowlive/internal/pkg/utils"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
)

var (
	builders      sync.Map
	limitKeyLocks sync.Map
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
		title = utils.T(ctx, "ext.push.live.start.title")
		content = utils.Tf(ctx, "ext.push.live.start.content", anchor)
	} else {
		title = utils.T(ctx, "ext.push.live.end.title")
		content = utils.Tf(ctx, "ext.push.live.end.content", anchor)
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
		g.Log().Error(ctx, err)
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

func CheckLimit(ctx context.Context, key string) error {
	mu := getLimitKeyLock(key)
	mu.Lock()
	defer mu.Unlock()
	limit, err := gcache.Get(ctx, key)
	if err != nil {
		return gerror.New("获取限制缓存失败")
	}
	if limit == nil {
		gcache.Set(ctx, key, 1, time.Minute)
	} else {
		if limit.Int() > consts.WebhookPushLimitPerMinute {
			return gerror.New("消息推送超出限制")
		}
		gcache.Update(ctx, key, limit.Int()+1)
	}
	return nil
}

func getLimitKeyLock(key string) *sync.Mutex {
	if l, ok := limitKeyLocks.Load(key); ok {
		return l.(*sync.Mutex)
	}
	mu := &sync.Mutex{}
	actual, _ := limitKeyLocks.LoadOrStore(key, mu)
	return actual.(*sync.Mutex)
}
