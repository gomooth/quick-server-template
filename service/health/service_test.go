package health

import (
	"context"
	"testing"

	"server-api/global"
	"server-api/testhelper"

	"github.com/stretchr/testify/assert"
)

func TestGetBuildInfo(t *testing.T) {
	err := testhelper.SetupTest()
	assert.Nil(t, err)
	orig := global.BuildParam
	defer func() { global.BuildParam = orig }()

	global.BuildParam = global.AppBuildParam{
		Version:   "v1.0",
		BuildTime: "2026-07-03",
		GitCommit: "abc1234",
	}

	svc := Service{}
	info := svc.getBuildInfo()
	assert.Equal(t, "v1.0", info.Version)
	assert.Equal(t, "2026-07-03", info.BuildTime)
	assert.Equal(t, "abc1234", info.GitCommit)
}

func TestGetBuildInfo_EmptyFields(t *testing.T) {
	err := testhelper.SetupTest()
	assert.Nil(t, err)
	orig := global.BuildParam
	defer func() { global.BuildParam = orig }()

	global.BuildParam = global.AppBuildParam{}

	svc := Service{}
	info := svc.getBuildInfo()
	assert.Equal(t, "unknown", info.Version)
	assert.Equal(t, "unknown", info.BuildTime)
	assert.Equal(t, "unknown", info.GitCommit)
}

func TestCheckSystem_ReportsRuntimeMetrics(t *testing.T) {
	err := testhelper.SetupTest()
	assert.Nil(t, err)

	svc := Service{}
	result := svc.checkSystem()
	assert.Equal(t, StatusOK, result.Status)
	assert.GreaterOrEqual(t, result.Goroutines, 1)
	assert.Less(t, result.Goroutines, 10000)
	assert.Greater(t, result.MemoryMB, float64(0))
}

func TestCheckDB_Disabled(t *testing.T) {
	err := testhelper.SetupTest()
	assert.Nil(t, err)
	orig := global.Config.Data.Persistent.Enabled
	defer func() { global.Config.Data.Persistent.Enabled = orig }()

	global.Config.Data.Persistent.Enabled = false

	svc := Service{}
	result := svc.checkDB(context.Background())
	assert.Equal(t, StatusSkip, result.Status)
	assert.Equal(t, "disabled", result.Detail)
}

func TestCheckRedis_Disabled(t *testing.T) {
	err := testhelper.SetupTest()
	assert.Nil(t, err)
	orig := global.Config.Data.Redis.Enabled
	defer func() { global.Config.Data.Redis.Enabled = orig }()

	global.Config.Data.Redis.Enabled = false

	svc := Service{}
	result := svc.checkRedis(context.Background())
	assert.Equal(t, StatusSkip, result.Status)
	assert.Equal(t, "disabled", result.Detail)
}

func TestLiveness(t *testing.T) {
	err := testhelper.SetupTest()
	assert.Nil(t, err)
	orig := global.BuildParam
	defer func() { global.BuildParam = orig }()

	global.BuildParam = global.AppBuildParam{
		Version:   "v1.0",
		BuildTime: "2026-07-03",
		GitCommit: "abc1234",
	}

	svc := Service{}
	resp := svc.Liveness()
	assert.Equal(t, HealthOK, resp.Status)
	assert.Equal(t, "v1.0", resp.Build.Version)
}

func TestReadiness_DBDisabled(t *testing.T) {
	err := testhelper.SetupTest()
	assert.Nil(t, err)
	origDB := global.Config.Data.Persistent.Enabled
	origRedis := global.Config.Data.Redis.Enabled
	defer func() {
		global.Config.Data.Persistent.Enabled = origDB
		global.Config.Data.Redis.Enabled = origRedis
	}()

	global.Config.Data.Persistent.Enabled = false
	global.Config.Data.Redis.Enabled = false

	svc := Service{}
	resp := svc.Readiness(context.Background())
	assert.Equal(t, HealthOK, resp.Status) // system always ok
	assert.Equal(t, StatusSkip, resp.Checks["database"].Status)
	assert.Equal(t, StatusSkip, resp.Checks["redis"].Status)
	assert.Equal(t, StatusOK, resp.Checks["system"].Status)
}

func TestDetermineOverallStatus_AllOK(t *testing.T) {
	checks := map[string]CheckResult{
		"database": {Status: StatusOK},
		"redis":    {Status: StatusOK},
		"system":   {Status: StatusOK},
	}
	assert.Equal(t, HealthOK, determineOverallStatus(checks))
}

func TestDetermineOverallStatus_WithFail(t *testing.T) {
	checks := map[string]CheckResult{
		"database": {Status: StatusFail, Detail: "connection refused"},
		"redis":    {Status: StatusSkip, Detail: "disabled"},
		"system":   {Status: StatusOK},
	}
	assert.Equal(t, HealthFail, determineOverallStatus(checks))
}

func TestDetermineOverallStatus_AllSkip(t *testing.T) {
	checks := map[string]CheckResult{
		"database": {Status: StatusSkip, Detail: "disabled"},
		"redis":    {Status: StatusSkip, Detail: "disabled"},
		"system":   {Status: StatusOK},
	}
	assert.Equal(t, HealthOK, determineOverallStatus(checks))
}

func TestPing(t *testing.T) {
	svc := Service{}
	resp := svc.Ping()
	assert.Equal(t, "pong", resp.Message)
}
