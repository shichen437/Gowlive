// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// LiveManage is the golang structure for table live_manage.
type LiveManage struct {
	Id             int         `json:"id"             orm:"id"               description:""`
	RoomUrl        string      `json:"roomUrl"        orm:"room_url"         description:""`
	Interval       int         `json:"interval"       orm:"interval"         description:""`
	Format         string      `json:"format"         orm:"format"           description:""`
	MonitorType    int         `json:"monitorType"    orm:"monitor_type"     description:""`
	MonitorStartAt string      `json:"monitorStartAt" orm:"monitor_start_at" description:""`
	MonitorStopAt  string      `json:"monitorStopAt"  orm:"monitor_stop_at"  description:""`
	Remark         string      `json:"remark"         orm:"remark"           description:""`
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"       description:""`
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"       description:""`
}
