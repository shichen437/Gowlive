package manager

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/service"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

const (
	maxTasks   = 1000
	workerNums = 1
)

var (
	fileCheckManager *FileCheckManager
	once             sync.Once
)

type FileCheckManager struct {
	tasks  chan int
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	mu       sync.Mutex
	inFlight int
}

func GetFileCheckManager() *FileCheckManager {
	once.Do(func() {
		ctx, cancel := context.WithCancel(gctx.New())
		fileCheckManager = &FileCheckManager{
			tasks:  make(chan int, maxTasks),
			ctx:    ctx,
			cancel: cancel,
		}
		for i := 0; i < workerNums; i++ {
			fileCheckManager.wg.Add(1)
			go fileCheckManager.autoExec()
		}
	})
	return fileCheckManager
}

func (f *FileCheckManager) Add(tid int) {
	f.tasks <- tid
}

func (f *FileCheckManager) BatchAdd(tids []int) {
	for _, tid := range tids {
		f.Add(tid)
	}
}

func (f *FileCheckManager) Count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.inFlight + len(f.tasks)
}

func (f *FileCheckManager) autoExec() {
	defer f.wg.Done()
	for {
		select {
		case <-f.ctx.Done():
			return
		case tid := <-f.tasks:
			f.mu.Lock()
			f.inFlight++
			f.mu.Unlock()

			checkMedia(f.ctx, tid)

			f.mu.Lock()
			if f.inFlight > 0 {
				f.inFlight--
			}
			f.mu.Unlock()
		}
	}
}

func (f *FileCheckManager) Close() {
	f.cancel()
	f.wg.Wait()
	close(f.tasks)
}

func checkMedia(ctx context.Context, tid int) {
	if tid <= 0 {
		return
	}
	current := utils.Now()
	task := service.GetTaskById(ctx, tid)
	if task == nil || tid != task.Id || task.FileStatus != consts.MediaCheckFileStatusInit {
		service.UpdateCheckTaskFileStatus(ctx, tid, consts.MediaCheckFileStatusNotExists)
		return
	}
	absPath, err := utils.FileAbsPath(task.Path, task.Filename)
	if err != nil || absPath == "" {
		service.UpdateCheckTaskFileStatus(ctx, task.Id, consts.MediaCheckFileStatusNotExists)
		return
	}
	service.UpdateCheckTaskProgress(ctx, task.Id, consts.MediaCheckProgressQcing)
	err = utils.QuickCheckFile(ctx, absPath)
	if err != nil {
		service.UpdateCheckTask(ctx, task.Id, consts.MediaCheckProgressQcError, consts.MediaCheckFileStatusError, current)
		return
	}
	service.UpdateCheckTaskProgress(ctx, task.Id, consts.MediaCheckProgressCcing)
	err = utils.CompletedCheckFile(ctx, absPath)
	if err != nil {
		service.UpdateCheckTask(ctx, task.Id, consts.MediaCheckProgressCcError, consts.MediaCheckFileStatusError, current)
		return
	}
	service.UpdateCheckTask(ctx, task.Id, consts.MediaCheckProgressCcSuccess, consts.MediaCheckFileStatusSuccess, current)
}
