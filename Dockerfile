
FROM golang:1.18-alpine3.16 AS builder

RUN set -eux \
    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk add upx make gcc musl-dev

COPY . /src
WORKDIR /src

RUN set -ex \
    && GOPROXY=https://goproxy.cn go build -ldflags="-s -w" -o ./bin/main \
    && upx ./bin/main


FROM alpine:3.16
RUN set -eux \
    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    # 设置时区为上海
    && apk add tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata \
    # 创建相关目录
    && mkdir -p /app/config /app/storage/logs

EXPOSE 8000

ENV CONF="config/config.toml"
ENV MODE="all"

COPY --from=builder /src/bin /app
WORKDIR /app

RUN set -eux \
    # 创建运行用户
    && addgroup -S www \
    && adduser -u 65530 -S www -G www \
    && chown -R www:www /app

USER 65530

CMD ["./main", "-conf", "${CONF}", "-mode", "${MODE}"]
