// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysUser is the golang structure for table sys_user.
type SysUser struct {
	Id        int         `json:"id"        orm:"id"         description:""`
	Username  string      `json:"username"  orm:"username"   description:""`
	Password  string      `json:"password"  orm:"password"   description:""`
	Nickname  string      `json:"nickname"  orm:"nickname"   description:""`
	Sex       int         `json:"sex"       orm:"sex"        description:""`
	Status    int         `json:"status"    orm:"status"     description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`
}
