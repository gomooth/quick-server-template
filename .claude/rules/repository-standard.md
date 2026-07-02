---
paths:
  - "/repository/**/*"
---

# Repository 层开发规则

## 目录结构

```
repository/
├── platform/              # 按数据库名隔离，每个数据库一个顶层目录
│   ├── pmodel/            # 表实体和视图实体
│   ├── pdao/              # 数据访问对象
│   ├── pcache/            # 缓存层
│   ├── pattr/             # 业务属性类型（枚举）
│   └── pfilter/           # 列表查询过滤器
└── shared/                # 跨库共享类型
```

子目录前缀为库名首字母（`platform` → `p`）。新增数据库时创建对应前缀的目录组（如 `analytics` → `amodel/`、`adao/`、`acache/`、`aattr/`、`afilter/`）。

## 数据库设计约定

1. 表名均为实体复数形式（如 `users`, `langs`, `failed_jobs`）
2. 每个表除自增 id 外，必须按顺序包含 `created_at, updated_at, deleted_at` 三个字段，放在字段末尾
3. 禁止外键约束，引用外部表 id 做软约束（仅存 id，不建 FK）
4. 所有表必须考虑索引建立和查询索引优化
5. 表示是否等二义性的字段类型为 `tinyint(1) unsigned NOT NULL`
6. 表示枚举的字段类型为 `tinyint(2) NOT NULL`

## pmodel — 实体模型

1. 每张表一个文件，文件名 `snake_case`，表名的单数形式（如 `表user_login_logs` -> `user_login_log`）
2. 嵌入 `gorm.Model`；视图实体嵌入基础实体（而非 `gorm.Model`）
3. 列名不同时标注 `gorm:"column:xxx"`，表名不同时实现 `TableName()`
4. 枚举字段用 `pattr` 业务类型（如 `Gender`、`OpenAPPState`），不用原生 `int8`/`string`
5. 可选时间用 `*time.Time`；模型可含业务方法（如 `HasPassword()`），但不含 DAO 操作

## pattr — 业务属性类型

1. 每个业务概念一个文件
2. 自定义类型 + `iota` 常量，负值用 `iota - N` 偏移
3. 必须实现 `String()` 返回中文文本，内部 map 命名 `_xxxTitleMap`
4. 必须提供 `ParseXxx(val int8) (Xxx, error)` 安全解析；字符串解析提供 `NewXxx(str) Xxx`

## pfilter — 过滤器

1. 零值表示"不筛选"；仅当零值本身有查询意义时用指针（如 `*int`、`*pattr.OpenAPPState`）
2. 批量 ID 用 `IDs []uint`；模糊匹配字段名加 `Like` 后缀
3. 必须提供数据库表中索引字段的过滤

## pdao — DAO

- DAO 接口以 `I` 开头（`IUser`、`IVWUser`、`ILang`）
- 获取连接：`global.Database().Get("platform")`

### 标准模式（推荐）

`dbrepo.IDAO[T]` + `dbrepo.ISearcher[T, F]`，适用于 CRUD + 列表查询：

```go
type vwUser struct {
    db       *gorm.DB
    dao      dbrepo.IDAO[pmodel.VWUser]
    searcher dbrepo.ISearcher[pmodel.VWUser, pfilter.User]
}

func NewVWUser() IVWUser {
    result := &vwUser{}  // 先创建零值，让 buildFilter/getSortMapping 可引用
    db, _ := global.Database().Get("platform")
    dao, _ := dbrepo.NewDAO[pmodel.VWUser](db)
    searcher, _ := dbrepo.NewSearcher[pmodel.VWUser, pfilter.User](db,
        dbrepo.WithFilterTransfer[pmodel.VWUser, pfilter.User](result.buildFilter),
        dbrepo.WithSortMapping[pmodel.VWUser, pfilter.User](result.getSortMapping()),
    )
    return &vwUser{db: db, dao: dao, searcher: searcher}
}
```

### 原始模式

直接 `*gorm.DB`，适用于复杂事务。Service 层禁止直接操作 DB。

### buildFilter

```go
func (u *vwUser) buildFilter(filter *pfilter.User, db *gorm.DB) *gorm.DB {
    if db == nil {
        db = u.db.Model(pmodel.VWUser{})  // 必须，保证 Count 查询正确
    }
    if v := filter.AccountLike; len(v) > 0 {
        db = db.Where("account like ?", fmt.Sprintf("%%%s%%", v))
    }
    if v := filter.State; v != nil {      // 指针：v != nil，不判断 *v
        db = db.Where("state = ?", v)
    }
    if v := filter.IDs; len(v) > 0 {
        db = db.Where("id in (?)", v)
    }
    return db
}
```

字符串判空用 `len(v) > 0`，不用 `!= ""`。模糊匹配在 DAO 层拼接 `%`。

### getSortMapping

```go
func (u *vwUser) getSortMapping() *dbquery.SortMapping {
    return dbquery.NewSortMapping(
        dbquery.WithSortKeyMap(map[string]string{
            "created_at":   "created_at",
            "created_time": "created_at",  // 前端别名
        }),
        dbquery.WithDefaultSort("created_at"),  // 默认 DESC，需 ASC 传第二参数
    )
}
```

### DAO 错误处理

使用 `xerror` + `xcode` 组合区分错误类型：

```go
// 参数错误
return nil, xerror.NewXCode(xcode.DBRequestParamError)

// 记录不存在
return nil, xerror.NewXCode(xcode.DBRecordNotFound)

// 数据库操作失败
return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
```

GORM 错误判断模式：

```go
if err := db.First(&record).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, xerror.WithXCode(xcode.DBRecordNotFound)
    }
    return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
}
```

所有错误都必须被处理：用 `xerror.Wrap` 给用户友好提示，或通过日志记录，或用 `_` 拦截。**只有完全不需要关注的错误才能用 `_`**。

### 自定义查询模板

1. 参数校验 → `xcode.DBRequestParamError`
2. 构建查询 → 必须 `u.db.WithContext(ctx)`
3. 错误处理：`gorm.ErrRecordNotFound` → `xcode.DBRecordNotFound`，其他 → `xerror.WrapWithXCode(err, xcode.DBFailed)`

### 事务

```go
return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
    // 事务内用 tx，禁止用 u.db
    // 错误都 xerror.WrapWithXCode(err, xcode.DBFailed) 包装返回
})
```

也支持 `dbrepo.RunInTx(ctx, u.db, fn)`，自动处理嵌套 Savepoint。

### 敏感数据原子操作

货币、余额等敏感数据必须用原子操作防止并发竞态：

```go
// 方式一：DB 函数运算（推荐，适合增减场景）
tx.Model(&platform.Wallet{}).Where("user_id = ?", userID).
    Update("balance", gorm.Expr("balance + ?", amount))

// 方式二：乐观锁（适合先读后改场景）
result := tx.Model(&platform.Wallet{}).
    Where("user_id = ? AND balance = ?", userID, oldBalance).
    Update("balance", newBalance)
if result.RowsAffected == 0 {
    return xerror.New("余额已变更，请重试")
}
```

### IDAO / ISearcher 方法速查

**IDAO**：`Create` / `Creates` / `Save` / `First` / `FirstWith` / `Delete`(软) / `Remove`(硬) / `Update`(指定字段) / `WithTx`

**ISearcher**：`FindAll` / `List`(无总数) / `Paginate`(有总数) / `ListByCursor` / `Find` / `FirstWith` / `CountBy` / `ExistsBy` / `WithTx`

## 依赖方向

`pdao → pmodel, pattr, pfilter`（不反向依赖）。获取连接：`global.Database().Get("platform")`
