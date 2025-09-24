// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PushChannelEmail is the golang structure of table push_channel_email for DAO operations like Where/Data.
type PushChannelEmail struct {
	g.Meta    `orm:"table:push_channel_email, do:true"`
	Id        interface{} //
	ChannelId interface{} //
	Sender    interface{} //
	Receiver  interface{} //
	Server    interface{} //
	Port      interface{} //
	AuthCode  interface{} //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
