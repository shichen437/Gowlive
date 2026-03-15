// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PushChannelCustomWebhook is the golang structure of table push_channel_custom_webhook for DAO operations like Where/Data.
type PushChannelCustomWebhook struct {
	g.Meta         `orm:"table:push_channel_custom_webhook, do:true"`
	Id             any         //
	ChannelId      any         //
	WebhookUrl     any         //
	RequestMethod  any         //
	RequestHeaders any         //
	RequestBody    any         //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
}
