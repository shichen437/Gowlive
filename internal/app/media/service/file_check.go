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
	IFileCheck interface {
		List(ctx context.Context, req *v1.GetFileCheckListReq) (res *v1.GetFileCheckListRes, err error)
		Post(ctx context.Context, req *v1.PostFileCheckReq) (res *v1.PostFileCheckRes, err error)
		Delete(ctx context.Context, req *v1.DeleteFileCheckReq) (res *v1.DeleteFileCheckRes, err error)
	}
)

var (
	localFileCheck IFileCheck
)

func FileCheck() IFileCheck {
	if localFileCheck == nil {
		panic("implement not found for interface IFileCheck, forgot register?")
	}
	return localFileCheck
}

func RegisterFileCheck(i IFileCheck) {
	localFileCheck = i
}
