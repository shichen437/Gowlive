// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysLogs is the golang structure for table sys_logs.
type SysLogs struct {
	Id        int         `json:"id"        orm:"id"         description:""`
	Type      int         `json:"type"      orm:"type"       description:""`
	Content   string      `json:"content"   orm:"content"    description:""`
	Status    int         `json:"status"    orm:"status"     description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
}
