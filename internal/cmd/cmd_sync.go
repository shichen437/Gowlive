package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shichen437/gowlive/internal/app/media/dao"
	"github.com/shichen437/gowlive/internal/app/media/model/entity"
	"github.com/shichen437/gowlive/internal/pkg/manager"
	"github.com/shichen437/gowlive/internal/pkg/third/openlist"
)

func SyncFile() {
	ctx := gctx.GetInitCtx()
	if openlist.CheckLoginParams() {
		err := openlist.Login()
		if err != nil {
			g.Log().Info(ctx, err)
			return
		}
		g.Log().Info(ctx, "Openlist Login successful")
	}
	m := manager.GetFileSyncManager()
	tasks := getSyncTasks(ctx)
	if tasks == nil {
		return
	}
	var ids []int
	for _, task := range tasks {
		ids = append(ids, task.Id)
	}
	m.BatchAdd(ids)
}

func getSyncTasks(ctx context.Context) []*entity.FileSyncTask {
	var tasks []*entity.FileSyncTask
	err := dao.FileSyncTask.Ctx(ctx).
		WhereIn(dao.FileSyncTask.Columns().Status, []int{0, 1}).
		OrderAsc(dao.FileSyncTask.Columns().Id).
		Limit(100).Scan(&tasks)
	if err != nil {
		return nil
	}
	if len(tasks) == 0 {
		return nil
	}
	return tasks
}
