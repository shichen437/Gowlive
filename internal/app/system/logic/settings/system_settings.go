package logic

import (
	"context"

	v1 "github.com/shichen437/gowlive/api/v1/system"
	"github.com/shichen437/gowlive/internal/app/system/service"
	"github.com/shichen437/gowlive/internal/pkg/crons/system"
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
