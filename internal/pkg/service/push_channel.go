package service

import (
	"context"
	"slices"

	"github.com/shichen437/gowlive/internal/app/system/dao"
	"github.com/shichen437/gowlive/internal/app/system/model"
	"github.com/shichen437/gowlive/internal/app/system/model/entity"
)

var (
	webhookTypeArr = []string{"lark", "dingTalk", "weCom"}
)

func GetAllPushChannel(ctx context.Context) []*model.PushChannel {
	var list []*model.PushChannel
	m := dao.PushChannel.Ctx(ctx)
	m = m.Where(dao.PushChannel.Columns().Status, 1)
	err := m.Scan(&list)
	if len(list) == 0 || err != nil {
		return list
	}
	earr, warr := []int{}, []int{}
	cMap := make(map[int]*model.PushChannel)
	for _, v := range list {
		if v.Type == "email" {
			earr = append(earr, v.Id)
		}
		if slices.Contains(webhookTypeArr, v.Type) {
			warr = append(warr, v.Id)
		}
		cMap[v.Id] = v
	}
	if len(earr) > 0 {
		emails := []*entity.PushChannelEmail{}
		err = dao.PushChannelEmail.Ctx(ctx).
			WhereIn(dao.PushChannelEmail.Columns().ChannelId, earr).Scan(&emails)
		for _, v := range emails {
			cMap[v.ChannelId].Email = v
		}
	}
	if len(warr) > 0 {
		webhooks := []*entity.PushChannelWebhook{}
		err = dao.PushChannelWebhook.Ctx(ctx).
			WhereIn(dao.PushChannelWebhook.Columns().ChannelId, warr).Scan(&webhooks)
		for _, v := range webhooks {
			cMap[v.ChannelId].Webhook = v
		}
	}
	return list
}
