package conf

import (
	"github.com/pelletier/go-toml"
	"go.uber.org/zap"
	"orangeadd.com/clipboard-client/resource"
	"orangeadd.com/clipboard-client/utils"
	"os"
	"runtime"
)

type ClipboardConfig struct {
	SecretKey  string
	DeviceName string
	SystemName string
	ServerUrl  string
}

var GlobalConfig ClipboardConfig

const (
	ConfigFilePath   = "./conf.toml"
	DefaultServerUrl = "ws://localhost:8080/ws"
)

func InitConf() {
	conf, err := toml.LoadFile(ConfigFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			resource.Logger.Info("conf.toml配置文件不存在,尝试创建并生成相关信息")
			createConf()
			return
		}
		resource.Logger.Error("加载配置文件失败", zap.Error(err))
		return
	}
	err = conf.Unmarshal(&GlobalConfig)
	if err != nil {
		resource.Logger.Error("解析配置文件出错", zap.Error(err))
		return
	}
	resource.Logger.Info("配置信息", zap.String("secretKey", GlobalConfig.SecretKey), zap.String("deviceName", GlobalConfig.DeviceName))
}

func createConf() {
	confFile, err := os.OpenFile(ConfigFilePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		resource.Logger.Error("初始化配置文件失败", zap.Error(err))
	}
	defer confFile.Close()
	deviceName := "None"
	hostname, err := os.Hostname()
	if err != nil {
		resource.Logger.Warn("初始化配置获得主机名失败", zap.Error(err))
	} else {
		deviceName = hostname
	}
	GlobalConfig.SystemName = runtime.GOOS
	GlobalConfig.SecretKey = utils.GenerateRandomBytes()
	GlobalConfig.DeviceName = deviceName
	GlobalConfig.ServerUrl = DefaultServerUrl
	encoder := toml.NewEncoder(confFile)
	encoder.Encode(GlobalConfig)
}
