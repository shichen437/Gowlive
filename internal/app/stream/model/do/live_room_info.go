// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// LiveRoomInfo is the golang structure of table live_room_info for DAO operations like Where/Data.
type LiveRoomInfo struct {
	g.Meta    `orm:"table:live_room_info, do:true"`
	Id        interface{} //
	LiveId    interface{} //
	RoomName  interface{} //
	Anchor    interface{} //
	Platform  interface{} //
	Status    interface{} //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
