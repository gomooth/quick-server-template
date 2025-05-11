package file

import (
	"server-api/service/lang"

	"github.com/gomooth/pkg/http/restful"
	"github.com/gomooth/utils/valutil"

	"github.com/gin-gonic/gin"

	"github.com/pkg/errors"
)

type Controller struct {
}

func (c *Controller) UploadPublic(ctx *gin.Context) {
	rru := restful.NewResponse(
		ctx,
		restful.WithResponseErrorMsgHandler(lang.Handler()),
	)

	file, _ := ctx.FormFile("file")

	in := uploadRequest{
		Genre:    ctx.Param("genre"),
		Business: ctx.Param("business"),
		File:     file,
	}

	url, err := new(service).UploadPublic(ctx, &in)
	if nil != err {
		rru.WithError(err)
		return
	}

	rru.Retrieve(map[string]string{
		"url": url,
	})
}

func (c *Controller) UploadPublicBase64(ctx *gin.Context) {
	rru := restful.NewResponse(
		ctx,
		restful.WithResponseErrorMsgHandler(lang.Handler()),
	)

	var in base64Request
	if err := ctx.ShouldBindJSON(&in); nil != err {
		rru.WithError(err)
		return
	}

	in.Genre = ctx.Param("genre")
	in.Business = ctx.Param("business")
	url, err := new(service).UploadPublicBase64(ctx, &in)
	if nil != err {
		rru.WithError(err)
		return
	}

	rru.Retrieve(map[string]string{
		"url": url,
	})
}

func (c *Controller) UploadPublicChunk(ctx *gin.Context) {
	rru := restful.NewResponse(
		ctx,
		restful.WithResponseErrorMsgHandler(lang.Handler()),
	)

	err := ctx.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		rru.WithError(errors.New("分块文件太大"))
		return
	}

	file, _ := ctx.FormFile("file")

	in := chunkRequest{
		Genre:       ctx.Param("genre"),
		Business:    ctx.Param("business"),
		Params:      nil,
		File:        file,
		FileId:      valutil.Int(ctx.GetHeader("uploader-file-id")),
		ChunksTotal: valutil.Int(ctx.GetHeader("uploader-chunks-total")),
		ChunkNumber: valutil.Int(ctx.GetHeader("uploader-chunk-number")),
	}

	res, err := new(service).UploadPublicChunk(ctx, &in)
	if nil != err {
		rru.WithError(err)
		return
	}

	rru.Retrieve(res)
}
