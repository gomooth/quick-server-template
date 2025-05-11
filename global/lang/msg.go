package lang

import "github.com/save95/xerror/xcode"

var (
	TestContent = xcode.NewWithMessage(100000, "测试语言包内容")

	None      = xcode.NewWithMessage(100001, "")
	Unknown   = xcode.NewWithMessage(100002, "未知")
	Undefined = xcode.NewWithMessage(100003, "未定义")
	Valid     = xcode.NewWithMessage(100004, "有效")
	Invalid   = xcode.NewWithMessage(100005, "无效")
	Expired   = xcode.NewWithMessage(100006, "已过期")

	// NOTE: 其他业务语言码

)
