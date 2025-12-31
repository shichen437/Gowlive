package controller

import (
	"context"

	v1 "github.com/shichen437/gowlive/api/v1/third"
	"github.com/shichen437/gowlive/internal/app/third/service"
)

type openlistController struct {
}

var Openlist = openlistController{}

func (s *openlistController) Status(ctx context.Context, req *v1.GetOpenlistStatusReq) (res *v1.GetOpenlistStatusRes, err error) {
	return service.Openlist().Status(ctx, req)
}
