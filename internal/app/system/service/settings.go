// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	v1 "github.com/shichen437/gowlive/api/v1/system"
)

type (
	ISystemSettings interface {
		LatestVersion(ctx context.Context, req *v1.GetLatestVersionReq) (res *v1.GetLatestVersionRes, err error)
		PutSettings(ctx context.Context, req *v1.PutSysSettingsReq) (res *v1.PutSysSettingsRes, err error)
		GetSettings(ctx context.Context, req *v1.GetSysSettingsReq) (res *v1.GetSysSettingsRes, err error)
	}
)

var (
	localSystemSettings ISystemSettings
)

func SystemSettings() ISystemSettings {
	if localSystemSettings == nil {
		panic("implement not found for interface ISystemSettings, forgot register?")
	}
	return localSystemSettings
}

func RegisterSystemSettings(i ISystemSettings) {
	localSystemSettings = i
}
