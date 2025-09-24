package logic

import (
	"context"

	v1 "github.com/shichen437/gowlive/api/v1/system"
	"github.com/shichen437/gowlive/internal/app/stream/dao"
	"github.com/shichen437/gowlive/internal/app/system/model"
	"github.com/shichen437/gowlive/internal/app/system/service"
	"github.com/shichen437/gowlive/internal/pkg/registry"
)

type sSystemOverview struct {
}

func init() {
	service.RegisterSystemOverview(New())
}

func New() service.ISystemOverview {
	return &sSystemOverview{}
}

func (s *sSystemOverview) Overview(ctx context.Context, req *v1.GetOverviewReq) (res *v1.GetOverviewRes, err error) {
	res = &v1.GetOverviewRes{}
	data := &model.Overview{}
	genRoomData(ctx, data)
	genRecordingTime(ctx, data)
	res.Data = data
	return
}

func genRoomData(ctx context.Context, data *model.Overview) {
	m := dao.LiveRoomInfo.Ctx(ctx)
	data.LiveRoomCount, _ = m.Count()
	data.RecordingRoomCount = registry.Get().RecordingCount()
}

func genRecordingTime(ctx context.Context, data *model.Overview) {
	timeCount, err := dao.LiveHistory.Ctx(ctx).Sum(dao.LiveHistory.Columns().Duration)
	if err != nil {
		return
	}
	data.RecordTimeCount = float64(timeCount)
}
