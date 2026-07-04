package testhelper

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"server-api/global"
)

// SetupOptions 定义测试初始化的选项
type SetupOptions struct {
	// ConfigPath 指定配置文件路径，如果为空则自动发现
	ConfigPath string

	// WithDB 是否初始化数据库连接
	WithDB bool

	// WithRedis 是否初始化Redis连接
	WithRedis bool

	// WithCache 是否初始化缓存
	WithCache bool

	// WithLocker 是否初始化分布式锁
	WithLocker bool

	// LogLevel 测试日志级别，默认 "debug"
	LogLevel string

	// Env 测试环境，默认 "test"
	Env string
}

// DefaultSetupOptions 返回默认的测试初始化选项
func DefaultSetupOptions() SetupOptions {
	return SetupOptions{
		WithDB:     false, // 单元测试默认不连接数据库
		WithRedis:  false,
		WithCache:  false,
		WithLocker: false,
		LogLevel:   "debug",
		Env:        "test",
	}
}

// SetupTest 使用默认选项初始化测试环境
// 自动发现配置文件并初始化全局配置和日志
func SetupTest() error {
	return SetupTestWithOptions(DefaultSetupOptions())
}

// SetupTestWithDB 初始化测试环境并包含数据库支持
// 适用于需要数据库访问的集成测试
func SetupTestWithDB() error {
	opts := DefaultSetupOptions()
	opts.WithDB = true
	opts.WithRedis = true
	opts.WithCache = true
	opts.WithLocker = true
	return SetupTestWithOptions(opts)
}

// SetupTestWithOptions 使用自定义选项初始化测试环境
func SetupTestWithOptions(opts SetupOptions) error {
	// 1. 查找配置文件
	configPath := opts.ConfigPath
	if configPath == "" {
		var err error
		configPath, err = findConfigFile()
		if err != nil {
			return err
		}
	}

	// 2. 解析配置
	if err := global.ParseConfig(configPath); err != nil {
		return err
	}

	// 3. 覆盖环境设置（如果指定）
	if opts.Env != "" {
		global.Config.App.Env = opts.Env
	}

	// 4. 设置测试日志级别
	if opts.LogLevel != "" {
		global.Config.App.Log.Level = opts.LogLevel
	}

	// 5. 修复配置中的相对路径，使其相对于项目根目录
	if err := fixRelativePaths(); err != nil {
		return err
	}

	// 6. 初始化日志，使用固定的 "test" 作为日志目录
	if err := global.InitLogger("test"); err != nil {
		return err
	}

	// 7. 初始化数据库（如果请求）
	if opts.WithDB {
		if err := global.InitDataBase(); err != nil {
			return err
		}
	}

	return nil
}

// FindProjectRoot 通过定位go.mod文件查找项目根目录
func FindProjectRoot() (string, error) {
	// 从当前目录开始
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// 向上遍历目录树查找go.mod
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir, nil
		}

		// 移动到父目录
		parent := filepath.Dir(dir)
		if parent == dir {
			// 到达根目录
			break
		}
		dir = parent
	}

	return "", errors.New("go.mod not found (not a Go module?)")
}

// findProjectRoot 是FindProjectRoot的内部版本
func findProjectRoot() (string, error) {
	return FindProjectRoot()
}

// fixRelativePaths 确保配置中的所有相对路径都是相对于项目根目录的绝对路径
func fixRelativePaths() error {
	projectRoot, err := findProjectRoot()
	if err != nil {
		return err
	}

	// 修复日志目录
	if !filepath.IsAbs(global.Config.App.Log.Dir) {
		global.Config.App.Log.Dir = filepath.Join(projectRoot, global.Config.App.Log.Dir)
		// 确保目录存在
		if err := os.MkdirAll(global.Config.App.Log.Dir, 0755); err != nil {
			return fmt.Errorf("failed to create log directory %s: %w", global.Config.App.Log.Dir, err)
		}
	}

	// 修复静态资源路径（如果存在）
	if global.Config.Server.HTTP.Resource.Path != "" && !filepath.IsAbs(global.Config.Server.HTTP.Resource.Path) {
		global.Config.Server.HTTP.Resource.Path = filepath.Join(projectRoot, global.Config.Server.HTTP.Resource.Path)
	}

	// 注意：当前版本的OpenAPI配置没有Resource字段
	// 如果将来需要，在此添加类似的修复

	return nil
}

// findConfigFile 使用默认查找器查找配置文件
func findConfigFile() (string, error) {
	finder := NewDefaultConfigFinder()
	return finder.Find()
}

// Cleanup 在测试后执行清理操作
// 应该在TestMain中的m.Run()之后调用
func Cleanup() {
	// 释放已注册的全局基础设施资源（Redis/Cache/DB/Producer 等），LIFO 顺序。
	// 测试环境记录错误但不影响清理流程。
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := global.Release(ctx); err != nil {
		// 测试环境记录但不影响清理流程
		_ = err
	}
	log.Println("test cleanup completed")
}

// GetTestConfigPath 返回测试中使用的配置文件路径
// 这对于需要知道配置位置的测试很有用
func GetTestConfigPath() (string, error) {
	return findConfigFile()
}

// IsCI 检查是否在CI环境中运行
func IsCI() bool {
	return os.Getenv("CI") == "true" ||
		os.Getenv("GITHUB_ACTIONS") == "true" ||
		os.Getenv("GITLAB_CI") == "true"
}

// SkipIfCI 如果在CI环境中运行则跳过测试
// 用于需要本地资源的测试
func SkipIfCI(t testingT) {
	if IsCI() {
		t.Skip("Skipping test in CI environment")
	}
}

// testingT 是testing.T的最小接口，避免直接依赖
type testingT interface {
	Skip(args ...interface{})
	Skipf(format string, args ...interface{})
}
