package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/gowlive/api/v1/common"
	"github.com/shichen437/gowlive/internal/app/stream/model"
	"github.com/shichen437/gowlive/internal/app/stream/model/entity"
)

type GetRoomListReq struct {
	g.Meta `path:"/live/room/list" method:"get" tags:"直播管理" summary:"获取直播房间列表"`
	common.PageReq
	Anchor   string `json:"anchor"`
	RoomName string `json:"roomName"`
	Platform string `json:"platform"`
	Sort     string `json:"sort"`
}
type GetRoomListRes struct {
	g.Meta `mime:"application/json"`
	Rows   []*model.RoomInfo `json:"rows"`
	Total  int               `json:"total"`
}

type GetLiveManageReq struct {
	g.Meta `path:"/live/room/manage/{liveId}" method:"get" tags:"直播管理" summary:"获取直播房间管理"`
	LiveId int `json:"liveId" dc:"直播房间ID" v:"min:1#直播房间ID不能为空"`
}
type GetLiveManageRes struct {
	g.Meta `mime:"application/json"`
	Data   *entity.LiveManage `json:"data"`
}

type PostLiveManageReq struct {
	g.Meta         `path:"/live/room/manage" method:"post" tags:"直播管理" summary:"创建直播房间"`
	RoomUrl        string `json:"roomUrl"  v:"required|url#房间地址不能为空|房间地址必须为有效URL"`
	Interval       int    `json:"interval" v:"required|min:30#间隔时间不能为空|间隔时间最小为30s"`
	Format         string `json:"format" v:"required|in:flv,mp4,mp3#录制格式不能为空|录制格式不合法"`
	MonitorType    int    `json:"monitorType" v:"required|in:0,1,2#监控类型不能为空|监控类型参数不合法"`
	MonitorStartAt string `json:"monitorStartAt" v:"required-if:monitorType,2#定时监控开始时间不能为空"`
	MonitorStopAt  string `json:"monitorStopAt" v:"required-if:monitorType,2|not-eq:monitorStartAt#定时监控结束时间不能为空|定时监控结束时间不能与开始时间相同"`
	Remark         string `json:"remark" v:"max-length:45#备注最大长度为45"`
}
type PostLiveManageRes struct {
	g.Meta `mime:"application/json"`
}

type PutLiveManageReq struct {
	g.Meta         `path:"/live/room/manage" method:"put" tags:"直播管理" summary:"更新直播房间"`
	Id             int    `json:"id" dc:"直播房间ID" v:"min:1#直播房间ID不能为空"`
	Interval       int    `json:"interval" v:"required|min:30#间隔时间不能为空|间隔时间最小为30s"`
	Format         string `json:"format" v:"required|in:flv,mp4,mp3#录制格式不能为空|录制格式不合法"`
	MonitorType    int    `json:"monitorType" v:"required|in:0,1,2#监控类型不能为空|监控类型参数不合法"`
	MonitorStartAt string `json:"monitorStartAt" v:"required-if:monitorType,2#定时监控开始时间不能为空"`
	MonitorStopAt  string `json:"monitorStopAt" v:"required-if:monitorType,2|not-eq:monitorStartAt#定时监控结束时间不能为空|定时监控结束时间不能与开始时间相同"`
	Remark         string `json:"remark" v:"max-length:45#备注最大长度为45"`
}
type PutLiveManageRes struct {
	g.Meta `mime:"application/json"`
}

type DeleteLiveManageReq struct {
	g.Meta `path:"/live/room/manage/{liveId}" method:"delete" tags:"直播管理" summary:"删除直播房间"`
	LiveId int `json:"liveId" dc:"直播房间ID" v:"min:1#直播房间ID不能为空"`
}
type DeleteLiveManageRes struct {
	g.Meta `mime:"application/json"`
}
