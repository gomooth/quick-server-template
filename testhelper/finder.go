package testhelper

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

// ConfigFinder 定义查找配置文件的接口
type ConfigFinder interface {
	Find() (string, error)
}

// EnvConfigFinder 从环境变量查找配置文件
type envConfigFinder struct{}

func (f *envConfigFinder) Find() (string, error) {
	if envPath := os.Getenv("TEST_CONFIG_PATH"); envPath != "" {
		if _, err := os.Stat(envPath); err == nil {
			return envPath, nil
		}
	}
	return "", os.ErrNotExist
}

// GoModConfigFinder 相对于go.mod（项目根目录）查找配置文件
type goModConfigFinder struct{}

func (f *goModConfigFinder) Find() (string, error) {
	projectRoot, err := findProjectRoot()
	if err != nil {
		return "", err
	}

	// 相对于项目根目录的常见配置位置
	configPaths := []string{
		filepath.Join(projectRoot, "storage/config/config.toml"),
		filepath.Join(projectRoot, "config/config.toml"),
		filepath.Join(projectRoot, "testdata/config.toml"),
	}

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", os.ErrNotExist
}

// FixedPathConfigFinder 尝试从测试文件的固定相对路径查找
type fixedPathConfigFinder struct{}

func (f *fixedPathConfigFinder) Find() (string, error) {
	// 获取调用者目录（测试文件位置）
	_, filename, _, ok := runtime.Caller(2) // 跳过findConfigFile和SetupTestWithOptions两层
	if !ok {
		return "", errors.New("无法获取调用者信息")
	}
	callerDir := filepath.Dir(filename)

	// 尝试从测试文件的常见相对路径
	searchPaths := []string{
		filepath.Join(callerDir, "../../../storage/config/config.toml"), // app/xxx/internal/xxx_test.go
		filepath.Join(callerDir, "../../storage/config/config.toml"),    // app/xxx/xxx_test.go
		filepath.Join(callerDir, "../config/config.toml"),               // 根目录测试文件
		filepath.Join(callerDir, "config/config.toml"),
	}

	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", os.ErrNotExist
}

// CompositeConfigFinder 按顺序尝试多个查找器
type compositeConfigFinder struct {
	finders []ConfigFinder
}

func (f *compositeConfigFinder) Find() (string, error) {
	for _, finder := range f.finders {
		if path, err := finder.Find(); err == nil {
			return path, nil
		}
	}
	return "", os.ErrNotExist
}

// NewDefaultConfigFinder 创建默认的配置查找器
// 搜索顺序：
// 1. 环境变量 (TEST_CONFIG_PATH)
// 2. 项目根目录 (基于go.mod)
// 3. 测试文件的固定相对路径
func NewDefaultConfigFinder() ConfigFinder {
	return &compositeConfigFinder{
		finders: []ConfigFinder{
			&envConfigFinder{},       // 最高优先级：显式环境变量
			&goModConfigFinder{},     // 基于项目根目录（最可靠）
			&fixedPathConfigFinder{}, // 回退方案：相对路径
		},
	}
}
