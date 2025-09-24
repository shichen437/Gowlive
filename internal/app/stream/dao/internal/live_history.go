// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// LiveHistoryDao is the data access object for the table live_history.
type LiveHistoryDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  LiveHistoryColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// LiveHistoryColumns defines and stores column names for the table live_history.
type LiveHistoryColumns struct {
	Id        string //
	LiveId    string //
	Anchor    string //
	StartedAt string //
	EndedAt   string //
	Duration  string //
	CreatedAt string //
	UpdatedAt string //
}

// liveHistoryColumns holds the columns for the table live_history.
var liveHistoryColumns = LiveHistoryColumns{
	Id:        "id",
	LiveId:    "live_id",
	Anchor:    "anchor",
	StartedAt: "started_at",
	EndedAt:   "ended_at",
	Duration:  "duration",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewLiveHistoryDao creates and returns a new DAO object for table data access.
func NewLiveHistoryDao(handlers ...gdb.ModelHandler) *LiveHistoryDao {
	return &LiveHistoryDao{
		group:    "default",
		table:    "live_history",
		columns:  liveHistoryColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *LiveHistoryDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *LiveHistoryDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *LiveHistoryDao) Columns() LiveHistoryColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *LiveHistoryDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *LiveHistoryDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *LiveHistoryDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
