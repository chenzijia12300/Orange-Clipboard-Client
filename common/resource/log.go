package resource

import (
	"fmt"
	"go.uber.org/zap"
	"os"
)

var (
	Logger      *zap.Logger
	SugarLogger *zap.SugaredLogger
	DebugFlag   bool
)

func InitLog() {
	var err error
	Logger, err = zap.NewDevelopment()
	if err != nil {
		fmt.Println("初始化日志化组件失败")
		os.Exit(0)
	}
	SugarLogger = Logger.Sugar()
}

func Debug(msg string) {
	fmt.Println(msg)
}
