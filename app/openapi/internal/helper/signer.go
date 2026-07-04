package helper

import (
	"bytes"
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
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

// GetRequestPath 从上下文中提取请求路径
func GetRequestPath(ctx context.Context) string {
	gtx, ok := ctx.(*gin.Context)
	if !ok {
		return ""
	}
	return gtx.Request.URL.Path
}

// ExtractRequestParams 从上下文中提取 query 参数和 body，并回写 body 供下游读取
func ExtractRequestParams(ctx context.Context) (qs map[string]string, body string) {
	gtx, ok := ctx.(*gin.Context)
	if !ok {
		return nil, ""
	}

	// 提取 query 参数
	qs = make(map[string]string)
	for key := range gtx.Request.URL.Query() {
		v := gtx.Request.URL.Query().Get(key)
		if len(v) > 0 {
			qs[key] = v
		}
	}

	// 读取 body 并回写
	bodyRaw, _ := io.ReadAll(gtx.Request.Body)
	gtx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyRaw))
	body = strings.TrimSpace(string(bodyRaw))

	return qs, body
}
