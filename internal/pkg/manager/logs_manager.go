package manager

import (
	"sync"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shichen437/gowlive/internal/app/system/dao"
	"github.com/shichen437/gowlive/internal/app/system/model/do"
	"github.com/shichen437/gowlive/internal/app/system/model/entity"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

type LogManager struct {
	queue    chan entity.SysLogs
	wg       sync.WaitGroup
	mu       sync.Mutex
	stopChan chan struct{}
}

var (
	logManagerInstance *LogManager
	once               sync.Once
)

// 单例获取方法
func GetLogManager() *LogManager {
	once.Do(func() {
		logManagerInstance = &LogManager{
			queue:    make(chan entity.SysLogs, 30),
			stopChan: make(chan struct{}),
		}
		logManagerInstance.wg.Add(1)
		go logManagerInstance.processLogs()
	})
	return logManagerInstance
}

func (lm *LogManager) AddSuccessLog(logType int, content string) {
	lm.addLog(logType, 1, content)
}

func (lm *LogManager) AddErrorLog(logType int, content string) {
	lm.addLog(logType, 0, content)
}

func (lm *LogManager) addLog(logType, status int, content string) {
	entry := entity.SysLogs{
		Type:      logType,
		Status:    status,
		Content:   content,
		CreatedAt: utils.Now(),
	}
	lm.queue <- entry
}

func (lm *LogManager) processLogs() {
	defer lm.wg.Done()
	for {
		select {
		case entry := <-lm.queue:
			lm.mu.Lock()
			dao.SysLogs.Ctx(gctx.GetInitCtx()).Insert(do.SysLogs{
				Type:      entry.Type,
				Status:    entry.Status,
				Content:   entry.Content,
				CreatedAt: entry.CreatedAt,
			})
			lm.mu.Unlock()
		case <-lm.stopChan:
			return
		}
	}
}

func (lm *LogManager) Stop() {
	close(lm.stopChan)
	lm.wg.Wait()
	close(lm.queue)
}
