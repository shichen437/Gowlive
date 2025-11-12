// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// FileCheckTask is the golang structure for table file_check_task.
type FileCheckTask struct {
	Id         int         `json:"id"         orm:"id"          description:""`
	Path       string      `json:"path"       orm:"path"        description:""`
	Filename   string      `json:"filename"   orm:"filename"    description:""`
	Duration   int         `json:"duration"   orm:"duration"    description:""`
	Progress   int         `json:"progress"   orm:"progress"    description:""`
	FileStatus int         `json:"fileStatus" orm:"file_status" description:""`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:""`
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  description:""`
}
