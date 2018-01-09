package main

import (
	"flag"
	"fmt"
	"github.com/zhifeichen/bluesky-protocol/agent/cfg"
	"github.com/zhifeichen/bluesky-protocol/common/xlogger"
	"os"
	"runtime"
	"github.com/zhifeichen/bluesky-protocol/agent/servers/blueSkyProtocol"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	config := flag.String("c", "./agent/cfg.json", "配置文件")
	version := flag.Bool("v", false, "显示版本")
	if *version {
		fmt.Println(cfg.VERSION)
		os.Exit(0)
	}

	// 解析配置文件
	err := cfg.ParseConfig(*config)
	if err != nil {
		os.Exit(1)
	}

	// 启动系统日志
	logFilePath := cfg.Config().LogFile
	level := xlogger.INFO
	if cfg.Config().Debug {
		level = xlogger.DEBUG
	}
	xlogger.New(logFilePath, level, cfg.Config().Debug)
	defer xlogger.Close()

	// 初始化配置
	if err = cfg.InitConfig(); err != nil {
		os.Exit(1)
	}

	xlogger.Info("启动服务ip:", cfg.Config().Ip, " port:", cfg.Config().Port, " ... [ok]")

	servers.Start(cfg.Config().Ip,cfg.Config().Port)
	// TODO 接收信号?
	select {}
}
