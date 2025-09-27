// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// LiveHistory is the golang structure for table live_history.
type LiveHistory struct {
	Id        int         `json:"id"        orm:"id"         description:""`
	LiveId    int         `json:"liveId"    orm:"live_id"    description:""`
	Anchor    string      `json:"anchor"    orm:"anchor"     description:""`
	StartedAt *gtime.Time `json:"startedAt" orm:"started_at" description:""`
	EndedAt   *gtime.Time `json:"endedAt"   orm:"ended_at"   description:""`
	Duration  float64     `json:"duration"  orm:"duration"   description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`
	IsDelete  int         `json:"isDelete"  orm:"is_delete"  description:""`
}
