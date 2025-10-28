package manager

import (
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shichen437/gowlive/internal/app/system/dao"
	"github.com/shichen437/gowlive/internal/app/system/model/do"
	"github.com/shichen437/gowlive/internal/app/system/model/entity"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/sse"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

type NotifyManager struct {
	queue    chan entity.SysNotify
	wg       sync.WaitGroup
	mu       sync.Mutex
	stopChan chan struct{}
}

var (
	notifyManagerInstance *NotifyManager
	notifyOnce            sync.Once
)

// 单例获取方法
func GetNotifyManager() *NotifyManager {
	notifyOnce.Do(func() {
		notifyManagerInstance = &NotifyManager{
			queue:    make(chan entity.SysNotify, 30),
			stopChan: make(chan struct{}),
		}
		notifyManagerInstance.wg.Add(1)
		go notifyManagerInstance.processNotify()
	})
	return notifyManagerInstance
}

func (nm *NotifyManager) AddInfoNotify(title, content string) {
	nm.addNotify("info", title, content)
}

func (nm *NotifyManager) AddWarningNotify(title, content string) {
	nm.addNotify("warning", title, content)
}

func (nm *NotifyManager) addNotify(level, title, content string) {
	entry := entity.SysNotify{
		Level:     level,
		Title:     title,
		Content:   content,
		CreatedAt: utils.Now(),
	}
	nm.queue <- entry
	msg := sse.GetSseMsgStr(consts.SSE_EVENT_TYPE_GLOBAL, g.MapStrAny{
		"notify": entry,
	})
	sse.BroadcastMessage(consts.SSE_CHANNEL_GLOBAL, msg)
}

func (nm *NotifyManager) processNotify() {
	defer nm.wg.Done()
	for {
		select {
		case entry := <-nm.queue:
			nm.mu.Lock()
			dao.SysNotify.Ctx(gctx.GetInitCtx()).Insert(do.SysNotify{
				Level:     entry.Level,
				Title:     entry.Title,
				Content:   entry.Content,
				Status:    0,
				CreatedAt: entry.CreatedAt,
			})
			nm.mu.Unlock()
		case <-nm.stopChan:
			return
		}
	}
}

func (nm *NotifyManager) Stop() {
	close(nm.stopChan)
	nm.wg.Wait()
	close(nm.queue)
}
