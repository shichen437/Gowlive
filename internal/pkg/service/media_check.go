package service

import (
	"context"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shichen437/gowlive/internal/app/media/dao"
	"github.com/shichen437/gowlive/internal/app/media/model/do"
	"github.com/shichen437/gowlive/internal/app/media/model/entity"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

func UpdateCheckTask(ctx context.Context, taskId, progress, status int, createdAt *gtime.Time) error {
	_, err := dao.FileCheckTask.Ctx(ctx).WherePri(taskId).Update(do.FileCheckTask{
		Progress:   progress,
		FileStatus: status,
		Duration:   utils.DiffNowSeconds(createdAt),
		UpdatedAt:  utils.Now(),
	})
	return err
}

func UpdateCheckTaskProgress(ctx context.Context, taskId, progress int) error {
	_, err := dao.FileCheckTask.Ctx(ctx).WherePri(taskId).Update(do.FileCheckTask{
		Progress:  progress,
		UpdatedAt: utils.Now(),
	})
	return err
}

func UpdateCheckTaskFileStatus(ctx context.Context, taskId, status int) error {
	_, err := dao.FileCheckTask.Ctx(ctx).WherePri(taskId).Update(do.FileCheckTask{
		FileStatus: status,
		UpdatedAt:  utils.Now(),
	})
	return err
}

func GetTaskById(ctx context.Context, tid int) *entity.FileCheckTask {
	var task entity.FileCheckTask
	err := dao.FileCheckTask.Ctx(ctx).WherePri(tid).Scan(&task)
	if err != nil {
		return nil
	}
	return &task
}
