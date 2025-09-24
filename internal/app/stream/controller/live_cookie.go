package controller

import (
	"context"

	v1 "github.com/shichen437/gowlive/api/v1/stream"
	"github.com/shichen437/gowlive/internal/app/stream/service"
)

type liveCookieController struct {
}

var LiveCookie = liveCookieController{}

func (c *liveCookieController) All(ctx context.Context, req *v1.GetAllCookieReq) (res *v1.GetAllCookieRes, err error) {
	return service.LiveCookie().All(ctx, req)
}

func (c *liveCookieController) Add(ctx context.Context, req *v1.PostLiveCookieReq) (res *v1.PostLiveCookieRes, err error) {
	return service.LiveCookie().Add(ctx, req)
}

func (c *liveCookieController) Update(ctx context.Context, req *v1.PutLiveCookieReq) (res *v1.PutLiveCookieRes, err error) {
	return service.LiveCookie().Update(ctx, req)
}

func (c *liveCookieController) Delete(ctx context.Context, req *v1.DeleteLiveCookieReq) (res *v1.DeleteLiveCookieRes, err error) {
	return service.LiveCookie().Delete(ctx, req)
}
