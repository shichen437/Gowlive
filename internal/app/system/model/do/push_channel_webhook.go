// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PushChannelWebhook is the golang structure of table push_channel_webhook for DAO operations like Where/Data.
type PushChannelWebhook struct {
	g.Meta      `orm:"table:push_channel_webhook, do:true"`
	Id          any         //
	ChannelId   any         //
	WebhookUrl  any         //
	MessageType any         //
	Sign        any         //
	At          any         //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
}
