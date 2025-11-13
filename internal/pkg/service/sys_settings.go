package service

import (
	"context"

	"github.com/shichen437/gowlive/internal/app/system/dao"
	"github.com/shichen437/gowlive/internal/app/system/model/entity"
)

func GetSettings(ctx context.Context, key string) int {
	var result *entity.SysSettings
	err := dao.SysSettings.Ctx(ctx).Where(dao.SysSettings.Columns().SKey, key).Limit(1).Scan(&result)
	if err != nil || result == nil {
		return 0
	}
	return result.SValue
}
