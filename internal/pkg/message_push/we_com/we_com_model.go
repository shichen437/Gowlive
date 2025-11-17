package we_com

import (
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
)

type weComMessageReq struct {
	Timestamp   string `json:"timestamp"`
	Sign        string `json:"sign"`
	MessageType string `json:"msgtype"`
	Text        any    `json:"text"`
	Markdown    any    `json:"markdown"`
}
type weComMessageRes struct {
	Code    int    `json:"errcode"`
	Message string `json:"errmsg"`
}

func genTextMessage(content string) *weComMessageReq {
	lmq := &weComMessageReq{
		MessageType: "text",
		Text: g.Map{
			"content": content,
		},
	}
	return lmq
}

func genRichTextMessage(title, content string) *weComMessageReq {
	lmq := &weComMessageReq{
		MessageType: "markdown",
		Markdown: g.Map{
			"content": fmt.Sprintf("# %s\n>%s", title, content),
		},
	}
	return lmq
}
