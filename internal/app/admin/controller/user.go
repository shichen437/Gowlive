package controller

import (
	"context"

	v1 "github.com/shichen437/gowlive/api/v1/admin"
	"github.com/shichen437/gowlive/internal/app/admin/service"
)

type sysUserController struct {
}

var SysUser = sysUserController{}

func (s *sysUserController) GetUserInfo(ctx context.Context, req *v1.GetUserInfoReq) (res *v1.GetUserInfoRes, err error) {
	return service.SysUser().GetUserInfo(ctx, req)
}

func (s *sysUserController) UpdateProfile(ctx context.Context, req *v1.PutProfileReq) (res *v1.PutProfileRes, err error) {
	return service.SysUser().UpdateProfile(ctx, req)
}

func (s *sysUserController) UpdatePassword(ctx context.Context, req *v1.PutPasswordReq) (res *v1.PutPasswordRes, err error) {
	return service.SysUser().UpdatePassword(ctx, req)
}
