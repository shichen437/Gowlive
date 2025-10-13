package logic

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	v1 "github.com/shichen437/gowlive/api/v1/stream"
	"github.com/shichen437/gowlive/internal/app/stream/dao"
	"github.com/shichen437/gowlive/internal/app/stream/model"
	"github.com/shichen437/gowlive/internal/app/stream/model/do"
	"github.com/shichen437/gowlive/internal/app/stream/model/entity"
	"github.com/shichen437/gowlive/internal/app/stream/service"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/crons"
	"github.com/shichen437/gowlive/internal/pkg/lives"
	"github.com/shichen437/gowlive/internal/pkg/registry"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

type sLiveManage struct {
	lock sync.RWMutex
}

func init() {
	service.RegisterLiveManage(New())
}

func New() service.ILiveManage {
	return &sLiveManage{}
}

func (s *sLiveManage) List(ctx context.Context, req *v1.GetRoomListReq) (res *v1.GetRoomListRes, err error) {
	res = &v1.GetRoomListRes{}
	var list []*model.RoomInfo
	m := dao.LiveRoomInfo.Ctx(ctx)
	if req.Anchor != "" {
		m = m.WhereLike(dao.LiveRoomInfo.Columns().Anchor, "%"+req.Anchor+"%")
	}
	if req.RoomName != "" {
		m = m.WhereLike(dao.LiveRoomInfo.Columns().RoomName, "%"+req.RoomName+"%")
	}
	if req.Platform != "" {
		m = m.Where(dao.LiveRoomInfo.Columns().Platform, req.Platform)
	}
	m = dealSortParams(m, req.Sort)
	res.Total, err = m.Count()
	if err != nil || res.Total <= 0 {
		return
	}
	if err = m.Page(req.PageNum, req.PageSize).Scan(&list); err != nil {
		return
	}
	for _, item := range list {
		if item.Status != 0 {
			item.IsRecording = registry.Get().IsRecording(item.LiveId)
		}
	}
	res.Rows = list
	return res, nil
}

func (s *sLiveManage) Get(ctx context.Context, req *v1.GetLiveManageReq) (res *v1.GetLiveManageRes, err error) {
	res = &v1.GetLiveManageRes{}
	var entity *entity.LiveManage
	err = dao.LiveManage.Ctx(ctx).WherePri(req.LiveId).Scan(&entity)
	if err != nil {
		return
	}
	if entity == nil || entity.Id == 0 {
		err = gerror.New("直播间不存在")
		return
	}
	res.Data = entity
	return res, nil
}

func (s *sLiveManage) Add(ctx context.Context, req *v1.PostLiveManageReq) (res *v1.PostLiveManageRes, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	liveApi, err := lives.New(req.RoomUrl)
	if err != nil {
		return
	}
	info, err := liveApi.GetInfo()
	if err != nil || info == nil {
		err = gerror.New("获取直播间信息失败")
		return
	}
	var liveId int64
	err = saveLiveConfig(ctx, req, &liveId, info)
	if err != nil {
		return nil, gerror.Wrap(err, "添加直播间失败")
	}
	go listenerForAdd(int(liveId), req.MonitorType, req.MonitorStartAt, req.MonitorStopAt)
	return
}

func (s *sLiveManage) Update(ctx context.Context, req *v1.PutLiveManageReq) (res *v1.PutLiveManageRes, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	var tempData *entity.LiveManage
	err = dao.LiveManage.Ctx(ctx).WherePri(req.Id).Scan(&tempData)
	if err != nil {
		return nil, gerror.New("获取直播间失败")
	}
	if tempData == nil || tempData.Id == 0 {
		return nil, gerror.New("直播间不存在")
	}
	validNeedUpdate := false
	if tempData.MonitorType != req.MonitorType || tempData.Interval != req.Interval || tempData.Format != req.Format ||
		tempData.MonitorStartAt != req.MonitorStartAt || tempData.MonitorStopAt != req.MonitorStopAt || tempData.Remark != req.Remark {
		validNeedUpdate = true
	}
	if !validNeedUpdate {
		return
	}
	_, err = dao.LiveManage.Ctx(ctx).WherePri(req.Id).Update(do.LiveManage{
		Format:         req.Format,
		Interval:       req.Interval,
		MonitorType:    req.MonitorType,
		MonitorStartAt: req.MonitorStartAt,
		MonitorStopAt:  req.MonitorStopAt,
		Remark:         req.Remark,
		UpdatedAt:      utils.Now(),
	})
	if tempData.MonitorType != req.MonitorType {
		dao.LiveRoomInfo.Ctx(ctx).Where(dao.LiveRoomInfo.Columns().LiveId, req.Id).Update(do.LiveRoomInfo{
			Status:    req.MonitorType,
			UpdatedAt: utils.Now(),
		})
	}
	if err != nil {
		return nil, gerror.New("更新直播间失败")
	}
	go listenerForUpdate(req, tempData)
	return
}

func (s *sLiveManage) Delete(ctx context.Context, req *v1.DeleteLiveManageReq) (res *v1.DeleteLiveManageRes, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err = dao.LiveManage.Ctx(ctx).Where(dao.LiveManage.Columns().Id, req.LiveId).Delete()
		if err != nil {
			return err
		}
		_, err = dao.LiveRoomInfo.Ctx(ctx).Where(dao.LiveRoomInfo.Columns().LiveId, req.LiveId).Delete()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return
	}
	go listenerForDelete(req.LiveId)
	return
}

func (s *sLiveManage) Start(ctx context.Context, req *v1.PutLiveManageStartReq) (res *v1.PutLiveManageStartRes, err error) {
	var tempData *entity.LiveManage
	err = dao.LiveManage.Ctx(ctx).WherePri(req.Id).Scan(&tempData)
	if err != nil {
		return nil, gerror.New("获取直播间失败")
	}
	if tempData == nil || tempData.Id == 0 {
		return nil, gerror.New("直播间不存在")
	}
	if tempData.MonitorType != consts.MonitorTypeStop {
		return nil, gerror.New("直播间已在监控中")
	}
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err = dao.LiveManage.Ctx(ctx).WherePri(req.Id).Update(do.LiveManage{
			MonitorType: consts.MonitorTypeStart,
			UpdatedAt:   utils.Now(),
		})
		if err != nil {
			return err
		}
		_, err = dao.LiveRoomInfo.Ctx(ctx).
			Where(dao.LiveRoomInfo.Columns().LiveId, req.Id).
			Update(do.LiveRoomInfo{
				Status:    consts.MonitorTypeStart,
				UpdatedAt: utils.Now(),
			})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, gerror.New("开始监控直播间失败")
	}

	go listenerForQuickAdd(req.Id)
	return
}

func (s *sLiveManage) Stop(ctx context.Context, req *v1.PutLiveManageStopReq) (res *v1.PutLiveManageStopRes, err error) {
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err = dao.LiveManage.Ctx(ctx).WherePri(req.Id).Update(do.LiveManage{
			MonitorType: consts.MonitorTypeStop,
			UpdatedAt:   utils.Now(),
		})
		if err != nil {
			return err
		}
		_, err = dao.LiveRoomInfo.Ctx(ctx).
			Where(dao.LiveRoomInfo.Columns().LiveId, req.Id).
			Update(do.LiveRoomInfo{
				Status:    consts.MonitorTypeStop,
				UpdatedAt: utils.Now(),
			})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, gerror.New("停止监控直播间失败")
	}
	go listenerForDelete(req.Id)
	return
}

func dealSortParams(m *gdb.Model, sort string) *gdb.Model {
	switch sort {
	case "id:asc":
		m = m.OrderAsc(dao.LiveRoomInfo.Columns().LiveId)
	case "status:asc":
		m = m.OrderAsc(dao.LiveRoomInfo.Columns().Status)
	case "status:desc":
		m = m.OrderDesc(dao.LiveRoomInfo.Columns().Status)
	default:
		m = m.OrderDesc(dao.LiveRoomInfo.Columns().LiveId)
	}
	return m
}

func saveLiveConfig(ctx context.Context, req *v1.PostLiveManageReq, liveId *int64, info *lives.LiveState) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		var txErr error
		result, txErr := dao.LiveManage.Ctx(ctx).Insert(do.LiveManage{
			RoomUrl:        req.RoomUrl,
			Format:         req.Format,
			Interval:       req.Interval,
			MonitorType:    req.MonitorType,
			MonitorStartAt: req.MonitorStartAt,
			MonitorStopAt:  req.MonitorStopAt,
			Remark:         req.Remark,
			CreatedAt:      utils.Now(),
		})
		if txErr != nil {
			g.Log().Error(ctx, txErr)
			return txErr
		}
		*liveId, txErr = result.LastInsertId()
		if txErr != nil {
			g.Log().Error(ctx, txErr)
			return txErr
		}
		_, txErr = dao.LiveRoomInfo.Ctx(ctx).Insert(do.LiveRoomInfo{
			LiveId:    liveId,
			Platform:  info.Platform,
			RoomName:  info.RoomName,
			Anchor:    info.Anchor,
			Status:    req.MonitorType,
			CreatedAt: utils.Now(),
		})
		if txErr != nil {
			g.Log().Error(ctx, txErr)
			return txErr
		}
		return nil
	})
}

func listenerForQuickAdd(liveId int) {
	listenerForAdd(liveId, consts.MonitorTypeStart, "", "")
}

func listenerForAdd(liveId, monitorType int, monitorStartAt, monitorStopAt string) {
	ctx := gctx.GetInitCtx()
	if monitorType == 1 {
		registry.Get().Add(ctx, liveId)
	}
	if monitorType == 2 {
		crons.AddStreamCron(ctx, liveId, monitorStartAt, monitorStopAt)
	}
}

func listenerForUpdate(req *v1.PutLiveManageReq, tempData *entity.LiveManage) {
	ctx := gctx.GetInitCtx()
	if req.MonitorType == 0 {
		listenerForDelete(req.Id)
	}
	if req.MonitorType == 1 {
		if req.MonitorType != tempData.MonitorType || req.Format != tempData.Format || req.Interval != tempData.Interval {
			listenerForDelete(req.Id)
			registry.Get().Add(ctx, req.Id)
			return
		}
	}
	if req.MonitorType == 2 {
		if req.MonitorType != tempData.MonitorType || req.Format != tempData.Format ||
			req.Interval != tempData.Interval || req.MonitorStartAt != tempData.MonitorStartAt ||
			req.MonitorStopAt != tempData.MonitorStopAt {
			listenerForDelete(req.Id)
			crons.AddStreamCron(ctx, req.Id, req.MonitorStartAt, req.MonitorStopAt)
		}
	}
}

func listenerForDelete(liveId int) {
	ctx := gctx.GetInitCtx()
	registry.Get().Remove(ctx, liveId)
	crons.RemoveStreamCron(liveId)
	g.Log().Info(ctx, "listenerForDelete", liveId)
}
