// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// LiveManage is the golang structure of table live_manage for DAO operations like Where/Data.
type LiveManage struct {
	g.Meta         `orm:"table:live_manage, do:true"`
	Id             any         //
	RoomUrl        any         //
	Interval       any         //
	Format         any         //
	MonitorType    any         //
	MonitorStartAt any         //
	MonitorStopAt  any         //
	Remark         any         //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
	Quality        any         //
}
