package controller

import (
	"context"

	v1 "github.com/shichen437/gowlive/api/v1/media"
	"github.com/shichen437/gowlive/internal/app/media/service"
)

type fileCheckController struct {
}

var FileCheck = fileCheckController{}

func (f *fileCheckController) List(ctx context.Context, req *v1.GetFileCheckListReq) (res *v1.GetFileCheckListRes, err error) {
	return service.FileCheck().List(ctx, req)
}

func (f *fileCheckController) Post(ctx context.Context, req *v1.PostFileCheckReq) (res *v1.PostFileCheckRes, err error) {
	return service.FileCheck().Post(ctx, req)
}

func (f *fileCheckController) Delete(ctx context.Context, req *v1.DeleteFileCheckReq) (res *v1.DeleteFileCheckRes, err error) {
	return service.FileCheck().Delete(ctx, req)
}
