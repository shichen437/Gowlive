// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// FileCheckTask is the golang structure of table file_check_task for DAO operations like Where/Data.
type FileCheckTask struct {
	g.Meta     `orm:"table:file_check_task, do:true"`
	Id         any         //
	Path       any         //
	Filename   any         //
	Duration   any         //
	Progress   any         //
	FileStatus any         //
	CreatedAt  *gtime.Time //
	UpdatedAt  *gtime.Time //
}
