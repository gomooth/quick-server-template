package config

import "strings"

// OriginMatcher 根据预编译的 origin 模式列表构造匹配器。
//
// 模式形如 "https://admin.example.com"（精确匹配）或
// "https://*.example.com"（单层子域名通配，"*" 仅匹配一个不含点的标签）。
// 协议（scheme）必须完全一致，host 区分大小写按原样比较。
//
// 返回的匹配器在 patterns 为空时拒绝所有 origin。
func OriginMatcher(patterns []string) func(string) bool {
	compiled := make([]originPattern, 0, len(patterns))
	for _, p := range patterns {
		compiled = append(compiled, compileOriginPattern(p))
	}

	return func(origin string) bool {
		if origin == "" {
			return false
		}
		scheme, host := splitOrigin(origin)
		if scheme == "" || host == "" {
			return false
		}
		for _, pat := range compiled {
			if pat.match(scheme, host) {
				return true
			}
		}
		return false
	}
}

// originPattern 是单条 origin 模式编译后的结构。
type originPattern struct {
	scheme string
	host   string
	wild   bool // host 是否以 "*." 开头（单层通配）
}

func compileOriginPattern(pattern string) originPattern {
	scheme, host := splitOrigin(pattern)
	wild := strings.HasPrefix(host, "*.")
	if wild {
		host = host[1:] // 去掉 '*'，保留 ".example.com"
	}
	return originPattern{scheme: scheme, host: host, wild: wild}
}

func (p originPattern) match(scheme, host string) bool {
	if p.scheme != scheme {
		return false
	}
	if p.wild {
		// host 必须以 ".example.com" 结尾，且前缀恰为一个不含点的标签
		if !strings.HasSuffix(host, p.host) {
			return false
		}
		prefix := strings.TrimSuffix(host, p.host)
		return prefix != "" && !strings.Contains(prefix, ".")
	}
	return p.host == host
}

// splitOrigin 将 "scheme://host[:port]" 拆分为 scheme 与 host。
// 不含 "://" 时返回空串，视为非法 origin。
func splitOrigin(origin string) (scheme, host string) {
	idx := strings.Index(origin, "://")
	if idx < 0 {
		return "", ""
	}
	scheme = origin[:idx]
	rest := origin[idx+3:]
	if scheme == "" || rest == "" {
		return "", ""
	}
	// 去掉 path/query，并分离端口；本匹配仅关心 host
	if i := strings.IndexAny(rest, "/?#"); i >= 0 {
		rest = rest[:i]
	}
	// 暂不处理 userinfo；标准 origin 头不含 userinfo
	host = rest
	return scheme, host
}
