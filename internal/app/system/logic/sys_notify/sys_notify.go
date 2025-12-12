package logic

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	v1 "github.com/shichen437/gowlive/api/v1/system"
	"github.com/shichen437/gowlive/internal/app/system/dao"
	"github.com/shichen437/gowlive/internal/app/system/model/do"
	"github.com/shichen437/gowlive/internal/app/system/service"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

type sSysNotify struct {
}

func init() {
	service.RegisterSysNotify(New())
}

func New() service.ISysNotify {
	return &sSysNotify{}
}

func (s *sSysNotify) List(ctx context.Context, req *v1.GetNotifyListReq) (res *v1.GetNotifyListRes, err error) {
	res = &v1.GetNotifyListRes{}
	m := dao.SysNotify.Ctx(ctx)
	if req.Status != nil {
		m = m.Where(dao.SysNotify.Columns().Status, req.Status)
	}
	res.Total, err = m.Count()
	if err != nil {
		return nil, utils.TError(ctx, "system.push.error.GetList")
	}
	if res.Total <= 0 {
		return
	}
	m = m.OrderAsc(dao.SysNotify.Columns().Status).OrderDesc(dao.SysNotify.Columns().Id)
	err = m.Page(req.PageNum, req.PageSize).Scan(&res.Rows)
	if err != nil {
		return nil, utils.TError(ctx, "system.push.error.GetList")
	}
	return
}

func (s *sSysNotify) MarkRead(ctx context.Context, req *v1.PutMarkNotifyReadReq) (res *v1.PutMarkNotifyReadRes, err error) {
	_, err = dao.SysNotify.Ctx(ctx).WherePri(req.Id).Update(do.SysNotify{
		Status:    1,
		UpdatedAt: utils.Now(),
	})
	if err != nil {
		g.Log().Error(ctx, err)
		err = utils.TError(ctx, "system.push.error.Mark")
	}
	return
}

func (s *sSysNotify) MarkAll(ctx context.Context, req *v1.PutMarkNotifyAllReadReq) (res *v1.PutMarkNotifyAllReadRes, err error) {
	_, err = dao.SysNotify.Ctx(ctx).Where(dao.SysNotify.Columns().Status, 0).Update(do.SysNotify{
		Status:    1,
		UpdatedAt: utils.Now(),
	})
	if err != nil {
		g.Log().Error(ctx, err)
		err = utils.TError(ctx, "system.push.error.MarkAll")
	}
	return
}

func (s *sSysNotify) Delete(ctx context.Context, req *v1.DeleteNotifyReq) (res *v1.DeleteNotifyRes, err error) {
	_, err = dao.SysNotify.Ctx(ctx).WherePri(req.Id).Delete()
	if err != nil {
		g.Log().Error(ctx, err)
		err = utils.TError(ctx, "system.push.error.Delete")
	}
	return
}

func (s *sSysNotify) DeleteAll(ctx context.Context, req *v1.DeleteAllNotifyReq) (res *v1.DeleteAllNotifyRes, err error) {
	_, err = dao.SysNotify.Ctx(ctx).WhereGT(dao.SysLogs.Columns().Id, 0).Delete()
	if err != nil {
		g.Log().Error(ctx, err)
		err = utils.TError(ctx, "system.push.error.DeleteAll")
	}
	return
}
