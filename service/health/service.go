package health

import (
	"context"
	"runtime"
	"time"

	"server-api/global"
)

type Service struct{}

// Liveness 返回存活状态和构建信息
func (s *Service) Liveness() LivenessResponse {
	return LivenessResponse{
		Status: HealthOK,
		Build:  s.getBuildInfo(),
	}
}

// Readiness 检查所有依赖并返回就绪状态
func (s *Service) Readiness(ctx context.Context) ReadinessResponse {
	checks := map[string]CheckResult{
		"database": s.checkDB(ctx),
		"redis":    s.checkRedis(ctx),
		"system":   s.checkSystem(),
	}

	return ReadinessResponse{
		Status: determineOverallStatus(checks),
		Checks: checks,
	}
}

// Ping 兼容旧 /ping 端点
func (s *Service) Ping() PongResponse {
	return PongResponse{Message: "pong"}
}

// determineOverallStatus 根据各项检查结果确定整体健康状态
func determineOverallStatus(checks map[string]CheckResult) HealthStatus {
	for _, c := range checks {
		if c.Status == StatusFail {
			return HealthFail
		}
	}
	return HealthOK
}

// checkDB 检查数据库连接
func (s *Service) checkDB(ctx context.Context) CheckResult {
	if !global.Config.Data.Persistent.Enabled {
		return CheckResult{Status: StatusSkip, Detail: "disabled"}
	}

	db, err := global.Database().Get("platform")
	if err != nil {
		return CheckResult{Status: StatusFail, Detail: err.Error()}
	}

	checkCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlDB, err := db.DB()
	if err != nil {
		return CheckResult{Status: StatusFail, Detail: err.Error()}
	}

	if err := sqlDB.PingContext(checkCtx); err != nil {
		return CheckResult{Status: StatusFail, Detail: err.Error()}
	}

	return CheckResult{Status: StatusOK}
}

// checkRedis 检查 Redis 连接
func (s *Service) checkRedis(ctx context.Context) CheckResult {
	if !global.Config.Data.Redis.Enabled {
		return CheckResult{Status: StatusSkip, Detail: "disabled"}
	}

	client, err := global.RedisClient()
	if err != nil {
		return CheckResult{Status: StatusFail, Detail: err.Error()}
	}

	checkCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	pong, err := client.Ping(checkCtx).Result()
	if err != nil {
		return CheckResult{Status: StatusFail, Detail: err.Error()}
	}

	return CheckResult{Status: StatusOK, Detail: pong}
}

// checkSystem 检查系统资源（始终执行）
func (s *Service) checkSystem() CheckResult {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return CheckResult{
		Status:     StatusOK,
		Goroutines: runtime.NumGoroutine(),
		MemoryMB:   bytesToMB(m.Alloc),
	}
}

// getBuildInfo 返回构建信息，空值返回 "unknown"
func (s *Service) getBuildInfo() BuildInfo {
	return BuildInfo{
		Version:   unknownIfEmpty(global.BuildParam.Version),
		BuildTime: unknownIfEmpty(global.BuildParam.BuildTime),
		GitCommit: unknownIfEmpty(global.BuildParam.GitCommit),
	}
}

func unknownIfEmpty(v string) string {
	if v == "" {
		return "unknown"
	}
	return v
}

func bytesToMB(b uint64) float64 {
	return float64(b) / 1024 / 1024
}
