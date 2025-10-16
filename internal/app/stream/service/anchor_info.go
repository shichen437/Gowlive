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
	IAnchorInfo interface {
		List(ctx context.Context, req *v1.GetAnchorListReq) (res *v1.GetAnchorListRes, err error)
		Add(ctx context.Context, req *v1.PostAnchorReq) (res *v1.PostAnchorRes, err error)
		Delete(ctx context.Context, req *v1.DeleteAnchorReq) (res *v1.DeleteAnchorRes, err error)
		StatInfo(ctx context.Context, req *v1.GetAnchorStatInfoReq) (res *v1.GetAnchorStatInfoRes, err error)
	}
)

var (
	localAnchorInfo IAnchorInfo
)

func AnchorInfo() IAnchorInfo {
	if localAnchorInfo == nil {
		panic("implement not found for interface IAnchorInfo, forgot register?")
	}
	return localAnchorInfo
}

func RegisterAnchorInfo(i IAnchorInfo) {
	localAnchorInfo = i
}
