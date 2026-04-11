# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

基于Go的快速服务器模板，Gin框架，支持多种运行模式（HTTP服务、OpenAPI、定时任务、消息消费者）。

## 常用命令

```bash
make build                    # 本地构建
make build:linux              # 构建Linux版本（UPX压缩）
make build:linux version=v1.0 # 指定版本号
make clean                    # 清理构建产物

./main                        # 启动所有服务
./main -mode=web              # 只启动HTTP服务（端口8000）
./main -mode=openapi          # 只启动OpenAPI服务（端口8001）
./main -mode=cronjob          # 只启动定时任务
./main -mode=consumer         # 只启动消息消费者
./main -mode=cmd -cmd.name=example-simple          # 执行一次性命令
./main -mode=cmd -cmd.name=example-simple -cmd.args=key1:value1
./main -version               # 查看版本

# 配置
# 配置文件：storage/config/config.toml
# 远程配置：./main -conf=https://www.domain.com/app/config.toml
# 热重载：app.watch_config_enabled
```

## 架构

```
app/              → 服务实现（HTTP/OpenAPI/Job/Consumer）
boot/             → 启动框架
testhelper/       → 测试工具
global/           → 基础设施（配置/日志/数据库/环境/错误码/Kafka常量）
repository/       → 数据访问（Model/DAO/Cache/Attr/Filter）
service/          → 公共业务服务（lang/i18n）
storage/          → 配置/日志/语言/静态资源
```

依赖方向：`app/ → boot/ → global/ → repository/ → service/`

## 文档（按需加载）

写代码前根据场景读取对应文档：

- 架构细节 → `docs/architecture.md`
- 配置参数 → `docs/configuration.md`
- 部署运维 → `docs/deployment.md`
- 安全/性能/注意事项 → `docs/notes.md`

## 编码规则（路径自动加载）

`.claude/rules/` 下的规则文件带 `paths:` frontmatter，编辑匹配路径的代码时自动注入上下文，无需手动读取；新增/修改规则时以各文件 frontmatter 为准。

## 快速开始

1. `cp storage/config/config.example.toml storage/config/config.toml`
2. 修改数据库连接等参数
3. `make build && ./main`
