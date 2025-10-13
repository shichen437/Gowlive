package controller

import (
	"context"

	v1 "github.com/shichen437/gowlive/api/v1/stream"
	"github.com/shichen437/gowlive/internal/app/stream/service"
)

type liveManageController struct {
}

var LiveManage = liveManageController{}

func (c *liveManageController) List(ctx context.Context, req *v1.GetRoomListReq) (res *v1.GetRoomListRes, err error) {
	return service.LiveManage().List(ctx, req)
}

func (c *liveManageController) Get(ctx context.Context, req *v1.GetLiveManageReq) (res *v1.GetLiveManageRes, err error) {
	return service.LiveManage().Get(ctx, req)
}

func (c *liveManageController) Add(ctx context.Context, req *v1.PostLiveManageReq) (res *v1.PostLiveManageRes, err error) {
	return service.LiveManage().Add(ctx, req)
}

func (c *liveManageController) Update(ctx context.Context, req *v1.PutLiveManageReq) (res *v1.PutLiveManageRes, err error) {
	return service.LiveManage().Update(ctx, req)
}

func (c *liveManageController) Delete(ctx context.Context, req *v1.DeleteLiveManageReq) (res *v1.DeleteLiveManageRes, err error) {
	return service.LiveManage().Delete(ctx, req)
}

func (c *liveManageController) Start(ctx context.Context, req *v1.PutLiveManageStartReq) (res *v1.PutLiveManageStartRes, err error) {
	return service.LiveManage().Start(ctx, req)
}

func (c *liveManageController) Stop(ctx context.Context, req *v1.PutLiveManageStopReq) (res *v1.PutLiveManageStopRes, err error) {
	return service.LiveManage().Stop(ctx, req)
}
