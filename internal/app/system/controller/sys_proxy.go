package controller

import (
	"context"

	v1 "github.com/shichen437/gowlive/api/v1/system"
	"github.com/shichen437/gowlive/internal/app/system/service"
)

type sysProxyontroller struct {
}

var SysProxy = sysProxyontroller{}

func (c *sysProxyontroller) All(ctx context.Context, req *v1.GetAllProxyReq) (res *v1.GetAllProxyRes, err error) {
	return service.SysProxy().All(ctx, req)
}

func (c *sysProxyontroller) Add(ctx context.Context, req *v1.PostSysProxyReq) (res *v1.PostSysProxyRes, err error) {
	return service.SysProxy().Add(ctx, req)
}

func (c *sysProxyontroller) Update(ctx context.Context, req *v1.PutSysProxyReq) (res *v1.PutSysProxyRes, err error) {
	return service.SysProxy().Update(ctx, req)
}

func (c *sysProxyontroller) Delete(ctx context.Context, req *v1.DeleteSysProxyReq) (res *v1.DeleteSysProxyRes, err error) {
	return service.SysProxy().Delete(ctx, req)
}
