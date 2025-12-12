package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/gowlive/api/v1/common"
	"github.com/shichen437/gowlive/internal/app/system/model"
	"github.com/shichen437/gowlive/internal/app/system/model/entity"
)

type GetPushChannelListReq struct {
	g.Meta `path:"/system/push/channel/list" method:"get" tags:"消息推送" summary:"渠道列表"`
	common.PageReq
	Name string `p:"name"`
	Type string `p:"type"`
}
type GetPushChannelListRes struct {
	g.Meta `mime:"application/json"`
	Rows   []*entity.PushChannel `json:"rows"`
	Total  int                   `json:"total"`
}

type PostPushChannelReq struct {
	g.Meta  `path:"/system/push/channel" method:"post" tags:"消息推送" summary:"添加渠道"`
	Name    string                   `json:"name" v:"required#system.push.valid.ChannelNameRequired"`
	Type    string                   `p:"type" json:"type" v:"required|in:gotify,email,lark,dingTalk,weCom#system.push.valid.ChannelTypeRequired|system.push.valid.ChannelTypeUnsupported"`
	Status  int                      `json:"status" v:"required#system.push.valid.ChannelStatusRequired"`
	Url     string                   `json:"url" v:"required-if:type,gotify#system.push.valid.UrlRequired"`
	Remark  string                   `json:"remark" v:"max-length:45#system.push.valid.RemarkLength"`
	Email   *PushChannelEmailModel   `json:"email" v:"required-if:type,email"`
	Webhook *PushChannelWebhookModel `json:"webhook" v:"required-if:type,lark,type,dingTalk,type,weCom"`
}
type PostPushChannelRes struct {
	g.Meta `mime:"application/json"`
}

type PutPushChannelReq struct {
	g.Meta  `path:"/system/push/channel" method:"put" tags:"消息推送" summary:"修改渠道"`
	Id      int                      `json:"id" v:"required#system.push.valid.ChannelIDRequired"`
	Name    string                   `json:"name" v:"required#system.push.valid.ChannelNameRequired"`
	Type    string                   `p:"type" json:"type" v:"required|in:gotify,email,lark,dingTalk,weCom#system.push.valid.ChannelTypeRequired|system.push.valid.ChannelTypeUnsupported"`
	Status  int                      `json:"status" v:"required#system.push.valid.ChannelStatusRequired"`
	Url     string                   `json:"url" v:"required-if:type,gotify#system.push.valid.UrlRequired"`
	Remark  string                   `json:"remark" v:"max-length:45#system.push.valid.RemarkLength"`
	Email   *PushChannelEmailModel   `json:"email" v:"required-if:type,email"`
	Webhook *PushChannelWebhookModel `json:"webhook" v:"required-if:type,lark,type,dingTalk,type,weCom"`
}
type PutPushChannelRes struct {
	g.Meta `mime:"application/json"`
}

type GetPushChannelReq struct {
	g.Meta `path:"/system/push/channel/{id}" method:"get" tags:"消息推送" summary:"获取渠道信息"`
	Id     int `p:"id"  v:"required#system.push.valid.ChannelIDRequired"`
}
type GetPushChannelRes struct {
	g.Meta `mime:"application/json"`
	*model.PushChannel
}

type DeletePushChannelReq struct {
	g.Meta `path:"/system/push/channel/{id}" method:"delete" tags:"消息推送" summary:"删除渠道"`
	Id     string `p:"id"  v:"required#system.push.valid.ChannelIDRequired"`
}
type DeletePushChannelRes struct {
	g.Meta `mime:"application/json"`
}

type PushChannelEmailModel struct {
	Sender   string `json:"sender" v:"required|email#system.push.valid.EmailSenderRequired|system.push.valid.EmailFormatInvalid"`
	Receiver string `json:"receiver" v:"required#system.push.valid.EmailReceiverRequired"`
	Server   string `json:"server" v:"required#system.push.valid.EmailServerRequired"`
	Port     int    `json:"port" v:"required#system.push.valid.EmailPortRequired"`
	AuthCode string `json:"authCode" v:"required#system.push.valid.EmailAuthCodeRequired"`
}

type PushChannelWebhookModel struct {
	WebhookUrl  string `json:"webhookUrl" v:"required|url#system.push.valid.WebhookUrlRequired|system.push.valid.WebhookUrlFormatInvalid"`
	MessageType int    `json:"messageType" v:"required|in:0,1,2#system.push.valid.MessageTypeRequired|system.push.valid.MessageTypeInvalid"`
	Sign        string `json:"sign"`
	At          string `json:"at"`
}
