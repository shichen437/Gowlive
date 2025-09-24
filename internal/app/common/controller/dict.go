package controller

import (
	"context"

	"github.com/shichen437/gowlive/api/v1/common"
	"github.com/shichen437/gowlive/internal/app/common/service"
)

type internalDictController struct {
}

var InternalDict = internalDictController{}

func (c *internalDictController) GetDictDataByType(ctx context.Context, req *common.GetInternalDictByTypeReq) (res *common.GetInternalDictByTypeRes, err error) {
	return service.InternalDict().GetDictDataByType(ctx, req)
}
