package helper

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"strings"
)

type SignResult struct {
	Path        string            `json:"path"`       // 请求路径
	Global      map[string]string `json:"global"`     // 公共参数
	QueryString map[string]string `json:"qs"`         // 查询参数
	Body        string            `json:"body"`       // body 参数
	SignedText  string            `json:"signedText"` // 最后组成的签名字符串
	Signature   string            `json:"signature"`  // 最终签名
	Input       string            `json:"input"`      // 输入的签名
	Success     bool              `json:"success"`    // 是否正确
}

func Sign(appID, appSecret, ts string, path string, qs map[string]string, body string) *SignResult {
	//ts := fmt.Sprintf("%d", time.Now().Unix())

	result := &SignResult{
		Global: map[string]string{
			"appId":      appID,
			"appSecret":  appSecret,
			"apiVersion": Version,
			"signType":   SignType,
			"timestamp":  ts,
		},
		Path:        path,
		QueryString: qs,
		Body:        body,
	}

	// 所有参数按 key 的字典排序
	keys := []string{"path"}
	params := map[string]string{
		"path": path,
	}
	for key, val := range result.Global {
		params[key] = val
		keys = append(keys, key)
	}
	if result.QueryString != nil {
		for key, val := range result.QueryString {
			if len(val) > 0 {
				params[key] = val
				keys = append(keys, key)
			}
		}
	}
	if len(body) > 0 {
		params["body"] = body
		keys = append(keys, "body")
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	values := make([]string, 0)
	for _, key := range keys {
		if val, ok := params[key]; ok {
			values = append(values, val)
		}
	}

	result.SignedText = strings.Join(values, "&")
	result.Signature = fmt.Sprintf("%x", sha1.Sum([]byte(result.SignedText)))

	return result
}
