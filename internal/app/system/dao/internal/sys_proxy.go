// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysProxyDao is the data access object for the table sys_proxy.
type SysProxyDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  SysProxyColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// SysProxyColumns defines and stores column names for the table sys_proxy.
type SysProxyColumns struct {
	Id        string //
	Platform  string //
	Proxy     string //
	Remark    string //
	CreatedAt string //
	UpdatedAt string //
}

// sysProxyColumns holds the columns for the table sys_proxy.
var sysProxyColumns = SysProxyColumns{
	Id:        "id",
	Platform:  "platform",
	Proxy:     "proxy",
	Remark:    "remark",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewSysProxyDao creates and returns a new DAO object for table data access.
func NewSysProxyDao(handlers ...gdb.ModelHandler) *SysProxyDao {
	return &SysProxyDao{
		group:    "default",
		table:    "sys_proxy",
		columns:  sysProxyColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SysProxyDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SysProxyDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SysProxyDao) Columns() SysProxyColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SysProxyDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SysProxyDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *SysProxyDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
