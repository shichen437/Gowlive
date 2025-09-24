// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// LiveHistory is the golang structure of table live_history for DAO operations like Where/Data.
type LiveHistory struct {
	g.Meta    `orm:"table:live_history, do:true"`
	Id        interface{} //
	LiveId    interface{} //
	Anchor    interface{} //
	StartedAt *gtime.Time //
	EndedAt   *gtime.Time //
	Duration  interface{} //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
