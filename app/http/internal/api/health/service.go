package health

import (
	"context"
	"runtime"
	"time"

	"server-api/global"
)

type service struct{}

// liveness 返回存活状态和构建信息
func (s *service) liveness() livenessResponse {
	return livenessResponse{
		Status: healthOK,
		Build:  s.getBuildInfo(),
	}
}

// readiness 检查所有依赖并返回就绪状态
func (s *service) readiness(ctx context.Context) readinessResponse {
	checks := map[string]checkResult{
		"database": s.checkDB(ctx),
		"redis":    s.checkRedis(ctx),
		"system":   s.checkSystem(),
	}

	return readinessResponse{
		Status: determineOverallStatus(checks),
		Checks: checks,
	}
}

// determineOverallStatus 根据各项检查结果确定整体健康状态
func determineOverallStatus(checks map[string]checkResult) healthStatus {
	for _, c := range checks {
		if c.Status == statusFail {
			return healthFail
		}
	}
	return healthOK
}

// ping 兼容旧 /ping 端点
func (s *service) ping() pongResponse {
	return pongResponse{Message: "pong"}
}

// checkDB 检查数据库连接
func (s *service) checkDB(ctx context.Context) checkResult {
	if !global.Config.Data.Persistent.Enabled {
		return checkResult{Status: statusSkip, Detail: "disabled"}
	}

	db, err := global.Database().Get("platform")
	if err != nil {
		return checkResult{Status: statusFail, Detail: err.Error()}
	}

	checkCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlDB, err := db.DB()
	if err != nil {
		return checkResult{Status: statusFail, Detail: err.Error()}
	}

	if err := sqlDB.PingContext(checkCtx); err != nil {
		return checkResult{Status: statusFail, Detail: err.Error()}
	}

	return checkResult{Status: statusOK}
}

// checkRedis 检查 Redis 连接
func (s *service) checkRedis(ctx context.Context) checkResult {
	if !global.Config.Data.Redis.Enabled {
		return checkResult{Status: statusSkip, Detail: "disabled"}
	}

	client, err := global.RedisClient()
	if err != nil {
		return checkResult{Status: statusFail, Detail: err.Error()}
	}

	checkCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	pong, err := client.Ping(checkCtx).Result()
	if err != nil {
		return checkResult{Status: statusFail, Detail: err.Error()}
	}

	return checkResult{Status: statusOK, Detail: pong}
}

// checkSystem 检查系统资源（始终执行）
func (s *service) checkSystem() checkResult {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return checkResult{
		Status:     statusOK,
		Goroutines: runtime.NumGoroutine(),
		MemoryMB:   bytesToMB(m.Alloc),
	}
}

// getBuildInfo 返回构建信息，空值返回 "unknown"
func (s *service) getBuildInfo() buildInfo {
	return buildInfo{
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
