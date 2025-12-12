package logic

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	v1 "github.com/shichen437/gowlive/api/v1/stream"
	"github.com/shichen437/gowlive/internal/app/stream/dao"
	"github.com/shichen437/gowlive/internal/app/stream/model"
	"github.com/shichen437/gowlive/internal/app/stream/model/do"
	"github.com/shichen437/gowlive/internal/app/stream/model/entity"
	"github.com/shichen437/gowlive/internal/app/stream/service"
	"github.com/shichen437/gowlive/internal/pkg/anchor"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

type (
	sAnchorInfo struct{}
)

func init() {
	service.RegisterAnchorInfo(New())
}

func New() service.IAnchorInfo {
	return &sAnchorInfo{}
}

func (c *sAnchorInfo) List(ctx context.Context, req *v1.GetAnchorListReq) (res *v1.GetAnchorListRes, err error) {
	res = &v1.GetAnchorListRes{}
	m := dao.AnchorInfo.Ctx(ctx)
	if req.Platform != "" {
		m = m.Where(dao.AnchorInfo.Columns().Platform, req.Platform)
	}
	if req.Nickname != "" {
		m = m.WhereLike(dao.AnchorInfo.Columns().AnchorName, "%"+req.Nickname+"%")
	}
	res.Total, err = m.Count()
	if err != nil {
		return
	}
	if res.Total <= 0 {
		return
	}
	var list []*entity.AnchorInfo
	m = m.OrderDesc(dao.AnchorInfo.Columns().Id)
	err = m.Page(req.PageNum, req.PageSize).Scan(&list)
	if err != nil {
		return
	}
	res.Rows = list
	return
}

func (c *sAnchorInfo) Add(ctx context.Context, req *v1.PostAnchorReq) (res *v1.PostAnchorRes, err error) {
	aApi, err := anchor.New(req.Url)
	if err != nil {
		g.Log().Errorf(ctx, "获取解析 api 失败，错误信息：%v", err)
		return
	}
	anchorInfo, err := aApi.ParseAnchorInfo(ctx)
	if err != nil {
		g.Log().Errorf(ctx, "获取主播数据失败，错误信息：%v", err)
		return
	}
	if anchorInfo == nil || anchorInfo.AnchorName == "" {
		err = utils.TError(ctx, "stream.anchor.error.Parse")
		return
	}
	count, _ := dao.AnchorInfo.Ctx(ctx).
		Where(dao.AnchorInfo.Columns().Platform, anchorInfo.Platform).
		Where(dao.AnchorInfo.Columns().UniqueId, anchorInfo.UniqueId).Count()
	if count > 0 {
		err = utils.TError(ctx, "stream.anchor.error.Repeated")
		return
	}
	err = saveAnchorInfo(ctx, req, anchorInfo)
	if err != nil {
		g.Log().Errorf(ctx, "添加主播数据失败,错误信息：%v", err)
		err = utils.TError(ctx, "stream.anchor.error.Add")
		return
	}
	return
}

func (c *sAnchorInfo) Delete(ctx context.Context, req *v1.DeleteAnchorReq) (res *v1.DeleteAnchorRes, err error) {
	_, err = dao.AnchorInfo.Ctx(ctx).Where(dao.AnchorInfo.Columns().Id, req.Id).Delete()
	if err != nil {
		err = utils.TError(ctx, "stream.anchor.error.Delete")
		return
	}
	_, err = dao.AnchorInfoHistory.Ctx(ctx).Where(dao.AnchorInfoHistory.Columns().AnchorId, req.Id).Delete()
	if err != nil {
		err = utils.TError(ctx, "stream.anchor.error.DeleteHistory")
		return
	}
	return
}

func (c *sAnchorInfo) StatInfo(ctx context.Context, req *v1.GetAnchorStatInfoReq) (res *v1.GetAnchorStatInfoRes, err error) {
	res = &v1.GetAnchorStatInfoRes{}
	info := &model.AnchorStatInfo{
		WeekFollowersIncr:  0,
		WeekLikeNumIncr:    0,
		MonthFollowersIncr: 0,
	}
	var list []*entity.AnchorInfoHistory
	m := dao.AnchorInfoHistory.Ctx(ctx).Where(dao.AnchorInfoHistory.Columns().AnchorId, req.Id)
	count, err := m.Count()
	if err != nil || count <= 0 {
		return
	}

	thirtyOneDaysAgo := gtime.Now().AddDate(0, 0, -31).Format("Y-m-d")
	m = m.WhereGTE(dao.AnchorInfoHistory.Columns().CollectedDate, thirtyOneDaysAgo)
	err = m.OrderAsc(dao.AnchorInfoHistory.Columns().CollectedDate).Scan(&list)
	if err != nil || len(list) <= 0 {
		return
	}

	info.HistoryData = make([]*model.AnchorStatData, 0, len(list))
	historyMap := make(map[string]*entity.AnchorInfoHistory, len(list))
	for _, h := range list {
		info.HistoryData = append(info.HistoryData, &model.AnchorStatData{
			RecordDate: h.CollectedDate,
			Followers:  h.FollowerCount,
			LikeCount:  h.LikeCount,
		})
		historyMap[h.CollectedDate] = h
	}

	if len(list) < 2 {
		res.Data = info
		return
	}

	latestRecord := list[len(list)-1]
	latestDate := gtime.NewFromStr(latestRecord.CollectedDate)

	sevenDaysAgoDateStr := latestDate.AddDate(0, 0, -7).Format("Y-m-d")
	if sevenDaysAgoRecord, ok := historyMap[sevenDaysAgoDateStr]; ok {
		increase := latestRecord.FollowerCount - sevenDaysAgoRecord.FollowerCount
		info.WeekFollowersIncr = increase
		increase = latestRecord.LikeCount - sevenDaysAgoRecord.LikeCount
		info.WeekLikeNumIncr = increase
	}

	// 30-day average
	thirtyDaysAgoDateStr := latestDate.AddDate(0, 0, -30).Format("Y-m-d")
	if thirtyDaysAgoRecord, ok := historyMap[thirtyDaysAgoDateStr]; ok {
		increase := latestRecord.FollowerCount - thirtyDaysAgoRecord.FollowerCount
		info.MonthFollowersIncr = increase
	}

	res.Data = info
	return
}

func saveAnchorInfo(ctx context.Context, req *v1.PostAnchorReq, anchorInfo *anchor.AnchorInfo) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		result, txErr := dao.AnchorInfo.Ctx(ctx).Insert(do.AnchorInfo{
			Url:            req.Url,
			Platform:       anchorInfo.Platform,
			AnchorName:     anchorInfo.AnchorName,
			UniqueId:       anchorInfo.UniqueId,
			Signature:      anchorInfo.Signature,
			FollowerCount:  anchorInfo.FollowerCount,
			FollowingCount: anchorInfo.FollowingCount,
			LikeCount:      anchorInfo.LikeCount,
			VideoCount:     anchorInfo.VideoCount,
			CreatedAt:      utils.Now(),
		})
		if txErr != nil {
			return txErr
		}
		anchorId, txErr := result.LastInsertId()
		if txErr != nil {
			return txErr
		}
		_, txErr = dao.AnchorInfoHistory.Ctx(ctx).Insert(do.AnchorInfoHistory{
			AnchorId:       anchorId,
			CollectedDate:  gtime.Now().Format("Y-m-d"),
			AnchorName:     anchorInfo.AnchorName,
			Signature:      anchorInfo.Signature,
			FollowerCount:  anchorInfo.FollowerCount,
			FollowingCount: anchorInfo.FollowingCount,
			LikeCount:      anchorInfo.LikeCount,
			VideoCount:     anchorInfo.VideoCount,
			CreatedAt:      utils.Now(),
		})
		if txErr != nil {
			return txErr
		}
		return nil
	})
}
