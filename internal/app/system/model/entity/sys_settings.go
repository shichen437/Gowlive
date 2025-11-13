// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysSettings is the golang structure for table sys_settings.
type SysSettings struct {
	Id        int         `json:"id"        orm:"id"         description:""`
	SKey      string      `json:"sKey"      orm:"s_key"      description:""`
	SValue    int         `json:"sValue"    orm:"s_value"    description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`
}
