// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// FileCheckTaskDao is the data access object for the table file_check_task.
type FileCheckTaskDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  FileCheckTaskColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// FileCheckTaskColumns defines and stores column names for the table file_check_task.
type FileCheckTaskColumns struct {
	Id         string //
	Path       string //
	Filename   string //
	Duration   string //
	Progress   string //
	FileStatus string //
	CreatedAt  string //
	UpdatedAt  string //
}

// fileCheckTaskColumns holds the columns for the table file_check_task.
var fileCheckTaskColumns = FileCheckTaskColumns{
	Id:         "id",
	Path:       "path",
	Filename:   "filename",
	Duration:   "duration",
	Progress:   "progress",
	FileStatus: "file_status",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
}

// NewFileCheckTaskDao creates and returns a new DAO object for table data access.
func NewFileCheckTaskDao(handlers ...gdb.ModelHandler) *FileCheckTaskDao {
	return &FileCheckTaskDao{
		group:    "default",
		table:    "file_check_task",
		columns:  fileCheckTaskColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *FileCheckTaskDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *FileCheckTaskDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *FileCheckTaskDao) Columns() FileCheckTaskColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *FileCheckTaskDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *FileCheckTaskDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *FileCheckTaskDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
