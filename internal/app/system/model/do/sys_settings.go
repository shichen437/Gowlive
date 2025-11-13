// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysSettings is the golang structure of table sys_settings for DAO operations like Where/Data.
type SysSettings struct {
	g.Meta    `orm:"table:sys_settings, do:true"`
	Id        any         //
	SKey      any         //
	SValue    any         //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
