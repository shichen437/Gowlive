package recorders

import (
	"context"
	"fmt"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shichen437/gowlive/internal/app/stream/dao"
	"github.com/shichen437/gowlive/internal/app/stream/model/do"
	"github.com/shichen437/gowlive/internal/pkg/lives"
	"github.com/shichen437/gowlive/internal/pkg/service"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

var (
	liveStartTimes = &sync.Map{}
)

func liveStartBiz(ctx context.Context, liveId int) {
	g.Log().Info(ctx, "liveStartBiz", liveId)
	liveStartTimes.Store(liveId, utils.Now())
}

func liveEndBiz(ctx context.Context, liveId int, anchor string) {
	g.Log().Info(ctx, "liveEndBiz", liveId)
	startTime := getStartTime(ctx, liveId)
	addHistory(ctx, liveId, anchor, startTime, utils.Now())
}

func addHistory(ctx context.Context, liveId int, anchor string, startTime, endTime *gtime.Time) {
	if startTime != nil {
		liveStartTimes.Delete(liveId)
	}
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

func getStartTime(ctx context.Context, liveId int) *gtime.Time {
	iStartTime, ok := liveStartTimes.Load(liveId)
	if !ok {
		return nil
	}

	startTime, ok := iStartTime.(*gtime.Time)
	if !ok {
		g.Log().Errorf(ctx, "Invalid start time type in map for liveId %d", liveId)
		return nil
	}

	return startTime
}

func (*manager) updateName(ctx context.Context, session *lives.LiveSession) {
	if session == nil || session.Id == 0 {
		return
	}
	service.UpdateRoomInfo(ctx, session)
}
