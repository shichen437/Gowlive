package logic

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	v1 "github.com/shichen437/gowlive/api/v1/system"
	"github.com/shichen437/gowlive/internal/app/system/dao"
	"github.com/shichen437/gowlive/internal/app/system/model/do"
	"github.com/shichen437/gowlive/internal/app/system/model/entity"
	"github.com/shichen437/gowlive/internal/app/system/service"
	"github.com/shichen437/gowlive/internal/pkg/manager"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

type sSysProxy struct {
}

func init() {
	service.RegisterSysProxy(New())
}

func New() service.ISysProxy {
	return &sSysProxy{}
}

func (s *sSysProxy) All(ctx context.Context, req *v1.GetAllProxyReq) (res *v1.GetAllProxyRes, err error) {
	res = &v1.GetAllProxyRes{}
	var list []*entity.SysProxy
	err = dao.SysProxy.Ctx(ctx).OrderDesc(dao.SysProxy.Columns().Id).Scan(&list)
	if err != nil {
		return
	}
	res.Rows = list
	return
}

func (s *sSysProxy) Add(ctx context.Context, req *v1.PostSysProxyReq) (res *v1.PostSysProxyRes, err error) {
	count, err := dao.SysProxy.Ctx(ctx).Where(dao.SysProxy.Columns().Platform, req.Platform).Count()
	if count > 0 {
		err = utils.TError(ctx, "system.proxy.error.Added")
		return
	}
	vp, err := validProxyAndGet(ctx, req.Proxy)
	if err != nil {
		return
	}
	_, err = dao.SysProxy.Ctx(ctx).Insert(do.SysProxy{
		Platform:  req.Platform,
		Proxy:     vp,
		Remark:    req.Remark,
		CreatedAt: utils.Now(),
	})
	if err != nil {
		g.Log().Errorf(ctx, "Failed to add system proxy: %v", err)
		return
	}
	manager.GetProxyManager().SaveProxy(ctx, req.Platform, vp)
	return
}

func (s *sSysProxy) Update(ctx context.Context, req *v1.PutSysProxyReq) (res *v1.PutSysProxyRes, err error) {
	oldData, err := getSysProxyById(ctx, req.Id)
	if err != nil || oldData == nil {
		g.Log().Errorf(ctx, "Failed to get system proxy by id: %v", err)
		return
	}
	vp, err := validProxyAndGet(ctx, req.Proxy)
	if err != nil {
		return
	}
	if oldData.Proxy == vp && oldData.Remark == req.Remark {
		return
	}
	_, err = dao.SysProxy.Ctx(ctx).WherePri(req.Id).Update(do.SysProxy{
		Proxy:     vp,
		Remark:    req.Remark,
		UpdatedAt: utils.Now(),
	})
	if err != nil {
		g.Log().Errorf(ctx, "Failed to update system proxy: %v", err)
		return
	}
	manager.GetProxyManager().SaveProxy(ctx, oldData.Platform, vp)
	return
}

func (s *sSysProxy) Delete(ctx context.Context, req *v1.DeleteSysProxyReq) (res *v1.DeleteSysProxyRes, err error) {
	oldData, err := getSysProxyById(ctx, req.Id)
	if err != nil || oldData == nil {
		g.Log().Errorf(ctx, "Failed to get system proxy by id: %v", err)
		return
	}
	_, err = dao.SysProxy.Ctx(ctx).WherePri(req.Id).Delete()
	if err != nil {
		g.Log().Errorf(ctx, "Failed to delete system proxy: %v", err)
		return
	}
	manager.GetProxyManager().RemoveProxy(ctx, oldData.Platform)
	return
}

func getSysProxyById(ctx context.Context, id int) (cookie *entity.SysProxy, err error) {
	var item *entity.SysProxy
	err = dao.SysProxy.Ctx(ctx).WherePri(id).Scan(&item)
	if err != nil {
		g.Log().Errorf(ctx, "Failed to get system proxy by id: %v", err)
		return
	}
	if item == nil {
		err = utils.TError(ctx, "system.proxy.error.NotFound")
		return
	}
	cookie = item
	return
}

func validProxyAndGet(ctx context.Context, proxies []string) (string, error) {
	if len(proxies) == 0 {
		return "", utils.TError(ctx, "system.proxy.error.ValidUrlBlank")
	}
	if len(proxies) > 30 {
		return "", utils.TError(ctx, "system.proxy.error.UrlLength")
	}

	seen := make(map[string]struct{})
	unique := make([]string, 0, len(proxies))
	for _, p := range proxies {
		trimmed := strings.TrimSpace(p)
		if trimmed == "" {
			continue
		}
		if !(strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "socks5://")) {
			return "", utils.TError(ctx, "system.proxy.error.UrlUnSupported")
		}
		if _, ok := seen[trimmed]; !ok {
			seen[trimmed] = struct{}{}
			unique = append(unique, trimmed)
		}
	}

	if len(unique) == 0 {
		return "", utils.TError(ctx, "system.proxy.error.ValidUrlBlank")
	}
	return strings.Join(unique, ","), nil
}
