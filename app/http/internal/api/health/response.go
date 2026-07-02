package health

// checkStatus 检查项状态枚举
type checkStatus string

const (
	statusOK   checkStatus = "ok"
	statusFail checkStatus = "fail"
	statusSkip checkStatus = "skip"
)

// healthStatus 整体健康状态
type healthStatus string

const (
	healthOK   healthStatus = "ok"
	healthFail healthStatus = "fail"
)

// pongResponse /ping 端点响应（兼容旧客户端）
type pongResponse struct {
	Message string `json:"message"`
}

// livenessResponse /healthz 端点响应
type livenessResponse struct {
	Status healthStatus `json:"status"`
	Build  buildInfo    `json:"build"`
}

// readinessResponse /readyz 端点响应
type readinessResponse struct {
	Status healthStatus           `json:"status"`
	Checks map[string]checkResult `json:"checks"`
}

// buildInfo 构建信息
type buildInfo struct {
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
	GitCommit string `json:"git_commit"`
}

// checkResult 单项检查结果
type checkResult struct {
	Status checkStatus `json:"status"`
	Detail string     `json:"detail,omitempty"`
	// system 专用字段
	Goroutines int     `json:"goroutines,omitempty"`
	MemoryMB   float64 `json:"memory_mb,omitempty"`
}
