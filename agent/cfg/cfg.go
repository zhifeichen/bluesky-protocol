package cfg

import (
	"encoding/json"
	"fmt"
	"github.com/zhifeichen/bluesky-protocol/common/xlogger"
	"os"
	"sync"
	"log"
	"github.com/zhifeichen/bluesky-protocol/common"
)

var (
	lock    = new(sync.RWMutex)
	cfgPath string
	config  *GlobalConfig
	Root    string
)

type TransferConfig struct {
	Interval int `json:"interval"`
}

type GlobalConfig struct {
	Debug    bool            `json:"debug"`
	Ip       string          `json:"ip"`
	Port     int             `json:"port"`
	Uuid     string          `json:"uuid"`
	LogFile  string          `json:"logFile"`
	UDPAddr  string          `json:"udpAddr"`
	Timeout  int             `json:"timeout"`
	EndPoint string          `json:"endPoint"`
	Transfer *TransferConfig `json:"transfer"`
}

func InitRootDir() {
	var err error
	Root, err = os.Getwd()
	if err != nil {
		log.Fatalln("getwd fail:", err)
	}
}
func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func IsDebug() bool {
	return config.Debug
}

func rewriteToCfg(cfgPath string, config *GlobalConfig) error {
	f, err := os.OpenFile(cfgPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	data, err := json.Marshal(config)
	if err == nil {
		if _, err = f.Write(data); err != nil {
			xlogger.Info("写入配置文件失败:", err)
		} else {
			xlogger.Info("写入配置文件成功 .... [ok]")
		}

	}

	return err
}

func genAgentUUid() string {
	//return uuid.TimeUUID().String()
	return fmt.Sprintf("G_%s:%d", config.Ip, config.Port)
}

/**
初始化配置等
*/
func InitConfig() (err error) {
	if config.Uuid == "" {
		config.Uuid = genAgentUUid()
		if err = rewriteToCfg(cfgPath, config); err != nil {
			xlogger.Error("写入配置文件:", cfgPath, " 错误:", err)
		}
	}
	return
}

/**
解析配置文件
*/
func ParseConfig(cfg string) error {
	if cfg == "" {
		fmt.Println("使用 -c 指定配置文件")
		return common.NewError(common.CONFIG_ARG_ERROR)
	}

	if !common.IsExist(cfg) {
		fmt.Println("配置文件:", cfg, "不存在,使用 -c 指定配置文件")
		return common.NewError(common.CONFIG_NOT_FOUND)
	}

	cfgPath = cfg
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
