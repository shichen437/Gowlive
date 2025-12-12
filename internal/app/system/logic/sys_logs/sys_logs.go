package logic

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	v1 "github.com/shichen437/gowlive/api/v1/system"
	"github.com/shichen437/gowlive/internal/app/system/dao"
	"github.com/shichen437/gowlive/internal/app/system/model/entity"
	"github.com/shichen437/gowlive/internal/app/system/service"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

type sSysLogs struct {
}

func init() {
	service.RegisterSysLogs(New())
}

func New() service.ISysLogs {
	return &sSysLogs{}
}

func (s *sSysLogs) List(ctx context.Context, req *v1.GetLogsListReq) (res *v1.GetLogsListRes, err error) {
	res = &v1.GetLogsListRes{}
	var list []*entity.SysLogs
	m := dao.SysLogs.Ctx(ctx)
	if req.Type != nil {
		m = m.Where(dao.SysLogs.Columns().Type, req.Type)
	}
	if req.Status != nil {
		m = m.Where(dao.SysLogs.Columns().Status, req.Status)
	}
	res.Total, err = m.Count()
	if res.Total <= 0 || err != nil {
		return
	}
	m = dealSortParams(m, req.Sort)
	err = m.Page(req.PageNum, req.PageSize).Scan(&list)
	if err != nil {
		return nil, utils.TError(ctx, "system.logs.error.GetList")
	}
	res.Rows = list
	return
}

func (s *sSysLogs) Delete(ctx context.Context, req *v1.DeleteLogsReq) (res *v1.DeleteLogsRes, err error) {
	_, err = dao.SysLogs.Ctx(ctx).WherePri(req.Id).Delete()
	if err != nil {
		err = utils.TError(ctx, "system.logs.error.Delete")
	}
	return
}

func (s *sSysLogs) DeleteAll(ctx context.Context, req *v1.DeleteAllLogsReq) (res *v1.DeleteAllLogsRes, err error) {
	_, err = dao.SysLogs.Ctx(ctx).WhereGT(dao.SysLogs.Columns().Id, 0).Delete()
	if err != nil {
		err = utils.TError(ctx, "system.logs.error.DeleteAll")
	}
	return
}

func dealSortParams(m *gdb.Model, sort string) *gdb.Model {
	switch sort {
	case "id:asc":
		m.OrderAsc(dao.SysLogs.Columns().Id)
	case "type:asc":
		m.OrderAsc(dao.SysLogs.Columns().Type)
	case "type:desc":
		m.OrderDesc(dao.SysLogs.Columns().Type)
	case "status:asc":
		m.OrderAsc(dao.SysLogs.Columns().Status)
	case "status:desc":
		m.OrderDesc(dao.SysLogs.Columns().Status)
	default:
		m = m.OrderDesc(dao.SysLogs.Columns().Id)
	}
	return m
}
