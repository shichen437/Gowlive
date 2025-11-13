package v1

import "github.com/gogf/gf/v2/frame/g"

type GetLatestVersionReq struct {
	g.Meta `path:"/system/latestVersion" method:"get" tags:"系统管理" summary:"获取最新版本"`
}
type GetLatestVersionRes struct {
	g.Meta        `mime:"application/json"`
	LatestVersion string `json:"latestVersion"`
}

type GetSysSettingsReq struct {
	g.Meta `path:"/system/settings" method:"get" tags:"系统管理" summary:"获取系统配置"`
	Key    string `json:"key" v:"required#键名不能为空"`
}
type GetSysSettingsRes struct {
	g.Meta `mime:"application/json"`
	Data   map[string]int `json:"data"`
}

type PutSysSettingsReq struct {
	g.Meta `path:"/system/settings" method:"put" tags:"系统管理" summary:"修改系统配置"`
	Key    string `json:"key" v:"required#键名不能为空"`
	Value  int    `json:"value" v:"required#值不能为空"`
}
type PutSysSettingsRes struct {
	g.Meta `mime:"application/json"`
}
