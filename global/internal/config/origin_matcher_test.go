package config

import "testing"

func TestOriginMatcher(t *testing.T) {
	patterns := []string{
		"https://admin.example.com",
		"https://*.example.com",
	}
	matcher := OriginMatcher(patterns)

	cases := []struct {
		origin string
		want   bool
	}{
		// 精确匹配
		{"https://admin.example.com", true},
		// 单层子域名通配
		{"https://api.example.com", true},
		{"https://www.example.com", true},
		// 多层子域名不匹配（避免误配扩大范围）
		{"https://a.b.example.com", false},
		// 协议不同不匹配
		{"http://admin.example.com", false},
		// 其他域名不匹配
		{"https://evil.com", false},
		{"https://evil.example.com.attacker.com", false},
		// 空字符串
		{"", false},
	}

	for _, c := range cases {
		if got := matcher(c.origin); got != c.want {
			t.Errorf("OriginMatcher(%q) = %v, want %v", c.origin, got, c.want)
		}
	}
}

func TestOriginMatcherEmpty(t *testing.T) {
	matcher := OriginMatcher(nil)
	if matcher("https://anywhere.com") {
		t.Error("empty patterns should reject all origins")
	}
}
