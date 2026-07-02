---
paths:
  - "/app/**/*"
  - "/repository/**/*"
  - "/service/**/*"
---

# 通用编码约定

跨切面的核心约定，适用于应用层、repository、service。各路径作用域的具体规则见同目录其他文件。

## 工厂函数模式

`NewXxx()` 工厂函数**不接受** `ctx` 参数，上下文通过各方法参数传递。接口公开、实现私有，工厂返回接口类型（`func NewUser() IUser`）。

## 全局组件

配置、日志、数据库、Redis、分布式锁、缓存等全局组件的访问方式和禁用行为见 `docs/architecture.md`「全局组件访问」。所有组件（Redis/Locker/Cache）未启用时访问返回错误，使用前务必检查。

## 错误码编号规则

错误码按模块/类型分段，新增错误码在对应文件 `// NOTE:` 之后按分类递增：

| 范围 | 定义位置 | 用途 |
|------|---------|------|
| `1xxx` | `app/openapi/internal/helper/ecode/vars.go` | OpenAPI 错误（签名/参数/权限） |
| `3xxx` | `global/ecode/error_code.go` | 请求/数据错误 |
| `4xxx` | `global/ecode/error_code.go` | 认证/授权错误 |
| `100xxx` | `global/lang/msg.go` | 语言包消息码 |

## 语言包与国际化

翻译文件按语言标识（如 `zh_CN`、`en_US`）放在 `storage/langs/`，错误 KEY 保持字典顺序；支持的语言标识以该目录下实际文件为准。

消息码用于国际化，编号从 `100000` 开始，定义在 `global/lang/msg.go` 的 `// NOTE:` 之后递增。启用方式：

- 通用：`restful.NewResponse(ctx, restful.WithResponseErrorMsgHandler(lang.Handler()))`
- `app/http` 下可直接 `helper.NewResponse(ctx, true)`

## 设计模式

- **分层架构**：`Controller → Service → Repository(DAO + Cache) → Database`，依赖单向 — 详见 `docs/architecture.md`
- **私有结构体 + 公开工厂函数**：`type user struct{}` + `func NewUser() IUser`
