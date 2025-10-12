// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AnchorInfoHistory is the golang structure of table anchor_info_history for DAO operations like Where/Data.
type AnchorInfoHistory struct {
	g.Meta         `orm:"table:anchor_info_history, do:true"`
	Id             interface{} //
	AnchorId       interface{} //
	AnchorName     interface{} //
	Signature      interface{} //
	FollowerCount  interface{} //
	FollowingCount interface{} //
	LikeCount      interface{} //
	VideoCount     interface{} //
	CollectedDate  interface{} //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
}
