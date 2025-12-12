package logic

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	v1 "github.com/shichen437/gowlive/api/v1/stream"
	"github.com/shichen437/gowlive/internal/app/stream/dao"
	"github.com/shichen437/gowlive/internal/app/stream/model/do"
	"github.com/shichen437/gowlive/internal/app/stream/model/entity"
	"github.com/shichen437/gowlive/internal/app/stream/service"
	"github.com/shichen437/gowlive/internal/pkg/manager"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

type (
	sLiveCookie struct{}
)

func init() {
	service.RegisterLiveCookie(New())
}

func New() service.ILiveCookie {
	return &sLiveCookie{}
}

func (s *sLiveCookie) All(ctx context.Context, req *v1.GetAllCookieReq) (res *v1.GetAllCookieRes, err error) {
	res = &v1.GetAllCookieRes{}
	var list []*entity.LiveCookie
	err = dao.LiveCookie.Ctx(ctx).OrderDesc(dao.LiveCookie.Columns().Id).Scan(&list)
	if err != nil {
		return
	}
	res.Rows = list
	return
}

func (s *sLiveCookie) Add(ctx context.Context, req *v1.PostLiveCookieReq) (res *v1.PostLiveCookieRes, err error) {
	count, err := dao.LiveCookie.Ctx(ctx).Where(dao.LiveCookie.Columns().Platform, req.Platform).Count()
	if count > 0 {
		err = utils.TError(ctx, "stream.cookie.error.Added")
		return
	}
	_, err = dao.LiveCookie.Ctx(ctx).Insert(do.LiveCookie{
		Platform:  req.Platform,
		Cookie:    req.Cookie,
		Remark:    req.Remark,
		CreatedAt: utils.Now(),
	})
	if err != nil {
		g.Log().Errorf(ctx, "Failed to add live cookie: %v", err)
		return
	}
	manager.GetCookieManager().Save(ctx, req.Platform, req.Cookie)
	return
}

func (s *sLiveCookie) Update(ctx context.Context, req *v1.PutLiveCookieReq) (res *v1.PutLiveCookieRes, err error) {
	oldData, err := getLiveCookieById(ctx, req.Id)
	if err != nil || oldData == nil {
		g.Log().Errorf(ctx, "Failed to get live cookie by id: %v", err)
		return
	}
	if oldData.Cookie == req.Cookie && oldData.Remark == req.Remark {
		return
	}
	_, err = dao.LiveCookie.Ctx(ctx).WherePri(req.Id).Update(do.LiveCookie{
		Cookie:    req.Cookie,
		Remark:    req.Remark,
		UpdatedAt: utils.Now(),
	})
	if err != nil {
		g.Log().Errorf(ctx, "Failed to update live cookie: %v", err)
		return
	}
	manager.GetCookieManager().Save(ctx, oldData.Platform, req.Cookie)
	return
}

func (s *sLiveCookie) Delete(ctx context.Context, req *v1.DeleteLiveCookieReq) (res *v1.DeleteLiveCookieRes, err error) {
	oldData, err := getLiveCookieById(ctx, req.Id)
	if err != nil || oldData == nil {
		g.Log().Errorf(ctx, "Failed to get live cookie by id: %v", err)
		return
	}
	_, err = dao.LiveCookie.Ctx(ctx).WherePri(req.Id).Delete()
	if err != nil {
		g.Log().Errorf(ctx, "Failed to delete live cookie: %v", err)
		return
	}
	manager.GetCookieManager().Remove(ctx, oldData.Platform)
	return
}

func getLiveCookieById(ctx context.Context, id int) (cookie *entity.LiveCookie, err error) {
	var item *entity.LiveCookie
	err = dao.LiveCookie.Ctx(ctx).WherePri(id).Scan(&item)
	if err != nil {
		g.Log().Errorf(ctx, "Failed to get live cookie by id: %v", err)
		return
	}
	if item == nil {
		err = utils.TError(ctx, "stream.cookie.error.NotFound")
		return
	}
	cookie = item
	return
}
