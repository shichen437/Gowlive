package utils

import (
	"context"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	instance     *SingleI18n
	once         sync.Once
	assist_langs = g.MapStrStr{
		"en":    "en",
		"zh":    "zh-CN",
		"zh-cn": "zh-CN",
		"zh-tw": "zh-TW",
	}
)

type SingleI18n struct {
	m *gi18n.Manager
}

func getInstance() *SingleI18n {
	once.Do(func() {
		g.Log().Info(gctx.GetInitCtx(), "init i18n")
		instance = &SingleI18n{}
		i18n := g.I18n()
		lang := strings.ToLower(LANG)
		if lang == "" || assist_langs[lang] == "" {
			i18n.SetLanguage("zh-CN")
		} else {
			i18n.SetLanguage(assist_langs[lang])
		}
		instance.m = i18n
		g.Validator().I18n(instance.m)
	})
	return instance
}

func GetDefaultLang() string {
	getInstance()
	lang := strings.ToLower(LANG)
	if lang == "" || assist_langs[lang] == "" {
		return "zh-CN"
	} else {
		return assist_langs[lang]
	}
}

func T(ctx context.Context, key string) string {
	return getInstance().m.T(ctx, key)
}

func Tf(ctx context.Context, key string, format any) string {
	return getInstance().m.Tf(ctx, key, format)
}

func TError(ctx context.Context, key string) error {
	return gerror.New(T(ctx, key))
}

func TfError(ctx context.Context, key string, format any) error {
	return gerror.New(Tf(ctx, key, format))
}
