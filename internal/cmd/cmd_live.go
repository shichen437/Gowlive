package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shichen437/gowlive/internal/app/stream/dao"
	"github.com/shichen437/gowlive/internal/app/stream/model/entity"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/crons"
	"github.com/shichen437/gowlive/internal/pkg/manager"
	"github.com/shichen437/gowlive/internal/pkg/registry"
)

func LiveMonitor() {
	ctx := gctx.GetInitCtx()
	defer g.Log().Info(ctx, "LiveMonitor Started!")
	initCookieRegistry(ctx)
	sessionIds := getLiveSessionIds4Init(ctx)
	if len(sessionIds) > 0 {
		reg := registry.Get()
		reg.AddAll(ctx, sessionIds)
	}
	sessions4Job := getSessionIds4Job(ctx)
	if len(sessions4Job) > 0 {
		for _, session := range sessions4Job {
			crons.AddStreamCron(ctx, session.Id, session.MonitorStartAt, session.MonitorStopAt)
		}
	}
}

func getLiveSessionIds4Init(ctx context.Context) []int {
	var list []*entity.LiveManage
	err := dao.LiveManage.Ctx(ctx).
		WhereIn(dao.LiveManage.Columns().MonitorType, []int{consts.MonitorTypeStart, consts.MonitorTypeIntelligent}).
		Scan(&list)
	if err != nil {
		g.Log().Errorf(ctx, "Failed to get live session ids from database: %v", err)
		return nil
	}
	var ids []int
	if len(list) <= 0 {
		return ids
	}
	for _, v := range list {
		ids = append(ids, v.Id)
	}
	return ids
}

func getSessionIds4Job(ctx context.Context) []*entity.LiveManage {
	var sessionIds []*entity.LiveManage
	dao.LiveManage.Ctx(ctx).
		Where(dao.LiveManage.Columns().MonitorType, 2).
		Scan(&sessionIds)
	return sessionIds
}

func initCookieRegistry(ctx context.Context) {
	var cookies []*entity.LiveCookie
	err := dao.LiveCookie.Ctx(ctx).Scan(&cookies)
	if err != nil {
		g.Log().Errorf(ctx, "Failed to get live cookies from database: %v", err)
		return
	}
	if len(cookies) <= 0 {
		return
	}
	cookieReg := manager.GetCookieManager()
	for _, cookie := range cookies {
		cookieReg.Save(ctx, cookie.Platform, cookie.Cookie)
	}
}
