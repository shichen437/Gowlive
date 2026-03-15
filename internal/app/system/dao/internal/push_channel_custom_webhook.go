// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PushChannelCustomWebhookDao is the data access object for the table push_channel_custom_webhook.
type PushChannelCustomWebhookDao struct {
	table    string                          // table is the underlying table name of the DAO.
	group    string                          // group is the database configuration group name of the current DAO.
	columns  PushChannelCustomWebhookColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler              // handlers for customized model modification.
}

// PushChannelCustomWebhookColumns defines and stores column names for the table push_channel_custom_webhook.
type PushChannelCustomWebhookColumns struct {
	Id             string //
	ChannelId      string //
	WebhookUrl     string //
	RequestMethod  string //
	RequestHeaders string //
	RequestBody    string //
	CreatedAt      string //
	UpdatedAt      string //
}

// pushChannelCustomWebhookColumns holds the columns for the table push_channel_custom_webhook.
var pushChannelCustomWebhookColumns = PushChannelCustomWebhookColumns{
	Id:             "id",
	ChannelId:      "channel_id",
	WebhookUrl:     "webhook_url",
	RequestMethod:  "request_method",
	RequestHeaders: "request_headers",
	RequestBody:    "request_body",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewPushChannelCustomWebhookDao creates and returns a new DAO object for table data access.
func NewPushChannelCustomWebhookDao(handlers ...gdb.ModelHandler) *PushChannelCustomWebhookDao {
	return &PushChannelCustomWebhookDao{
		group:    "default",
		table:    "push_channel_custom_webhook",
		columns:  pushChannelCustomWebhookColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PushChannelCustomWebhookDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PushChannelCustomWebhookDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PushChannelCustomWebhookDao) Columns() PushChannelCustomWebhookColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PushChannelCustomWebhookDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PushChannelCustomWebhookDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PushChannelCustomWebhookDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
