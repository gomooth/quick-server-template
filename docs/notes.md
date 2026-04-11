# 注意事项

> 通用组件访问方式与禁用行为见 `docs/architecture.md`「全局组件访问」；Session 存储降级见 `rules/jwt.md`；健康检查端点见 `docs/deployment.md`。本文仅列项目特有的陷阱与警告。

## 安全警告

- **JWT 密钥**：`app.secret` 必须改为强随机字符串，不要用默认值 `go-quick-server-template`
- **示例文件**：生产环境设置 `app.clear_example_file=true` 清理示例代码
- **角色名称**：超级管理员角色为 `RoleSuper`（非 RoleAdmin）
- **OpenAPI 签名**：`appSecret` 存于 `open_apps` 表不可泄露；`sign-debugger` AppID 仅用于调试，生产环境确保不存在此记录；签名算法不加密传输内容，生产必须走 HTTPS

## 常见陷阱

- **Cache 禁用但代码使用**：DAO/Service 调用 `global.StringCacheManager()` 前未检查错误，导致运行时 panic
- **SessionStoreClient 为 nil**：Redis 未启用时用 JWT 有状态中间件会空指针（自动降级为内存存储，重启后 Session 丢失）
- **模式别名**：`cron`/`job`/`cronjob` 等效；`cmd`/`command` 等效
- **配置文件删除**：`clear_config_file=true` 会删除本地配置文件，Docker 挂载卷时需注意
- **热重载局限**：`watch_config_enabled` 开启后部分配置（如数据库连接）仍需重启生效

## 性能提醒

- 数据库连接池 `max_idle=10`、`max_open=100` 为默认值，高负载场景需要调整
- 避免 N+1 查询，使用 GORM Preload
- 热点数据用 `pcache` 层缓存，设置合理过期时间（默认 15 分钟，防止缓存永存）
- Job 防并发用 `global.Locker()`，锁键名格式 `jobTask:{业务域}:{动作}`
- OpenAPI 限流中间件 `IPRateLimit()` 按请求频率限制，保护后端

## 扩展点

- **新增数据库连接**：在 `data.persistent.connects` 加配置 → 在 `repository/{库名}/` 下建 `pmodel/pdao/pcache/pfilter/pattr` 子目录（前缀为库名首字母）→ `global.Database().Get("连接名")` 取连接
- **新增缓存/锁引擎**：配置引擎 → 在 `global/internal/cache/` 或 `global/internal/locker/` 下实现适配器
