---
paths:
  - "/**/*_test.go"
---

# 测试约定

## 基本约定

- 所有测试统一使用 `testhelper/`，避免硬编码路径
  - 单元测试使用 `SetupTest()` 初始化：配置、日志（不连接 DB/Redis）
  - 集成测试使用 `SetupTestWithDB()` 初始化：配置、日志、数据库、Redis、缓存、锁
  - 自定义场景 `SetupTestWithOptions()` 按选项初始化
  - 初始化时，必须断言 error 为 nil
- 使用标准库 `testing` + `github.com/stretchr/testify` 断言 assert(常规) / require(致命)
- Mock 接口用 `testify/mock`，单元测试避免对 DB/HTTP 真实调用

### SetupOptions 详解

```go
opts := testhelper.DefaultSetupOptions()
opts.ConfigPath = ""       // 配置文件路径，空则自动发现
opts.WithDB = false        // 是否初始化数据库
opts.WithRedis = false     // 是否初始化 Redis
opts.WithCache = false     // 是否初始化缓存
opts.WithLocker = false    // 是否初始化分布式锁
opts.LogLevel = "debug"    // 日志级别
opts.Env = "test"          // 环境标识
testhelper.SetupTestWithOptions(opts)
```

注意：`WithDB=true` 时会自动设置 `WithRedis=true`、`WithCache=true`、`WithLocker=true`。

## 覆盖要求（每个公开函数）
- **正常路径**：典型合法输入
- **边界条件**：
  - 数值 → 0、±1、min、max、越界(n-1/n/n+1)
  - slice/string/map → nil、空、单元素、满长度/超长
  - 指针 → nil 与有效实例
- **错误路径**：非法入参、返回 error 的场景至少一例

## 辅助工具

```go
// 查找项目根目录
root, err := testhelper.FindProjectRoot()

// 获取测试配置文件路径
path, err := testhelper.GetTestConfigPath()

// CI 环境检测
if testhelper.IsCI() {
// CI 环境逻辑
}
```

## 运行测试

```bash
go test ./...                                      # 全部测试
go test ./testhelper/...                           # testhelper 包测试
go test -run TestFunctionName ./path/to/package    # 指定测试
go test -v ./path/to/package                       # 详细输出
```
