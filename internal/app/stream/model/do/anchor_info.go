// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AnchorInfo is the golang structure of table anchor_info for DAO operations like Where/Data.
type AnchorInfo struct {
	g.Meta         `orm:"table:anchor_info, do:true"`
	Id             interface{} //
	AnchorName     interface{} //
	Url            interface{} //
	Signature      interface{} //
	Platform       interface{} //
	UniqueId       interface{} //
	FollowerCount  interface{} //
	FollowingCount interface{} //
	LikeCount      interface{} //
	VideoCount     interface{} //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
}
