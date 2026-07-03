package admin

import (
	"server-api/app/http/internal/api/admin/file"

	"github.com/gin-gonic/gin"
)

func registerFile(router gin.IRouter) {
	api := file.Controller{}

	ra := router.Group("/files")
	{
		ra.POST("/by-simple/:genre/:business", api.UploadPublic)
		ra.POST("/by-base64/:genre/:business", api.UploadPublicBase64)

		// 分块上传
		// 配合 npm install huge-uploader --save 食用
		// https://www.npmjs.com/package/huge-uploader
		ra.POST("/by-chunk/:genre/:business", api.UploadPublicChunk)
	}
}
