# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

这是一个基于Go的快速服务器模板，采用模块化架构设计，支持多种运行模式（HTTP服务、OpenAPI、定时任务、消息消费者等）。项目使用Gin框架，包含完整的认证、数据库、缓存、日志等基础设施。

## 常用命令

### 构建
```bash
# 本地构建
make build

# 构建Linux发行版本（使用UPX压缩）
make build:linux

# 指定版本号构建
make build:linux version=v1.0

# 清理构建产物
make clean
```

### 运行
```bash
# 默认启动所有服务（HTTP、OpenAPI、定时任务、消费者）
./main

# 启动特定模式
./main -mode=web           # 只启动HTTP服务
./main -mode=openapi       # 只启动OpenAPI服务
./main -mode=cronjob       # 只启动定时任务
./main -mode=consumer      # 只启动消息消费者

# 执行一次性脚本命令
./main -mode=cmd -cmd.name=example-simple
./main -mode=cmd -cmd.name=example-simple -cmd.args=key1:value1 -cmd.args=key2:value2

# 查看版本信息
./main -version
```

### 配置
- 配置文件路径：`storage/config/config.toml`
- 支持远程配置文件URL：`./main -conf=https://www.domain.com/app/config.toml`
- 配置文件热重载：通过`app.watch_config_enabled`配置启用

## 架构概览

### 目录结构
```
├── app/                    # 应用层
│   ├── http/              # HTTP服务（管理后台API）
│   │   ├── internal/
│   │   │   ├── api/       # API控制器、服务层、请求/响应
│   │   │   ├── middleware/# 中间件
│   │   │   ├── route/     # 路由注册
│   │   │   └── helper/    # 辅助函数
│   │   └── main.go        # HTTP服务入口
│   ├── openapi/           # OpenAPI服务（对外API）
│   │   └── internal/      # 类似http目录结构
│   ├── job/               # 定时任务
│   │   └── internal/handler/ # 任务处理器
│   └── consumer/          # 消息消费者
│       └── internal/consumer/ # 消费者实现
├── boot/                  # 启动引导
│   ├── internal/          # 各服务启动器
│   ├── app.go             # 应用初始化
│   ├── config.go          # 配置结构
│   └── command.go         # 命令模式
├── global/                # 全局组件
│   ├── internal/          # 内部实现
│   │   ├── database/      # 数据库连接和迁移
│   │   ├── config/        # 配置解析
│   │   ├── logger/        # 日志
│   │   └── locker/        # 分布式锁
│   ├── env.go             # 环境判断
│   ├── logger.go          # 日志接口
│   └── role.go            # 角色权限
├── repository/            # 数据访问层
│   ├── platform/          # 平台相关实体
│   │   ├── dao/           # 数据访问对象
│   │   └── cache/         # 缓存层
│   └── types/             # 类型定义
├── service/               # 公共业务服务层
├── storage/               # 存储目录
│   ├── config/            # 配置文件
│   ├── langs/             # 多语言文件
│   ├── logs/              # 日志文件
│   └── public/            # 静态文件
└── main.go                # 程序入口
```

### 核心组件

1. **配置管理** (`global/internal/config/`)
   - 支持TOML格式配置文件
   - 支持本地和远程配置文件
   - 支持配置文件热重载

2. **数据库** (`global/internal/database/`)
   - 支持多数据库连接
   - 自动迁移表结构
   - 初始化数据支持
   - 使用GORM作为ORM

3. **日志系统** (`global/logger.go`)
   - 结构化日志
   - 支持不同日志级别
   - HTTP请求日志
   - 文件轮转

4. **认证授权**
   - JWT认证
   - 角色权限控制 (`global/role.go`)
   - 验证码支持
   - 2FA支持

5. **缓存和锁**
   - Redis缓存 (`data.cache`)
   - 分布式锁 (`data.locker`)
   - 支持多种存储引擎

### 服务模式

项目支持多种运行模式，通过`-mode`参数控制：
- `all`: 启动所有服务（默认）
- `web`: HTTP管理后台服务（端口8000）
- `openapi`: 对外API服务（端口8001）
- `cronjob`: 定时任务服务
- `consumer`: 消息消费者服务
- `cmd`: 一次性命令执行模式

### API设计模式

HTTP和OpenAPI服务采用相似的分层架构：
```
Controller → Service → Repository → Database
```
- **Controller**: 处理HTTP请求/响应，参数验证
- **Service**: 业务逻辑实现
- **Repository**: 数据访问抽象
- **DAO**: 具体的数据访问实现

### 消息处理
- 支持Kafka消费者 (`app/consumer/kafka.go`)
- 支持Redis队列消费者
- 支持HTTPSQS队列消费者
- 消费者配置在`storage/config/config.toml`中

## 开发指南

### 添加新的API端点
1. 在`app/http/internal/api/`或`app/openapi/internal/api/`下创建模块目录
2. 创建`controller.go`、`service.go`、`request.go`、`response.go`
3. 在对应的路由文件中注册路由（`app/http/internal/route/`或`app/openapi/internal/route/`）

### 添加定时任务
1. 在`app/job/internal/handler/`下创建任务处理器
2. 在`app/job/job.go`中注册任务

### 添加消息消费者
1. 在`app/consumer/internal/consumer/`下创建消费者实现
2. 在`app/consumer/single.go`中注册消费者

### 数据库操作
- 实体定义在`repository/platform/`
- DAO实现在`repository/platform/dao/`
- 缓存实现在`repository/platform/cache/`
- 使用`global.DB()`获取数据库连接

### 配置说明
关键配置项：
- `server.http.addr`: HTTP服务地址（默认:8000）
- `server.openapi.addr`: OpenAPI服务地址（默认:8001）
- `data.persistent.connects`: 数据库连接配置
- `data.redis`: Redis配置
- `app.admin`: 默认管理员账号
- `app.auth_captcha_enabled`: 验证码开关
- `app.auth_2fa_enabled`: 2FA开关

## 注意事项

1. **环境变量**: 通过`app.env`配置环境（local/test/prod）
2. **多语言**: 语言文件在`storage/langs/`，通过`service/lang`包管理
3. **静态文件**: 静态资源在`storage/public/`，通过`/storage`路径访问
4. **示例文件**: 设置`app.clear_example_file=false`保留示例代码
5. **日志目录**: 日志文件输出到`storage/logs/`，按服务分类

## 测试
项目包含基本的测试文件，使用标准Go测试框架。运行测试：
```bash
go test ./...
```