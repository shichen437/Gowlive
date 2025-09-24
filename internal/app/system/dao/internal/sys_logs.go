// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysLogsDao is the data access object for the table sys_logs.
type SysLogsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  SysLogsColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// SysLogsColumns defines and stores column names for the table sys_logs.
type SysLogsColumns struct {
	Id        string //
	Type      string //
	Content   string //
	Status    string //
	CreatedAt string //
}

// sysLogsColumns holds the columns for the table sys_logs.
var sysLogsColumns = SysLogsColumns{
	Id:        "id",
	Type:      "type",
	Content:   "content",
	Status:    "status",
	CreatedAt: "created_at",
}

// NewSysLogsDao creates and returns a new DAO object for table data access.
func NewSysLogsDao(handlers ...gdb.ModelHandler) *SysLogsDao {
	return &SysLogsDao{
		group:    "default",
		table:    "sys_logs",
		columns:  sysLogsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SysLogsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SysLogsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SysLogsDao) Columns() SysLogsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SysLogsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SysLogsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SysLogsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
