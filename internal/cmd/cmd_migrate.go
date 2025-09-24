package cmd

import (
	"context"
	"os"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/golang-migrate/migrate/v4"

	"github.com/shichen437/gowlive/internal/app/admin/dao"
	"github.com/shichen437/gowlive/internal/app/admin/model/do"
	"github.com/shichen437/gowlive/internal/app/admin/model/entity"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/utils"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	DB_URL     = "sqlite3://resources/data/gowlive.db"
	SOURCE_URL = "file://manifest/migrate"
)

func Migrate() (err error) {
	ctx := gctx.New()
	g.Log().Info(ctx, "Starting Migrate...")
	os.MkdirAll(utils.DATA_PATH, os.ModePerm)
	m, err := migrate.New(SOURCE_URL, DB_URL)
	if err != nil {
		return
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		g.Log().Error(ctx, "Database Migrate Failed", err)
		return
	}
	err = initAdminInfo(ctx)
	if err != nil {
		g.Log().Error(ctx, "初始化管理员信息失败")
		return
	}
	g.Log().Info(ctx, "Migrate Done!")
	return nil
}

func initAdminInfo(ctx context.Context) (err error) {
	var users *entity.SysUser
	err = dao.SysUser.Ctx(ctx).WherePri(1).Scan(&users)
	if err != nil {
		g.Log().Error(ctx, "获取管理员账户信息失败")
		return
	}
	enc, _ := utils.Encrypt(ctx, consts.DefaultPassword)
	if users == nil {
		_, err = dao.SysUser.Ctx(ctx).Insert(do.SysUser{
			Id:        1,
			Username:  "admin",
			Nickname:  "admin",
			Password:  enc,
			CreatedAt: utils.Now(),
		})
	}
	return
}
