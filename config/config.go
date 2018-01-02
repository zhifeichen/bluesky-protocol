package config

import (
	"sync"
	"encoding/json"
	"github.com/zhifeichen/bluesky-protocol/common/utils"
)

// ServerConfig server config
type ServerConfig struct {
	IP	string	`json:"ip"`
	Port int		`json:"port"`
}

var (
	lock    = new(sync.RWMutex)
	serverConfig ServerConfig
)

// Config get the config
func Config() *ServerConfig {
	lock.Lock()
	defer lock.Unlock()
	return &serverConfig
}

// ParseConfig parse config file.
func ParseConfig(file string) error {
	if file == "" {
		return common.NewError(common.CONFIG_ARG_ERROR)
	}

	if !common.IsExist(file) {
		return common.NewError(common.CONFIG_NOT_FOUND)
	}

	content, err := common.ReadToTrimString(file)
	if err != nil {
		return common.NewErrorOfMsg(common.READ_FILE_ERROR, err.Error());
	}
	err = json.Unmarshal([]byte(content), &serverConfig)
	if err != nil {
		return common.NewError(common.CONFIG_PARSE_ERROR)
	}
	return nil
}
