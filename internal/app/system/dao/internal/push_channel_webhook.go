// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PushChannelWebhookDao is the data access object for the table push_channel_webhook.
type PushChannelWebhookDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  PushChannelWebhookColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// PushChannelWebhookColumns defines and stores column names for the table push_channel_webhook.
type PushChannelWebhookColumns struct {
	Id          string //
	ChannelId   string //
	WebhookUrl  string //
	MessageType string //
	Sign        string //
	At          string //
	CreatedAt   string //
	UpdatedAt   string //
}

// pushChannelWebhookColumns holds the columns for the table push_channel_webhook.
var pushChannelWebhookColumns = PushChannelWebhookColumns{
	Id:          "id",
	ChannelId:   "channel_id",
	WebhookUrl:  "webhook_url",
	MessageType: "message_type",
	Sign:        "sign",
	At:          "at",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewPushChannelWebhookDao creates and returns a new DAO object for table data access.
func NewPushChannelWebhookDao(handlers ...gdb.ModelHandler) *PushChannelWebhookDao {
	return &PushChannelWebhookDao{
		group:    "default",
		table:    "push_channel_webhook",
		columns:  pushChannelWebhookColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PushChannelWebhookDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PushChannelWebhookDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PushChannelWebhookDao) Columns() PushChannelWebhookColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PushChannelWebhookDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PushChannelWebhookDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PushChannelWebhookDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
