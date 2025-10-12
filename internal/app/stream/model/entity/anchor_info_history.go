// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AnchorInfoHistory is the golang structure for table anchor_info_history.
type AnchorInfoHistory struct {
	Id             int         `json:"id"             orm:"id"              description:""`
	AnchorId       int         `json:"anchorId"       orm:"anchor_id"       description:""`
	AnchorName     string      `json:"anchorName"     orm:"anchor_name"     description:""`
	Signature      string      `json:"signature"      orm:"signature"       description:""`
	FollowerCount  int         `json:"followerCount"  orm:"follower_count"  description:""`
	FollowingCount int         `json:"followingCount" orm:"following_count" description:""`
	LikeCount      int         `json:"likeCount"      orm:"like_count"      description:""`
	VideoCount     int         `json:"videoCount"     orm:"video_count"     description:""`
	CollectedDate  string      `json:"collectedDate"  orm:"collected_date"  description:""`
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      description:""`
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      description:""`
}
