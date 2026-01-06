package logic

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	v1 "github.com/shichen437/gowlive/api/v1/stream"
	"github.com/shichen437/gowlive/internal/app/stream/dao"
	"github.com/shichen437/gowlive/internal/app/stream/model"
	"github.com/shichen437/gowlive/internal/app/stream/model/do"
	"github.com/shichen437/gowlive/internal/app/stream/model/entity"
	"github.com/shichen437/gowlive/internal/app/stream/service"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/crons"
	"github.com/shichen437/gowlive/internal/pkg/lives"
	"github.com/shichen437/gowlive/internal/pkg/manager"
	"github.com/shichen437/gowlive/internal/pkg/registry"
	"github.com/shichen437/gowlive/internal/pkg/utils"
	"github.com/xuri/excelize/v2"
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
			item.IsLiving = registry.Get().IsLiving(item.LiveId)
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
		err = utils.TError(ctx, "stream.live.error.NotExists")
		return
	}
	res.Data = entity
	return res, nil
}

func (s *sLiveManage) Add(ctx context.Context, req *v1.PostLiveManageReq) (res *v1.PostLiveManageRes, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if checkExistsRoomUrl(ctx, req.RoomUrl) {
		err = utils.TError(ctx, "stream.anchor.error.Repeated")
		return
	}
	liveApi, err := lives.New(req.RoomUrl)
	if err != nil {
		g.Log().Errorf(ctx, "获取解析 api 失败，错误信息：%v", err)
		return
	}
	info, err := liveApi.GetInfo()
	if err != nil || info == nil {
		g.Log().Errorf(ctx, "获取直播数据失败，错误信息：%v", err)
		err = utils.TError(ctx, "stream.live.error.GetRoomInfo")
		return
	}
	var liveId int64
	err = saveLiveConfig(ctx, req, &liveId, info)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, utils.TError(ctx, "stream.live.error.Add")
	}
	go listenerForAdd(int(liveId), req.MonitorType, req.MonitorStartAt, req.MonitorStopAt)
	return
}

func (s *sLiveManage) BatchAdd(ctx context.Context, req *v1.PostLiveManageBatchReq) (res *v1.PostLiveManageBatchRes, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if len(req.RoomUrls) < 1 || len(req.RoomUrls) > 30 {
		err = utils.TError(ctx, "stream.live.error.BatchAddNumLimit")
		return
	}
	go func() {
		newCtx := gctx.GetInitCtx()
		errUrls := make([]string, 0)
		for _, roomUrl := range req.RoomUrls {
			if checkExistsRoomUrl(ctx, roomUrl) {
				g.Log().Errorf(newCtx, "链接 %s 重复添加", roomUrl)
				errUrls = append(errUrls, roomUrl)
				continue
			}
			liveApi, err := lives.New(roomUrl)
			if err != nil {
				g.Log().Errorf(newCtx, "链接 %s 获取解析 API 失败，错误信息：%v", roomUrl, err)
				errUrls = append(errUrls, roomUrl)
				continue
			}
			info, err := liveApi.GetInfo()
			if err != nil || info == nil {
				g.Log().Errorf(newCtx, "链接 %s 获取解析信息失败，错误信息：%v", roomUrl, err)
				errUrls = append(errUrls, roomUrl)
				continue
			}
			var liveId int64
			buildReq := &v1.PostLiveManageReq{
				RoomUrl:     roomUrl,
				MonitorType: consts.MonitorTypeStop,
				Interval:    req.Interval,
				Format:      req.Format,
				Remark:      req.Remark,
			}
			err = saveLiveConfig(newCtx, buildReq, &liveId, info)
			if err != nil {
				g.Log().Errorf(newCtx, "链接 %s 保存数据失败，错误信息：%v", roomUrl, err)
				errUrls = append(errUrls, roomUrl)
				continue
			}
		}
		if len(errUrls) > 0 {
			manager.GetNotifyManager().AddWarningNotify(utils.T(ctx, "stream.live.error.BatchAdd"), utils.T(ctx, "stream.live.error.BatchAddWithFailedLink")+strings.Join(errUrls, ","))
			return
		}
		manager.GetNotifyManager().AddInfoNotify(utils.T(ctx, "stream.live.success.BatchAdd"), utils.T(ctx, "stream.live.success.BatchAddWithLinkNum")+strconv.Itoa(len(req.RoomUrls)))
	}()
	return
}

func (s *sLiveManage) Update(ctx context.Context, req *v1.PutLiveManageReq) (res *v1.PutLiveManageRes, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	var tempData *entity.LiveManage
	err = dao.LiveManage.Ctx(ctx).WherePri(req.Id).Scan(&tempData)
	if err != nil {
		return nil, utils.TError(ctx, "stream.live.error.Get")
	}
	if tempData == nil || tempData.Id == 0 {
		return nil, utils.TError(ctx, "stream.live.error.NotExists")
	}
	validNeedUpdate := false
	if tempData.MonitorType != req.MonitorType || tempData.Interval != req.Interval || tempData.Format != req.Format ||
		tempData.MonitorStartAt != req.MonitorStartAt || tempData.MonitorStopAt != req.MonitorStopAt || tempData.Remark != req.Remark ||
		tempData.Quality != req.Quality || tempData.SegmentTime != req.SegmentTime || tempData.MonitorOnly != req.MonitorOnly || tempData.SyncPath != req.SyncPath {
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
		Quality:        req.Quality,
		SegmentTime:    req.SegmentTime,
		MonitorOnly:    req.MonitorOnly,
		Remark:         req.Remark,
		SyncPath:       req.SyncPath,
		UpdatedAt:      utils.Now(),
	})
	if tempData.MonitorType != req.MonitorType {
		dao.LiveRoomInfo.Ctx(ctx).Where(dao.LiveRoomInfo.Columns().LiveId, req.Id).Update(do.LiveRoomInfo{
			Status:    req.MonitorType,
			UpdatedAt: utils.Now(),
		})
	}
	if err != nil {
		g.Log().Errorf(ctx, "更新直播间失败，错误信息：%v", err)
		return nil, utils.TError(ctx, "stream.live.error.Update")
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
		g.Log().Errorf(ctx, "删除直播间失败，错误信息：%v", err)
		err = utils.TError(ctx, "stream.live.error.Delete")
		return
	}
	go listenerForDelete(req.LiveId)
	return
}

func (s *sLiveManage) Start(ctx context.Context, req *v1.PutLiveManageStartReq) (res *v1.PutLiveManageStartRes, err error) {
	var tempData *entity.LiveManage
	err = dao.LiveManage.Ctx(ctx).WherePri(req.Id).Scan(&tempData)
	if err != nil {
		return nil, utils.TError(ctx, "stream.live.error.Get")
	}
	if tempData == nil || tempData.Id == 0 {
		return nil, utils.TError(ctx, "stream.live.error.NotExists")
	}
	if tempData.MonitorType != consts.MonitorTypeStop {
		return nil, utils.TError(ctx, "stream.live.error.AlreadyInMonitor")
	}
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err = dao.LiveManage.Ctx(ctx).WherePri(req.Id).Update(do.LiveManage{
			MonitorType: consts.MonitorTypeIntelligent,
			UpdatedAt:   utils.Now(),
		})
		if err != nil {
			return err
		}
		_, err = dao.LiveRoomInfo.Ctx(ctx).
			Where(dao.LiveRoomInfo.Columns().LiveId, req.Id).
			Update(do.LiveRoomInfo{
				Status:    consts.MonitorTypeIntelligent,
				UpdatedAt: utils.Now(),
			})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, utils.TError(ctx, "stream.live.error.Start")
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
		return nil, utils.TError(ctx, "stream.live.error.Stop")
	}
	go listenerForDelete(req.Id)
	return
}

func (s *sLiveManage) Top(ctx context.Context, req *v1.PutLiveManageTopReq) (res *v1.PutLiveManageTopRes, err error) {
	count, err := dao.LiveRoomInfo.Ctx(ctx).Where(dao.LiveRoomInfo.Columns().IsTop, 1).Count()
	if err != nil {
		return nil, utils.TError(ctx, "stream.live.error.GetTopNum")
	}
	if count >= consts.MaxTopCount {
		return nil, utils.TError(ctx, "stream.live.error.TopNumLimit")
	}
	_, err = dao.LiveRoomInfo.Ctx(ctx).Where(dao.LiveRoomInfo.Columns().LiveId, req.Id).Update(do.LiveRoomInfo{
		IsTop:     1,
		ToppedAt:  utils.Now(),
		UpdatedAt: utils.Now(),
	})
	if err != nil {
		return nil, utils.TError(ctx, "stream.live.error.Top")
	}
	return
}

func (s *sLiveManage) UnTop(ctx context.Context, req *v1.PutLiveManageUnTopReq) (res *v1.PutLiveManageUnTopRes, err error) {
	_, err = dao.LiveRoomInfo.Ctx(ctx).Data(g.Map{
		"is_top":     0,
		"topped_at":  nil,
		"updated_at": utils.Now(),
	}).Where(dao.LiveRoomInfo.Columns().LiveId, req.Id).Update()
	if err != nil {
		return nil, utils.TError(ctx, "stream.live.error.UnTop")
	}
	return
}

func (s *sLiveManage) Export(ctx context.Context, req *v1.ExportRoomInfoReq) (res *v1.ExportRoomInfoRes, err error) {
	list, err := getExportData(ctx, req)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, utils.TError(ctx, "stream.live.error.NoDataToExport")
	}
	r := g.RequestFromCtx(ctx)
	r.Response.ClearBuffer()
	r.Response.Header().Set("Cache-Control", "no-store")
	ustr := gtime.Now().UnixMilli()
	switch req.ExportType {
	case 1:
		exportExcel(r, ustr, list)
	default:
		exportTxt(r, ustr, list)
	}
	return
}

func (s *sLiveManage) Preview(ctx context.Context, req *v1.PreviewRoomReq) (res *v1.PreviewRoomRes, err error) {
	res = &v1.PreviewRoomRes{}
	info := registry.Get().GetPreviewInfo(req.Id)
	if info == nil {
		return nil, utils.TError(ctx, "stream.live.error.Preview")
	}
	gconv.Struct(info, &res.PreviewInfo)
	return
}

func (s *sLiveManage) PreviewList(ctx context.Context, req *v1.PreviewRoomListReq) (res *v1.PreviewRoomListRes, err error) {
	res = &v1.PreviewRoomListRes{}
	list := registry.Get().GetPreviewList()
	if list == nil {
		return
	}
	gconv.Struct(list, &res.PreviewList)
	return
}

func (s *sLiveManage) QuickAdd(ctx context.Context, req *v1.PostQuickLinkReq) (res *v1.PostQuickLinkRes, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	liveApi, err := lives.New(req.Url)
	if err == nil {
		if checkExistsRoomUrl(ctx, req.Url) {
			err = utils.TError(ctx, "stream.anchor.error.Repeated")
			return
		}
		info, liveErr := liveApi.GetInfo()
		if liveErr != nil || info == nil {
			g.Log().Errorf(ctx, "获取直播数据失败，错误信息：%v", liveErr)
			err = utils.TError(ctx, "stream.live.error.GetRoomInfo")
			return
		}
		var liveId int64
		liveReq := &v1.PostLiveManageReq{
			RoomUrl:     req.Url,
			MonitorType: consts.MonitorTypeStop,
			Interval:    30,
			Format:      "flv",
		}
		err = saveLiveConfig(ctx, liveReq, &liveId, info)
		if err != nil {
			g.Log().Errorf(ctx, "保存直播配置失败，错误信息：%v", err)
			err = utils.TError(ctx, "stream.live.error.Add")
			return
		}
		return
	}
	_, err = service.AnchorInfo().Add(ctx, &v1.PostAnchorReq{
		Url: req.Url,
	})
	if err != nil {
		g.Log().Errorf(ctx, "添加主播信息失败，错误信息：%v", err)
		err = utils.TError(ctx, "stream.live.error.QuickAdd")
		return
	}
	return
}

func checkExistsRoomUrl(ctx context.Context, url string) bool {
	qUrl := utils.UrlRemoveParams(url)
	count, _ := dao.LiveManage.Ctx(ctx).WhereLike(dao.LiveManage.Columns().RoomUrl, qUrl+"%").Count()
	return count > 0
}

func dealSortParams(m *gdb.Model, sort string) *gdb.Model {
	m = m.OrderDesc(dao.LiveRoomInfo.Columns().IsTop)
	m = m.OrderDesc(dao.LiveRoomInfo.Columns().ToppedAt)
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
			SegmentTime:    req.SegmentTime,
			Quality:        req.Quality,
			MonitorOnly:    req.MonitorOnly,
			SyncPath:       req.SyncPath,
			Remark:         req.Remark,
			CreatedAt:      utils.Now(),
		})
		if txErr != nil {
			return txErr
		}
		*liveId, txErr = result.LastInsertId()
		if txErr != nil {
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
			return txErr
		}
		return nil
	})
}

func listenerForQuickAdd(liveId int) {
	listenerForAdd(liveId, consts.MonitorTypeIntelligent, "", "")
}

func listenerForAdd(liveId, monitorType int, monitorStartAt, monitorStopAt string) {
	ctx := gctx.GetInitCtx()
	if monitorType == consts.MonitorTypeStart || monitorType == consts.MonitorTypeIntelligent {
		registry.Get().Add(ctx, liveId)
	}
	if monitorType == consts.MonitorTypeCron {
		crons.AddStreamCron(ctx, liveId, monitorStartAt, monitorStopAt)
	}
}

func listenerForUpdate(req *v1.PutLiveManageReq, tempData *entity.LiveManage) {
	ctx := gctx.GetInitCtx()
	if req.MonitorType == consts.MonitorTypeStop {
		listenerForDelete(req.Id)
	}
	if req.MonitorType == consts.MonitorTypeStart || req.MonitorType == consts.MonitorTypeIntelligent {
		if req.MonitorType != tempData.MonitorType || req.Format != tempData.Format || req.Interval != tempData.Interval ||
			req.Quality != tempData.Quality {
			listenerForDelete(req.Id)
			registry.Get().Add(ctx, req.Id)
			return
		}
	}
	if req.MonitorType == consts.MonitorTypeCron {
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

func getExportData(ctx context.Context, req *v1.ExportRoomInfoReq) ([]*model.ExportRoomInfo, error) {
	var list []*model.ExportRoomInfo
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
	count, err := m.Count()
	if err != nil || count <= 0 {
		return nil, utils.TError(ctx, "stream.live.error.GetRoomInfo")
	}
	err = m.Scan(&list)
	if err != nil {
		return nil, utils.TError(ctx, "stream.live.error.GetRoomInfo")
	}
	ids := make([]int, count)
	for i, item := range list {
		ids[i] = item.LiveId
	}
	var mList []*entity.LiveManage
	err = dao.LiveManage.Ctx(ctx).WhereIn(dao.LiveManage.Columns().Id, ids).Scan(&mList)
	if err != nil {
		return nil, utils.TError(ctx, "stream.live.error.GetConfig")
	}
	mMap := make(map[int]*entity.LiveManage, count)
	for _, item := range mList {
		mMap[item.Id] = item
	}
	for _, item := range list {
		if m, ok := mMap[item.LiveId]; ok {
			item.Url = m.RoomUrl
		}
	}
	return list, nil
}

func exportTxt(r *ghttp.Request, ustr int64, list []*model.ExportRoomInfo) {
	r.Response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	filename := utils.Tf(r.Context(), "stream.live.info.ExportTxtTitle", ustr)
	disposition := fmt.Sprintf("attachment; filename=%s; filename*=UTF-8''%s", filename, url.QueryEscape(filename))
	r.Response.Header().Set("Content-Disposition", disposition)

	var builder strings.Builder
	for _, item := range list {
		if item.Url == "" {
			continue
		}
		builder.WriteString(item.Url)
		builder.WriteString("\n")
	}
	r.Response.Write([]byte(builder.String()))
}

func exportExcel(r *ghttp.Request, ustr int64, list []*model.ExportRoomInfo) {
	r.Response.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	filename := utils.Tf(r.Context(), "stream.live.info.ExportExcelTitle", ustr)
	disposition := fmt.Sprintf("attachment; filename=%s; filename*=UTF-8''%s", filename, url.QueryEscape(filename))
	r.Response.Header().Set("Content-Disposition", disposition)
	f := excelize.NewFile()

	const sheet = "房间信息"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)
	_ = f.SetSheetRow(sheet, "A1", &[]string{utils.T(r.Context(), "stream.live.info.ExportURL"), utils.T(r.Context(), "stream.live.info.ExportPlatform"), utils.T(r.Context(), "stream.live.info.ExportAnchor"), utils.T(r.Context(), "stream.live.info.ExportRoomName")})
	for i, row := range list {
		cell := fmt.Sprintf("A%d", i+2)
		_ = f.SetSheetRow(sheet, cell, &[]string{
			row.Url,
			row.Platform,
			row.Anchor,
			row.RoomName,
		})
	}
	_ = f.SetColWidth(sheet, "A", "D", 18)
	_ = f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})
	f.Write(r.Response.Writer)
}
