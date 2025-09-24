package cmd

import (
	"context"
	"fmt"
	"os"

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
	"github.com/shichen437/gowlive/internal/pkg/manager"
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
					bindRoute(group)
				})
			})
			JobInit()
			LiveMonitor()
			s.Run()
			go func() {
				gproc.AddSigHandlerShutdown(shutdown)
				gproc.Listen()
			}()
			return nil
		},
	}
)

func GetGtoken(ctx context.Context) (gfToken *gtoken.GfToken, err error) {
	gfToken = &gtoken.GfToken{
		AuthAfterFunc:    AuthAfterFunc,
		AuthExcludePaths: g.SliceStr{"/logout"},
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

func LoginFunc(r *ghttp.Request) (string, interface{}) {
	ctx := context.TODO()
	username := r.Get("username").String()
	password := r.Get("password").String()
	if username == "" || password == "" {
		r.Response.WriteJson(gtoken.Fail("用户名或密码不能为空"))
		r.ExitAll()
	}
	var users *entity.SysUser
	enc, _ := utils.Encrypt(ctx, password)
	err := dao.SysUser.Ctx(ctx).Where(do.SysUser{
		Username: username,
		Password: enc,
	}).Scan(&users)
	if err != nil || users == nil {
		manager.GetLogManager().AddErrorLog(consts.LogTypeUser, "用户名或密码错误")
		r.Response.WriteJson(gtoken.Fail("用户名或密码错误"))
		r.ExitAll()
	}
	if users.Status == consts.StatusDisable {
		r.Response.WriteJson(gtoken.Fail("用户已被禁用"))
		r.ExitAll()
	}
	manager.GetLogManager().AddSuccessLog(consts.LogTypeUser, "用户登录成功")
	return fmt.Sprintf("%s%d", consts.GTokenAdminPrefix, users.Id), users
}

func LogoutAfterFunc(r *ghttp.Request, respData gtoken.Resp) {
	service.Session().RemoveUser(r.Context())
	manager.GetLogManager().AddSuccessLog(consts.LogTypeUser, "用户退出登录")
	r.Middleware.Next()
}

func AuthAfterFunc(r *ghttp.Request, respData gtoken.Resp) {
	var users entity.SysUser
	err := gconv.Struct(respData.GetString("data"), &users)
	if err != nil {
		r.Response.WriteJson(gtoken.Unauthorized("未授权", nil))
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
		Media.FileManage,
		Stream.LiveManage, Stream.LiveHistory, Stream.LiveCookie,
		System.SystemOverview, System.SystemSettings, System.SysLogs)
}

func shutdown(sig os.Signal) {
	manager.GetLogManager().Stop()
	g.Log().Info(gctx.GetInitCtx(), "live monitor shutdown!")
}
