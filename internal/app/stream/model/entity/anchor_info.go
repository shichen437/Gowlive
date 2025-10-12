// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AnchorInfo is the golang structure for table anchor_info.
type AnchorInfo struct {
	Id             int         `json:"id"             orm:"id"              description:""`
	AnchorName     string      `json:"anchorName"     orm:"anchor_name"     description:""`
	Url            string      `json:"url"            orm:"url"             description:""`
	Signature      string      `json:"signature"      orm:"signature"       description:""`
	Platform       string      `json:"platform"       orm:"platform"        description:""`
	UniqueId       string      `json:"uniqueId"       orm:"unique_id"       description:""`
	FollowerCount  int         `json:"followerCount"  orm:"follower_count"  description:""`
	FollowingCount int         `json:"followingCount" orm:"following_count" description:""`
	LikeCount      int         `json:"likeCount"      orm:"like_count"      description:""`
	VideoCount     int         `json:"videoCount"     orm:"video_count"     description:""`
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      description:""`
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      description:""`
}
