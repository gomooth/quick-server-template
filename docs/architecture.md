# 架构说明

## 目录结构

```
├── app/                   # 应用层（服务实现）
│   ├── http/              # HTTP服务（端口8000，toB/toC管理后台）
│   │   ├── internal/
│   │   │   ├── api/            # API控制器、服务层、请求/响应
│   │   │   │   ├── admin/      #   管理端（ToB）：auth, user, file
│   │   │   │   ├── user/       #   用户端（ToC）：auth
│   │   │   │   └── health/     #   健康检查：Ping, Liveness, Readiness
│   │   │   ├── middleware/     # 中间件（CORS、XSS过滤、HTTP缓存等）
│   │   │   ├── route/          # 路由注册
│   │   │   └── helper/         # 辅助函数（JWT配置、用户解析）
│   │   ├── middleware.go       # 中间件包装
│   │   └── main.go             # HTTP服务入口
│   ├── openapi/           # OpenAPI服务（三方对接API，端口8001）
│   │   ├── internal/
│   │   │   ├── api/            # API实现（ping, system）
│   │   │   ├── middleware/     # 中间件（签名认证、限流、缓存）
│   │   │   ├── route/          # 路由注册
│   │   │   └── helper/         # 辅助工具（签名、Header解析、错误码）
│   │   ├── middleware.go       # 中间件包装
│   │   └── main.go             # OpenAPI服务入口
│   ├── job/               # 定时任务服务
│   │   ├── internal/handler/   # 任务处理器（按业务域分子目录）
│   │   ├── job.go              # Cron 注册
│   │   └── command.go          # CMD 注册
│   └── consumer/          # 消息消费者服务
│       ├── internal/consumer/  # 消费者实现（按业务域分子目录）
│       ├── internal/helper/    # 辅助工具（Kafka Producer）
│       ├── kafka.go            # Kafka 消费者注册
│       ├── single.go           # 单消费者注册（Redis/HTTPSQS）
│       └── main.go             # 消费者服务入口
├── boot/                  # 启动引导框架
│   ├── internal/          # 各服务启动器
│   ├── app.go             # 应用初始化（Boot + initialize）
│   ├── config.go          # 启动参数结构（Param, ServerType）
│   └── command.go         # 命令模式实现
├── testhelper/            # 测试辅助工具包
│   ├── setup.go           # 测试初始化逻辑
│   └── finder.go          # 配置文件查找器
├── global/                # 全局组件（基础设施层，公开包）
│   ├── internal/          # 内部实现
│   │   ├── database/      # 数据库连接和迁移
│   │   ├── config/        # 配置解析和管理
│   │   ├── logger/        # 日志系统实现
│   │   └── locker/        # 分布式锁实现
│   ├── ecode/             # 业务错误码定义
│   ├── lang/              # 语言包消息码
│   ├── kafka/             # Kafka 常量（consumer group, topic）
│   ├── config.go          # 全局配置变量
│   ├── database.go        # 数据库/Redis/缓存/锁初始化
│   ├── env.go             # 环境判断（local/test/prod）
│   ├── logger.go          # 日志初始化
│   ├── role.go            # 角色权限定义（RoleUser/RoleSuper）
│   ├── vars.go            # 构建参数
│   └── var_const.go       # API 版本常量
├── repository/            # 数据访问层（公共业务模型，公开包）
│   ├── platform/          # 平台db相关
│   │   ├── pattr/         # 业务数据属性类型
│   │   ├── pcache/        # 缓存层实现
│   │   ├── pdao/          # 数据访问对象实现
│   │   ├── pfilter/       # 数据过滤器类型
│   │   └── pmodel/        # 表实体模型
│   └── shared/            # 共享类型
├── service/               # 公共业务服务层（公开包）
│   └── lang/              # 语言包初始化和 i18n 处理
├── storage/               # 存储目录（公共资源）
│   ├── config/            # 配置文件
│   ├── langs/             # 多语言文件
│   ├── logs/              # 日志文件目录
│   └── public/            # 静态文件目录
├── main.go               # 程序入口
├── go.mod                # Go模块定义
├── Makefile              # 构建脚本
└── README.md             # 项目说明文档
```

## 依赖方向

```
app/ → boot/ → global/ → repository/ → service/
```

- 上层可以导入下层
- 下层不应导入上层
- 同级包可以相互导入

## 核心组件

1. **配置管理** (`global/internal/config/`) — TOML格式，多源配置，热重载
2. **数据库** (`global/internal/database/`) — 多数据库支持，GORM集成，自动迁移
3. **日志系统** (`global/logger.go`) — 结构化日志，slog 适配，文件轮转
4. **认证授权** — JWT认证（HTTP），AppID+签名认证（OpenAPI），角色权限（RBAC）
5. **缓存系统** — Redis 缓存，分布式锁
6. **消息队列** (`app/consumer/`) — Kafka、Redis、HTTPSQS
7. **定时任务** (`app/job/`) — Cron表达式，分布式调度

## 启动流程

### 正常启动（Boot 模式）

```
main.go → step() [解析命令行参数]
  → boot.Boot(cnf)
    → initialize(cnf)
      → global.ParseConfig()     # 1. 加载配置
      → global.InitLogger()      # 2. 初始化日志（传入 logCategory）
      → global.InitDataBase()    # 3. 初始化基础设施（DB → Locker → Redis → Cache）
    → global.MigrateData()       # 4. 数据库迁移
    → lang.Init(ctx)             # 5. 初始化语言包
    → app.NewManager()           # 6. 创建应用管理器
    → 注册配置文件监听器（可选）
    → 注册服务到 Manager
    → apps.Run(ctx)              # 7. 启动所有服务（阻塞）
```

### 命令模式（Command 模式）

```
main.go → boot.Command(cnf)
  → initialize(cnf)             # 同上
  → jobapp.CMDRegister(cmd)     # 注册所有命令
  → cmd.Execute(name, args...)  # 执行指定命令
```

### 初始化顺序

`global.InitDataBase()` 内部按以下顺序初始化，有依赖关系：

```
initDB()      → 注册数据库连接
initLocker()  → 分布式锁（可独立于 Redis）
initRedis()   → Redis 连接 + SessionStoreClient（DB=10）
initCache()   → 缓存管理器（依赖 Redis）
```

每个组件都检查 `enabled` 配置，禁用时跳过并打印日志。

## 模式切换

通过 `-mode` 参数控制启动哪些服务：

| mode 值 | 注册的服务 | 日志分类 |
|---------|-----------|---------|
| `all` | Web + OpenAPI + Cronjob + Consumer | `all` |
| `web` | HTTP 服务 | `default` |
| `openapi` | OpenAPI 服务 | `openapi` |
| `cron` / `job` / `cronjob` | 定时任务 | `cronjob` |
| `consumer` | 消息消费者 | `consumer` |
| `cmd` / `command` | 一次性命令 | `command` |

## 服务生命周期

每个服务实现 `app.IServer` 接口：

```
New(ctx) → 创建服务实例
  → 启动 HTTP/Cron/Consumer 等
  → 运行中（阻塞）
  → 收到信号 → Release() 释放资源
```

各服务的 `Release()` 函数用于清理资源（关闭连接、停止消费者等），当前模板中大部分为空实现，按需填充。

资源释放分两层：

- app 层 `Release()`/`Shutdown()`：释放 app 私有资源（如 consumer 的 producer、各 app 独占连接），由 `app.Manager` 在收到信号后逆序调用。
- `global.Release`：释放全局共享基础设施（DB 连接池、Redis、Cache、Locker），LIFO 顺序，由 `boot.Boot` 在 `apps.Run` 返回后统一调用。注册通过 `global.RegisterRelease`，在各 init 阶段完成。

注意：`Release` 执行后再调用 `RegisterRelease` 会 panic（生命周期倒置）；`Release` 本身幂等，重复调用返回 nil。

## 全局组件访问

| 组件 | 访问方式 | 禁用时行为 |
|------|---------|-----------|
| 配置 | `global.Config` | — |
| 日志 | `global.Log`（`*slog.Logger`） | — |
| 环境 | `global.Env()` | — |
| 数据库 | `global.Database().Get("platform")` | 返回未注册错误 |
| Redis | `global.RedisClient()` | 返回 `("redis client is disabled" error)` |
| 分布式锁 | `global.Locker()` | 返回 `("locker is disabled" error)` |
| 缓存 | `global.StringCacheManager()` | 返回 `("cache manager is disabled" error)` |
| Session | `global.SessionStoreClient` | `nil`（需 Redis 启用） |
| 构建参数 | `global.BuildParam` | — |
| 全局释放钩子注册 | `global.RegisterRelease(fn)` | Release 后调用将 panic |
| 全局资源释放 | `global.Release(ctx)` | 幂等，重复调用返回 nil |
