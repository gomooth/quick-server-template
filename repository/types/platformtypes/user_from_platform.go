package platformtypes

import "strings"

type UserFromPlatform uint

const (
	UserFromPlatformUndefined UserFromPlatform = iota
	UserFromPlatformAccount                    // 账号
	UserFromPlatformEmail                      // Email
	UserFromPlatformMobile                     // 手机号
	UserFromPlatformWechat                     // 微信
	UserFromPlatformAlipay                     // 支付宝
	UserFromPlatformWeibo                      // 新浪微博
	UserFromPlatformDouyin                     // 抖音
	UserFromPlatformFacebook                   // Facebook
	UserFromPlatformApple                      // Apple
	UserFromPlatformGoogle                     // 谷歌
	UserFromPlatformTwitter                    // Twitter
	UserFromPlatformTiktok                     // Tiktok
)

func NewUserFromPlatform(str string) UserFromPlatform {
	str = strings.TrimSpace(strings.ToLower(str))
	switch str {
	case "email":
		return UserFromPlatformEmail
	case "mobile":
		return UserFromPlatformMobile
	case "wechat", "wechat-mp":
		return UserFromPlatformWechat
	case "alipay", "alipay-mp":
		return UserFromPlatformAlipay
	default:
		return UserFromPlatformUndefined
	}
}
