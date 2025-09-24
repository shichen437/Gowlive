package v1

import "github.com/gogf/gf/v2/frame/g"

type GetLatestVersionReq struct {
	g.Meta `path:"/system/latestVersion" method:"get" tags:"系统" summary:"获取最新版本"`
}
type GetLatestVersionRes struct {
	g.Meta        `mime:"application/json"`
	LatestVersion string `json:"latestVersion"`
}
