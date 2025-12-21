// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysProxy is the golang structure for table sys_proxy.
type SysProxy struct {
	Id        int         `json:"id"        orm:"id"         description:""`
	Platform  string      `json:"platform"  orm:"platform"   description:""`
	Proxy     string      `json:"proxy"     orm:"proxy"      description:""`
	Remark    string      `json:"remark"    orm:"remark"     description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`
}
