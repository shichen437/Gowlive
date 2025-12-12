package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/goflyfox/gtoken/gtoken"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/util/gconv"

	Admin "github.com/shichen437/gowlive/internal/app/admin/controller"
	Common "github.com/shichen437/gowlive/internal/app/common/controller"
	Media "github.com/shichen437/gowlive/internal/app/media/controller"
	Stream "github.com/shichen437/gowlive/internal/app/stream/controller"
	System "github.com/shichen437/gowlive/internal/app/system/controller"

	"github.com/shichen437/gowlive/internal/app/admin/dao"
	"github.com/shichen437/gowlive/internal/app/admin/model/do"
	"github.com/shichen437/gowlive/internal/app/admin/model/entity"
	"github.com/shichen437/gowlive/internal/app/common/service"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/lives"
	"github.com/shichen437/gowlive/internal/pkg/manager"
	"github.com/shichen437/gowlive/internal/pkg/registry"
	"github.com/shichen437/gowlive/internal/pkg/sse"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "Gowlive backend api",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			g.Log().Info(ctx, "Starting HTTP Server...")
			gfToken, err := GetGtoken(ctx)
			if err != nil {
				g.Log().Error(ctx, err)
				return err
			}
			initDir()
			s := g.Server()
			s.SetGraceful(true)
			s.SetIndexFolder(true)
			s.SetServerRoot(utils.DATA_PATH)
			s.AddSearchPath(utils.DATA_PATH)
			s.SetSwaggerUITemplate(consts.MySwaggerUITemplate)
			s.Use(service.Middleware().HandlerResponse)
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(
					service.Middleware().Ctx,
					ghttp.MiddlewareCORS,
				)
				group.Group("/", func(group *ghttp.RouterGroup) {
					err = gfToken.Middleware(ctx, group)
					if err != nil {
						panic(err)
					}
					group.ALL("/sse", sse.HandleSSE)
					bindRoute(group)
				})
			})
			if err = s.Start(); err != nil {
				g.Log().Errorf(ctx, "server 启动异常，错误信息：%v", err)
				return err
			}
			g.Go(ctx, func(c context.Context) {
				time.Sleep(1500 * time.Millisecond)
				CheckFile()
				JobInit()
				LiveMonitor()
			}, func(c context.Context, err error) {
				g.Log().Errorf(c, "starter 启动异常，错误信息：%v", err)
			})
			gproc.AddSigHandlerShutdown(shutdown)
			gproc.Listen()
			return nil
		},
	}
)

func GetGtoken(ctx context.Context) (gfToken *gtoken.GfToken, err error) {
	gfToken = &gtoken.GfToken{
		AuthAfterFunc:    AuthAfterFunc,
		AuthExcludePaths: g.SliceStr{"/logout", "/system/lang"},
		AuthPaths:        g.SliceStr{"/*"},
		CacheMode:        3,
		Timeout:          30 * 86400 * 1000,
		LoginBeforeFunc:  LoginFunc,
		LoginPath:        "/login",
		LogoutAfterFunc:  LogoutAfterFunc,
		LogoutPath:       "post:/logout",
		MultiLogin:       consts.MultiLogin,
		ServerName:       consts.ServerName,
	}
	err = gfToken.Start()
	return
}

func LoginFunc(r *ghttp.Request) (string, any) {
	ctx := context.TODO()
	username := r.Get("username").String()
	password := r.Get("password").String()
	if username == "" || password == "" {
		r.Response.WriteJson(gtoken.Fail(utils.T(r.Context(), "user.login.blank")))
		r.ExitAll()
	}
	var users *entity.SysUser
	enc, _ := utils.Encrypt(ctx, password)
	err := dao.SysUser.Ctx(ctx).Where(do.SysUser{
		Username: username,
		Password: enc,
	}).Scan(&users)
	if err != nil || users == nil {
		manager.GetLogManager().AddErrorLog(consts.LogTypeUser, utils.T(r.Context(), "user.login.error"))
		r.Response.WriteJson(gtoken.Fail(utils.T(r.Context(), "user.login.error")))
		r.ExitAll()
	}
	if users.Status == consts.StatusDisable {
		r.Response.WriteJson(gtoken.Fail(utils.T(r.Context(), "user.login.banned")))
		r.ExitAll()
	}
	manager.GetLogManager().AddSuccessLog(consts.LogTypeUser, utils.T(r.Context(), "user.login.success"))
	return fmt.Sprintf("%s%d", consts.GTokenAdminPrefix, users.Id), users
}

func LogoutAfterFunc(r *ghttp.Request, respData gtoken.Resp) {
	service.Session().RemoveUser(r.Context())
	manager.GetLogManager().AddSuccessLog(consts.LogTypeUser, utils.T(r.Context(), "user.login.exit"))
	r.Middleware.Next()
}

func AuthAfterFunc(r *ghttp.Request, respData gtoken.Resp) {
	var users entity.SysUser
	err := gconv.Struct(respData.GetString("data"), &users)
	if err != nil {
		r.Response.WriteJson(gtoken.Unauthorized(utils.T(r.Context(), "user.login.unauthorized"), nil))
		r.ExitAll()
	}
	service.Session().SetUser(r.Context(), &users)
	r.SetCtxVar(consts.CtxAdminId, users.Id)
	r.SetCtxVar(consts.CtxAdminName, users.Username)
	r.Middleware.Next()
}

func initDir() {
	os.MkdirAll(utils.DATA_PATH, os.ModePerm)
	os.MkdirAll(utils.STREAM_PATH, os.ModePerm)
	os.MkdirAll(utils.DOWNLOAD_PATH, os.ModePerm)
}

func bindRoute(group *ghttp.RouterGroup) {
	group.Bind(Admin.SysUser,
		Common.InternalDict,
		Media.FileManage, Media.FileCheck,
		Stream.LiveManage, Stream.LiveHistory, Stream.LiveCookie, Stream.AnchorInfo,
		System.SystemOverview, System.SystemSettings, System.SysLogs, System.PushChannel, System.SysNotify)
}

func shutdown(sig os.Signal) {
	registry.Get().StopAll(gctx.GetInitCtx())
	manager.GetLogManager().Stop()
	manager.GetNotifyManager().Stop()
	manager.GetFileCheckManager().Close()
	lives.GetBucketManager().Stop()
	g.Log().Info(gctx.GetInitCtx(), "all monitor shutdown!")
}
