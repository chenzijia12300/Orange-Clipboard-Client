package server

import (
	"github.com/pelletier/go-toml"
	"go.uber.org/zap"
	"orange-clipboard/common/resource"
	"orange-clipboard/common/utils"
	"os"
)

type ClipboardServerConfig struct {
	Ipv4          string
	Port          int
	MaxConnectNum string
}

var GlobalServerConfig ClipboardServerConfig

const (
	ConfigFilePath = "./conf-server.toml"
)

func InitConf() {
	conf, err := toml.LoadFile(ConfigFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			resource.Logger.Info("conf-server.toml配置文件不存在,尝试创建并生成相关信息")
			createConf()
			return
		}
		resource.Logger.Error("加载配置文件失败", zap.Error(err))
		return
	}
	err = conf.Unmarshal(&GlobalServerConfig)
	if err != nil {
		resource.Logger.Error("解析配置文件出错", zap.Error(err))
		return
	}
	resource.Logger.Info("服务器配置信息", zap.String("Ipv4", GlobalServerConfig.Ipv4), zap.Int("Port", GlobalServerConfig.Port))
}

func createConf() {
	confFile, err := os.OpenFile(ConfigFilePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		resource.Logger.Error("初始化配置文件失败", zap.Error(err))
	}
	defer confFile.Close()
	GlobalServerConfig.Ipv4 = utils.GetLocalIP()
	GlobalServerConfig.Port = 8900
	encoder := toml.NewEncoder(confFile)
	encoder.Encode(GlobalServerConfig)
}
