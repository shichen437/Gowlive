package controller

import (
	"context"

	v1 "github.com/shichen437/gowlive/api/v1/stream"
	"github.com/shichen437/gowlive/internal/app/stream/service"
)

type liveHistoryController struct {
}

var LiveHistory = liveHistoryController{}

func (c *liveHistoryController) List(ctx context.Context, req *v1.GetLiveHistoryListReq) (res *v1.GetLiveHistoryListRes, err error) {
	return service.LiveHistory().List(ctx, req)
}

func (c *liveHistoryController) Delete(ctx context.Context, req *v1.DeleteLiveHistoryReq) (res *v1.DeleteLiveHistoryRes, err error) {
	return service.LiveHistory().Delete(ctx, req)
}
