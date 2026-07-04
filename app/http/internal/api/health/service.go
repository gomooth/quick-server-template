package health

import (
	healthsvc "server-api/service/health"
)

// Service 健康检查服务（薄代理，实际逻辑在 service/health 包）
type Service = healthsvc.Service
