// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// LiveRoomInfo is the golang structure for table live_room_info.
type LiveRoomInfo struct {
	Id        int         `json:"id"        orm:"id"         description:""`
	LiveId    int         `json:"liveId"    orm:"live_id"    description:""`
	RoomName  string      `json:"roomName"  orm:"room_name"  description:""`
	Anchor    string      `json:"anchor"    orm:"anchor"     description:""`
	Platform  string      `json:"platform"  orm:"platform"   description:""`
	Status    int         `json:"status"    orm:"status"     description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`
}
