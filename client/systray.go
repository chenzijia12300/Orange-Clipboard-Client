package client

import (
	"fmt"
	"github.com/getlantern/systray"
	"go.uber.org/zap"
	"orangeadd.com/clipboard-client/common/resource"
	"os"
)

var SysTrayConnectStatusCh = make(chan bool)

const (
	SuccessTitle = "断开连接"
	FailureTitle = "重新连接"
)

func InitSystemTray() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("Orange Clipboard")
	systray.SetTooltip("Orange Clipboard")
	addConnectStatusMenu()
}

func addConnectStatusMenu() {
	menuItem := systray.AddMenuItem("连接状态", "初始化")
	clickedCh := menuItem.ClickedCh
	go func() {
		defer func() {
			if err := recover(); err != nil {
				resource.Logger.Error("初始化systray组件失败", zap.Error(err.(error)))
			}
		}()
		for {
			select {
			case connectErrorFlag := <-SysTrayConnectStatusCh:
				title := ""
				if connectErrorFlag {
					title = FailureTitle
				} else {
					title = SuccessTitle
				}
				menuItem.SetTitle(title)
				menuItem.SetTooltip(title)
			case <-clickedCh:

			}
		}
	}()
}

func onExit() {
	fmt.Println("退出本系统")
	os.Exit(0)
}
