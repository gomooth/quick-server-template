package file

import (
	"server-api/global/ecode"
	"server-api/service/lang"

	"github.com/gin-gonic/gin"
	"github.com/gomooth/pkg/http/restful"
	"github.com/gomooth/utils/valutil"
	"github.com/gomooth/xerror"
)

type Controller struct{}

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

	targetPath, url, err := new(service).UploadPublic(ctx.Request.Context(), &in)
	if nil != err {
		rru.WithError(err)
		return
	}

	if err := ctx.SaveUploadedFile(in.File, targetPath); nil != err {
		rru.WithError(xerror.WrapWithXCode(err, ecode.ErrorFileSave))
		return
	}

	rru.Retrieve(&uploadEntity{URL: url})
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
	url, err := new(service).UploadPublicBase64(ctx.Request.Context(), &in)
	if nil != err {
		rru.WithError(err)
		return
	}

	rru.Retrieve(&uploadEntity{URL: url})
}

func (c *Controller) UploadPublicChunk(ctx *gin.Context) {
	rru := restful.NewResponse(
		ctx,
		restful.WithResponseErrorMsgHandler(lang.Handler()),
	)

	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil {
		rru.WithError(xerror.New("分块文件太大"))
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

	svc := new(service)

	tempPath, err := svc.PrepareChunk(ctx.Request.Context(), &in)
	if nil != err {
		rru.WithError(err)
		return
	}

	if err := ctx.SaveUploadedFile(in.File, tempPath); nil != err {
		rru.WithError(xerror.WrapWithXCode(err, ecode.ErrorFileSave))
		return
	}

	res, err := svc.FinalizeChunk(ctx.Request.Context(), &in)
	if nil != err {
		rru.WithError(err)
		return
	}

	rru.Retrieve(res)
}
