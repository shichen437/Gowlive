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
	IFileManage interface {
		List(ctx context.Context, req *v1.GetFileListReq) (res *v1.GetFileListRes, err error)
		Delete(ctx context.Context, req *v1.DeleteFileReq) (res *v1.DeleteFileRes, err error)
		Empty(ctx context.Context, req *v1.GetEmptyFolderReq) (res *v1.GetEmptyFolderRes, err error)
		AnchorFilePath(ctx context.Context, req *v1.GetAnchorFilePathReq) (res *v1.GetAnchorFilePathRes, err error)
		Play(ctx context.Context, req *v1.GetFilePlayReq) (res *v1.GetFilePlayRes, err error)
		Clip(ctx context.Context, req *v1.PostFileClipReq) (res *v1.PostFileClipRes, err error)
	}
)

var (
	localFileManage IFileManage
)

func FileManage() IFileManage {
	if localFileManage == nil {
		panic("implement not found for interface IFileManage, forgot register?")
	}
	return localFileManage
}

func RegisterFileManage(i IFileManage) {
	localFileManage = i
}
