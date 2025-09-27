// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysNotify is the golang structure for table sys_notify.
type SysNotify struct {
	Id        int         `json:"id"        orm:"id"         description:""`
	Title     string      `json:"title"     orm:"title"      description:""`
	Content   string      `json:"content"   orm:"content"    description:""`
	Level     string      `json:"level"     orm:"level"      description:""`
	Status    int         `json:"status"    orm:"status"     description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`
}
