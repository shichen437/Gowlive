// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	v1 "github.com/shichen437/gowlive/api/v1/system"
)

type (
	ISysProxy interface {
		All(ctx context.Context, req *v1.GetAllProxyReq) (res *v1.GetAllProxyRes, err error)
		Add(ctx context.Context, req *v1.PostSysProxyReq) (res *v1.PostSysProxyRes, err error)
		Update(ctx context.Context, req *v1.PutSysProxyReq) (res *v1.PutSysProxyRes, err error)
		Delete(ctx context.Context, req *v1.DeleteSysProxyReq) (res *v1.DeleteSysProxyRes, err error)
	}
)

var (
	localSysProxy ISysProxy
)

func SysProxy() ISysProxy {
	if localSysProxy == nil {
		panic("implement not found for interface ISysProxy, forgot register?")
	}
	return localSysProxy
}

func RegisterSysProxy(i ISysProxy) {
	localSysProxy = i
}
