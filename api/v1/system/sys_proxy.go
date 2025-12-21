package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/gowlive/internal/app/system/model/entity"
)

type GetAllProxyReq struct {
	g.Meta `path:"/system/proxy/list" method:"get" tags:"代理管理" summary:"获取所有代理"`
}
type GetAllProxyRes struct {
	g.Meta `mime:"application/json"`
	Rows   []*entity.SysProxy `json:"rows"`
}

type PostSysProxyReq struct {
	g.Meta   `path:"/system/proxy" method:"post" tags:"代理管理" summary:"添加代理"`
	Platform string   `json:"platform" dc:"平台" v:"required#system.proxy.valid.PlatformRequired"`
	Proxy    []string `json:"proxy" dc:"代理" v:"required#system.proxy.valid.ProxyRequired"`
	Remark   string   `json:"remark" dc:"备注" v:"max-length:45#system.proxy.valid.RemarkMaxLength"`
}
type PostSysProxyRes struct {
	g.Meta `mime:"application/json"`
}

type PutSysProxyReq struct {
	g.Meta `path:"/system/proxy" method:"put" tags:"代理管理" summary:"修改代理"`
	Id     int      `json:"id" dc:"id" v:"required#system.proxy.valid.IDRequired"`
	Proxy  []string `json:"proxy" dc:"代理" v:"required#system.proxy.valid.ProxyRequired"`
	Remark string   `json:"remark" dc:"备注" v:"max-length:45#system.proxy.valid.RemarkMaxLength"`
}
type PutSysProxyRes struct {
	g.Meta `mime:"application/json"`
}

type DeleteSysProxyReq struct {
	g.Meta `path:"/system/proxy/{id}" method:"delete" tags:"代理管理" summary:"删除代理"`
	Id     int `json:"id" dc:"id" v:"required#system.proxy.valid.IDRequired"`
}
type DeleteSysProxyRes struct {
	g.Meta `mime:"application/json"`
}
