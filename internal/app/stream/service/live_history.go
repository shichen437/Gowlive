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
	ILiveHistory interface {
		List(ctx context.Context, req *v1.GetLiveHistoryListReq) (res *v1.GetLiveHistoryListRes, err error)
		Delete(ctx context.Context, req *v1.DeleteLiveHistoryReq) (res *v1.DeleteLiveHistoryRes, err error)
		DeleteAll(ctx context.Context, req *v1.DeleteAllHistoryReq) (res *v1.DeleteAllHistoryRes, err error)
	}
)

var (
	localLiveHistory ILiveHistory
)

func LiveHistory() ILiveHistory {
	if localLiveHistory == nil {
		panic("implement not found for interface ILiveHistory, forgot register?")
	}
	return localLiveHistory
}

func RegisterLiveHistory(i ILiveHistory) {
	localLiveHistory = i
}
