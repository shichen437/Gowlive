// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// LiveManageDao is the data access object for the table live_manage.
type LiveManageDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  LiveManageColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// LiveManageColumns defines and stores column names for the table live_manage.
type LiveManageColumns struct {
	Id             string //
	RoomUrl        string //
	Interval       string //
	Format         string //
	MonitorType    string //
	MonitorStartAt string //
	MonitorStopAt  string //
	Remark         string //
	CreatedAt      string //
	UpdatedAt      string //
	Quality        string //
}

// liveManageColumns holds the columns for the table live_manage.
var liveManageColumns = LiveManageColumns{
	Id:             "id",
	RoomUrl:        "room_url",
	Interval:       "interval",
	Format:         "format",
	MonitorType:    "monitor_type",
	MonitorStartAt: "monitor_start_at",
	MonitorStopAt:  "monitor_stop_at",
	Remark:         "remark",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	Quality:        "quality",
}

// NewLiveManageDao creates and returns a new DAO object for table data access.
func NewLiveManageDao(handlers ...gdb.ModelHandler) *LiveManageDao {
	return &LiveManageDao{
		group:    "default",
		table:    "live_manage",
		columns:  liveManageColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *LiveManageDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *LiveManageDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *LiveManageDao) Columns() LiveManageColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *LiveManageDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *LiveManageDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *LiveManageDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
