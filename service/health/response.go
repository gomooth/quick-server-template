package health

// CheckStatus 检查项状态枚举
type CheckStatus string

const (
	StatusOK   CheckStatus = "ok"
	StatusFail CheckStatus = "fail"
	StatusSkip CheckStatus = "skip"
)

// HealthStatus 整体健康状态
type HealthStatus string

const (
	HealthOK   HealthStatus = "ok"
	HealthFail HealthStatus = "fail"
)

// PongResponse /ping 端点响应（兼容旧客户端）
type PongResponse struct {
	Message string `json:"message"`
}

// LivenessResponse /healthz 端点响应
type LivenessResponse struct {
	Status HealthStatus `json:"status"`
	Build  BuildInfo    `json:"build"`
}

// ReadinessResponse /readyz 端点响应
type ReadinessResponse struct {
	Status HealthStatus           `json:"status"`
	Checks map[string]CheckResult `json:"checks"`
}

// BuildInfo 构建信息
type BuildInfo struct {
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
	GitCommit string `json:"git_commit"`
}

// CheckResult 单项检查结果
type CheckResult struct {
	Status CheckStatus `json:"status"`
	Detail string     `json:"detail,omitempty"`
	// system 专用字段
	Goroutines int     `json:"goroutines,omitempty"`
	MemoryMB   float64 `json:"memory_mb,omitempty"`
}
