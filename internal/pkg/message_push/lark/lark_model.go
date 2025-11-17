package lark

import (
	"github.com/gogf/gf/v2/frame/g"
)

type larkMessageReq struct {
	Timestamp   string `json:"timestamp"`
	Sign        string `json:"sign"`
	MessageType string `json:"msg_type"`
	Content     any    `json:"content"`
	Card        any    `json:"card"`
}
type larkMessageRes struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

func genTextMessage(content string) *larkMessageReq {
	lmq := &larkMessageReq{
		MessageType: "text",
		Content: g.Map{
			"text": content,
		},
	}
	return lmq
}

func genRichTextMessage(title, content string) *larkMessageReq {
	contentRows := g.Slice{
		g.Slice{
			g.Map{
				"tag":  "text",
				"text": content,
			},
		},
	}
	lmq := &larkMessageReq{
		MessageType: "post",
		Content: g.Map{
			"post": g.Map{
				"zh-cn": g.Map{
					"title":   title,
					"content": contentRows,
				},
			},
		},
	}
	return lmq
}

func genCardMessage(title, content string) *larkMessageReq {
	lmq := &larkMessageReq{
		MessageType: "interactive",
		Card: g.Map{
			"schema": "2.0",
			"body": g.Map{
				"elements": g.Slice{
					g.Map{
						"tag":     "markdown",
						"content": content,
					},
				},
			},
			"header": g.Map{
				"title": g.Map{
					"tag":     "plain_text",
					"content": title,
				},
				"template": "blue",
			},
		},
	}
	return lmq
}
