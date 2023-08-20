package main

import (
	"fmt"
	"go.uber.org/zap"
)

var (
	Logger      *zap.Logger
	SugarLogger *zap.SugaredLogger
)

func InitLog() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("初始化日志化组件失败")
		return
	}
	Logger = logger
	SugarLogger = logger.Sugar()
}
