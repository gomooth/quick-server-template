# 测试辅助包

`testhelper` 包提供测试初始化的工具函数，解决了重复初始化代码和硬编码配置路径的问题。

## Features

1. **Automatic configuration file discovery** - Searches for `config.toml` in common locations
2. **Unified initialization** - Single function call for all test setup
3. **Flexible options** - Customize initialization with `SetupOptions`
4. **CI awareness** - Detect CI environment and skip tests if needed
5. **Environment variable support** - Use `TEST_CONFIG_PATH` to specify config location

## Usage

### Basic Setup (Unit Tests)

For unit tests that don't need database access:

```go
import "server-api/testhelper"

func TestMain(m *testing.M) {
    if err := testhelper.SetupTest(); err != nil {
        panic(err)
    }
    
    code := m.Run()
    testhelper.Cleanup()
    os.Exit(code)
}
```

### Setup with Database (Integration Tests)

For integration tests that need database access:

```go
func TestMain(m *testing.M) {
    if err := testhelper.SetupTestWithDB(); err != nil {
        panic(err)
    }
    
    code := m.Run()
    testhelper.Cleanup()
    os.Exit(code)
}
```

### Custom Setup Options

For advanced configuration:

```go
opts := testhelper.SetupOptions{
    ConfigPath: "/path/to/config.toml", // Optional, auto-discovered if empty
    WithDB:     true,                   // Initialize database
    WithRedis:  true,                   // Initialize Redis
    LogLevel:   "info",                 // Log level
    Env:        "test",                 // Environment
}

if err := testhelper.SetupTestWithOptions(opts); err != nil {
    t.Fatal(err)
}
```

### Skipping Tests in CI

For tests that require local resources:

```go
func TestRequiresLocalDB(t *testing.T) {
    testhelper.SkipIfCI(t) // Skips test in CI environment
    
    // Test code that needs local database
}
```

## Configuration File Discovery

The package uses a smart, multi-strategy approach to find configuration files:

### Search Strategy (in order):

1. **Environment Variable** (`TEST_CONFIG_PATH`) - Highest priority
   - Explicit path specified by user
   - Useful for CI/CD or custom setups

2. **Project Root Based** (go.mod location) - Most reliable
   - Finds `go.mod` file to determine project root
   - Searches from project root:
     - `storage/config/config.toml`
     - `config/config.toml`
     - `testdata/config.toml`

3. **Relative Path Fallback** - For edge cases
   - Tries common relative paths from test file location
   - Handles deeply nested test files

### Why This Approach is Better:

1. **No Enumeration Problems** - Doesn't rely on hardcoded path lists
2. **Works Anywhere** - Test files can be in any subdirectory
3. **Project-Agnostic** - Works across different project structures
4. **Predictable** - Always finds config relative to project root
5. **Override Support** - Environment variable for special cases

## Environment Variables

- `TEST_CONFIG_PATH`: Override configuration file path
- `CI`, `GITHUB_ACTIONS`, `GITLAB_CI`: Detect CI environment

## Best Practices

1. **Use `TestMain` for package-level setup** - Avoid duplicate `init()` functions
2. **Separate unit and integration tests** - Use `SetupTest()` for unit tests, `SetupTestWithDB()` for integration tests
3. **Skip resource-intensive tests in CI** - Use `SkipIfCI(t)` for tests that need local resources
4. **Set `TEST_CONFIG_PATH` in CI** - Ensure consistent config location

## Migration from Old Pattern

**Before:**
```go
func init() {
    // Hardcoded path
    if err := global.ParseConfig("../../config/config.toml"); err != nil {
        log.Fatal(err)
    }
    // More initialization...
}
```

**After:**
```go
func TestMain(m *testing.M) {
    if err := testhelper.SetupTestWithDB(); err != nil {
        panic(err)
    }
    code := m.Run()
    testhelper.Cleanup()
    os.Exit(code)
}
```

## File Structure

```
testhelper/
├── setup.go          # 主要测试初始化逻辑
├── finder.go         # 配置文件查找器接口和实现
├── example_test.go   # 使用示例
└── README.md         # 本文档
```

### 模块职责分离

- **setup.go**: 包含测试初始化、选项配置、路径修复等核心逻辑
- **finder.go**: 专门处理配置文件查找，包含多种查找策略的实现
- **example_test.go**: 提供使用示例和最佳实践

这种分离使得代码更加模块化，便于维护和测试。

## Example

See `example_test.go` for complete usage examples.