package health

import (
	"context"
	"testing"

	"server-api/global"
	"server-api/testhelper"

	"github.com/stretchr/testify/assert"
)

func TestGetBuildInfo(t *testing.T) {
	testhelper.SetupTest()
	orig := global.BuildParam
	defer func() { global.BuildParam = orig }()

	global.BuildParam = global.AppBuildParam{
		Version:   "v1.0",
		BuildTime: "2026-07-03",
		GitCommit: "abc1234",
	}

	svc := service{}
	info := svc.getBuildInfo()
	assert.Equal(t, "v1.0", info.Version)
	assert.Equal(t, "2026-07-03", info.BuildTime)
	assert.Equal(t, "abc1234", info.GitCommit)
}

func TestGetBuildInfo_EmptyFields(t *testing.T) {
	testhelper.SetupTest()
	orig := global.BuildParam
	defer func() { global.BuildParam = orig }()

	global.BuildParam = global.AppBuildParam{}

	svc := service{}
	info := svc.getBuildInfo()
	assert.Equal(t, "unknown", info.Version)
	assert.Equal(t, "unknown", info.BuildTime)
	assert.Equal(t, "unknown", info.GitCommit)
}

func TestCheckSystem_ReportsRuntimeMetrics(t *testing.T) {
	testhelper.SetupTest()

	svc := service{}
	result := svc.checkSystem()
	assert.Equal(t, statusOK, result.Status)
	assert.GreaterOrEqual(t, result.Goroutines, 1)
	assert.Less(t, result.Goroutines, 10000)
	assert.Greater(t, result.MemoryMB, float64(0))
}

func TestCheckDB_Disabled(t *testing.T) {
	testhelper.SetupTest()
	orig := global.Config.Data.Persistent.Enabled
	defer func() { global.Config.Data.Persistent.Enabled = orig }()

	global.Config.Data.Persistent.Enabled = false

	svc := service{}
	result := svc.checkDB(context.Background())
	assert.Equal(t, statusSkip, result.Status)
	assert.Equal(t, "disabled", result.Detail)
}

func TestCheckRedis_Disabled(t *testing.T) {
	testhelper.SetupTest()
	orig := global.Config.Data.Redis.Enabled
	defer func() { global.Config.Data.Redis.Enabled = orig }()

	global.Config.Data.Redis.Enabled = false

	svc := service{}
	result := svc.checkRedis(context.Background())
	assert.Equal(t, statusSkip, result.Status)
	assert.Equal(t, "disabled", result.Detail)
}

func TestLiveness(t *testing.T) {
	testhelper.SetupTest()
	orig := global.BuildParam
	defer func() { global.BuildParam = orig }()

	global.BuildParam = global.AppBuildParam{
		Version:   "v1.0",
		BuildTime: "2026-07-03",
		GitCommit: "abc1234",
	}

	svc := service{}
	resp := svc.liveness()
	assert.Equal(t, healthOK, resp.Status)
	assert.Equal(t, "v1.0", resp.Build.Version)
}

func TestReadiness_DBDisabled(t *testing.T) {
	testhelper.SetupTest()
	origDB := global.Config.Data.Persistent.Enabled
	origRedis := global.Config.Data.Redis.Enabled
	defer func() {
		global.Config.Data.Persistent.Enabled = origDB
		global.Config.Data.Redis.Enabled = origRedis
	}()

	global.Config.Data.Persistent.Enabled = false
	global.Config.Data.Redis.Enabled = false

	svc := service{}
	resp := svc.readiness(context.Background())
	assert.Equal(t, healthOK, resp.Status) // system always ok
	assert.Equal(t, statusSkip, resp.Checks["database"].Status)
	assert.Equal(t, statusSkip, resp.Checks["redis"].Status)
	assert.Equal(t, statusOK, resp.Checks["system"].Status)
}

func TestDetermineOverallStatus_AllOK(t *testing.T) {
	checks := map[string]checkResult{
		"database": {Status: statusOK},
		"redis":    {Status: statusOK},
		"system":   {Status: statusOK},
	}
	assert.Equal(t, healthOK, determineOverallStatus(checks))
}

func TestDetermineOverallStatus_WithFail(t *testing.T) {
	checks := map[string]checkResult{
		"database": {Status: statusFail, Detail: "connection refused"},
		"redis":    {Status: statusSkip, Detail: "disabled"},
		"system":   {Status: statusOK},
	}
	assert.Equal(t, healthFail, determineOverallStatus(checks))
}

func TestDetermineOverallStatus_AllSkip(t *testing.T) {
	checks := map[string]checkResult{
		"database": {Status: statusSkip, Detail: "disabled"},
		"redis":    {Status: statusSkip, Detail: "disabled"},
		"system":   {Status: statusOK},
	}
	assert.Equal(t, healthOK, determineOverallStatus(checks))
}

func TestPing(t *testing.T) {
	svc := service{}
	resp := svc.ping()
	assert.Equal(t, "pong", resp.Message)
}
