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
		Status: OK,
		Build:  s.getBuildInfo(),
	}
}

// Readiness 检查所有依赖并返回就绪状态
func (s *Service) Readiness(ctx context.Context) ReadinessResponse {
	checks := map[string]any{
		"database": s.checkAllDBs(ctx),
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
// checks 中的值应为 CheckResult 或 GroupCheckResult，未知类型将被忽略
func determineOverallStatus(checks map[string]any) Status {
	for _, c := range checks {
		switch v := c.(type) {
		case CheckResult:
			if v.Status == StatusFail {
				return Fail
			}
		case GroupCheckResult:
			if v.Status == StatusFail {
				return Fail
			}
		}
	}
	return OK
}

// checkAllDBs 检查所有已注册数据库连接
func (s *Service) checkAllDBs(ctx context.Context) GroupCheckResult {
	if !global.Config.Data.Persistent.Enabled {
		return GroupCheckResult{Status: StatusSkip, Detail: "disabled"}
	}

	names := global.Database().List()
	if len(names) == 0 {
		return GroupCheckResult{Status: StatusSkip, Detail: "no connections registered"}
	}

	items := make(map[string]CheckResult, len(names))
	groupStatus := StatusOK

	for _, name := range names {
		items[name] = s.checkSingleDB(ctx, name)
		if items[name].Status == StatusFail {
			groupStatus = StatusFail
		}
	}

	return GroupCheckResult{Status: groupStatus, Items: items}
}

// checkSingleDB 检查单个数据库连接
func (s *Service) checkSingleDB(ctx context.Context, name string) CheckResult {
	db, err := global.Database().Get(name)
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
