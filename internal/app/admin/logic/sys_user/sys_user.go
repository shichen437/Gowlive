package logic

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
	v1 "github.com/shichen437/gowlive/api/v1/admin"
	"github.com/shichen437/gowlive/internal/app/admin/dao"
	"github.com/shichen437/gowlive/internal/app/admin/model"
	"github.com/shichen437/gowlive/internal/app/admin/model/do"
	"github.com/shichen437/gowlive/internal/app/admin/model/entity"
	"github.com/shichen437/gowlive/internal/app/admin/service"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/manager"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

type (
	sSysUser struct{}
)

func init() {
	service.RegisterSysUser(New())
}

func New() service.ISysUser {
	return &sSysUser{}
}

func (s *sSysUser) GetUserInfo(ctx context.Context, req *v1.GetUserInfoReq) (res *v1.GetUserInfoRes, err error) {
	res = &v1.GetUserInfoRes{}
	uid := gconv.Int(ctx.Value(consts.CtxAdminId))
	if uid == 0 {
		return
	}
	user := getUserById(ctx, uid)
	if user == nil {
		return
	}
	var userInfo *model.UserInfo
	gconv.Struct(user, &userInfo)
	res.UserInfo = userInfo
	return
}

func (s *sSysUser) UpdateProfile(ctx context.Context, req *v1.PutProfileReq) (res *v1.PutProfileRes, err error) {
	uid := gconv.Int(ctx.Value(consts.CtxAdminId))
	if uid == 0 {
		return
	}
	user := getUserById(ctx, uid)
	if user == nil {
		return
	}
	_, err = dao.SysUser.Ctx(ctx).WherePri(uid).Update(do.SysUser{
		Nickname:  req.Nickname,
		Sex:       req.Sex,
		UpdatedAt: utils.Now(),
	})
	if err != nil {
		err = gerror.New("修改用户信息失败")
		manager.GetLogManager().AddErrorLog(consts.LogTypeUser, "修改信息失败")
		return
	}
	manager.GetLogManager().AddSuccessLog(consts.LogTypeUser, "修改用户信息成功")
	return
}

func (s *sSysUser) UpdatePassword(ctx context.Context, req *v1.PutPasswordReq) (res *v1.PutPasswordRes, err error) {
	uid := gconv.Int(ctx.Value(consts.CtxAdminId))
	if uid == 0 {
		return
	}
	op, err := utils.Encrypt(ctx, req.OldPwd)
	if err != nil {
		err = gerror.New("密码加密失败")
		return
	}
	user := getUserById(ctx, uid)
	if user == nil || user.Password != op {
		err = gerror.New("旧密码错误")
		return
	}
	np, err := utils.Encrypt(ctx, req.NewPwd)
	if err != nil {
		err = gerror.New("密码加密失败")
		return
	}
	_, err = dao.SysUser.Ctx(ctx).WherePri(uid).Update(do.SysUser{
		Password:  np,
		UpdatedAt: utils.Now(),
	})
	if err != nil {
		err = gerror.New("修改用户密码失败")
		manager.GetLogManager().AddErrorLog(consts.LogTypeUser, "修改用户密码失败")
		return
	}
	manager.GetLogManager().AddErrorLog(consts.LogTypeUser, "修改用户密码成功")
	return
}

func getUserById(ctx context.Context, uid int) (user *entity.SysUser) {
	user = &entity.SysUser{}
	err := dao.SysUser.Ctx(ctx).WherePri(uid).Scan(&user)
	if err != nil {
		return nil
	}
	return
}
