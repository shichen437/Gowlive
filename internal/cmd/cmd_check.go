package cmd

import (
	"context"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shichen437/gowlive/internal/app/media/dao"
	"github.com/shichen437/gowlive/internal/app/media/model/entity"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/manager"
)

func CheckFile() {
	ctx := gctx.GetInitCtx()
	tasks := getTaskIds4Check(ctx)
	if len(tasks) == 0 {
		return
	}
	ids := make([]int, 0, len(tasks))
	for _, task := range tasks {
		ids = append(ids, task.Id)
	}

	go manager.GetFileCheckManager().BatchAdd(ids)
}

func getTaskIds4Check(ctx context.Context) []*entity.FileCheckTask {
	var tasks []*entity.FileCheckTask
	dao.FileCheckTask.Ctx(ctx).
		Where(dao.FileCheckTask.Columns().FileStatus, consts.MediaCheckFileStatusInit).
		Scan(&tasks)
	return tasks
}
