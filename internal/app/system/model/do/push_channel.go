// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PushChannel is the golang structure of table push_channel for DAO operations like Where/Data.
type PushChannel struct {
	g.Meta    `orm:"table:push_channel, do:true"`
	Id        interface{} //
	Name      interface{} //
	Type      interface{} //
	Status    interface{} //
	Url       interface{} //
	Remark    interface{} //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
