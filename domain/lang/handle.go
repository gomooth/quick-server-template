package lang

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"server-api/global"

	"github.com/gin-gonic/gin"

	"github.com/gomooth/pkg/http/restful"

	"github.com/redis/go-redis/v9"

	"github.com/gomooth/xerror/xcode"

	"golang.org/x/sync/singleflight"

	"golang.org/x/text/language"
)

var single singleflight.Group

var supportedLanguages = make([]language.Tag, 0)

const (
	singleKey    = "lang"
	msgKeyFormat = "lang:%s:%d"
)

func getMsgKey(language language.Tag, code int) string {
	return fmt.Sprintf(msgKeyFormat, language.String(), code)
}

func Handler() ([]language.Tag, func(code int, lang language.Tag) string) {
	return supportedLanguages, langHandle()
}

func langHandle() func(code int, lang language.Tag) string {
	ctx := context.Background()

	return func(code int, lang language.Tag) string {
		v, err, _ := single.Do(singleKey, func() (interface{}, error) {
			key := getMsgKey(lang, code)
			cacheManager, err := global.StringCacheManager()
			if nil != err {
				slog.Warn("get Cache Manager failed", "lang", lang, "code", code, "err", err)
				return "", nil
			}
			return cacheManager.Get(ctx, key)
		})

		if nil != err {
			if !errors.Is(err, redis.Nil) {
				slog.Error("get lang failed", "lang", lang, "err", err)
			}
			return ""
		}

		return v.(string)
	}
}

// Content 获得语言包对应文字
// 请求头中必须包含 `X-Language` 才可以
func Content(ctx context.Context, code xcode.XCode) string {
	gtx, ok := ctx.(*gin.Context)
	if !ok {
		slog.Error("parse header failed in lang.Content")
		return ""
	}

	lang := detectLanguage(gtx)
	msg := langHandle()(code.Code(), lang)
	if len(msg) == 0 {
		return code.String()
	}
	return msg
}

func detectLanguage(ctx *gin.Context) language.Tag {
	supported := supportedLanguages
	if len(supported) == 0 {
		return language.Chinese
	}

	lang := ctx.GetHeader(restful.LangHeaderKey)
	if len(lang) == 0 {
		lang = ctx.GetHeader("Accept-Language")
	}
	if len(lang) == 0 {
		return language.Chinese
	}

	tags, _, err := language.ParseAcceptLanguage(lang)
	if err != nil || len(tags) == 0 {
		return language.Chinese
	}

	// 匹配支持的语言
	matcher := language.NewMatcher(supported)

	tag, _, _ := matcher.Match(tags...)
	return tag
}
