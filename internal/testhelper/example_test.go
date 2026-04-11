package testhelper_test

import (
	"os"
	"testing"

	"server-api/internal/testhelper"
)

func TestExampleSetup(t *testing.T) {
	// 简单设置，不包含数据库
	if err := testhelper.SetupTest(); err != nil {
		t.Fatalf("SetupTest失败: %v", err)
	}
}

func TestExampleSetupWithDB(t *testing.T) {
	// 如果数据库不可用，在CI中跳过
	testhelper.SkipIfCI(t)

	// 为集成测试设置包含数据库的环境
	if err := testhelper.SetupTestWithDB(); err != nil {
		t.Fatalf("SetupTestWithDB失败: %v", err)
	}
}

func TestExampleCustomOptions(t *testing.T) {
	// 使用特定选项的自定义设置
	opts := testhelper.SetupOptions{
		ConfigPath: os.Getenv("TEST_CONFIG_PATH"), // 如果设置了环境变量则使用
		WithDB:     false,
		LogLevel:   "info",  // 设置日志级别为info
		Env:        "local", // 设置环境为local
	}

	if err := testhelper.SetupTestWithOptions(opts); err != nil {
		t.Fatalf("SetupTestWithOptions失败: %v", err)
	}
}

func TestMain(m *testing.M) {
	// 全局测试设置
	if err := testhelper.SetupTest(); err != nil {
		panic(err)
	}

	// 运行测试
	code := m.Run()

	// 全局测试清理
	testhelper.Cleanup()

	os.Exit(code)
}
