package system

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/manager"
	mp "github.com/shichen437/gowlive/internal/pkg/message_push"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

func StorageWarning(ctx context.Context) {
	usage := utils.GetDiskUsage()
	if usage > consts.StorageThreshold {
		title := utils.T(ctx, "ext.storage.waning.title")
		content := utils.Tf(ctx, "ext.storage.waning.desc", gconv.String(usage))
		manager.GetNotifyManager().AddWarningNotify(title, content)
		err := mp.PushMessage(gctx.New(), &mp.MessageModel{
			Title:   title,
			Content: content,
		})
		if err != nil {
			g.Log().Error(ctx, "空间预警发送通知失败：", err)
		}
	}
}
