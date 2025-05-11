package platformtypes

import "github.com/save95/xerror"

// Gender 性别
type Gender int8

const (
	GenderUnknown Gender = iota // 保密
	GenderMale                  // 男性
	GenderFemale                // 女性
)

func (m Gender) String() string {
	switch m {
	case GenderUnknown:
		return "保密"
	case GenderMale:
		return "男"
	case GenderFemale:
		return "女"
	default:
		return "未定义"
	}
}

func MustUserGender(val int8) (Gender, error) {
	switch val {
	case int8(GenderUnknown):
		return GenderUnknown, nil
	case int8(GenderMale):
		return GenderMale, nil
	case int8(GenderFemale):
		return GenderFemale, nil
	default:
		return GenderUnknown, xerror.Errorf("未定义的性别[%d]", val)
	}
}
