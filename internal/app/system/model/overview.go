package model

type Overview struct {
	LiveRoomCount      int     `json:"liveRoomCount"`      // 直播间总数
	RecordingRoomCount int     `json:"recordingRoomCount"` // 正在录制的直播间总数
	RecordTimeCount    float64 `json:"recordTimeCount"`    // 总录制时长（单位：时）
	UnreadMessageCount int     `json:"unreadMessageCount"` // 未读消息总数
	ParseMediaCount    int     `json:"parseMediaCount"`    // 已解析的媒体文件总数
}
