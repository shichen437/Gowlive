package controller

import (
	"context"

	v1 "github.com/shichen437/gowlive/api/v1/system"
	"github.com/shichen437/gowlive/internal/app/system/service"
)

type systemOverviewController struct {
}

var SystemOverview = systemOverviewController{}

func (s *systemOverviewController) Overview(ctx context.Context, req *v1.GetOverviewReq) (res *v1.GetOverviewRes, err error) {
	return service.SystemOverview().Overview(ctx, req)
}

func (s *systemOverviewController) GetLang(ctx context.Context, req *v1.GetLangReq) (res *v1.GetLangRes, err error) {
	return service.SystemOverview().GetLang(ctx, req)
}
