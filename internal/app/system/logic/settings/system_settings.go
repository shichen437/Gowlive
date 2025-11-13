package logic

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	v1 "github.com/shichen437/gowlive/api/v1/system"
	"github.com/shichen437/gowlive/internal/app/system/dao"
	"github.com/shichen437/gowlive/internal/app/system/model/do"
	"github.com/shichen437/gowlive/internal/app/system/service"
	"github.com/shichen437/gowlive/internal/pkg/crons/system"
	"github.com/shichen437/gowlive/internal/pkg/manager"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

type sSystemSettings struct {
}

func init() {
	service.RegisterSystemSettings(New())
}

func New() service.ISystemSettings {
	return &sSystemSettings{}
}

func (s *sSystemSettings) LatestVersion(ctx context.Context, req *v1.GetLatestVersionReq) (res *v1.GetLatestVersionRes, err error) {
	res = &v1.GetLatestVersionRes{}
	res.LatestVersion = system.GetLatestVersion(ctx, false)
	return
}

func (s *sSystemSettings) PutSettings(ctx context.Context, req *v1.PutSysSettingsReq) (res *v1.PutSysSettingsRes, err error) {
	count, err := dao.SysSettings.Ctx(ctx).Where(dao.SysSettings.Columns().SKey, req.Key).Count()
	if err != nil {
		return nil, gerror.New("获取设置数据失败")
	}
	if count > 0 {
		_, err = dao.SysSettings.Ctx(ctx).
			Where(dao.SysSettings.Columns().SKey, req.Key).
			Update(do.SysSettings{
				SValue:    req.Value,
				UpdatedAt: utils.Now(),
			})
		if err != nil {
			return nil, gerror.New("修改设置失败")
		}
	} else {
		_, err = dao.SysSettings.Ctx(ctx).
			Insert(do.SysSettings{
				SKey:      req.Key,
				SValue:    req.Value,
				CreatedAt: utils.Now(),
			})
		if err != nil {
			return nil, gerror.New("添加设置失败")
		}
	}
	manager.GetSettingsManager().SaveSetting(req.Key, req.Value)
	return
}

func (s *sSystemSettings) GetSettings(ctx context.Context, req *v1.GetSysSettingsReq) (res *v1.GetSysSettingsRes, err error) {
	res = &v1.GetSysSettingsRes{}
	rm := make(map[string]int)
	for key := range strings.SplitSeq(req.Key, ",") {
		value := manager.GetSettingsManager().GetSetting(key)
		rm[key] = value
	}
	res.Data = rm
	return
}
