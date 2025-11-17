// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PushChannelWebhook is the golang structure for table push_channel_webhook.
type PushChannelWebhook struct {
	Id          int         `json:"id"          orm:"id"           description:""`
	ChannelId   int         `json:"channelId"   orm:"channel_id"   description:""`
	WebhookUrl  string      `json:"webhookUrl"  orm:"webhook_url"  description:""`
	MessageType int         `json:"messageType" orm:"message_type" description:""`
	Sign        string      `json:"sign"        orm:"sign"         description:""`
	At          string      `json:"at"          orm:"at"           description:""`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:""`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   description:""`
}
