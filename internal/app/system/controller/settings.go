package controller

import (
	"context"

	v1 "github.com/shichen437/gowlive/api/v1/system"
	"github.com/shichen437/gowlive/internal/app/system/service"
)

type systemSettingsController struct {
}

var SystemSettings = systemSettingsController{}

func (s *systemSettingsController) LatestVersion(ctx context.Context, req *v1.GetLatestVersionReq) (res *v1.GetLatestVersionRes, err error) {
	return service.SystemSettings().LatestVersion(ctx, req)
}
