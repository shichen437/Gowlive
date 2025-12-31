package controller

import (
	"context"

	v1 "github.com/shichen437/gowlive/api/v1/media"
	"github.com/shichen437/gowlive/internal/app/media/service"
)

type fileSyncController struct {
}

var FileSync = fileSyncController{}

func (f *fileSyncController) List(ctx context.Context, req *v1.GetFileSyncListReq) (res *v1.GetFileSyncListRes, err error) {
	return service.FileSync().List(ctx, req)
}

func (f *fileSyncController) Resync(ctx context.Context, req *v1.ResyncFileReq) (res *v1.ResyncFileRes, err error) {
	return service.FileSync().Resync(ctx, req)
}

func (f *fileSyncController) Delete(ctx context.Context, req *v1.DeleteFileSyncReq) (res *v1.DeleteFileSyncRes, err error) {
	return service.FileSync().Delete(ctx, req)
}

func (f *fileSyncController) DeleteAll(ctx context.Context, req *v1.DeleteAllSyncReq) (res *v1.DeleteAllSyncRes, err error) {
	return service.FileSync().DeleteAll(ctx, req)
}
