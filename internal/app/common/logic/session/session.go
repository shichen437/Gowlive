package session

import (
	"context"

	"github.com/shichen437/gowlive/internal/app/admin/model/entity"
	"github.com/shichen437/gowlive/internal/app/common/service"
	"github.com/shichen437/gowlive/internal/pkg/consts"
)

type (
	sSession struct{}
)

func init() {
	service.RegisterSession(New())
}

func New() service.ISession {
	return &sSession{}
}

func (s *sSession) SetUser(ctx context.Context, user *entity.SysUser) error {
	return service.BizCtx().Get(ctx).Session.Set(consts.UserSessionKey, user)
}

func (s *sSession) GetUser(ctx context.Context) *entity.SysUser {
	customCtx := service.BizCtx().Get(ctx)
	if customCtx != nil {
		if v := customCtx.Session.MustGet(consts.UserSessionKey); !v.IsNil() {
			var user *entity.SysUser
			_ = v.Struct(&user)
			return user
		}
	}
	return nil
}

func (s *sSession) RemoveUser(ctx context.Context) error {
	customCtx := service.BizCtx().Get(ctx)
	if customCtx != nil {
		return customCtx.Session.Remove(consts.UserSessionKey)
	}
	return nil
}
