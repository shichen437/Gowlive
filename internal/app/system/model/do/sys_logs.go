// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysLogs is the golang structure of table sys_logs for DAO operations like Where/Data.
type SysLogs struct {
	g.Meta    `orm:"table:sys_logs, do:true"`
	Id        interface{} //
	Type      interface{} //
	Content   interface{} //
	Status    interface{} //
	CreatedAt *gtime.Time //
}
