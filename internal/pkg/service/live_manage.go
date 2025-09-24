package service

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/gowlive/internal/app/stream/dao"
	"github.com/shichen437/gowlive/internal/app/stream/model/do"
	"github.com/shichen437/gowlive/internal/pkg/lives"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

func GenLiveSessionById(ctx context.Context, liveId int) *lives.LiveSession {
	liveSession := &lives.LiveSession{
		Id: liveId,
	}
	var liveConfig *lives.LiveConfig
	err := dao.LiveManage.Ctx(ctx).WherePri(liveId).Scan(&liveConfig)
	if err != nil || liveConfig == nil {
		g.Log().Error(ctx, "获取直播配置失败", err)
		return nil
	}
	liveSession.Config = *liveConfig
	var liveState *lives.LiveState
	err = dao.LiveRoomInfo.Ctx(ctx).Where(dao.LiveRoomInfo.Columns().LiveId, liveId).Scan(&liveState)
	if err != nil || liveState == nil {
		g.Log().Error(ctx, "获取直播房间信息失败", err)
		return nil
	}
	liveSession.State = *liveState
	return liveSession
}

func GenLiveSessionsByIds(ctx context.Context, liveIds []int) []*lives.LiveSession {
	liveSessions := make([]*lives.LiveSession, 0, len(liveIds))
	var liveConfigs []*lives.LiveConfig
	err := dao.LiveManage.Ctx(ctx).WhereIn(dao.LiveManage.Columns().Id, liveIds).Scan(&liveConfigs)
	if err != nil || liveConfigs == nil || len(liveConfigs) == 0 {
		g.Log().Error(ctx, "获取直播配置失败", err)
		return nil
	}

	var liveRoomInfos []*lives.LiveState
	err = dao.LiveRoomInfo.Ctx(ctx).WhereIn(dao.LiveRoomInfo.Columns().LiveId, liveIds).Scan(&liveRoomInfos)
	if err != nil || liveRoomInfos == nil || len(liveRoomInfos) == 0 {
		g.Log().Error(ctx, "获取直播房间信息失败", err)
		return nil
	}
	liveStateMap := make(map[int]*lives.LiveState)
	for _, liveRoomInfo := range liveRoomInfos {
		liveStateMap[liveRoomInfo.LiveId] = liveRoomInfo
	}
	for _, liveConfig := range liveConfigs {
		liveState, ok := liveStateMap[liveConfig.Id]
		if !ok {
			g.Log().Error(ctx, "live state not found, liveId: %d", liveConfig.Id)
			continue
		}
		liveSessions = append(liveSessions, &lives.LiveSession{
			Id:     liveConfig.Id,
			Config: *liveConfig,
			State:  *liveState,
		})
	}
	return liveSessions
}

func UpdateRoomInfo(ctx context.Context, liveSession *lives.LiveSession) {
	state := liveSession.State
	dao.LiveRoomInfo.Ctx(ctx).Where(dao.LiveRoomInfo.Columns().LiveId, liveSession.Id).
		Update(do.LiveRoomInfo{
			RoomName:  state.RoomName,
			Anchor:    state.Anchor,
			UpdatedAt: utils.Now(),
		})
}
