package pattr

import "github.com/gomooth/xerror"

type OpenAPPState int8

const (
	OpenAPPStateForbidden OpenAPPState = iota - 2
	OpenAPPStateStopped
	OpenAPPStateNoOpened
	OpenAPPStateApplying
	OpenAPPStateNormal
)

var _openAPIStateTitleMap = map[OpenAPPState]string{
	OpenAPPStateForbidden: "禁用",
	OpenAPPStateStopped:   "停用",
	OpenAPPStateNoOpened:  "未开通",
	OpenAPPStateApplying:  "申请中",
	OpenAPPStateNormal:    "正常",
}

func (x OpenAPPState) String() string {
	if s, ok := _openAPIStateTitleMap[x]; ok {
		return s
	}
	return "未定义"
}

func ParseOpenAPPState(state int8) (OpenAPPState, error) {
	val := OpenAPPState(state)
	if _, ok := _openAPIStateTitleMap[val]; ok {
		return val, nil
	}
	return 0, xerror.New("不支持的 OpenAPP 状态")
}
