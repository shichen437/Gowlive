package manager

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/service"
	"github.com/shichen437/gowlive/internal/pkg/third/openlist"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

const (
	syncMaxTasks   = 1000
	syncWorkerNums = 3
)

var (
	fileSyncManager *FileSyncManager
	syncOnce        sync.Once
)

type FileSyncManager struct {
	tasks  chan int
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	mu       sync.Mutex
	inFlight int
}

func GetFileSyncManager() *FileSyncManager {
	syncOnce.Do(func() {
		ctx, cancel := context.WithCancel(gctx.New())
		fileSyncManager = &FileSyncManager{
			tasks:  make(chan int, syncMaxTasks),
			ctx:    ctx,
			cancel: cancel,
		}
		for range make([]int, syncWorkerNums) {
			fileSyncManager.wg.Add(1)
			go fileSyncManager.autoExec()
		}
	})
	return fileSyncManager
}

func (f *FileSyncManager) Add(sid int) {
	f.tasks <- sid
}

func (f *FileSyncManager) BatchAdd(sids []int) {
	for _, sid := range sids {
		f.Add(sid)
	}
}

func (f *FileSyncManager) Count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.inFlight + len(f.tasks)
}

func (f *FileSyncManager) Close() {
	f.cancel()
	f.wg.Wait()
	close(f.tasks)
}

func (f *FileSyncManager) autoExec() {
	defer f.wg.Done()
	for {
		select {
		case <-f.ctx.Done():
			return
		case tid := <-f.tasks:
			f.mu.Lock()
			f.inFlight++
			f.mu.Unlock()

			syncFile(f.ctx, tid)

			f.mu.Lock()
			if f.inFlight > 0 {
				f.inFlight--
			}
			f.mu.Unlock()
		}
	}
}

func syncFile(ctx context.Context, sid int) {
	if sid <= 0 {
		return
	}
	current := utils.Now()
	task := service.GetSyncTaskById(ctx, sid)
	if task == nil || sid != task.Id || task.Status != consts.FileSyncStatusInit {
		return
	}
	absPath, err := utils.FileAbsPath(task.Path, task.Filename)
	if err != nil || absPath == "" {
		service.UpdateSyncTaskStatus(ctx, task.Id, consts.FileSyncStatusFileNotFound, current)
		return
	}
	if !gfile.Exists(absPath) {
		service.UpdateSyncTaskStatus(ctx, task.Id, consts.FileSyncStatusFileNotFound, current)
		return
	}
	service.UpdateSyncTaskStatus(ctx, task.Id, consts.FileSyncStatusUploading, current)
	err = openlist.Upload(task.SyncPath, task.Path, task.Filename)
	if err != nil {
		g.Log().Error(ctx, "file sync error : ", err)
		service.UpdateSyncTaskStatus(ctx, task.Id, consts.FileSyncStatusError, current)
		return
	}
	service.UpdateSyncTaskStatus(ctx, task.Id, consts.FileSyncStatusSuccess, current)
}
