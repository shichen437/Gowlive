package utils

import (
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/genv"
)

const (
	E_PROJECT_DATA   = "PROJECT_DATA"
	E_PROJECT_SM4KEY = "PROJECT_SM4KEY"
)

var (
	DATA_PATH     = getEnvWithDefault(E_PROJECT_DATA)
	SM4_KEY       = getEnvWithDefault(E_PROJECT_SM4KEY)
	STREAM_PATH   = DATA_PATH + "/stream"
	DOWNLOAD_PATH = DATA_PATH + "/download"
)

func getEnvWithDefault(envKey string) string {
	if envStr := genv.Get(envKey); envStr != nil {
		return envStr.String()
	}
	r, _ := g.Cfg().Get(gctx.GetInitCtx(), convertEnvToConfig(envKey))
	return r.String()
}

func convertEnvToConfig(env string) string {
	return strings.ToLower(strings.ReplaceAll(env, "_", "."))
}
