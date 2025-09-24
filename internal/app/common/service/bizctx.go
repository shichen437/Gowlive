// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/shichen437/gowlive/internal/app/common/model"
)

type (
	IBizCtx interface {
		Init(r *ghttp.Request, customCtx *model.Context)
		Get(ctx context.Context) *model.Context
		SetUser(ctx context.Context, ctxUser *model.ContextUser)
	}
)

var (
	localBizCtx IBizCtx
)

func BizCtx() IBizCtx {
	if localBizCtx == nil {
		panic("implement not found for interface IBizCtx, forgot register?")
	}
	return localBizCtx
}

func RegisterBizCtx(i IBizCtx) {
	localBizCtx = i
}
