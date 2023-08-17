package client

import (
	"fmt"
	"golang.design/x/clipboard"
)

func MustInit() error {
	err := clipboard.Init()
	if err != nil {
		fmt.Println("初始化组件失败:", err)
		return err
	}
	InitSystemTray()
	InitClipboard()
	go InitServer()
	// TODO 初始化日志组件
	return nil
}
