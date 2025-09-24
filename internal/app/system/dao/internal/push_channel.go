// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PushChannelDao is the data access object for the table push_channel.
type PushChannelDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  PushChannelColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// PushChannelColumns defines and stores column names for the table push_channel.
type PushChannelColumns struct {
	Id        string //
	Name      string //
	Type      string //
	Status    string //
	Url       string //
	Remark    string //
	CreatedAt string //
	UpdatedAt string //
}

// pushChannelColumns holds the columns for the table push_channel.
var pushChannelColumns = PushChannelColumns{
	Id:        "id",
	Name:      "name",
	Type:      "type",
	Status:    "status",
	Url:       "url",
	Remark:    "remark",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewPushChannelDao creates and returns a new DAO object for table data access.
func NewPushChannelDao(handlers ...gdb.ModelHandler) *PushChannelDao {
	return &PushChannelDao{
		group:    "default",
		table:    "push_channel",
		columns:  pushChannelColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PushChannelDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PushChannelDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PushChannelDao) Columns() PushChannelColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PushChannelDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PushChannelDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PushChannelDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
