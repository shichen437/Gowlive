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

var (
	title = "空间预警"
)

func StorageWarning(ctx context.Context) {
	if utils.GetDiskUsage() > consts.StorageThreshold {
		content := "存储空间已达到" + gconv.String(consts.StorageThreshold) + "%"
		manager.GetNotfiyManager().AddWarningNotify(title, content)
		err := mp.PushMessage(gctx.New(), &mp.MessageModel{
			Title:   title,
			Content: content,
		})
		if err != nil {
			g.Log().Error(ctx, "空间预警发送通知失败：", err)
		}
	}
}
