package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/gowlive/internal/app/admin/model"
)

type GetUserInfoReq struct {
	g.Meta `path:"/user/getInfo" method:"get" tags:"个人信息" summary:"获取个人信息"`
}
type GetUserInfoRes struct {
	g.Meta          `mime:"application/json"`
	*model.UserInfo `json:"userInfo"`
}

type PutProfileReq struct {
	g.Meta   `path:"/user/profile" method:"put" tags:"个人信息" summary:"修改个人信息"`
	Nickname string `v:"required|length:2,12#请输入昵称|昵称长度为 2-12 位" json:"nickname"`
	Sex      int    `v:"required#请选择性别" json:"sex"`
}
type PutProfileRes struct {
	g.Meta `mime:"application/json"`
}

type PutPasswordReq struct {
	g.Meta `path:"/user/password" method:"put" tags:"个人信息" summary:"修改密码"`
	OldPwd string `v:"required|length:6,20#请输入旧密码|旧密码长度为 6-20 位" json:"oldPwd"`
	NewPwd string `v:"required|length:6,20|not-eq:OldPwd#请输入新密码|新密码长度为 6-20 位|新密码不能与旧密码相同" json:"newPwd"`
}
type PutPasswordRes struct {
	g.Meta `mime:"application/json"`
}
