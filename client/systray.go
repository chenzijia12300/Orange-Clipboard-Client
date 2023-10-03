package client

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"fyne.io/systray"
	"github.com/eiannone/keyboard"
	"go.uber.org/zap"
	"orangeadd.com/clipboard-client/client/conf"
	"orangeadd.com/clipboard-client/client/db"
	"orangeadd.com/clipboard-client/common/resource"
	"os"
	"time"
)

var SysTrayConnectStatusCh = make(chan bool)
var (
	GlobalWindow fyne.Window
)

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

func InitUI() {
	a := app.New()
	a.Lifecycle().SetOnEnteredForeground(func() {
		fmt.Println("集中焦点")
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		fmt.Println("退出焦点")
		GlobalWindow.Hide()
	})
	drv := fyne.CurrentApp().Driver()
	if drv, ok := drv.(desktop.Driver); ok {
		GlobalWindow = drv.CreateSplashWindow()
		AddShortcuts(GlobalWindow)
		GlobalWindow.Resize(fyne.NewSize(400, 400))
		GlobalWindow.SetContent(makeListTab(GlobalWindow))
		GlobalWindow.ShowAndRun()
		GlobalWindow.SetMainMenu(nil)
		GlobalWindow.Title()
	}
}

func makeListTab(window fyne.Window) fyne.CanvasObject {
	clipboardModels := db.Query(0, 20)
	fmt.Printf("数据:%+v\n", clipboardModels)
	icon := widget.NewIcon(nil)
	list := widget.NewList(
		func() int {
			return len(clipboardModels)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel("未知时间"), widget.NewLabel("未知内容"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			model := clipboardModels[id]
			container := item.(*fyne.Container)
			timeLabel := container.Objects[0].(*widget.Label)
			msgLabel := container.Objects[1].(*widget.Label)
			timeLabel.SetText(time.Unix(model.CreateTime, 0).Format(conf.DateTime))
			msgLabel.SetText(model.Msg)
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		model := clipboardModels[id]
		window.Hide()
		WriteClipboard([]byte(model.Msg))
	}
	list.OnUnselected = func(id widget.ListItemID) {
		icon.SetResource(nil)
	}
	list.SetItemHeight(5, 50)
	list.SetItemHeight(6, 50)

	return list
}

func AddShortcuts(window fyne.Window) {
	go func() {
		err := keyboard.Open()
		if err != nil {
			resource.Logger.Error("初始化全局键盘监听失败", zap.Error(err))
		}
		defer keyboard.Close()
		for {

			_, key, err := keyboard.GetKey()
			if err != nil {
				resource.Logger.Error("GetKey() failure", zap.Error(err))
			}
			if key == keyboard.KeyCtrlQ {
				resource.Logger.Debug("唤醒剪贴板")
				window.Show()
				window.RequestFocus()

			}
		}
	}()
}
