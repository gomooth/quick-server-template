package file

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	"math/rand/v2"
	"os"
	"path"
	"server-api/global"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gomooth/pkg/storage"
	"github.com/gomooth/utils/fsutil"

	"github.com/pkg/errors"

	"github.com/gomooth/xerror"
	"github.com/gomooth/xerror/xcode"
)

type service struct {
}

func (s *service) UploadPublic(ctx *gin.Context, in *uploadRequest) (string, error) {
	if err := in.Validate(); nil != err {
		return "", xerror.NewXCode(xcode.RequestParamError, err.Error())
	}

	// 生成随机文件名
	ext := strings.ToLower(path.Ext(in.File.Filename))
	name := strings.ToLower(fmt.Sprintf("%s%s", s.makeFilename(in.Genre), ext))

	now := time.Now()
	store := storage.Public()

	// 不同业务存储逻辑
	switch in.Business {
	case "articles", "icons", "banners":
		store = store.AppendDir(in.Business)
	default:
		return "", xerror.New("不支持的上传业务")
	}

	// 存储规则：业务/文件类型/年/月/文件
	store = store.AppendDir(in.Genre, now.Format("2006"), now.Format("01")).
		SetName(name)

	// 创建目录
	dir, err := store.Dir()
	if err != nil {
		return "", xerror.Wrap(err, "获取存储目录失败")
	}
	_ = os.MkdirAll(dir, os.ModePerm)

	filePath, err := store.Path()
	if err != nil {
		return "", xerror.Wrap(err, "获取存储路径失败")
	}
	if err := ctx.SaveUploadedFile(in.File, filePath); nil != err {
		return "", errors.Wrap(err, "文件保存失败")
	}

	url, err := store.URLWithHost(global.Config.Server.HTTP.Resource.Host)
	if err != nil {
		return "", xerror.Wrap(err, "获取文件URL失败")
	}
	return url, nil
}

func (s *service) UploadPublicBase64(_ *gin.Context, in *base64Request) (string, error) {
	if err := in.Validate(); nil != err {
		return "", xerror.NewXCode(xcode.RequestParamError, err.Error())
	}

	// 生成随机文件名
	ext := strings.ToLower(path.Ext(in.Filename))
	name := strings.ToLower(fmt.Sprintf("%s%s", s.makeFilename(in.Genre), ext))

	now := time.Now()
	store := storage.Public()

	// 不同业务存储逻辑
	switch in.Business {
	case "articles", "icons", "banners":
		store = store.AppendDir(in.Business)
	default:
		return "", xerror.New("不支持的上传业务")
	}

	// 存储规则：业务/文件类型/年/月/文件
	store = store.AppendDir(in.Genre, now.Format("2006"), now.Format("01")).
		SetName(name)

	// 创建目录
	dir, err := store.Dir()
	if err != nil {
		return "", xerror.Wrap(err, "获取存储目录失败")
	}
	_ = os.MkdirAll(dir, os.ModePerm)

	switch in.Genre {
	case "pictures":
		// 转成图片文件并把文件写入到 buffer
		ib, err := base64.StdEncoding.DecodeString(in.Data)
		if nil != err {
			return "", xerror.Wrap(err, "不是一个合法的图片")
		}

		img, _, err := image.Decode(bytes.NewBuffer(ib))
		if nil != err {
			return "", xerror.Wrap(err, "图片解析失败")
		}

		filePath, pathErr := store.Path()
		if pathErr != nil {
			return "", xerror.Wrap(pathErr, "获取存储路径失败")
		}
		out, err := os.Create(filePath)
		if nil != err {
			return "", xerror.Wrap(err, "写入文件失败")
		}
		defer func() {
			_ = out.Close()
		}()

		if err := jpeg.Encode(out, img, &jpeg.Options{Quality: 100}); nil != err {
			return "", xerror.Wrap(err, "写入文件失败")
		}
	}

	url, err := store.URLWithHost(global.Config.Server.HTTP.Resource.Host)
	if err != nil {
		return "", xerror.Wrap(err, "获取文件URL失败")
	}
	return url, nil
}

func (s *service) UploadPublicChunk(gtx *gin.Context, in *chunkRequest) (*chunkResponse, error) {
	if err := in.Validate(); nil != err {
		return nil, xerror.NewXCode(xcode.RequestParamError, err.Error())
	}

	now := time.Now()
	tmpStore := storage.Temp()
	store := storage.Public()

	// 不同业务存储逻辑
	switch in.Business {
	case "articles", "icons", "banners":
		tmpStore = tmpStore.AppendDir(in.Business)
		store = store.AppendDir(in.Business)
	default:
		return nil, xerror.New("不支持的上传业务")
	}

	// 分块文件临时存储规则：业务/文件类型/文件id/文件块
	chunkNameFormat := "chunk-%d"
	chunkName := fmt.Sprintf(chunkNameFormat, in.ChunkNumber)
	tmpStore = tmpStore.AppendDir(in.Genre, strconv.Itoa(in.FileId)).SetName(chunkName)

	// 创建目录
	tmpDir, err := tmpStore.Dir()
	if err != nil {
		return nil, xerror.Wrap(err, "获取临时目录失败")
	}
	_ = os.MkdirAll(tmpDir, os.ModePerm)

	tmpPath, err := tmpStore.Path()
	if err != nil {
		return nil, xerror.Wrap(err, "获取临时路径失败")
	}
	if err := gtx.SaveUploadedFile(in.File, tmpPath); nil != err {
		return nil, err
	}

	res := &chunkResponse{
		Over: false,
		Url:  "",
	}

	// 最后一个分块时，拼接数据
	idx := in.ChunkNumber
	total := in.ChunksTotal
	if idx == total-1 {
		// 生成随机文件名
		ext := strings.ToLower(path.Ext(in.File.Filename))
		name := strings.ToLower(fmt.Sprintf("%s%s", s.makeFilename(in.Genre), ext))

		// 存储子目录规则：业务/文件类型/年/月/文件
		store = store.AppendDir(in.Genre, now.Format("2006"), now.Format("01")).
			SetName(name)

		// 创建目录
		storeDir, dirErr := store.Dir()
		if dirErr != nil {
			return nil, xerror.Wrap(dirErr, "获取存储目录失败")
		}
		_ = os.MkdirAll(storeDir, os.ModePerm)

		storePath, pathErr := store.Path()
		if pathErr != nil {
			return nil, xerror.Wrap(pathErr, "获取存储路径失败")
		}
		fd, err := os.OpenFile(storePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return nil, xerror.Wrap(err, "存储文件时，打开失败")
		}
		defer func() {
			_ = fd.Close()
		}()

		for i := 0; i < total; i++ {
			chunkName := fmt.Sprintf(chunkNameFormat, i)
			tmpFile := path.Join(tmpDir, chunkName)
			err := fsutil.BlockRead(tmpFile, func(data []byte) error {
				_, err := fd.Write(data)
				return err
			})
			if nil != err {
				return nil, xerror.Wrap(err, "文件合并时，写入失败")
			}
		}

		// 删除临时目录
		_ = os.RemoveAll(tmpDir)

		res.Over = true
		url, urlErr := store.URLWithHost(global.Config.Server.HTTP.Resource.Host)
		if urlErr != nil {
			return nil, xerror.Wrap(urlErr, "获取文件URL失败")
		}
		res.Url = url
	}

	return res, nil
}

func (s *service) makeFilename(secret string) string {
	h := md5.New()
	h.Write([]byte(strconv.FormatInt(int64(rand.Int32()), 10)))
	h.Write([]byte("-"))
	h.Write([]byte(strconv.FormatInt(time.Now().UnixNano(), 10)))
	h.Write([]byte("-"))
	h.Write([]byte(strconv.FormatInt(rand.Int64(), 10)))

	name := hex.EncodeToString(h.Sum([]byte(secret)))

	r := strconv.FormatInt(rand.Int64(), 10)[0:6]

	return fmt.Sprintf("%s_%s", name, r)
}
