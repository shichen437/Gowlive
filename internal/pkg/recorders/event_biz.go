package recorders

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shichen437/gowlive/internal/app/stream/dao"
	"github.com/shichen437/gowlive/internal/app/stream/model/do"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/lives"
	mr "github.com/shichen437/gowlive/internal/pkg/manager"
	"github.com/shichen437/gowlive/internal/pkg/message_push"
	"github.com/shichen437/gowlive/internal/pkg/service"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

func liveStartBiz(ctx context.Context, session *lives.LiveSession) {
	g.Log().Info(ctx, "liveStartBiz", session.Id)
	session.StartedAt = utils.Now()
	go message_push.LivePush(gctx.GetInitCtx(), session.State.Anchor, true)
}

func liveEndBiz(ctx context.Context, session *lives.LiveSession) {
	g.Log().Info(ctx, "liveEndBiz", session.Id)
	startTime := session.StartedAt
	addHistory(ctx, session.Id, session.State.Anchor, startTime, utils.Now())
	enable := mr.GetSettingsManager().GetSetting(consts.SKLiveEndNotify)
	if enable == 1 {
		go message_push.LivePush(gctx.GetInitCtx(), session.State.Anchor, false)
	}
}

func addHistory(ctx context.Context, liveId int, anchor string, startTime, endTime *gtime.Time) {
	if startTime == nil || endTime == nil {
		g.Log().Warningf(ctx, "Invalid start or end time for liveId %d.", liveId)
		return
	}
	_, err := dao.LiveHistory.Ctx(ctx).Insert(do.LiveHistory{
		LiveId:    liveId,
		Anchor:    anchor,
		StartedAt: startTime,
		EndedAt:   endTime,
		Duration:  fmt.Sprintf("%.2f", endTime.Sub(startTime).Hours()),
		CreatedAt: utils.Now(),
	})
	if err != nil {
		g.Log().Errorf(ctx, "Failed to save live history for liveId %d: %v", liveId, err)
	}
}

func (*manager) updateName(ctx context.Context, session *lives.LiveSession) {
	if session == nil || session.Id == 0 {
		return
	}
	service.UpdateRoomInfo(ctx, session)
}
