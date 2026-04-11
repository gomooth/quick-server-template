package lang

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"server-api/global"
	"server-api/repository/platform/pdao"
	"server-api/repository/platform/pfilter"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	"github.com/gomooth/pkg/framework/dbquery"

	"github.com/gomooth/pkg/storage"
	"github.com/gomooth/utils/sliceutil"
	"github.com/gomooth/utils/valutil"

	"github.com/gomooth/xerror"

	"golang.org/x/text/language"
)

func Init(ctx context.Context) error {
	if !global.Config.Data.Cache.Enabled {
		slog.Warn("lang init skip, because cache is disabled")
		return nil
	}

	cacheManager, err := global.StringCacheManager()
	if nil != err {
		return xerror.Wrap(err, "get cache manager failed")
	}
	// 清空旧缓存
	_ = cacheManager.Invalidate(ctx, store.WithInvalidateTags([]string{"lang"}))

	// 先缓存本地文件
	if err := initForLocal(ctx, cacheManager); nil != err {
		return err
	}

	// 再使用 db 覆盖
	if err := initForMysql(ctx, cacheManager); nil != err {
		return err
	}

	// 去重
	supportedLanguages = sliceutil.Unique(supportedLanguages)

	return nil
}

func initForMysql(ctx context.Context, cacheManager *cache.Cache[string]) error {
	if !global.Config.Data.Persistent.Enabled {
		slog.Debug("lang init for mysql, but database disabled, skip")
		return nil
	}

	start := 0
	limit := 100
	hasMore := true

	repo := pdao.NewLang()
	supportedSet := map[language.Tag]struct{}{}
	for hasMore {
		records, err := repo.List(ctx, dbquery.NewQuery(pfilter.Lang{}, dbquery.WithOffsetPage[pfilter.Lang](start, limit)))
		if nil != err {
			return err
		}
		if len(records) == 0 {
			hasMore = false
			continue
		}

		start += limit
		hasMore = len(records) >= limit

		for _, record := range records {
			msg := record.Content
			tag, err := language.Parse(record.Lang)
			if err != nil {
				tag = language.English
			}

			supportedSet[tag] = struct{}{}

			key := getMsgKey(tag, record.Code)

			if err := cacheManager.Set(ctx, key, msg,
				store.WithTags([]string{"lang", tag.String()}),
				store.WithExpiration(0), // 永不过期
			); err != nil {
				//global.Log.Errorf("lang set failed: lang=%s, code=%d, msg=%s", language, k, msg)
				return xerror.Wrapf(err, "lang set failed: lang=%s, code=%d, msg=%s", tag, record.Code, msg)
			}
		}
	}
	for tag := range supportedSet {
		supportedLanguages = append(supportedLanguages, tag)
	}

	return nil
}

func initForLocal(ctx context.Context, cacheManager *cache.Cache[string]) error {
	langPath := storage.Disk("langs")
	langDir, err := langPath.Path()
	if err != nil {
		return xerror.Wrap(err, "get lang path failed")
	}
	files, err := filepath.Glob(langDir + "/*.json")
	if err != nil {
		return err
	}

	supportedSet := map[language.Tag]struct{}{}
	for _, file := range files {
		lang := filepath.Base(file[:len(file)-5]) // 移除 .json
		tag, err := language.Parse(lang)
		if err != nil {
			continue
		}

		supportedSet[tag] = struct{}{}

		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		var contents map[string]string
		if err := json.Unmarshal(data, &contents); err != nil {
			return err
		}

		for codeStr, msg := range contents {
			code := valutil.Int(codeStr)
			key := getMsgKey(tag, code)
			if err := cacheManager.Set(ctx, key, msg,
				store.WithTags([]string{"lang", tag.String()}),
				store.WithExpiration(0), // 永不过期
			); err != nil {
				//global.Log.Errorf("lang set failed: lang=%s, code=%d, msg=%s", language, k, msg)
				return xerror.Wrapf(err, "lang set failed: lang=%s, code=%d, msg=%s", tag, code, msg)
			}
		}
	}
	for tag := range supportedSet {
		supportedLanguages = append(supportedLanguages, tag)
	}

	return nil
}
