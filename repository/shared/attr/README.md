# shared/attr — 跨库共享业务属性类型

存放被多个数据库共享的业务属性类型（枚举），区别于 `repository/platform/pattr/` 仅服务于 `platform` 库。

## 何时放在这里

- 该枚举被 **两个及以上** 库的 model/dao/filter 引用（如 `platform` 与 `analytics` 共用同一套状态码）。
- 仅单个库使用的枚举应放在该库的 `xattr/` 下（如 `platform` → `pattr/`），避免过早抽象。

## 约定（与 `pattr/` 一致）

1. 每个业务概念一个文件，文件名 `snake_case`，类型名 `PascalCase`。
2. 自定义类型 + `iota` 常量；负值用 `iota - N` 偏移。
3. 必须实现 `String()` 返回中文文本。
4. 必须提供 `ParseXxx(val int8) (Xxx, error)` 安全解析；字符串解析提供 `NewXxx(str) Xxx`。

## 示例

```go
package attr

import "github.com/gomooth/xerror"

// OrderState 订单状态（跨库共享）
type OrderState int8

const (
    OrderStatePending OrderState = iota // 待支付
    OrderStatePaid                      // 已支付
    OrderStateClosed                    // 已关闭
)

func (s OrderState) String() string {
    switch s {
    case OrderStatePending:
        return "待支付"
    case OrderStatePaid:
        return "已支付"
    case OrderStateClosed:
        return "已关闭"
    default:
        return "未定义"
    }
}

func ParseOrderState(val int8) (OrderState, error) {
    switch val {
    case int8(OrderStatePending):
        return OrderStatePending, nil
    case int8(OrderStatePaid):
        return OrderStatePaid, nil
    case int8(OrderStateClosed):
        return OrderStateClosed, nil
    default:
        return OrderStatePending, xerror.Errorf("未定义的订单状态[%d]", val)
    }
}
```

## 依赖方向

`shared/attr` 不依赖任何 `xmodel`/`xdao`/`xfilter`，仅被它们引用。禁止反向依赖。
