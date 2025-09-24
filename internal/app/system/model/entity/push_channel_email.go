// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PushChannelEmail is the golang structure for table push_channel_email.
type PushChannelEmail struct {
	Id        int         `json:"id"        orm:"id"         description:""`
	ChannelId int         `json:"channelId" orm:"channel_id" description:""`
	Sender    string      `json:"sender"    orm:"sender"     description:""`
	Receiver  string      `json:"receiver"  orm:"receiver"   description:""`
	Server    string      `json:"server"    orm:"server"     description:""`
	Port      int         `json:"port"      orm:"port"       description:""`
	AuthCode  string      `json:"authCode"  orm:"auth_code"  description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`
}
