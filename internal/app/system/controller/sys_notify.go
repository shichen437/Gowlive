package controller

import (
	"context"

	v1 "github.com/shichen437/gowlive/api/v1/system"
	"github.com/shichen437/gowlive/internal/app/system/service"
)

type sysNotifyController struct {
}

var SysNotify = sysNotifyController{}

func (c *sysNotifyController) List(ctx context.Context, req *v1.GetNotifyListReq) (res *v1.GetNotifyListRes, err error) {
	return service.SysNotify().List(ctx, req)
}

func (c *sysNotifyController) MarkRead(ctx context.Context, req *v1.PutMarkNotifyReadReq) (res *v1.PutMarkNotifyReadRes, err error) {
	return service.SysNotify().MarkRead(ctx, req)
}

func (c *sysNotifyController) MarkAll(ctx context.Context, req *v1.PutMarkNotifyAllReadReq) (res *v1.PutMarkNotifyAllReadRes, err error) {
	return service.SysNotify().MarkAll(ctx, req)
}

func (c *sysNotifyController) Delete(ctx context.Context, req *v1.DeleteNotifyReq) (res *v1.DeleteNotifyRes, err error) {
	return service.SysNotify().Delete(ctx, req)
}

func (c *sysNotifyController) DeleteAll(ctx context.Context, req *v1.DeleteAllNotifyReq) (res *v1.DeleteAllNotifyRes, err error) {
	return service.SysNotify().DeleteAll(ctx, req)
}
