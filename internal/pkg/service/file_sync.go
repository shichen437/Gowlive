package service

import (
	"context"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shichen437/gowlive/internal/app/media/dao"
	"github.com/shichen437/gowlive/internal/app/media/model/do"
	"github.com/shichen437/gowlive/internal/app/media/model/entity"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

func UpdateSyncTaskStatus(ctx context.Context, taskId, status int, createdAt *gtime.Time) error {
	_, err := dao.FileSyncTask.Ctx(ctx).WherePri(taskId).Update(do.FileSyncTask{
		Status:    status,
		Duration:  utils.DiffNowSeconds(createdAt),
		UpdatedAt: utils.Now(),
	})
	return err
}

func UpdateSyncTaskStatusWithError(ctx context.Context, taskId, status int, errMsg string, createdAt *gtime.Time) error {
	_, err := dao.FileSyncTask.Ctx(ctx).WherePri(taskId).Update(do.FileSyncTask{
		Status:    status,
		Duration:  utils.DiffNowSeconds(createdAt),
		Remark:    errMsg,
		UpdatedAt: utils.Now(),
	})
	return err
}

func AddSyncTask(ctx context.Context, path, filename, syncPath string) int {
	if path == "" || filename == "" || syncPath == "" {
		return 0
	}
	task := do.FileSyncTask{
		Path:      path,
		Filename:  filename,
		SyncPath:  syncPath,
		CreatedAt: utils.Now(),
	}
	id, err := dao.FileSyncTask.Ctx(ctx).InsertAndGetId(task)
	if err != nil {
		return 0
	}
	return int(id)
}

func GetSyncTaskById(ctx context.Context, sid int) *entity.FileSyncTask {
	var task entity.FileSyncTask
	err := dao.FileSyncTask.Ctx(ctx).WherePri(sid).Scan(&task)
	if err != nil {
		return nil
	}
	return &task
}

func GetTasksWithStatus(ctx context.Context, status int) []*entity.FileSyncTask {
	var tasks []*entity.FileSyncTask
	err := dao.FileSyncTask.Ctx(ctx).
		Where(dao.FileSyncTask.Columns().Status, status).
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

func BatchUpdateStatus(ctx context.Context, status int, ids []int) {
	dao.FileSyncTask.Ctx(ctx).WhereIn(dao.FileSyncTask.Columns().Id, ids).Update(do.FileSyncTask{
		Status:    status,
		UpdatedAt: utils.Now(),
	})
}
