---
paths:
  - "/app/http/**/*"
  - "/app/openapi/**/*"
---

# API 开发规则

## API 分层规则

```
internal/api/{端口}/{模块}
  controller.go   # 路由定义
  request.go      # 请求参数模型
  response.go     # 响应模型
  service.go      # 业务逻辑
```

- 严格按照 `controller -> service -> dao` 三层隔离
- `controller` 仅负责接收请求参数，封装响应，基础参数校验，以及流程流转
- `service` 负责业务逻辑实现
- 请求参数必须实现 `validator.IValidator` 接口

### 端口选择规则：HTTP vs OpenAPI

新接口必须先判断归属端口，避免业务实现错位：

| 判断维度 | HTTP 服务（`app/http`，端口 8000） | OpenAPI 服务（`app/openapi`，端口 8001） |
|---------|-----------------------------------|-----------------------------------------|
| **服务对象** | 自有产品的前端/内部系统 | 外部第三方应用（OpenAPP） |
| **认证方式** | JWT + 角色权限 | AppID + SHA1 签名 |
| **典型场景** | 管理后台、C 端用户操作 | 三方数据推送、开放能力对接 |

**选择原则：**
- 面向**登录用户**（人）的接口 → HTTP 服务
- 面向**第三方应用**（机器）的接口 → OpenAPI 服务
- 犹豫时问自己：调用方是"人登录后操作"还是"应用拿 AppID 签名后调用"？前者走 HTTP，后者走 OpenAPI

### 端口命名规则：
  - `admin` - ToB，管理后台
  - `user`  - ToC，用户端口
  - `health` - 健康检查

## 分层校验规则

- **简单类型校验**（格式、范围、枚举等纯语法/类型检查）在 `request` 和 `controller` 中完成
- **数据依赖校验**（需要通过 DAO 查库才能判定的校验）在 `service` 中完成
- `controller` 向 `service` 传递已解析的参数，`service` 不做格式解析
- `controller`，`service` 禁止有状态，应该通过 func args 显式传递
- 多端口共用的 `service` 逻辑应该提升到根 `service` 下

## Controller 规则

- `Controller` 是空结构体 `type Controller struct{}`，方法使用指针接收者 `(c *Controller)`
- 参数解析绑定方式：`GET` 用 `ShouldBindQuery`，`POST/PUT` 用 `ShouldBindJSON`
- 请求参数校验委托给 `Request` 的 `Validate()` 方法，并使用 `xerror.NewXCode(xcode.RequestParamError, err.Error())` 包装
- 路径参数用 `github.com/gomooth/utils` 包工具函数，如 `valutil.Int(ctx.Param("id"))` 等
- 使用 `new(service)` 实例化 `Service`，不用构造函数
- 推荐使用 `restful.NewResponse(ctx)` 处理响应
    - 使用 `restful.WithResponseErrorMsgHandler(lang.Handler())` 支持国际化错误提示
    - 响应方法：`WithError`(错误) / `WithMessage`(简单消息) / `Retrieve`(单个资源) / `Post`(创建/更新) / `ListWithPagination`(分页列表) / `Delete`(删除) 等
    - `/app/http` 下响应优先使用 `helper.NewResponse` 包装
- 禁止直接调用 DAO 层

## Request 规则

- Request 结构体私有（小写开头），并以 `函数名 + Request` 命名（如：`createTokenRequest`），
- GET 参数用 `form` tag，POST/PUT 参数用 `json` tag
- 列表请求嵌入可以直接 `httpmodel.SearchRequest`（含 Start/Limit/Sort 等）
- 业务校验写在 `Validate()` 方法中，用 `xerror.New("中文提示")`
- 修改请求可嵌入创建请求复用字段和校验

## Response 规则

- Response Entity 结构体私有（entity），通过 Controller 返回给 restful 框架
- 可通过嵌入 `httpmodel.ResponseModel` 提供 ID、CreatedAt 等通用字段
- `json` tag 用小驼峰命名（avatarUrl、currentRole）
- `copy` tag 对应 model 字段名，支持 `.String` 等方法调用形式
- `omitempty` 只用于可能为空的字段（如 profile）

## Service 规则

- `Service` 是空结构体 `type service struct{}`
- 方法第一个参数为 `context.Context`
- `Service` 不持有状态，每次通过 `new(service)` 实例化
- DAO 通过工厂函数获取：`pdao.NewUser()`、`pdao.NewVWUser()`
- DAO 错误必须用 `xerror.WrapWithXCode` 包装成业务错误码，禁止泄漏数据库细节
- 使用 `dbquery.NewQuery(filter, options...)` 构建查询选项
- `Model → Response Entity` 转换用 `xcopy.Copy`
- 禁止 `Service` 直接使用 `*gorm.DB`，必须通过 DAO 层

### Service 错误处理

Service 层**必须包装 DAO 层错误**，不对外暴露数据库细节，对用户提供友好提示。需要前端特殊处理的错误，在 `global/ecode/error_code.go` 中定义业务错误码：

```go
// global/ecode/error_code.go
var (
    // 请求相关 (3xxx)
    ErrorBadRequest   = xcode.NewWithMessage(3001, "请求数据错误或不存在")
    ErrorSavedData    = xcode.NewWithMessage(3004, "数据保存失败")
    ErrorRecordExist  = xcode.NewWithMessage(3005, "数据已存在")

    // 认证相关 (4xxx)
    ErrorAuthParams   = xcode.NewWithMessage(4000, "账号或密码错误")
    ErrorAuthFailed   = xcode.NewWithMessage(4001, "授权登录失败")
)
```

包装方式区分：

```go
// 包装 DAO 错误，隐藏数据库细节
record, err := pdao.NewUser().First(ctx, id)
if err != nil {
    return nil, xerror.WrapWithXCode(err, ecode.ErrorBadRequest)
}

// 纯业务逻辑错误，直接返回 ecode
return nil, xerror.WithXCode(ecode.ErrorAuthParams)

// 需要保留原始错误的场景（如日志排查），用 Wrap
return nil, xerror.WrapWithXCode(err, ecode.ErrorSavedData)
```

新增错误码在 `global/ecode/error_code.go` 的 `// NOTE:` 之后按分类递增（编号规则见 [coding-standard.md](coding-standard.md)）。

## 路由注册规则

- 业务接口路由在 `internal/route/` 中注册
- 每个端口子包暴露一个公开路由注册方法 `Register(router *gin.Engine)` 并在 `main.go` 中组合
- 内部各模块用私有函数 `registerXxx` 注册，不对外暴露
- API 端点命名严格遵循 RESTful 资源命名风格：
    - 资源名用复数名词（如：tokens, passwords, users, files）
    - 操作通过 HTTP Method + 资源名表达，不用动词
    - 路径参数用 `:id`、`:genre`、`:business` 等单数名词
