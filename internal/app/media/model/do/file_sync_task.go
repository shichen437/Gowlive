// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// FileSyncTask is the golang structure of table file_sync_task for DAO operations like Where/Data.
type FileSyncTask struct {
	g.Meta    `orm:"table:file_sync_task, do:true"`
	Id        any         //
	Path      any         //
	Filename  any         //
	SyncPath  any         //
	Duration  any         //
	Status    any         //
	Remark    any         //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
