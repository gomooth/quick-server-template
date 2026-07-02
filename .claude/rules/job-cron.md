---
paths:
  - "/app/job/**/*"
---

# 定时任务与命令约定

## Handler 组织

在 `app/job/internal/handler/` 下按业务域创建子目录：

```
handler/
└── example/       # 示例任务
    └── simple.go
```

## 开发约定

- 必须实现 `job.ICommandJob` 接口
  - 方法签名：`Run(args ...string) error`
- 定时任务在 `app/job/job.go` 中注册，新任务注册必须在 `// NOTE:` 标记之后。
- 需要防并发的定时任务必须使用 `global.Locker()` 加锁：
  - 锁键名格式：`jobTask:{业务域}:{动作}`
  - 获取锁失败时 `return nil`（跳过），不是返回错误
  - 用 `defer` 确保释放锁
- 在定时任务内部用 `context.Background()` 创建上下文

## 命令约定

- 一次性命令在 `app/job/command.go` 中注册，新命令注册必须在 `// NOTE:` 标记之后。
- 命令命名规则：用 `:` 分隔业务域和动作（如 `lang:reload`），简单命令用 `-` 分隔（如 `example-simple`）。
- 命令参数使用 `helper.NewCMDArgs` 解析，`ICMDArgsParser` 接口方法有：
  - `Get(key string, alias ...string) string` — 获取字符串值，支持别名
  - `GetInt(key string, alias ...string) int`
  - `GetBool(key string, alias ...string) bool`
