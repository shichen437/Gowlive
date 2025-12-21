// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysProxy is the golang structure of table sys_proxy for DAO operations like Where/Data.
type SysProxy struct {
	g.Meta    `orm:"table:sys_proxy, do:true"`
	Id        any         //
	Platform  any         //
	Proxy     any         //
	Remark    any         //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
