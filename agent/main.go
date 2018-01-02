package main

import (
	"flag"
	"fmt"
	"os"
	"bluesky-protocol/agent/cfg"
	"bluesky-protocol/common/logger"
	"bluesky-protocol/agent/receiver"
)

func main(){
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
	logger.New(logFilePath, cfg.Config().Debug)
	defer logger.Close()

	// 初始化配置
	if err = cfg.InitConfig(); err != nil {
		os.Exit(1)
	}



	logger.Info.Println("启动服务ip:",cfg.Config().Ip," port:",cfg.Config().Port," ... [ok]")

	go receiver.Start()

	// TODO 接收信号?
	select {}
}