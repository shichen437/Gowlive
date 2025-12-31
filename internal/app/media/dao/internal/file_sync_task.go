// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// FileSyncTaskDao is the data access object for the table file_sync_task.
type FileSyncTaskDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  FileSyncTaskColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// FileSyncTaskColumns defines and stores column names for the table file_sync_task.
type FileSyncTaskColumns struct {
	Id        string //
	Path      string //
	Filename  string //
	SyncPath  string //
	Duration  string //
	Status    string //
	Remark    string //
	CreatedAt string //
	UpdatedAt string //
}

// fileSyncTaskColumns holds the columns for the table file_sync_task.
var fileSyncTaskColumns = FileSyncTaskColumns{
	Id:        "id",
	Path:      "path",
	Filename:  "filename",
	SyncPath:  "sync_path",
	Duration:  "duration",
	Status:    "status",
	Remark:    "remark",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewFileSyncTaskDao creates and returns a new DAO object for table data access.
func NewFileSyncTaskDao(handlers ...gdb.ModelHandler) *FileSyncTaskDao {
	return &FileSyncTaskDao{
		group:    "default",
		table:    "file_sync_task",
		columns:  fileSyncTaskColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *FileSyncTaskDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *FileSyncTaskDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *FileSyncTaskDao) Columns() FileSyncTaskColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *FileSyncTaskDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *FileSyncTaskDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *FileSyncTaskDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
