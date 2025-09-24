package controller

import (
	"context"

	v1 "github.com/shichen437/gowlive/api/v1/media"
	"github.com/shichen437/gowlive/internal/app/media/service"
)

type fileManageController struct {
}

var FileManage = fileManageController{}

func (f *fileManageController) List(ctx context.Context, req *v1.GetFileListReq) (res *v1.GetFileListRes, err error) {
	return service.FileManage().List(ctx, req)
}

func (f *fileManageController) Delete(ctx context.Context, req *v1.DeleteFileReq) (res *v1.DeleteFileRes, err error) {
	return service.FileManage().Delete(ctx, req)
}

func (f *fileManageController) Play(ctx context.Context, req *v1.GetFilePlayReq) (res *v1.GetFilePlayRes, err error) {
	return service.FileManage().Play(ctx, req)
}
