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
	ISysNotify interface {
		List(ctx context.Context, req *v1.GetNotifyListReq) (res *v1.GetNotifyListRes, err error)
		MarkRead(ctx context.Context, req *v1.PutMarkNotifyReadReq) (res *v1.PutMarkNotifyReadRes, err error)
		MarkAll(ctx context.Context, req *v1.PutMarkNotifyAllReadReq) (res *v1.PutMarkNotifyAllReadRes, err error)
		Delete(ctx context.Context, req *v1.DeleteNotifyReq) (res *v1.DeleteNotifyRes, err error)
		DeleteAll(ctx context.Context, req *v1.DeleteAllNotifyReq) (res *v1.DeleteAllNotifyRes, err error)
	}
)

var (
	localSysNotify ISysNotify
)

func SysNotify() ISysNotify {
	if localSysNotify == nil {
		panic("implement not found for interface ISysNotify, forgot register?")
	}
	return localSysNotify
}

func RegisterSysNotify(i ISysNotify) {
	localSysNotify = i
}
