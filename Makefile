# 版本号设置：优先使用命令行 `version=v1.0` 传入的版本，否则使用 Git Tag 或默认值
VERSION ?= $(if $(filter-out undefined default,$(origin version)),$(version),$(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev"))
GIT_COMMIT ?= $(shell git rev-parse --short HEAD)
BUILD_TIME ?= $(shell date -u +"%Y-%m-%d %H:%M:%S")

GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_BUILD_COMPRESSION=$(GO_BUILD) -ldflags="-s -w -X 'main.version=${VERSION}' -X 'main.buildTime=${BUILD_TIME}' -X 'main.gitCommit=${GIT_COMMIT}'"
BINARY_NAME=main

# 平台检测和设置
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
    CURRENT_OS = darwin
    BINARY_NAME = main
else ifeq ($(UNAME_S),Linux)
    CURRENT_OS = linux
    BINARY_NAME = main
else
    CURRENT_OS = windows
    BINARY_NAME = main.exe
endif

build:
	$(GO_BUILD) -ldflags="-X 'main.version=${VERSION}' -X 'main.buildTime=${BUILD_TIME}' -X 'main.gitCommit=${GIT_COMMIT}'" -o $(BINARY_NAME)

# 构建当前平台
build\:linux:
	#CGO_ENABLED=1 CC=x86_64-linux-musl-gcc CGO_LDFLAGS="-static" GOOS=linux GOARCH=amd64 $(GO_BUILD_COMPRESSION) -o $(BINARY_NAME) && upx $(BINARY_NAME)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_BUILD_COMPRESSION) -o $(BINARY_NAME) && upx $(BINARY_NAME)
	mkdir -p .dist .qa.dist
	mv $(BINARY_NAME) .dist
	cp -p .dist/$(BINARY_NAME) .qa.dist

docker:
	docker build . -t server-api:0.0.1

# 清理构建产物
clean:
	rm -rf .dist .qa.dist $(BINARY_NAME) $(BINARY_NAME).exe
