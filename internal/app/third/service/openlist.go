// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	v1 "github.com/shichen437/gowlive/api/v1/third"
)

type (
	IOpenlist interface {
		Status(ctx context.Context, req *v1.GetOpenlistStatusReq) (res *v1.GetOpenlistStatusRes, err error)
	}
)

var (
	localOpenlist IOpenlist
)

func Openlist() IOpenlist {
	if localOpenlist == nil {
		panic("implement not found for interface IOpenlist, forgot register?")
	}
	return localOpenlist
}

func RegisterOpenlist(i IOpenlist) {
	localOpenlist = i
}
