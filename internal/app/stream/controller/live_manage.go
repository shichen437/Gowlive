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

func (c *liveManageController) BatchAdd(ctx context.Context, req *v1.PostLiveManageBatchReq) (res *v1.PostLiveManageBatchRes, err error) {
	return service.LiveManage().BatchAdd(ctx, req)
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

func (c *liveManageController) Top(ctx context.Context, req *v1.PutLiveManageTopReq) (res *v1.PutLiveManageTopRes, err error) {
	return service.LiveManage().Top(ctx, req)
}

func (c *liveManageController) UnTop(ctx context.Context, req *v1.PutLiveManageUnTopReq) (res *v1.PutLiveManageUnTopRes, err error) {
	return service.LiveManage().UnTop(ctx, req)
}

func (c *liveManageController) Export(ctx context.Context, req *v1.ExportRoomInfoReq) (res *v1.ExportRoomInfoRes, err error) {
	return service.LiveManage().Export(ctx, req)
}

func (c *liveManageController) Preview(ctx context.Context, req *v1.PreviewRoomReq) (res *v1.PreviewRoomRes, err error) {
	return service.LiveManage().Preview(ctx, req)
}

func (c *liveManageController) PreviewList(ctx context.Context, req *v1.PreviewRoomListReq) (res *v1.PreviewRoomListRes, err error) {
	return service.LiveManage().PreviewList(ctx, req)
}

func (c *liveManageController) QuickAdd(ctx context.Context, req *v1.PostQuickLinkReq) (res *v1.PostQuickLinkRes, err error) {
	return service.LiveManage().QuickAdd(ctx, req)
}
