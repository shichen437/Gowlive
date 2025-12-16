package crons

import (
	"context"
	"math/rand"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shichen437/gowlive/internal/app/stream/model/do"
	"github.com/shichen437/gowlive/internal/pkg/anchor"
	"github.com/shichen437/gowlive/internal/pkg/service"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

func AnchorInfoCron(ctx context.Context) {
	g.Log().Info(ctx, "anchor info cron start...")
	// 获取所有主播信息
	anchorList := service.GetAllAnchorInfo(ctx)
	if len(anchorList) == 0 {
		g.Log().Info(ctx, "anchor info cron interrupted, no anchor info.")
		return
	}
	currentDay := gtime.Now().Format("Y-m-d")
	// 遍历并更新信息
	lp := ""
	for _, v := range anchorList {
		if service.ExistsTodayHistory(ctx, v.Id, currentDay) || v.Url == "" {
			continue
		}
		if lp != "" && v.Platform == lp {
			time.Sleep(time.Duration(rand.Intn(2000)+1000) * time.Millisecond)
		}
		lp = v.Platform
		api, err := anchor.New(v.Url)
		if err != nil {
			g.Log().Errorf(ctx, "AIC build anchor api error: %v", err)
			continue
		}
		info, err := api.ParseAnchorInfo(ctx)
		if err != nil || info.AnchorName == "" {
			g.Log().Errorf(ctx, "AIC parse anchor info error: %v", err)
			continue
		}
		err = service.UpdateAnchorInfo(ctx, do.AnchorInfo{
			AnchorName:     info.AnchorName,
			UniqueId:       info.UniqueId,
			Signature:      info.Signature,
			FollowerCount:  info.FollowerCount,
			FollowingCount: info.FollowingCount,
			LikeCount:      info.LikeCount,
			VideoCount:     info.VideoCount,
			UpdatedAt:      utils.Now(),
		}, v.Id)
		if err != nil {
			g.Log().Errorf(ctx, "AIC update anchor error: %v", err)
			continue
		}
		service.SaveTodayHistory(ctx, do.AnchorInfoHistory{
			AnchorId:       v.Id,
			AnchorName:     info.AnchorName,
			Signature:      info.Signature,
			CollectedDate:  currentDay,
			FollowerCount:  info.FollowerCount,
			FollowingCount: info.FollowingCount,
			LikeCount:      info.LikeCount,
			VideoCount:     info.VideoCount,
			CreatedAt:      utils.Now(),
		})
	}
	g.Log().Info(ctx, "anchor info cron finished.")
}
