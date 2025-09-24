// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	v1 "github.com/shichen437/gowlive/api/v1/stream"
)

type (
	ILiveManage interface {
		List(ctx context.Context, req *v1.GetRoomListReq) (res *v1.GetRoomListRes, err error)
		Get(ctx context.Context, req *v1.GetLiveManageReq) (res *v1.GetLiveManageRes, err error)
		Add(ctx context.Context, req *v1.PostLiveManageReq) (res *v1.PostLiveManageRes, err error)
		Update(ctx context.Context, req *v1.PutLiveManageReq) (res *v1.PutLiveManageRes, err error)
		Delete(ctx context.Context, req *v1.DeleteLiveManageReq) (res *v1.DeleteLiveManageRes, err error)
	}
)

var (
	localLiveManage ILiveManage
)

func LiveManage() ILiveManage {
	if localLiveManage == nil {
		panic("implement not found for interface ILiveManage, forgot register?")
	}
	return localLiveManage
}

func RegisterLiveManage(i ILiveManage) {
	localLiveManage = i
}
