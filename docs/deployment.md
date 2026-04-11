# 部署指南

## 部署

构建与运行命令见 `CLAUDE.md`「常用命令」与「快速开始」。本文只记录部署形态相关的配置：进程托管、容器化、健康检查、生产清单与版本管理。

环境要求：Go 1.25+，MySQL 5.7+ 或 SQLite3 或 PostgreSQL，Redis 6.0+（可选）。

## Systemd 服务配置

```ini
# /etc/systemd/system/server-api.service
[Unit]
Description=Server API Service
After=network.target mysql.service redis.service

[Service]
Type=simple
User=www
Group=www
WorkingDirectory=/opt/server-api
ExecStart=/opt/server-api/main -mode=all
Restart=always
RestartSec=3
Environment="APP_ENV=prod"

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl start server-api
sudo systemctl enable server-api
sudo journalctl -u server-api -f    # 查看日志
```

## Docker 部署

```dockerfile
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY .dist/main .
COPY storage/config/config.toml ./storage/config/
RUN mkdir -p storage/logs && \
    addgroup -g 1000 app && adduser -u 1000 -G app -s /bin/sh -D app
USER app
EXPOSE 8000 8001
CMD ["./main", "-mode=all"]
```

```bash
docker build -t server-api:latest .
docker run -d --name server-api \
  -p 8000:8000 -p 8001:8001 \
  -v ./storage/config:/app/storage/config \
  -v ./storage/logs:/app/storage/logs \
  server-api:latest
```

## 健康检查

```bash
curl http://localhost:8000/ping       # 简单 pong 响应
curl http://localhost:8000/healthz    # 存活探针
curl http://localhost:8000/readyz     # 就绪探针（检查 DB/Redis/系统资源）
curl http://localhost:8001/ping       # OpenAPI 服务
```

## 生产环境清单

- [ ] 修改 `app.secret` 为强随机字符串
- [ ] 修改 `app.admin.password`
- [ ] 设置 `app.clear_example_file=true`
- [ ] 日志级别设为 `info`
- [ ] 通过反向代理配置 HTTPS
- [ ] 只开放必要端口（8000, 8001）
- [ ] 使用非 root 用户运行应用

## 版本管理

```bash
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
./main -version               # 查看版本
```
