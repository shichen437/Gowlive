package openlist

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

var (
	_token     string
	_status    int
	loginMutex sync.Mutex
)

func withRetry(f func(ctx context.Context) (*gclient.Response, error)) (*gclient.Response, error) {
	ctx := gctx.GetInitCtx()
	resp, err := f(ctx)

	if checkResultUnauthorized(resp) {
		g.Log().Warning(ctx, "Openlist token expired, attempting to refresh.")

		loginMutex.Lock()
		defer loginMutex.Unlock()
		if job := gcron.Search("openlist-login-retry"); job != nil {
			return resp, err
		}

		probeResp, probeErr := getClientWithToken().Get(ctx, utils.OPENLIST_DOMAIN+userInfoUrl)
		if probeErr == nil && probeResp.StatusCode == http.StatusOK {
			probeResp.Close()
			g.Log().Info(ctx, "Token already refreshed, retrying original request.")
			return f(ctx)
		}
		if probeResp != nil {
			probeResp.Close()
		}

		g.Log().Info(ctx, "Performing login to refresh Openlist token.")
		if loginErr := Login(); loginErr != nil {
			g.Log().Error(ctx, "Failed to refresh Openlist token, scheduling retry.", loginErr)
			scheduleLoginRetry(ctx)
			return nil, gerror.Wrap(loginErr, "Openlist token refresh failed")
		}

		g.Log().Info(ctx, "Openlist token refreshed successfully, retrying original request.")
		return f(ctx)
	}

	return resp, err
}

func scheduleLoginRetry(ctx context.Context) {
	jobName := "openlist-login-retry"
	if gcron.Search(jobName) != nil {
		g.Log().Debug(ctx, "Openlist login retry job already exists.")
		return
	}

	g.Log().Info(ctx, "Scheduling Openlist login retry job.")
	_, err := gcron.AddSingleton(ctx, "@every 5m", func(ctx context.Context) {
		g.Log().Info(ctx, "Executing scheduled Openlist login retry.")
		if err := Login(); err == nil {
			g.Log().Info(ctx, "Scheduled Openlist login successful, removing retry job.")
			gcron.Remove(jobName)
		} else {
			g.Log().Warning(ctx, "Scheduled Openlist login retry failed.", err)
		}
	}, jobName)

	if err != nil {
		g.Log().Error(ctx, "Failed to schedule Openlist login retry job.", err)
	}
}

func getClientWithToken() *gclient.Client {
	c := g.Client()
	c.SetHeader("Authorization", _token)
	return c
}

func checkResultUnauthorized(resp *gclient.Response) bool {
	if resp == nil {
		return false
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return true
	}
	if resp.Body == nil {
		return false
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	var cresp *CommonResp
	if err := json.Unmarshal(b, &cresp); err != nil {
		return false
	}
	if cresp.Code == http.StatusUnauthorized {
		resp.Body.Close()
		return true
	}
	return false
}
