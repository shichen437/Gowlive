// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"github.com/shichen437/gowlive/internal/app/admin/model/entity"
)

type (
	ISession interface {
		SetUser(ctx context.Context, user *entity.SysUser) error
		GetUser(ctx context.Context) *entity.SysUser
		RemoveUser(ctx context.Context) error
	}
)

var (
	localSession ISession
)

func Session() ISession {
	if localSession == nil {
		panic("implement not found for interface ISession, forgot register?")
	}
	return localSession
}

func RegisterSession(i ISession) {
	localSession = i
}
