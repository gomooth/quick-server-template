package ecode

import "github.com/gomooth/xerror/xcode"

var (
	ErrorBadRequest   = xcode.NewWithMessage(3001, "请求数据错误或不存在")
	ErrorVOConverted  = xcode.NewWithMessage(3002, "数据转换失败")
	ErrorRequestData  = xcode.NewWithMessage(3003, "请求数据错误")
	ErrorSavedData    = xcode.NewWithMessage(3004, "数据保存失败")
	ErrorRecordExist  = xcode.NewWithMessage(3005, "数据已存在")
	ErrorHandleFailed = xcode.NewWithMessage(3006, "操作失败")
	ErrorFileUpload     = xcode.NewWithMessage(3007, "文件上传失败")
	ErrorFileSave       = xcode.NewWithMessage(3008, "文件保存失败")
	ErrorFileInvalid    = xcode.NewWithMessage(3009, "文件格式不合法")
	ErrorPasswordFailed = xcode.NewWithMessage(3010, "密码处理失败")

	ErrorAuthParams      = xcode.NewWithMessage(4000, "账号或密码错误")
	ErrorAuthFailed      = xcode.NewWithMessage(4001, "授权登录失败")
	ErrorAuthUse2FA      = xcode.NewWithMessage(4002, "请输入2FA认证码")
	ErrorAuthCodeInvalid = xcode.NewWithMessage(4003, "验证码错误")
	ErrorAuthCodeExpired = xcode.NewWithMessage(4004, "验证码 已过期，请刷新验证码")

	// NOTE: 其他业务错误码

)
