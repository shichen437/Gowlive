package controller

import (
	"context"

	v1 "github.com/shichen437/gowlive/api/v1/stream"
	"github.com/shichen437/gowlive/internal/app/stream/service"
)

type anchorInfoController struct {
}

var AnchorInfo = anchorInfoController{}

func (c *anchorInfoController) List(ctx context.Context, req *v1.GetAnchorListReq) (res *v1.GetAnchorListRes, err error) {
	return service.AnchorInfo().List(ctx, req)
}

func (c *anchorInfoController) Add(ctx context.Context, req *v1.PostAnchorReq) (res *v1.PostAnchorRes, err error) {
	return service.AnchorInfo().Add(ctx, req)
}

func (c *anchorInfoController) Delete(ctx context.Context, req *v1.DeleteAnchorReq) (res *v1.DeleteAnchorRes, err error) {
	return service.AnchorInfo().Delete(ctx, req)
}

func (c *anchorInfoController) StatInfo(ctx context.Context, req *v1.GetAnchorStatInfoReq) (res *v1.GetAnchorStatInfoRes, err error) {
	return service.AnchorInfo().StatInfo(ctx, req)
}
