package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"server-api/boot"
	"server-api/global"
)

var (
	version   = "unknown" // 版本号
	buildTime = "unknown" // 构建时间
	gitCommit = "unknown" // Git提交哈希
)

var (
	flagVersion bool

	flagConf, flagMode string

	flagCMDName    string
	flagCMDTimeout int
	flagCMDArgs    boot.FlagSlice
)

func step() {
	flag.BoolVar(&flagVersion, "version", false, "show version")

	flag.StringVar(&flagConf, "conf", "storage/config/config.toml", "config path, support remote url")
	flag.StringVar(&flagMode, "mode", "all", "server mode: all, web, openapi, cronjob/cron/job, consumer, cmd/command")

	// cmd 参数解析
	flag.StringVar(&flagCMDName, "cmd.name", "", "command task name, only support `cmd` mode")
	flag.IntVar(&flagCMDTimeout, "cmd.timeout", 0, "command task run timeout, second")
	flag.Var(&flagCMDArgs, "cmd.args", "command task run args. default use `:` split key and value, e.g., `-cmd.args=ver:v1.1028` is key=`ver`, value=`v1.1028`")

	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
		_, _ = fmt.Fprint(os.Stderr, "OPTIONS:\n")
		flag.PrintDefaults()
	}

	flag.Parse()
}

// @title server API
// @version 1.0
// @description 接口文档.

func main() {
	step()

	global.BuildParam = global.AppBuildParam{
		Version:   version,
		BuildTime: buildTime,
		GitCommit: gitCommit,
	}
	if flagVersion {
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Build at: %s\n", buildTime)
		fmt.Printf("Git commit: %s\n", gitCommit)
		fmt.Println("\nThank you for your use.")
		os.Exit(0)
	}

	log.Println("launcher starting...")
	log.Printf("launcher flags: conf=%s, mode=%s\n", flagConf, flagMode)

	cnf := boot.Param{
		ConfigFilename:  flagConf,
		RegisterServers: make([]boot.ServerType, 0, 4),
	}

	// command 命令
	if flagMode == "cmd" || flagMode == "command" {
		cnf.CMDParam = &boot.CMDBootParam{
			Name:    flagCMDName,
			Timeout: flagCMDTimeout,
			Args:    flagCMDArgs,
		}
		if err := boot.Command(cnf); err != nil {
			log.Fatalf("boot command cmd failed: %+v\n", err.Error())
		}
		return
	}

	switch flagMode {
	case "all":
		cnf.RegisterServers = append(
			cnf.RegisterServers,
			boot.InitServerTypeWeb,
			boot.InitServerTypeOpenAPI,
			boot.InitServerTypeCronjob,
			boot.InitServerTypeConsumer,
		)
	case "web":
		cnf.RegisterServers = append(cnf.RegisterServers, boot.InitServerTypeWeb)
	case "openapi":
		cnf.RegisterServers = append(cnf.RegisterServers, boot.InitServerTypeOpenAPI)
	case "cron", "job", "cronjob":
		cnf.RegisterServers = append(cnf.RegisterServers, boot.InitServerTypeCronjob)
	case "consumer":
		cnf.RegisterServers = append(cnf.RegisterServers, boot.InitServerTypeConsumer)
	default:
		log.Fatalf("boot failed: mode err. is %+v\n", flagMode)
	}

	if err := boot.Boot(cnf); err != nil {
		log.Fatalf("boot failed: %+v\n", err.Error())
	}
}
