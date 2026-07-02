---
paths:
  - "/app/http/**/*"
---

# JWT 与角色约定

HTTP 服务（端口 8000）使用 **JWT + 有状态存储 + 角色** 的认证体系。OpenAPI 服务走 AppID + 签名，不用 JWT，见 [openapi.md](openapi.md)。

## 角色系统

角色定义在 `global/role.go`，类型为 `global.Role`（`uint8`）：

```go
global.RoleUser   // 用户（0）
global.RoleSuper  // 超级管理员（1）
```

角色转换（JWT 中间件解析 token 时需要 `NewRole`）：

```go
global.NewRole(str string) (httpcontext.IRole, error)           // 字符串 → Role
global.NewRoleFromGenre(genre uint8) (httpcontext.IRole, error) // uint8 → Role
global.Role2Genre(role httpcontext.IRole) uint8                 // Role → uint8
```

**注意**：超级管理员是 `RoleSuper`（非 RoleAdmin）。Admin 路由组用 `middleware.WithRole(global.RoleSuper)` 守卫，User 路由组用 `middleware.WithRole(global.RoleUser)`。

## Token 创建

用 `jwt.NewTokenBuilder` 构建有状态 token，存储后端为 Redis：

```go
store := jwtstore.NewMultiRedisStore(global.SessionStoreClient)
tk, err := jwt.NewTokenBuilder(
    []byte(global.Config.App.Secret),
    httpcontext.User{
        ID:      user.ID,
        Account: user.Account,
        Name:    user.Nickname,
        Roles:   roles,
    },
).
    WithIssuer(global.Config.App.ID).
    WithExpiration(24 * time.Hour).
    WithStatefulStore(store).
    Build()
if err != nil {
    return nil, xerror.WrapWithXCode(err, ecode.ErrorAuthFailed)
}

tokenStr, err := tk.ToString(ctx)
```

- 密钥来自 `global.Config.App.Secret`（生产环境必须改为强随机值）
- `WithStatefulStore` 使 token 可服务端注销；store 用 `jwtstore.NewMultiRedisStore(global.SessionStoreClient)`，依赖 Redis 启用
- token 字符串通过响应头返回：`rru.SetHeader(jwt.TokenHeaderKey, token.AccessToken)`

## 中间件选择

认证路由组用 `middleware.JWTStatefulWith` + `middleware.WithRole`：

```go
router.Group(
    "/admin/...",
    middleware.RESTFul(global.ApiVersionLatest),
    middleware.JWTStatefulWith(
        []byte(global.Config.App.Secret),
        global.NewRole,
        jwtstore.NewMultiRedisStore(global.SessionStoreClient),
    ),
    middleware.WithRole(global.RoleSuper),
)
```

- `JWTStatefulWith(secret, roleParser, store)` — 校验 token 并按 store 校验服务端状态（支持注销）
- `WithRole(role)` — 角色守卫，要求当前用户至少持有该角色
- 公开接口（如登录 `/admin/auth/tokens`）不加这两个中间件

## Session 存储

`helper.SessionStore(opt)` 根据是否启用 Redis 自动选择存储后端：

```go
// helper/jwt_option.go
func SessionStore(opt middleware.SessionOption) sessions.Store {
    if global.Config.Data.Redis.Enabled {
        return sessionRedisStore(opt)  // Redis 存储
    }
    return memstore.NewStore([]byte(global.Config.App.Secret))  // 内存降级
}
```

`SessionStoreClient` 使用 Redis DB=10，独立于主 Redis 连接。**Redis 未启用时降级为内存存储，重启后 Session 丢失**，生产环境务必启用 Redis。

## 自动续期

续期窗口由 `helper.JWTOption(refresh)` 控制：

```go
func JWTOption(refresh bool) *jwt.Option {
    refreshDuration := time.Duration(0)
    if refresh {
        refreshDuration = 12 * time.Hour
    }
    return jwt.NewOption(
        []byte(global.Config.App.Secret),
        global.NewRole,
        jwt.WithRefreshDuration(refreshDuration),
    )
}
```

`refresh=true` 时设置 12 小时续期窗口，`jwt.ParseJWTUser(c, jwtOpt)` 解析时据此判断是否续签（如 `middleware/http_cache.go` 中的用法）。

## 注销

注销通过清理服务端 store 实现，使已签发的 token 立即失效：

```go
func (s *service) Logout(ctx context.Context) error {
    owner := helper.UserFromContext(ctx)
    if owner.GetID() == 0 {
        return nil
    }
    if err := jwtstore.NewMultiRedisStore(global.SessionStoreClient).Clean(ctx, owner.GetID()); err != nil {
        return xerror.WrapWithXCode(err, ecode.ErrorHandleFailed)
    }
    return nil
}
```

## 当前用户解析

中间件注入后，用 `helper.UserFromContext(ctx)` 获取当前登录用户（`*httpcontext.User`）：

```go
owner := helper.UserFromContext(ctx)
if owner.GetID() == 0 {
    return xerror.NewXCode(xcode.Unauthorized)
}
```
