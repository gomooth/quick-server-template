package route

import (
	"server-api/global"
	"server-api/app/http/internal/api/file"

	"github.com/gin-gonic/gin"

	"github.com/gomooth/pkg/http/jwt"
	"github.com/gomooth/pkg/http/middleware"
)

func RegisterFile(rg *gin.Engine) {
	api := file.Controller{}

	v1 := rg.Group(
		"/file",
		middleware.RESTFul(global.ApiVersionLatest),
		middleware.JWTStatefulWithout(
			[]byte(global.Config.App.Secret),
			global.NewRole,
			jwt.WithSilentMode(true),
		),
		middleware.WithRole(global.RoleUser),
	)
	{
		v1.POST("/by-simple/:genre/:business", api.UploadPublic)
		v1.POST("/by-base64/:genre/:business", api.UploadPublicBase64)

		// 分块上传
		// 配合 npm install huge-uploader --save 食用
		// https://www.npmjs.com/package/huge-uploader
		v1.POST("/by-chunk/:genre/:business", api.UploadPublicChunk)
	}
}
