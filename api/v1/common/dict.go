package common

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/gowlive/internal/app/common/model"
)

type GetInternalDictByTypeReq struct {
	g.Meta   `path:"/dict/internal/type" method:"get" tags:"字典" summary:"根据类型获取字典数据"`
	DictType string `v:"required#请输入字典类型" json:"dictType"`
}
type GetInternalDictByTypeRes struct {
	g.Meta `mime:"application/json"`
	Data   []*model.InternalDict `json:"data"`
}
