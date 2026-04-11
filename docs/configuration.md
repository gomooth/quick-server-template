# 配置说明

配置字段清单及注释见 `storage/config/config.example.toml`（不存在时会自动复制为 `config.toml`）。本文只记录 toml 表达不出的行为特性。

## 多源配置

- **本地文件**：`storage/config/config.toml`
- **远程 URL**：`./main -conf=https://www.domain.com/app/config.toml`
- **自动复制**：`config.toml` 不存在时从 `config.example.toml` 自动复制

## 环境变量覆盖

优先级：环境变量 > 配置文件 > 默认值。

```bash
export TEST_CONFIG_PATH=/path/to/config.toml   # 指定配置文件路径
export APP_ENV=prod                             # 设置应用环境
```

## 三个清理/监听开关的交互

| 开关 | 作用 | 交互 |
|------|------|------|
| `app.clear_example_file` | 删除 `config.example.toml` | 独立生效 |
| `app.clear_config_file` | 删除 `config.toml` | 若同时开启 `watch_config_enabled` 则跳过删除 |
| `app.watch_config_enabled` | 监听配置文件变更，免重启热重载 | 部分配置（如数据库连接）仍需重启生效 |

## 配置验证

```bash
./main -conf=config.toml -mode=web
```
