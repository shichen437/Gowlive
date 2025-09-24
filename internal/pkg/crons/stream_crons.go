package crons

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shichen437/gowlive/internal/pkg/registry"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

var (
	CronStartPrefix = "Cron-Start-"
	CronStopPrefix  = "Cron-Stop-"
)

func AddStreamCron(ctx context.Context, liveId int, startAt, stopAt string) {
	if startAt == "" || stopAt == "" || startAt == stopAt {
		g.Log().Error(ctx, "添加定时任务失败, 时间不能为空或相同")
		return
	}
	cronStart, err := parseCronTime(startAt)
	if err != nil {
		g.Log().Error(ctx, "解析开始时间失败: ", err)
		return
	}
	cronStop, err := parseCronTime(stopAt)
	if err != nil {
		g.Log().Error(ctx, "解析停止时间失败: ", err)
		return
	}
	gcron.Add(ctx, strings.Trim(cronStart, " "), func(ctx context.Context) {
		g.Log().Info(ctx, "启动定时任务-", liveId)
		registry.Get().Add(ctx, liveId)
	}, CronStartPrefix+strconv.Itoa(liveId))
	gcron.Add(ctx, strings.Trim(cronStop, " "), func(ctx context.Context) {
		g.Log().Info(ctx, "停止定时任务-", liveId)
		registry.Get().Remove(ctx, liveId)
	}, CronStopPrefix+strconv.Itoa(liveId))
	g.Log().Info(ctx, "添加定时任务-", liveId)
	if utils.IsTimeRange(startAt, stopAt) {
		registry.Get().Add(ctx, liveId)
	}
}

func RemoveStreamCron(liveId int) {
	search := gcron.Search(CronStartPrefix + strconv.Itoa(liveId))
	if search != nil {
		gcron.Remove(search.Name)
	}
	search = gcron.Search(CronStopPrefix + strconv.Itoa(liveId))
	if search != nil {
		gcron.Remove(search.Name)
	}
	g.Log().Info(gctx.GetInitCtx(), "移除定时任务-", liveId)
}

func RestartStreamCron(ctx context.Context, liveId int, startAt, stopAt string) {
	RemoveStreamCron(liveId)
	AddStreamCron(ctx, liveId, startAt, stopAt)
}

func parseCronTime(t string) (string, error) {
	if len(t) == 0 {
		return "", gerror.New("传入时间不能为空")
	}
	s := strings.Split(t, ":")
	if len(s) != 2 {
		return "", gerror.New("传入时间不能为空")
	}
	return fmt.Sprintf("0 %s %s * * *", s[1], s[0]), nil
}
