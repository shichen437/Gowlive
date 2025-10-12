// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AnchorInfoDao is the data access object for the table anchor_info.
type AnchorInfoDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  AnchorInfoColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// AnchorInfoColumns defines and stores column names for the table anchor_info.
type AnchorInfoColumns struct {
	Id             string //
	AnchorName     string //
	Url            string //
	Signature      string //
	Platform       string //
	UniqueId       string //
	FollowerCount  string //
	FollowingCount string //
	LikeCount      string //
	VideoCount     string //
	CreatedAt      string //
	UpdatedAt      string //
}

// anchorInfoColumns holds the columns for the table anchor_info.
var anchorInfoColumns = AnchorInfoColumns{
	Id:             "id",
	AnchorName:     "anchor_name",
	Url:            "url",
	Signature:      "signature",
	Platform:       "platform",
	UniqueId:       "unique_id",
	FollowerCount:  "follower_count",
	FollowingCount: "following_count",
	LikeCount:      "like_count",
	VideoCount:     "video_count",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewAnchorInfoDao creates and returns a new DAO object for table data access.
func NewAnchorInfoDao(handlers ...gdb.ModelHandler) *AnchorInfoDao {
	return &AnchorInfoDao{
		group:    "default",
		table:    "anchor_info",
		columns:  anchorInfoColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *AnchorInfoDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *AnchorInfoDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *AnchorInfoDao) Columns() AnchorInfoColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *AnchorInfoDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *AnchorInfoDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *AnchorInfoDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
