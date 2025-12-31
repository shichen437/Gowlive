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
	E_PROJECT_LANG   = "PROJECT_LANG"

	E_OPENLIST_DOMAIN   = "OPENLIST_DOMAIN"
	E_OPENLIST_USERNAME = "OPENLIST_USERNAME"
	E_OPENLIST_PASSWORD = "OPENLIST_PASSWORD"
	E_OPENLIST_CODE     = "OPENLIST_CODE"
)

var (
	DATA_PATH     = getEnvWithDefault(E_PROJECT_DATA)
	SM4_KEY       = getEnvWithDefault(E_PROJECT_SM4KEY)
	LANG          = getEnvWithDefault(E_PROJECT_LANG)
	STREAM_PATH   = DATA_PATH + "/stream"
	DOWNLOAD_PATH = DATA_PATH + "/download"

	OPENLIST_DOMAIN   = getEnvWithDefault(E_OPENLIST_DOMAIN)
	OPENLIST_USERNAME = getEnvWithDefault(E_OPENLIST_USERNAME)
	OPENLIST_PASSWORD = getEnvWithDefault(E_OPENLIST_PASSWORD)
	OPENLIST_CODE     = getEnvWithDefault(E_OPENLIST_CODE)
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
