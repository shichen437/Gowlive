package logic

import (
	"context"

	v1 "github.com/shichen437/gowlive/api/v1/media"
	"github.com/shichen437/gowlive/internal/app/media/dao"
	"github.com/shichen437/gowlive/internal/app/media/model/do"
	"github.com/shichen437/gowlive/internal/app/media/model/entity"
	"github.com/shichen437/gowlive/internal/app/media/service"
	"github.com/shichen437/gowlive/internal/pkg/manager"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

type (
	sFileCheck struct{}
)

func init() {
	service.RegisterFileCheck(New())
}

func New() service.IFileCheck {
	return &sFileCheck{}
}

func (s *sFileCheck) List(ctx context.Context, req *v1.GetFileCheckListReq) (res *v1.GetFileCheckListRes, err error) {
	res = &v1.GetFileCheckListRes{}
	var list []*entity.FileCheckTask
	m := dao.FileCheckTask.Ctx(ctx)
	if req.Path != "" {
		m = m.WhereLike(dao.FileCheckTask.Columns().Path, req.Path+"%")
	}
	if req.Filename != "" {
		m = m.WhereLike(dao.FileCheckTask.Columns().Filename, req.Filename+"%")
	}
	res.Total, err = m.Count()
	if err != nil || res.Total <= 0 {
		return
	}
	m = m.OrderDesc(dao.FileCheckTask.Columns().Id)
	if err = m.Page(req.PageNum, req.PageSize).Scan(&list); err != nil {
		return
	}
	res.Rows = list
	res.Executing = manager.GetFileCheckManager().Count() > 0
	return
}

func (s *sFileCheck) Post(ctx context.Context, req *v1.PostFileCheckReq) (res *v1.PostFileCheckRes, err error) {
	result, err := dao.FileCheckTask.Ctx(ctx).Insert(do.FileCheckTask{
		Path:      req.Path,
		Filename:  req.Filename,
		CreatedAt: utils.Now(),
	})
	if err != nil {
		return nil, utils.TError(ctx, "media.check.error.Add")
	}
	resultId, err := result.LastInsertId()
	if err != nil {
		return nil, utils.TError(ctx, "media.check.error.GetID")
	}
	manager.GetFileCheckManager().Add(int(resultId))
	return
}

func (s *sFileCheck) Delete(ctx context.Context, req *v1.DeleteFileCheckReq) (res *v1.DeleteFileCheckRes, err error) {
	_, err = dao.FileCheckTask.Ctx(ctx).Where(dao.FileCheckTask.Columns().Id, req.Id).Delete()
	if err != nil {
		return
	}
	return
}

func (s *sFileCheck) DeleteAll(ctx context.Context, req *v1.DeleteAllCheckReq) (res *v1.DeleteAllCheckRes, err error) {
	_, err = dao.FileCheckTask.Ctx(ctx).WhereGT(dao.FileCheckTask.Columns().FileStatus, 0).Delete()
	if err != nil {
		return
	}
	return
}
