package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/gowlive/api/v1/common"
	"github.com/shichen437/gowlive/internal/app/system/model"
	"github.com/shichen437/gowlive/internal/app/system/model/entity"
)

type GetPushChannelListReq struct {
	g.Meta `path:"/system/push/channel/list" method:"get" tags:"消息推送" summary:"渠道列表"`
	common.PageReq
	Name string `p:"name"`
	Type string `p:"type"`
}
type GetPushChannelListRes struct {
	g.Meta `mime:"application/json"`
	Rows   []*entity.PushChannel `json:"rows"`
	Total  int                   `json:"total"`
}

type PostPushChannelReq struct {
	g.Meta `path:"/system/push/channel" method:"post" tags:"消息推送" summary:"添加渠道"`
	Name   string                 `json:"name" v:"required#渠道名称不能为空"`
	Type   string                 `p:"type" json:"type" v:"required|in:gotify,email#渠道类型不能为空|不支持的渠道类型"`
	Status int                    `json:"status" v:"required#渠道启用状态不能为空"`
	Url    string                 `json:"url" v:"required-if:type,gotify#Url不能为空"`
	Remark string                 `json:"remark" v:"max-length:45#备注最大长度为45"`
	Email  *PushChannelEmailModel `json:"email" v:"required-if:type,email"`
}
type PostPushChannelRes struct {
	g.Meta `mime:"application/json"`
}

type PutPushChannelReq struct {
	g.Meta `path:"/system/push/channel" method:"put" tags:"消息推送" summary:"修改渠道"`
	Id     int                    `json:"id" v:"required#渠道ID不能为空"`
	Name   string                 `json:"name" v:"required#渠道名称不能为空"`
	Type   string                 `p:"type" json:"type" v:"required|in:gotify,email#渠道类型不能为空|不支持的渠道类型"`
	Status int                    `json:"status" v:"required#渠道启用状态不能为空"`
	Url    string                 `json:"url" v:"required-if:type,gotify#Url不能为空"`
	Remark string                 `json:"remark" v:"max-length:45#备注最大长度为45"`
	Email  *PushChannelEmailModel `json:"email" v:"required-if:type,email"`
}
type PutPushChannelRes struct {
	g.Meta `mime:"application/json"`
}

type GetPushChannelReq struct {
	g.Meta `path:"/system/push/channel/{id}" method:"get" tags:"消息推送" summary:"获取渠道信息"`
	Id     int `p:"id"  v:"required#渠道ID不能为空"`
}
type GetPushChannelRes struct {
	g.Meta `mime:"application/json"`
	*model.PushChannel
}

type DeletePushChannelReq struct {
	g.Meta `path:"/system/push/channel/{id}" method:"delete" tags:"消息推送" summary:"删除渠道"`
	Id     string `p:"id"  v:"required#渠道ID不能为空"`
}
type DeletePushChannelRes struct {
	g.Meta `mime:"application/json"`
}

type PushChannelEmailModel struct {
	Sender   string `json:"sender" v:"required|email#发送邮箱不能为空|邮箱格式不正确"`
	Receiver string `json:"receiver" v:"required#接收邮箱不能为空"`
	Server   string `json:"server" v:"required#邮箱服务器不能为空"`
	Port     int    `json:"port" v:"required#邮箱端口不能为空"`
	AuthCode string `json:"authCode" v:"required#邮箱授权码不能为空"`
}
