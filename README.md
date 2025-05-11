# quick-server-template

## Usage

```shell
git clone --depth=1 -b main https://github.com/gomooth/quick-server-template.git server-api
#git clone --depth=1 -b develop https://github.com/gomooth/quick-server-template.git server-api
#git clone --depth=1 -b feature/v1.0-captcha-dev https://github.com/gomooth/quick-server-template.git server-api

cd server-api
rm -rf .git

git init
git checkout -b main
git remote add origin your_url
git add .
git commit -m "init"
git push -u origin main
```



## 编译

### 本地编译

```shell
make build
```

编译后在项目根目录生成可执行文件

### 编译 linux 发行版本

```shell
make build:linux

# 指定版本号编译
make build:linux version=v1.0
```

该编译会使用 [upx](https://upx.github.io) 压缩可执行文件，并将文件移动到目录 `.dist` 中


## 启动

```shell

# 默认启动所有服务
./main
# 等同于 ./main -conf=storage/config/config.toml -mode=all

# 启动 web
./main -conf=https://www.domain.com/app/config.toml -mode=web

# 启动 cronjob
./main -conf=config.toml -mode=cron

# 启动 consumer
./main -conf=config.toml -mode=consumer

# 执行一次性脚本命令
./main -conf=config.toml -mode=cmd -cmd.name=example-simple
# 携带自定义参数执行脚本命令
./main -conf=config.toml -mode=cmd -cmd.name=example-simple -cmd.args=key1:value1 -cmd.args=key2:value2

```


