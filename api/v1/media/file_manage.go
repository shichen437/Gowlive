package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/gowlive/internal/app/media/model"
)

type GetFileListReq struct {
	g.Meta   `path:"/media/file/list" method:"get" tags:"文件管理" summary:"获取文件列表"`
	Path     string `json:"path"`
	Filename string `json:"filename"`
}

type GetFileListRes struct {
	g.Meta `mime:"application/json"`
	Rows   []*model.FileInfo `json:"rows"`
}

type GetAnchorFilePathReq struct {
	g.Meta   `path:"/media/file/roomPath" method:"get" tags:"文件管理" summary:"获取主播录制文件路径"`
	Anchor   string `json:"anchor"`
	Platform string `json:"platform"`
}
type GetAnchorFilePathRes struct {
	g.Meta `mime:"application/json"`
	Path   string `json:"path"`
}

type DeleteFileReq struct {
	g.Meta   `path:"/media/file" method:"delete" tags:"文件管理" summary:"删除文件"`
	Filename string `json:"filename" v:"required-with:Path#路径为空时文件名必选"`
	Path     string `json:"path"`
}

type DeleteFileRes struct {
	g.Meta `mime:"application/json"`
}

type GetEmptyFolderReq struct {
	g.Meta `path:"/media/file/empty" method:"get" tags:"文件管理" summary:"是否空文件夹"`
	Path   string `json:"path" v:"required#路径不能为空"`
}

type GetEmptyFolderRes struct {
	g.Meta  `mime:"application/json"`
	IsEmpty bool `json:"isEmpty"`
}

type GetFilePlayReq struct {
	g.Meta `path:"/media/file/play" method:"get,post,head,options,trace" tags:"文件管理" summary:"媒体文件流式传输"`
	Path   string `p:"path" v:"required#文件路径不能为空"`
}

type GetFilePlayRes struct {
	g.Meta `mime:"application/octet-stream"`
}
