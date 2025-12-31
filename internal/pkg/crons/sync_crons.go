package crons

import (
	"context"

	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/manager"
	"github.com/shichen437/gowlive/internal/pkg/service"
	"github.com/shichen437/gowlive/internal/pkg/third/openlist"
)

func SyncFile(ctx context.Context) {
	if manager.GetSettingsManager().GetSetting(consts.SKDataSyncEnable) != 1 ||
		!openlist.CheckLoginParams() {
		return
	}
	m := manager.GetFileSyncManager()
	if m.Count() > 0 {
		return
	}
	pts := service.GetTasksWithStatus(ctx, consts.FileSyncStatusInit)
	if len(pts) != 0 {
		var ids []int
		for _, pt := range pts {
			ids = append(ids, pt.Id)
		}
		m.BatchAdd(ids)
	}

	// 失败重试任务
	if manager.GetSettingsManager().GetSetting(consts.SKDataSyncFailedRetry) != 1 {
		return
	}
	ets := service.GetTasksWithStatus(ctx, consts.FileSyncStatusError)
	if len(ets) == 0 {
		return
	}
	var eids []int
	for _, et := range ets {
		eids = append(eids, et.Id)
	}
	service.BatchUpdateStatus(ctx, consts.FileSyncStatusInit, eids)
	m.BatchAdd(eids)
}
