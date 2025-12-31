package logic

import (
	"context"

	v1 "github.com/shichen437/gowlive/api/v1/third"
	"github.com/shichen437/gowlive/internal/app/third/service"
	"github.com/shichen437/gowlive/internal/pkg/third/openlist"
)

type sOpenlist struct {
}

func init() {
	service.RegisterOpenlist(New())
}

func New() service.IOpenlist {
	return &sOpenlist{}
}

func (s *sOpenlist) Status(ctx context.Context, req *v1.GetOpenlistStatusReq) (res *v1.GetOpenlistStatusRes, err error) {
	res = &v1.GetOpenlistStatusRes{}
	res.Status = openlist.GetStatus()
	return
}
