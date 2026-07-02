---
paths:
  - "/repository/*/*cache/*"
---

# DB 缓存层开发规则

- 缓存用 `dbcache.IDBCache[T, F]` + `getCacher()` 模式

## 基本结构

接口命名 `IXxxCache`，结构体仅 `name string` 字段：

```go
type user struct {
    name string
}

func NewUser() IUserCache {
    return &user{name: "user"}
}
```

`name` 即缓存命名空间/键前缀。

## getCacher 辅助方法

每次调用创建新实例（轻量设计，通过 `name` 隔离命名空间）：

```go
func (s *user) getCacher() (dbcache.IDBCache[pmodel.VWUser, pfilter.User], error) {
    cacheManager, err := global.StringCacheManager()
    if err != nil {
        return nil, err
    }
    return dbcache.New[pmodel.VWUser, pfilter.User](
        s.name, cacheManager, dbcache.WithExpiration(time.Hour),
    ), nil
}
```

## 穿透查询模式

### 按 ID 查询（标准模式）

回调委托 DAO 查询，缓存命中直接返回：

```go
return cacher.First(ctx, id, func(ctx context.Context) (*pmodel.VWUser, error) {
    return pdao.NewVWUser().First(ctx, id)
})
```

### 自定义键（推荐 RememberOf）

免手动序列化，类型安全：

```go
return dbcache.RememberOf[*pmodel.OpenAPP](ctx, cacher, key, func(ctx context.Context) (*pmodel.OpenAPP, error) {
    return pdao.NewOpenAPP().FirstByAppID(ctx, appID)
})
```

### 自定义键（Remember）

需手动 `json.Marshal`/`json.Unmarshal`，仅在需要自定义序列化时使用：

```go
result, err := cacher.Remember(ctx, key, func(ctx context.Context) ([]byte, error) {
    return json.Marshal(record)
})
var record pmodel.OpenAPP
json.Unmarshal(result, &record)
```

缓存键命名格式：`方法标识:业务参数`（如 `byAppID:xxx`、`byUserID:123`）

## 缓存清理

```go
cacher.Clear(ctx, dbcache.ClearWithAll(true))     // 清理全部
cacher.Clear(ctx, dbcache.ClearWithID(id))         // 按主键
cacher.Clear(ctx, dbcache.ClearWithKey(key))       // 按缓存键
cacher.Clear(ctx, dbcache.ClearWithTags(tag))      // 按标签
```

## 依赖方向

`pcache → pdao`（穿透时委托 DAO 查询），不反向依赖。获取缓存管理器：`global.StringCacheManager()`
