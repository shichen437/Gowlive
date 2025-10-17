package sse

import (
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/google/uuid"
)

type SSEMessage struct {
	Event     string      `json:"event"`
	Data      interface{} `json:"data"`
	ID        string      `json:"id"`
	Retry     int         `json:"retry"`
	Timestamp int64       `json:"timestamp"`
}

func GetSseMsgStr(event string, data interface{}) string {
	msg := new(event, data)
	return gjson.MustEncodeString(msg)
}

func new(event string, data interface{}) *SSEMessage {
	return &SSEMessage{
		Event:     event,
		Data:      data,
		ID:        uuid.New().String(),
		Timestamp: gconv.Int64(time.Now().Unix()),
	}
}
