// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PushChannelCustomWebhook is the golang structure for table push_channel_custom_webhook.
type PushChannelCustomWebhook struct {
	Id             int         `json:"id"             orm:"id"              description:""`
	ChannelId      int         `json:"channelId"      orm:"channel_id"      description:""`
	WebhookUrl     string      `json:"webhookUrl"     orm:"webhook_url"     description:""`
	RequestMethod  int         `json:"requestMethod"  orm:"request_method"  description:""`
	RequestHeaders string      `json:"requestHeaders" orm:"request_headers" description:""`
	RequestBody    string      `json:"requestBody"    orm:"request_body"    description:""`
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      description:""`
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      description:""`
}
