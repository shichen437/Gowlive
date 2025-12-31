package logic

import (
	"context"

	v1 "github.com/shichen437/gowlive/api/v1/media"
	"github.com/shichen437/gowlive/internal/app/media/dao"
	"github.com/shichen437/gowlive/internal/app/media/model/do"
	"github.com/shichen437/gowlive/internal/app/media/model/entity"
	"github.com/shichen437/gowlive/internal/app/media/service"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/manager"
)

type (
	sFileSync struct{}
)

func init() {
	service.RegisterFileSync(New())
}

func New() service.IFileSync {
	return &sFileSync{}
}

func (s *sFileSync) List(ctx context.Context, req *v1.GetFileSyncListReq) (res *v1.GetFileSyncListRes, err error) {
	res = &v1.GetFileSyncListRes{}
	var list []*entity.FileSyncTask
	m := dao.FileSyncTask.Ctx(ctx)
	if req.Path != "" {
		m = m.WhereLike(dao.FileSyncTask.Columns().Path, req.Path+"%")
	}
	if req.Filename != "" {
		m = m.WhereLike(dao.FileSyncTask.Columns().Filename, req.Filename+"%")
	}
	res.Total, err = m.Count()
	if err != nil || res.Total <= 0 {
		return
	}
	m = m.OrderDesc(dao.FileSyncTask.Columns().Id)
	if err = m.Page(req.PageNum, req.PageSize).Scan(&list); err != nil {
		return
	}
	res.Rows = list
	return
}

func (s *sFileSync) Resync(ctx context.Context, req *v1.ResyncFileReq) (res *v1.ResyncFileRes, err error) {
	count, err := dao.FileSyncTask.Ctx(ctx).Where(dao.FileSyncTask.Columns().Id, req.Id).
		Where(dao.FileSyncTask.Columns().Status, consts.FileSyncStatusError).Count()
	if err != nil || count <= 0 {
		return
	}
	_, err = dao.FileSyncTask.Ctx(ctx).WherePri(req.Id).Update(do.FileSyncTask{
		Status: consts.FileSyncStatusInit,
	})
	if err != nil {
		return
	}
	manager.GetFileSyncManager().Add(req.Id)
	return
}

func (s *sFileSync) Delete(ctx context.Context, req *v1.DeleteFileSyncReq) (res *v1.DeleteFileSyncRes, err error) {
	_, err = dao.FileSyncTask.Ctx(ctx).Where(dao.FileSyncTask.Columns().Id, req.Id).Delete()
	if err != nil {
		return
	}
	return
}

func (s *sFileSync) DeleteAll(ctx context.Context, req *v1.DeleteAllSyncReq) (res *v1.DeleteAllSyncRes, err error) {
	_, err = dao.FileSyncTask.Ctx(ctx).WhereGT(dao.FileSyncTask.Columns().Status, 1).Delete()
	if err != nil {
		return
	}
	return
}
