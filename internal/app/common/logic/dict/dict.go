package dict

import (
	"context"

	"github.com/shichen437/gowlive/api/v1/common"
	"github.com/shichen437/gowlive/internal/app/common/model"
	"github.com/shichen437/gowlive/internal/app/common/service"
)

type sInternalDict struct {
}

func init() {
	service.RegisterInternalDict(New())
}

func New() service.IInternalDict {
	return &sInternalDict{}
}

func (s *sInternalDict) GetDictDataByType(ctx context.Context, req *common.GetInternalDictByTypeReq) (res *common.GetInternalDictByTypeRes, err error) {
	res = &common.GetInternalDictByTypeRes{}
	if data := model.GetDictDataByType(req.DictType); data != nil {
		res.Data = make([]*model.InternalDict, len(*data))
		for i, v := range *data {
			res.Data[i] = &v
		}
	}
	return res, nil
}
