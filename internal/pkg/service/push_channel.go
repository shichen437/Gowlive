package service

import (
	"context"

	"github.com/shichen437/gowlive/internal/app/system/dao"
	"github.com/shichen437/gowlive/internal/app/system/model"
	"github.com/shichen437/gowlive/internal/app/system/model/entity"
)

func GetAllPushChannel(ctx context.Context) []*model.PushChannel {
	var list []*model.PushChannel
	m := dao.PushChannel.Ctx(ctx)
	m = m.Where(dao.PushChannel.Columns().Status, 1)
	err := m.Scan(&list)
	if len(list) == 0 || err != nil {
		return list
	}
	arr := []int{}
	cMap := make(map[int]*model.PushChannel)
	for _, v := range list {
		arr = append(arr, v.Id)
		cMap[v.Id] = v
	}
	emails := []*entity.PushChannelEmail{}
	err = dao.PushChannelEmail.Ctx(ctx).
		WhereIn(dao.PushChannelEmail.Columns().ChannelId, arr).Scan(&emails)
	for _, v := range emails {
		cMap[v.ChannelId].Email = v
	}
	return list
}
