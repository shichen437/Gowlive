package ding_talk

import (
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
)

type dingTalkMessageReq struct {
	MessageType string `json:"msgtype"`
	Content     any    `json:"content"`
	Markdown    any    `json:"markdown"`
	ActionCard  any    `json:"actionCard"`
}

type dingTalkMessageRes struct {
	Code    int    `json:"errcode"`
	Message string `json:"errmsg"`
}

func genTextMessage(content string) *dingTalkMessageReq {
	lmq := &dingTalkMessageReq{
		MessageType: "text",
		Content: g.Map{
			"text": content,
		},
	}
	return lmq
}

func genRichTextMessage(title, content string) *dingTalkMessageReq {
	lmq := &dingTalkMessageReq{
		MessageType: "markdown",
		Markdown: g.Map{
			"title": title,
			"text":  fmt.Sprintf("### %s\n>%s", title, content),
		},
	}
	return lmq
}
