package file

import (
	"bytes"
	"context"
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
	"server-api/global/ecode"
	"strconv"
	"strings"
	"time"

	"github.com/gomooth/pkg/storage"
	"github.com/gomooth/utils/fsutil"

	"github.com/gomooth/xerror"
	"github.com/gomooth/xerror/xcode"
)

type service struct{}

func (s *service) UploadPublic(ctx context.Context, in *uploadRequest) (string, string, error) {
	if err := in.Validate(); nil != err {
		return "", "", xerror.NewXCode(xcode.RequestParamError, err.Error())
	}

	ext := strings.ToLower(path.Ext(in.File.Filename))
	name := strings.ToLower(fmt.Sprintf("%s%s", s.makeFilename(in.Genre), ext))

	now := time.Now()
	store := storage.Public()

	switch in.Business {
	case "articles", "icons", "banners":
		store = store.AppendDir(in.Business)
	default:
		return "", "", xerror.New("不支持的上传业务")
	}

	store = store.AppendDir(in.Genre, now.Format("2006"), now.Format("01")).
		SetName(name)

	dir, err := store.Dir()
	if err != nil {
		return "", "", xerror.WrapWithXCode(err, ecode.ErrorFileUpload)
	}
	_ = os.MkdirAll(dir, os.ModePerm)

	filePath, err := store.Path()
	if err != nil {
		return "", "", xerror.WrapWithXCode(err, ecode.ErrorFileUpload)
	}

	url, err := store.URLWithHost(global.Config.Server.HTTP.Resource.Host)
	if err != nil {
		return "", "", xerror.WrapWithXCode(err, ecode.ErrorFileUpload)
	}

	return filePath, url, nil
}

func (s *service) UploadPublicBase64(ctx context.Context, in *base64Request) (string, error) {
	if err := in.Validate(); nil != err {
		return "", xerror.NewXCode(xcode.RequestParamError, err.Error())
	}

	ext := strings.ToLower(path.Ext(in.Filename))
	name := strings.ToLower(fmt.Sprintf("%s%s", s.makeFilename(in.Genre), ext))

	now := time.Now()
	store := storage.Public()

	switch in.Business {
	case "articles", "icons", "banners":
		store = store.AppendDir(in.Business)
	default:
		return "", xerror.New("不支持的上传业务")
	}

	store = store.AppendDir(in.Genre, now.Format("2006"), now.Format("01")).
		SetName(name)

	dir, err := store.Dir()
	if err != nil {
		return "", xerror.WrapWithXCode(err, ecode.ErrorFileUpload)
	}
	_ = os.MkdirAll(dir, os.ModePerm)

	switch in.Genre {
	case "pictures":
		ib, err := base64.StdEncoding.DecodeString(in.Data)
		if nil != err {
			return "", xerror.WrapWithXCode(err, ecode.ErrorFileInvalid)
		}

		img, _, err := image.Decode(bytes.NewBuffer(ib))
		if nil != err {
			return "", xerror.WrapWithXCode(err, ecode.ErrorFileInvalid)
		}

		filePath, pathErr := store.Path()
		if pathErr != nil {
			return "", xerror.WrapWithXCode(pathErr, ecode.ErrorFileUpload)
		}
		out, err := os.Create(filePath)
		if nil != err {
			return "", xerror.WrapWithXCode(err, ecode.ErrorFileSave)
		}
		defer func() {
			_ = out.Close()
		}()

		if err := jpeg.Encode(out, img, &jpeg.Options{Quality: 100}); nil != err {
			return "", xerror.WrapWithXCode(err, ecode.ErrorFileSave)
		}
	}

	url, err := store.URLWithHost(global.Config.Server.HTTP.Resource.Host)
	if err != nil {
		return "", xerror.WrapWithXCode(err, ecode.ErrorFileUpload)
	}
	return url, nil
}

func (s *service) PrepareChunk(ctx context.Context, in *chunkRequest) (string, error) {
	if err := in.Validate(); nil != err {
		return "", xerror.NewXCode(xcode.RequestParamError, err.Error())
	}

	tmpStore := storage.Temp()

	switch in.Business {
	case "articles", "icons", "banners":
		tmpStore = tmpStore.AppendDir(in.Business)
	default:
		return "", xerror.New("不支持的上传业务")
	}

	chunkName := fmt.Sprintf("chunk-%d", in.ChunkNumber)
	tmpStore = tmpStore.AppendDir(in.Genre, strconv.Itoa(in.FileId)).SetName(chunkName)

	tmpDir, err := tmpStore.Dir()
	if err != nil {
		return "", xerror.WrapWithXCode(err, ecode.ErrorFileUpload)
	}
	_ = os.MkdirAll(tmpDir, os.ModePerm)

	tmpPath, err := tmpStore.Path()
	if err != nil {
		return "", xerror.WrapWithXCode(err, ecode.ErrorFileUpload)
	}

	return tmpPath, nil
}

func (s *service) FinalizeChunk(ctx context.Context, in *chunkRequest) (*chunkResponse, error) {
	res := &chunkResponse{
		Over: false,
		Url:  "",
	}

	idx := in.ChunkNumber
	total := in.ChunksTotal
	if idx != total-1 {
		return res, nil
	}

	ext := strings.ToLower(path.Ext(in.File.Filename))
	name := strings.ToLower(fmt.Sprintf("%s%s", s.makeFilename(in.Genre), ext))

	now := time.Now()
	store := storage.Public()
	tmpStore := storage.Temp()

	switch in.Business {
	case "articles", "icons", "banners":
		store = store.AppendDir(in.Business)
		tmpStore = tmpStore.AppendDir(in.Business)
	default:
		return nil, xerror.New("不支持的上传业务")
	}

	tmpStore = tmpStore.AppendDir(in.Genre, strconv.Itoa(in.FileId))
	tmpDir, err := tmpStore.Dir()
	if err != nil {
		return nil, xerror.WrapWithXCode(err, ecode.ErrorFileUpload)
	}

	store = store.AppendDir(in.Genre, now.Format("2006"), now.Format("01")).
		SetName(name)

	storeDir, dirErr := store.Dir()
	if dirErr != nil {
		return nil, xerror.WrapWithXCode(dirErr, ecode.ErrorFileUpload)
	}
	_ = os.MkdirAll(storeDir, os.ModePerm)

	storePath, pathErr := store.Path()
	if pathErr != nil {
		return nil, xerror.WrapWithXCode(pathErr, ecode.ErrorFileUpload)
	}
	fd, err := os.OpenFile(storePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, xerror.WrapWithXCode(err, ecode.ErrorFileSave)
	}
	defer func() {
		_ = fd.Close()
	}()

	chunkNameFormat := "chunk-%d"
	for i := 0; i < total; i++ {
		chunkName := fmt.Sprintf(chunkNameFormat, i)
		tmpFile := path.Join(tmpDir, chunkName)
		err := fsutil.BlockRead(tmpFile, func(data []byte) error {
			_, err := fd.Write(data)
			return err
		})
		if nil != err {
			return nil, xerror.WrapWithXCode(err, ecode.ErrorFileSave)
		}
	}

	_ = os.RemoveAll(tmpDir)

	res.Over = true
	url, urlErr := store.URLWithHost(global.Config.Server.HTTP.Resource.Host)
	if urlErr != nil {
		return nil, xerror.WrapWithXCode(urlErr, ecode.ErrorFileUpload)
	}
	res.Url = url

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
