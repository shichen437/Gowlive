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
	LiveId int `json:"liveId" dc:"直播房间ID" v:"min:1#stream.live.valid.RoomIDRequired"`
}
type GetLiveManageRes struct {
	g.Meta `mime:"application/json"`
	Data   *entity.LiveManage `json:"data"`
}

type PostLiveManageReq struct {
	g.Meta         `path:"/live/room/manage" method:"post" tags:"直播管理" summary:"创建直播房间"`
	RoomUrl        string `json:"roomUrl"  v:"required|url#stream.live.valid.RoomUrlRequired|stream.live.valid.RoomUrlFormat"`
	Interval       int    `json:"interval" v:"required|min:30#stream.live.valid.IntervalRequired|stream.live.valid.IntervalMin"`
	Format         string `json:"format" v:"required|in:flv,mp4,mkv,ts,mp3#stream.live.valid.RecordFormatRequired|stream.live.valid.RecordFormatInvalid"`
	MonitorType    int    `json:"monitorType" v:"required|in:0,1,2,3#stream.live.valid.MonitorTypeRequired|stream.live.valid.MonitorTypeInvalid"`
	MonitorStartAt string `json:"monitorStartAt" v:"required-if:monitorType,2#stream.live.valid.MonitorStartTimeRequired"`
	MonitorStopAt  string `json:"monitorStopAt" v:"required-if:monitorType,2|not-eq:monitorStartAt#stream.live.valid.MonitorEndTimeRequired|stream.live.valid.MonitorEndTimeSame"`
	Quality        int    `json:"quality" v:"required|in:0,1,2,3,4#stream.live.valid.QualityRequired|stream.live.valid.QualityInvalid"`
	SegmentTime    int    `json:"segmentTime" v:"required#stream.live.valid.SliceTimeRequired"`
	Remark         string `json:"remark" v:"max-length:45#stream.live.valid.RemarkMaxLength"`
}
type PostLiveManageRes struct {
	g.Meta `mime:"application/json"`
}

type PostLiveManageBatchReq struct {
	g.Meta   `path:"/live/room/manage/batch" method:"post" tags:"直播管理" summary:"批量添加直播间"`
	RoomUrls []string `json:"roomUrls" v:"required|array|foreach|url#stream.live.valid.RoomUrlRequired|stream.live.valid.DataFormatInvalid|foreach|stream.live.valid.RoomUrlFormat"`
	Interval int      `json:"interval" v:"required|min:30#stream.live.valid.IntervalRequired|stream.live.valid.IntervalMin"`
	Format   string   `json:"format" v:"required|in:flv,mp4,mkv,ts,mp3#stream.live.valid.RecordFormatRequired|stream.live.valid.RecordFormatInvalid"`
	Remark   string   `json:"remark" v:"max-length:45#stream.live.valid.RemarkMaxLength"`
}
type PostLiveManageBatchRes struct {
	g.Meta `mime:"application/json"`
}

type PutLiveManageReq struct {
	g.Meta         `path:"/live/room/manage" method:"put" tags:"直播管理" summary:"更新直播房间"`
	Id             int    `json:"id" dc:"直播房间ID" v:"min:1#stream.live.valid.RoomIDRequired"`
	Interval       int    `json:"interval" v:"required|min:30#stream.live.valid.IntervalRequired|stream.live.valid.IntervalMin"`
	Format         string `json:"format" v:"required|in:flv,mp4,mkv,ts,mp3#stream.live.valid.RecordFormatRequired|stream.live.valid.RecordFormatInvalid"`
	MonitorType    int    `json:"monitorType" v:"required|in:0,1,2,3#stream.live.valid.MonitorTypeRequired|stream.live.valid.MonitorTypeInvalid"`
	MonitorStartAt string `json:"monitorStartAt" v:"required-if:monitorType,2#stream.live.valid.MonitorStartTimeRequired"`
	MonitorStopAt  string `json:"monitorStopAt" v:"required-if:monitorType,2|not-eq:monitorStartAt#stream.live.valid.MonitorEndTimeRequired|stream.live.valid.MonitorEndTimeSame"`
	Quality        int    `json:"quality" v:"required|in:0,1,2,3,4#stream.live.valid.QualityRequired|stream.live.valid.QualityInvalid"`
	SegmentTime    int    `json:"segmentTime" v:"required#stream.live.valid.SliceTimeRequired"`
	Remark         string `json:"remark" v:"max-length:45#stream.live.valid.RemarkMaxLength"`
}
type PutLiveManageRes struct {
	g.Meta `mime:"application/json"`
}

type DeleteLiveManageReq struct {
	g.Meta `path:"/live/room/manage/{liveId}" method:"delete" tags:"直播管理" summary:"删除直播房间"`
	LiveId int `json:"liveId" dc:"直播房间ID" v:"min:1#stream.live.valid.RoomIDRequired"`
}
type DeleteLiveManageRes struct {
	g.Meta `mime:"application/json"`
}

type PutLiveManageStartReq struct {
	g.Meta `path:"/live/room/manage/start/{id}" method:"put" tags:"直播管理" summary:"开始监控"`
	Id     int `json:"id" dc:"直播房间ID" v:"min:1#stream.live.valid.RoomIDRequired"`
}
type PutLiveManageStartRes struct {
	g.Meta `mime:"application/json"`
}

type PutLiveManageStopReq struct {
	g.Meta `path:"/live/room/manage/stop/{id}" method:"put" tags:"直播管理" summary:"停止监控"`
	Id     int `json:"id" dc:"直播房间ID" v:"min:1#stream.live.valid.RoomIDRequired"`
}
type PutLiveManageStopRes struct {
	g.Meta `mime:"application/json"`
}

type PutLiveManageTopReq struct {
	g.Meta `path:"/live/room/manage/top/{id}" method:"put" tags:"直播管理" summary:"置顶房间"`
	Id     int `json:"id" dc:"直播房间ID" v:"min:1#stream.live.valid.RoomIDRequired"`
}
type PutLiveManageTopRes struct {
	g.Meta `mime:"application/json"`
}

type PutLiveManageUnTopReq struct {
	g.Meta `path:"/live/room/manage/unTop/{id}" method:"put" tags:"直播管理" summary:"取消置顶房间"`
	Id     int `json:"id" dc:"直播房间ID" v:"min:1#stream.live.valid.RoomIDRequired"`
}
type PutLiveManageUnTopRes struct {
	g.Meta `mime:"application/json"`
}

type ExportRoomInfoReq struct {
	g.Meta     `path:"/live/room/export" method:"get" tags:"直播管理" summary:"导出直播间信息"`
	ExportType int    `json:"exportType" dc:"导出类型" v:"required|in:0,1#stream.live.valid.ExportTypeRequired|stream.live.valid.ExportTypeInvalid"`
	Anchor     string `json:"anchor"`
	RoomName   string `json:"roomName"`
	Platform   string `json:"platform"`
	Sort       string `json:"sort"`
}
type ExportRoomInfoRes struct {
	g.Meta `mime:"application/json"`
}
