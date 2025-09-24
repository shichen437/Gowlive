// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// LiveCookie is the golang structure for table live_cookie.
type LiveCookie struct {
	Id        int         `json:"id"        orm:"id"         description:""`
	Platform  string      `json:"platform"  orm:"platform"   description:""`
	Cookie    string      `json:"cookie"    orm:"cookie"     description:""`
	Remark    string      `json:"remark"    orm:"remark"     description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`
}
