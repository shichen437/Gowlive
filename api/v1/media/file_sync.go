package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/gowlive/api/v1/common"
	"github.com/shichen437/gowlive/internal/app/media/model/entity"
)

type GetFileSyncListReq struct {
	g.Meta `path:"/media/sync/list" method:"get" tags:"文件同步" summary:"获取文件同步任务列表"`
	common.PageReq
	Path     string `json:"path"`
	Filename string `json:"filename"`
	Status   *int   `json:"status"`
}

type GetFileSyncListRes struct {
	g.Meta `mime:"application/json"`
	Rows   []*entity.FileSyncTask `json:"rows"`
	Total  int                    `json:"total"`
}

type ResyncFileReq struct {
	g.Meta `path:"/media/sync/resync/{id}" method:"put" tags:"文件同步" summary:"重新同步文件"`
	Id     int `json:"id" v:"required"`
}

type ResyncFileRes struct {
	g.Meta `mime:"application/json"`
}

type DeleteFileSyncReq struct {
	g.Meta `path:"/media/sync/{id}" method:"delete" tags:"文件同步" summary:"删除文件同步任务"`
	Id     int `json:"id" v:"required"`
}

type DeleteFileSyncRes struct {
	g.Meta `mime:"application/json"`
}

type DeleteAllSyncReq struct {
	g.Meta `path:"/media/sync/all" method:"delete" tags:"文件同步" summary:"清空同步记录"`
	Status *int `json:"status"`
}

type DeleteAllSyncRes struct {
	g.Meta `mime:"application/json"`
}
