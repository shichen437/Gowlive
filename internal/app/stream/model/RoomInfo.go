package model

import "github.com/shichen437/gowlive/internal/app/stream/model/entity"

type RoomInfo struct {
	*entity.LiveRoomInfo
	IsRecording bool `json:"isRecording"`
	IsLiving    bool `json:"isLiving"`
}
