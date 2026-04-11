package ecode

import "github.com/gomooth/xerror/xcode"

var (
	InternalError      = xcode.NewWithMessage(1001, "系统内部错误")
	RequestParamError  = xcode.NewWithMessage(1002, "请求参数错误")
	OpenAPIClosed      = xcode.NewWithMessage(1003, "未被授权访问")
	RequestExpired     = xcode.NewWithMessage(1004, "请求已过期")
	MethodNotAllow     = xcode.NewWithMessage(1005, "请求方式不支持")
	MissRequired       = xcode.NewWithMessage(1006, "必填参数未设置")
	RequiredParamError = xcode.NewWithMessage(1007, "必填参数错误")
	SignError          = xcode.NewWithMessage(1008, "签名错误")
	BusinessError      = xcode.NewWithMessage(1009, "业务系统错误")
)
