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
	ISystemOverview interface {
		Overview(ctx context.Context, req *v1.GetOverviewReq) (res *v1.GetOverviewRes, err error)
	}
)

var (
	localSystemOverview ISystemOverview
)

func SystemOverview() ISystemOverview {
	if localSystemOverview == nil {
		panic("implement not found for interface ISystemOverview, forgot register?")
	}
	return localSystemOverview
}

func RegisterSystemOverview(i ISystemOverview) {
	localSystemOverview = i
}
