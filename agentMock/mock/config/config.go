package config

import (
	"fmt"
	"encoding/json"
	"sync"
	"github.com/zhifeichen/bluesky-protocol/common/utils"
)

type MockConfig struct {
	ServerAddr string `json:"serverAddr"`
	Interval int	`json:"interval"`
}

var (
	config MockConfig
	lock = new(sync.RWMutex)
)

func Config() *MockConfig {
	lock.Lock()
	defer lock.Unlock()
	return &config
}

func ParseConfig(cfg string) error {
	if cfg == "" {
		fmt.Println("使用 -c 指定配置文件")
		return common.NewError(common.CONFIG_ARG_ERROR)
	}

	if !common.IsExist(cfg) {
		fmt.Println("配置文件:", cfg, "不存在,使用 -c 指定配置文件")
		return common.NewError(common.CONFIG_NOT_FOUND)
	}

	content, err := common.ReadToTrimString(cfg)

	if err != nil {
		fmt.Println("读取配置文件错误:", cfg, " err:", err)
		return common.NewErrorOfMsg(common.READ_FILE_ERROR, err.Error())
	}

	err = json.Unmarshal([]byte(content), &config)
	if err != nil {
		fmt.Println("解析配置文件失败,", cfg, " err:", err)
		return common.NewError(common.CONFIG_PARSE_ERROR)
	}

	lock.Lock()
	defer lock.Unlock()

	fmt.Println("读取配置文件", cfg, " ... [ok]")
	return nil
}
