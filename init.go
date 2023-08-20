package main

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
	InitLog()
	Parse()
	go InitSystemTray()
	go InitClipboard()
	go InitServer()
	return nil
}
