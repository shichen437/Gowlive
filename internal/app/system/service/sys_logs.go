// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	v1 "github.com/shichen437/gowlive/api/v1/system"
)

type (
	ISysLogs interface {
		List(ctx context.Context, req *v1.GetLogsListReq) (res *v1.GetLogsListRes, err error)
		Delete(ctx context.Context, req *v1.DeleteLogsReq) (res *v1.DeleteLogsRes, err error)
		DeleteAll(ctx context.Context, req *v1.DeleteAllLogsReq) (res *v1.DeleteAllLogsRes, err error)
	}
)

var (
	localSysLogs ISysLogs
)

func SysLogs() ISysLogs {
	if localSysLogs == nil {
		panic("implement not found for interface ISysLogs, forgot register?")
	}
	return localSysLogs
}

func RegisterSysLogs(i ISysLogs) {
	localSysLogs = i
}
