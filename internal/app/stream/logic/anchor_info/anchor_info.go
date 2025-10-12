package logic

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	v1 "github.com/shichen437/gowlive/api/v1/stream"
	"github.com/shichen437/gowlive/internal/app/stream/dao"
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
		return
	}
	anchorInfo, err := aApi.ParseAnchorInfo(ctx)
	if err != nil {
		return
	}
	if anchorInfo == nil || anchorInfo.AnchorName == "" {
		err = gerror.New("解析主播数据失败")
		return
	}
	count, _ := dao.AnchorInfo.Ctx(ctx).
		Where(dao.AnchorInfo.Columns().Platform, anchorInfo.Platform).
		Where(dao.AnchorInfo.Columns().UniqueId, anchorInfo.UniqueId).Count()
	if count > 0 {
		err = gerror.New("重复添加")
		return
	}
	err = saveAnchorInfo(ctx, req, anchorInfo)
	if err != nil {
		err = gerror.New("添加主播数据失败")
		return
	}
	return
}

func (c *sAnchorInfo) Delete(ctx context.Context, req *v1.DeleteAnchorReq) (res *v1.DeleteAnchorRes, err error) {
	_, err = dao.AnchorInfo.Ctx(ctx).Where(dao.AnchorInfo.Columns().Id, req.Id).Delete()
	if err != nil {
		err = gerror.New("删除主播数据失败")
		return
	}
	_, err = dao.AnchorInfoHistory.Ctx(ctx).Where(dao.AnchorInfoHistory.Columns().AnchorId, req.Id).Delete()
	if err != nil {
		err = gerror.New("删除主播历史数据失败")
		return
	}
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
			g.Log().Error(ctx, txErr)
			return txErr
		}
		anchorId, txErr := result.LastInsertId()
		if txErr != nil {
			g.Log().Error(ctx, txErr)
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
		return nil
	})
}
