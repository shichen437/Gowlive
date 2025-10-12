package service

import (
	"context"

	"github.com/shichen437/gowlive/internal/app/stream/dao"
	"github.com/shichen437/gowlive/internal/app/stream/model/do"
	"github.com/shichen437/gowlive/internal/app/stream/model/entity"
)

func GetAllAnchorInfo(ctx context.Context) []*entity.AnchorInfo {
	var list []*entity.AnchorInfo
	m := dao.AnchorInfo.Ctx(ctx)
	err := m.Scan(&list)
	if len(list) == 0 || err != nil {
		return list
	}
	return list
}

func UpdateAnchorInfo(ctx context.Context, info do.AnchorInfo, aid int) error {
	_, err := dao.AnchorInfo.Ctx(ctx).WherePri(aid).Update(info)
	return err
}

func ExistsTodayHistory(ctx context.Context, anchorId int, date string) bool {
	count, _ := dao.AnchorInfoHistory.Ctx(ctx).
		Where(dao.AnchorInfoHistory.Columns().AnchorId, anchorId).
		Where(dao.AnchorInfoHistory.Columns().CollectedDate, date).
		Count()
	if count > 0 {
		return true
	}
	return false
}

func SaveTodayHistory(ctx context.Context, history do.AnchorInfoHistory) error {
	_, err := dao.AnchorInfoHistory.Ctx(ctx).Insert(history)
	return err
}
