---
paths:
  - "/app/openapi/**/*"
---

# OpenAPI 开发规则

OpenAPI 服务（端口 8001）用于三方对接，认证体系与 HTTP 服务完全不同：使用 **AppID + 签名** 代替 JWT。端口选择规则（HTTP vs OpenAPI 对照表）见 [api-standard.md](api-standard.md)。

## OpenAPI 专属约定

| 维度 | 约定 |
|------|------|
| 中间件 | Auth / AuthWithoutSign |
| 错误码 | `1xxx`（定义于 `app/openapi/internal/helper/ecode/vars.go`，新增按分类递增加在 `// NOTE:` 之后） |
| 响应格式 | 直接使用 `restful.NewResponse`（不用 `helper.NewResponse` 包装） |
| 上下文解析 | `helper.ParseOpenAPP(ctx)` 从上下文解析 open app 信息 |
