package controller

import (
	"context"

	v1 "github.com/shichen437/gowlive/api/v1/system"
	"github.com/shichen437/gowlive/internal/app/system/service"
)

type pushChannelController struct {
}

var PushChannel = pushChannelController{}

func (c *pushChannelController) List(ctx context.Context, req *v1.GetPushChannelListReq) (res *v1.GetPushChannelListRes, err error) {
	return service.PushChannel().List(ctx, req)
}

func (c *pushChannelController) Post(ctx context.Context, req *v1.PostPushChannelReq) (res *v1.PostPushChannelRes, err error) {
	return service.PushChannel().Post(ctx, req)
}

func (c *pushChannelController) Put(ctx context.Context, req *v1.PutPushChannelReq) (res *v1.PutPushChannelRes, err error) {
	return service.PushChannel().Put(ctx, req)
}

func (c *pushChannelController) Get(ctx context.Context, req *v1.GetPushChannelReq) (res *v1.GetPushChannelRes, err error) {
	return service.PushChannel().Get(ctx, req)
}

func (c *pushChannelController) Delete(ctx context.Context, req *v1.DeletePushChannelReq) (res *v1.DeletePushChannelRes, err error) {
	return service.PushChannel().Delete(ctx, req)
}
