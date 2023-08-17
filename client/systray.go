package client

import (
	"fmt"
	"github.com/getlantern/systray"
	"os"
)

func InitSystemTray() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("Orange Clipboard")
	systray.SetTooltip("Orange Clipboard")
}

func onExit() {
	fmt.Println("退出本系统")
	os.Exit(0)
}
