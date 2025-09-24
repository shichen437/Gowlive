package controller

import (
	"context"

	v1 "github.com/shichen437/gowlive/api/v1/system"
	"github.com/shichen437/gowlive/internal/app/system/service"
)

type sysLogsController struct {
}

var SysLogs = sysLogsController{}

func (s *sysLogsController) List(ctx context.Context, req *v1.GetLogsListReq) (res *v1.GetLogsListRes, err error) {
	return service.SysLogs().List(ctx, req)
}

func (s *sysLogsController) Delete(ctx context.Context, req *v1.DeleteLogsReq) (res *v1.DeleteLogsRes, err error) {
	return service.SysLogs().Delete(ctx, req)
}

func (s *sysLogsController) DeleteAll(ctx context.Context, req *v1.DeleteAllLogsReq) (res *v1.DeleteAllLogsRes, err error) {
	return service.SysLogs().DeleteAll(ctx, req)
}
