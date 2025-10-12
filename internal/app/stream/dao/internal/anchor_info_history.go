// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AnchorInfoHistoryDao is the data access object for the table anchor_info_history.
type AnchorInfoHistoryDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  AnchorInfoHistoryColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// AnchorInfoHistoryColumns defines and stores column names for the table anchor_info_history.
type AnchorInfoHistoryColumns struct {
	Id             string //
	AnchorId       string //
	AnchorName     string //
	Signature      string //
	FollowerCount  string //
	FollowingCount string //
	LikeCount      string //
	VideoCount     string //
	CollectedDate  string //
	CreatedAt      string //
	UpdatedAt      string //
}

// anchorInfoHistoryColumns holds the columns for the table anchor_info_history.
var anchorInfoHistoryColumns = AnchorInfoHistoryColumns{
	Id:             "id",
	AnchorId:       "anchor_id",
	AnchorName:     "anchor_name",
	Signature:      "signature",
	FollowerCount:  "follower_count",
	FollowingCount: "following_count",
	LikeCount:      "like_count",
	VideoCount:     "video_count",
	CollectedDate:  "collected_date",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewAnchorInfoHistoryDao creates and returns a new DAO object for table data access.
func NewAnchorInfoHistoryDao(handlers ...gdb.ModelHandler) *AnchorInfoHistoryDao {
	return &AnchorInfoHistoryDao{
		group:    "default",
		table:    "anchor_info_history",
		columns:  anchorInfoHistoryColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *AnchorInfoHistoryDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *AnchorInfoHistoryDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *AnchorInfoHistoryDao) Columns() AnchorInfoHistoryColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *AnchorInfoHistoryDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *AnchorInfoHistoryDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *AnchorInfoHistoryDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
