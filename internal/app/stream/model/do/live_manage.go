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
	Id             interface{} //
	RoomUrl        interface{} //
	Interval       interface{} //
	Format         interface{} //
	MonitorType    interface{} //
	MonitorStartAt interface{} //
	MonitorStopAt  interface{} //
	Remark         interface{} //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
}
