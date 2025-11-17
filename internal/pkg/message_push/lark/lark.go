package lark

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
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
	channelType = "lark"
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
	var req *larkMessageReq
	switch channel.Webhook.MessageType {
	case 1:
		req = genRichTextMessage(model.Title, model.Content)
	case 2:
		req = genCardMessage(model.Title, model.Content)
	default:
		req = genTextMessage(model.Content)
	}
	if channel.Webhook.Sign != "" {
		ts := gtime.Now().Unix()
		req.Timestamp = strconv.FormatInt(ts, 10)
		sign, err := genSign(channel.Webhook.Sign, ts)
		if err != nil {
			return err
		}
		req.Sign = sign
	}
	data, err := gjson.Marshal(req)
	if err != nil {
		return err
	}
	err = mp.CheckLimit(ctx, consts.LarkPushLimitKey)
	if err != nil {
		return err
	}
	c := g.Client()
	resp, err := c.Post(ctx, channel.Webhook.WebhookUrl, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	body, err := utils.Text(resp.Response)
	if err != nil {
		return err
	}
	res := larkMessageRes{}
	err = gjson.DecodeTo(body, &res)
	if err != nil {
		return err
	}
	if res.Code != 0 {
		errMsg := fmt.Sprintf("飞书消息推送失败，错误码：%d，错误信息：%s", res.Code, res.Message)
		manager.GetLogManager().AddErrorLog(consts.LogTypePush, errMsg)
		return gerror.New(errMsg)
	}
	return
}

func genSign(secret string, timestamp int64) (string, error) {
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret
	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}
