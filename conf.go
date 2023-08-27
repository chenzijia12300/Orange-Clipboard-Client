package main

import (
	"github.com/pelletier/go-toml"
	"go.uber.org/zap"
	"os"
)

type ClipboardConfig struct {
	SecretKey  string
	DeviceName string
}

var clipboardConfig ClipboardConfig

const ConfigFilePath = "./conf.toml"

func InitConf() {
	conf, err := toml.LoadFile(ConfigFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			Logger.Info("conf.toml配置文件不存在,尝试创建并生成相关信息")
			createConf()
			return
		}
		Logger.Error("加载配置文件失败", zap.Error(err))
		return
	}
	err = conf.Unmarshal(&clipboardConfig)
	if err != nil {
		Logger.Error("解析配置文件出错", zap.Error(err))
		return
	}
	Logger.Info("配置信息", zap.String("secretKey", clipboardConfig.SecretKey), zap.String("deviceName", clipboardConfig.DeviceName))
}

func createConf() {
	confFile, err := os.OpenFile(ConfigFilePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		Logger.Error("初始化配置文件失败", zap.Error(err))
	}
	defer confFile.Close()
	clipboardConfig.SecretKey = GenerateRandomBytes()
	clipboardConfig.DeviceName = "Master"
	encoder := toml.NewEncoder(confFile)
	encoder.Encode(clipboardConfig)
}
