// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"github.com/shichen437/gowlive/api/v1/common"
)

type (
	IInternalDict interface {
		GetDictDataByType(ctx context.Context, req *common.GetInternalDictByTypeReq) (res *common.GetInternalDictByTypeRes, err error)
	}
)

var (
	localInternalDict IInternalDict
)

func InternalDict() IInternalDict {
	if localInternalDict == nil {
		panic("implement not found for interface IInternalDict, forgot register?")
	}
	return localInternalDict
}

func RegisterInternalDict(i IInternalDict) {
	localInternalDict = i
}
