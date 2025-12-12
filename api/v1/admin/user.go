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
	Nickname string `v:"required|length:2,12#user.profile.valid.NicknameRequired|user.profile.valid.NicknameLength" json:"nickname"`
	Sex      int    `v:"required#user.profile.valid.SexRequired" json:"sex"`
}
type PutProfileRes struct {
	g.Meta `mime:"application/json"`
}

type PutPasswordReq struct {
	g.Meta `path:"/user/password" method:"put" tags:"个人信息" summary:"修改密码"`
	OldPwd string `v:"required|length:6,20#user.profile.valid.OldPwdRequired|user.profile.valid.OldPwdLength" json:"oldPwd"`
	NewPwd string `v:"required|length:6,20|not-eq:OldPwd#user.profile.valid.NewPwdRequired|user.profile.valid.NewPwdLength|user.profile.valid.NewPwdSame" json:"newPwd"`
}
type PutPasswordRes struct {
	g.Meta `mime:"application/json"`
}
