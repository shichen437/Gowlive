package recorders

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shichen437/gowlive/internal/app/stream/dao"
	"github.com/shichen437/gowlive/internal/app/stream/model/do"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/lives"
	mr "github.com/shichen437/gowlive/internal/pkg/manager"
	"github.com/shichen437/gowlive/internal/pkg/message_push"
	"github.com/shichen437/gowlive/internal/pkg/service"
	"github.com/shichen437/gowlive/internal/pkg/third/openlist"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

func liveStartBiz(ctx context.Context, session *lives.LiveSession) {
	g.Log().Info(ctx, "liveStartBiz", session.Id)
	session.StartedAt = utils.Now()
	go message_push.LivePush(gctx.GetInitCtx(), session.State.Anchor, true)
}

func liveEndBiz(ctx context.Context, session *lives.LiveSession) {
	g.Log().Info(ctx, "liveEndBiz", session.Id)
	startTime := session.StartedAt
	addHistory(ctx, session.Id, session.State.Anchor, startTime, utils.Now())
	enable := mr.GetSettingsManager().GetSetting(consts.SKLiveEndNotify)
	if enable == 1 {
		go message_push.LivePush(gctx.GetInitCtx(), session.State.Anchor, false)
	}
	autoClean := mr.GetSettingsManager().GetSetting(consts.SKAutoCleanLittleFile)
	if autoClean > 0 && session.Config.MonitorOnly == 0 {
		cleanLittleFiles(session.Filename, autoClean)
	}
	syncEnable := mr.GetSettingsManager().GetSetting(consts.SKDataSyncEnable)
	if syncEnable == 1 && openlist.CheckLoginParams() {
		go syncTask(session.Filename, session.Config.SyncPath)
	}
}

func addHistory(ctx context.Context, liveId int, anchor string, startTime, endTime *gtime.Time) {
	if startTime == nil || endTime == nil {
		g.Log().Warningf(ctx, "Invalid start or end time for liveId %d.", liveId)
		return
	}
	_, err := dao.LiveHistory.Ctx(ctx).Insert(do.LiveHistory{
		LiveId:    liveId,
		Anchor:    anchor,
		StartedAt: startTime,
		EndedAt:   endTime,
		Duration:  fmt.Sprintf("%.2f", endTime.Sub(startTime).Hours()),
		CreatedAt: utils.Now(),
	})
	if err != nil {
		g.Log().Errorf(ctx, "Failed to save live history for liveId %d: %v", liveId, err)
	}
}

func (*manager) updateName(ctx context.Context, session *lives.LiveSession) {
	if session == nil || session.Id == 0 {
		return
	}
	service.UpdateRoomInfo(ctx, session)
}

func cleanLittleFiles(filename string, fileSize int) {
	if strings.TrimSpace(filename) == "" || fileSize <= 0 {
		return
	}

	ctx := gctx.GetInitCtx()
	thresholdBytes := int64(fileSize) * 1024 * 1024

	dir, name := filepath.Split(filename)
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)

	origPath := filepath.Join(dir, name)

	deleteIfSmall := func(path string, threshold int64) bool {
		sz, ok := getFileSize(path)
		if !ok {
			return false
		}
		if sz < threshold {
			if err := os.Remove(path); err != nil {
				g.Log().Infof(ctx, "remove failed: %s, err: %v\n", path, err)
				return false
			}
			g.Log().Infof(ctx, "removed small file: %s (size=%d bytes, threshold=%d)\n", path, sz, threshold)
			return true
		}
		return false
	}

	if _, exists := getFileSize(origPath); exists {
		_ = deleteIfSmall(origPath, thresholdBytes)
		return
	}

	const maxParts = 1000
	deletedAny := false

	for i := range maxParts {
		partName := fmt.Sprintf("%s_%03d%s", base, i, ext)
		partPath := filepath.Join(dir, partName)

		if _, ok := getFileSize(partPath); !ok {
			break
		}

		deleteIfSmall(partPath, thresholdBytes)
	}

	if !deletedAny {
		g.Log().Info(ctx, "no little files found!")
		return
	}
}

func syncTask(filename, syncPath string) {
	if strings.TrimSpace(filename) == "" || syncPath == "" {
		return
	}

	ctx := gctx.GetInitCtx()

	dir, name := filepath.Split(filename)
	syncManager := mr.GetFileSyncManager()

	if _, exists := getFileSize(filename); exists {
		id := service.AddSyncTask(ctx, dir, name, syncPath)
		if id > 0 {
			syncManager.Add(id)
			g.Log().Infof(ctx, "Added sync task for file: %s", filename)
		}
		return
	}

	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	foundAny := false
	const maxParts = 1000
	for i := range maxParts {
		partName := fmt.Sprintf("%s_%03d%s", base, i, ext)
		partPath := filepath.Join(dir, partName)

		if _, ok := getFileSize(partPath); !ok {
			break
		}

		foundAny = true
		id := service.AddSyncTask(ctx, dir, partName, syncPath)
		if id > 0 {
			syncManager.Add(id)
			g.Log().Infof(ctx, "Added sync task for file part: %s", partPath)
		}
	}

	if !foundAny {
		g.Log().Infof(ctx, "No files found to sync for: %s", filename)
	}
}

func getFileSize(path string) (int64, bool) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, false
		}
		return 0, false
	}
	return info.Size(), info.Mode().IsRegular()
}
