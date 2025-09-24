// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// LiveRoomInfoDao is the data access object for the table live_room_info.
type LiveRoomInfoDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  LiveRoomInfoColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// LiveRoomInfoColumns defines and stores column names for the table live_room_info.
type LiveRoomInfoColumns struct {
	Id        string //
	LiveId    string //
	RoomName  string //
	Anchor    string //
	Platform  string //
	Status    string //
	CreatedAt string //
	UpdatedAt string //
}

// liveRoomInfoColumns holds the columns for the table live_room_info.
var liveRoomInfoColumns = LiveRoomInfoColumns{
	Id:        "id",
	LiveId:    "live_id",
	RoomName:  "room_name",
	Anchor:    "anchor",
	Platform:  "platform",
	Status:    "status",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewLiveRoomInfoDao creates and returns a new DAO object for table data access.
func NewLiveRoomInfoDao(handlers ...gdb.ModelHandler) *LiveRoomInfoDao {
	return &LiveRoomInfoDao{
		group:    "default",
		table:    "live_room_info",
		columns:  liveRoomInfoColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *LiveRoomInfoDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *LiveRoomInfoDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *LiveRoomInfoDao) Columns() LiveRoomInfoColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *LiveRoomInfoDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *LiveRoomInfoDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *LiveRoomInfoDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
