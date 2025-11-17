package email

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/gowlive/internal/app/system/model"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/manager"
	mp "github.com/shichen437/gowlive/internal/pkg/message_push"

	"gopkg.in/mail.v2"
)

const (
	channelType = "email"
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
	m := mail.NewMessage()
	m.SetHeader("From", channel.Email.Sender)
	m.SetHeader("To", channel.Email.Receiver)
	m.SetHeader("Subject", model.Title)
	m.SetBody("text/html", model.Content)
	d := mail.NewDialer(channel.Email.Server, channel.Email.Port, channel.Email.Sender, channel.Email.AuthCode)
	d.StartTLSPolicy = mail.MandatoryStartTLS
	err = d.DialAndSend(m)
	if err != nil {
		g.Log().Error(ctx, "邮箱推送失败，错误信息:"+err.Error())
		manager.GetLogManager().AddErrorLog(consts.LogTypePush, "邮箱推送失败")
	}
	return err
}
