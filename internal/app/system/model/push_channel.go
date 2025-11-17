package model

import "github.com/shichen437/gowlive/internal/app/system/model/entity"

type PushChannel struct {
	*entity.PushChannel
	Email   *entity.PushChannelEmail   `json:"email" desc:"邮箱信息"`
	Webhook *entity.PushChannelWebhook `json:"webhook" desc:"Webhook信息"`
}
