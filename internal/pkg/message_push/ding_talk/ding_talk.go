package ding_talk

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shichen437/gowlive/internal/app/system/model"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/manager"
	mp "github.com/shichen437/gowlive/internal/pkg/message_push"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

const (
	channelType = "dingTalk"
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
	if channel.Webhook == nil {
		return gerror.New("未找到 webhook 配置信息")
	}
	var req *dingTalkMessageReq
	switch channel.Webhook.MessageType {
	case 1:
		req = genRichTextMessage(model.Title, model.Content)
	default:
		req = genTextMessage(model.Content)
	}
	wurl := channel.Webhook.WebhookUrl
	if channel.Webhook.Sign != "" {
		ts := gtime.Now().UnixMilli()
		sign, err := genSecret(channel.Webhook.Sign, ts)
		if err != nil {
			return err
		}
		wurl = fmt.Sprintf("%s&timestamp=%s&sign=%s", wurl, strconv.FormatInt(ts, 10), sign)
	}
	data, err := gjson.Marshal(req)
	if err != nil {
		return err
	}
	err = mp.CheckLimit(ctx, consts.DingTalkPushLimitKey)
	if err != nil {
		return err
	}
	c := g.Client()
	resp, err := c.Post(ctx, wurl, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	body, err := utils.Text(resp.Response)
	if err != nil {
		return err
	}
	res := dingTalkMessageRes{}
	err = gjson.DecodeTo(body, &res)
	if err != nil {
		return err
	}
	if res.Code != 0 {
		errMsg := fmt.Sprintf("钉钉消息推送失败，错误码：%d，错误信息：%s", res.Code, res.Message)
		manager.GetLogManager().AddErrorLog(consts.LogTypePush, errMsg)
		return gerror.New(errMsg)
	}
	return
}

func genSecret(secret string, timestamp int64) (string, error) {
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)

	mac := hmac.New(sha256.New, []byte(secret))
	if _, err := mac.Write([]byte(stringToSign)); err != nil {
		panic(err)
	}
	base64Sig := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	sign := url.QueryEscape(base64Sig)

	return sign, nil
}
