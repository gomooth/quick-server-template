package apptypes

import (
	"net/url"

	"golang.org/x/text/language"
)

// headerI18N 国际化相关参数
// 其值为标准的国际化标识，详细参考：[附：国际化语言标识](../docs/i18n.md)
type headerI18N struct {
	UseLanguage string `header:"X-Language" form:"lang"` // 使用语言，默认是英文
}

func (c headerI18N) Language() language.Tag {
	if len(c.UseLanguage) == 0 {
		return language.Chinese
	}

	return language.MustParse(c.UseLanguage)
}

// headerUTM 广告统计分析相关参数
// @see https://www.shangyexinzhi.com/article/4258360.html
// @see https://zhuanlan.zhihu.com/p/378091279
type headerUTM struct {
	UTMSourceEncode   string `header:"X-Utm-Source" form:"utm_source"`     // 广告投放来源
	UTMMediumEncode   string `header:"X-Utm-Medium" form:"utm_medium"`     // 广告投放媒介：cpc
	UTMCampaignEncode string `header:"X-Utm-Campaign" form:"utm_campaign"` // 广告投放名称
	UTMTermEncode     string `header:"X-Utm-Term" form:"utm_term"`         // 广告投放字词关键字
	UTMContentEncode  string `header:"X-Utm-Content" form:"utm_content"`   // 广告投放内容
}

func (c headerUTM) UTMSource() string {
	str, _ := url.QueryUnescape(c.UTMSourceEncode)
	return str
}

func (c headerUTM) UTMMedium() string {
	str, _ := url.QueryUnescape(c.UTMMediumEncode)
	return str
}

func (c headerUTM) UTMCampaign() string {
	str, _ := url.QueryUnescape(c.UTMCampaignEncode)
	return str
}

func (c headerUTM) UTMTerm() string {
	str, _ := url.QueryUnescape(c.UTMTermEncode)
	return str
}

func (c headerUTM) UTMContent() string {
	str, _ := url.QueryUnescape(c.UTMContentEncode)
	return str
}
