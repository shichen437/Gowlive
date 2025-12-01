package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/gowlive/api/v1/common"
	"github.com/shichen437/gowlive/internal/app/media/model/entity"
)

type GetFileCheckListReq struct {
	g.Meta `path:"/media/check/list" method:"get" tags:"文件检测" summary:"获取文件检测任务列表"`
	common.PageReq
	Path     string `json:"path"`
	Filename string `json:"filename"`
}

type GetFileCheckListRes struct {
	g.Meta    `mime:"application/json"`
	Rows      []*entity.FileCheckTask `json:"rows"`
	Total     int                     `json:"total"`
	Executing bool                    `json:"executing"`
}

type PostFileCheckReq struct {
	g.Meta   `path:"/media/check" method:"post" tags:"文件检测" summary:"新增文件检测任务"`
	Path     string `json:"path" v:"required#路径不能为空"`
	Filename string `json:"filename" v:"required#文件名称不能为空"`
}

type PostFileCheckRes struct {
	g.Meta `mime:"application/json"`
}

type DeleteFileCheckReq struct {
	g.Meta `path:"/media/check/{id}" method:"delete" tags:"文件检测" summary:"删除文件检测任务"`
	Id     int `json:"id" v:"required#任务 ID不能为空"`
}

type DeleteFileCheckRes struct {
	g.Meta `mime:"application/json"`
}

type DeleteAllCheckReq struct {
	g.Meta `path:"/media/check/all" method:"delete" tags:"文件检测" summary:"清空记录"`
}

type DeleteAllCheckRes struct {
	g.Meta `mime:"application/json"`
}
