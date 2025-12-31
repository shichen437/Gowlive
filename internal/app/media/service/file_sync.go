// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	v1 "github.com/shichen437/gowlive/api/v1/media"
)

type (
	IFileSync interface {
		List(ctx context.Context, req *v1.GetFileSyncListReq) (res *v1.GetFileSyncListRes, err error)
		Resync(ctx context.Context, req *v1.ResyncFileReq) (res *v1.ResyncFileRes, err error)
		Delete(ctx context.Context, req *v1.DeleteFileSyncReq) (res *v1.DeleteFileSyncRes, err error)
		DeleteAll(ctx context.Context, req *v1.DeleteAllSyncReq) (res *v1.DeleteAllSyncRes, err error)
	}
)

var (
	localFileSync IFileSync
)

func FileSync() IFileSync {
	if localFileSync == nil {
		panic("implement not found for interface IFileSync, forgot register?")
	}
	return localFileSync
}

func RegisterFileSync(i IFileSync) {
	localFileSync = i
}
