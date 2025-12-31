// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// FileSyncTask is the golang structure for table file_sync_task.
type FileSyncTask struct {
	Id        int         `json:"id"        orm:"id"         description:""`
	Path      string      `json:"path"      orm:"path"       description:""`
	Filename  string      `json:"filename"  orm:"filename"   description:""`
	SyncPath  string      `json:"syncPath"  orm:"sync_path"  description:""`
	Duration  int         `json:"duration"  orm:"duration"   description:""`
	Status    int         `json:"status"    orm:"status"     description:""`
	Remark    string      `json:"remark"    orm:"remark"     description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`
}
