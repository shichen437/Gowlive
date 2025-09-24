package logic

import (
	"context"

	v1 "github.com/shichen437/gowlive/api/v1/stream"
	"github.com/shichen437/gowlive/internal/app/stream/dao"
	"github.com/shichen437/gowlive/internal/app/stream/service"
)

type (
	sLiveHistory struct{}
)

func init() {
	service.RegisterLiveHistory(New())
}

func New() service.ILiveHistory {
	return &sLiveHistory{}
}

func (s *sLiveHistory) List(ctx context.Context, req *v1.GetLiveHistoryListReq) (res *v1.GetLiveHistoryListRes, err error) {
	res = &v1.GetLiveHistoryListRes{}
	m := dao.LiveHistory.Ctx(ctx)
	if req.LiveId != nil && *req.LiveId != 0 {
		m = m.Where(dao.LiveHistory.Columns().LiveId, *req.LiveId)
	}
	m = m.OrderDesc(dao.LiveHistory.Columns().Id)
	res.Total, err = m.Count()
	if err != nil {
		return
	}
	if res.Total <= 0 {
		return
	}
	err = m.Page(req.PageNum, req.PageSize).Scan(&res.Rows)
	if err != nil {
		return
	}
	return
}

func (s *sLiveHistory) Delete(ctx context.Context, req *v1.DeleteLiveHistoryReq) (res *v1.DeleteLiveHistoryRes, err error) {
	_, err = dao.LiveHistory.Ctx(ctx).WherePri(req.Id).Delete()
	return
}
