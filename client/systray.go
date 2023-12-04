package client

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"fyne.io/systray"
	hook "github.com/robotn/gohook"
	"go.uber.org/zap"
	"orange-clipboard/client/conf"
	"orange-clipboard/client/db"
	"orange-clipboard/common/resource"
	"os"
	"time"
	"unicode/utf8"
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
			box := container.NewVBox()
			box.Resize(fyne.NewSize(400, 50))
			timeLabel := widget.NewLabel("03/11 09:00:00")
			label := widget.NewLabel("")
			label.Wrapping = fyne.TextWrapBreak
			box.Add(timeLabel)
			box.Add(label)
			return box
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			model := clipboardModels[id]
			container := item.(*fyne.Container)
			timeLabel := container.Objects[0].(*widget.Label)
			contentLabel := container.Objects[1].(*widget.Label)
			timeLabel.SetText(time.Unix(model.CreateTime, 0).Format(conf.DateTime))
			msg := model.Msg
			if utf8.RuneCountInString(msg) >= 100 {
				msg = string([]rune(msg)[0:100])
			}
			contentLabel.SetText(msg)
		},
	)
	list.ScrollToBottom()
	list.OnSelected = func(id widget.ListItemID) {
		model := clipboardModels[id]
		window.Hide()
		WriteClipboard([]byte(model.Msg))
	}
	list.OnUnselected = func(id widget.ListItemID) {
		icon.SetResource(nil)
	}
	return list
}

func AddShortcuts(window fyne.Window) {
	go func() {
		hook.Register(hook.KeyDown, []string{"`", "ctrl"}, func(e hook.Event) {
			resource.Logger.Debug("唤醒剪贴板")
			window.Show()
			window.SetMaster()
			window.RequestFocus()
		})
		s := hook.Start()
		<-hook.Process(s)
	}()
}

type MaxHeightLayout struct {
	*fyne.Container
	MaxHeight float32
}

func (l MaxHeightLayout) MinSize() fyne.Size {
	minSize := l.Container.MinSize()
	//if minSize.Height > l.MaxHeight {
	//	minSize.Height = l.MaxHeight
	//}
	//minSize.Width = 400
	//fmt.Println(minSize)
	return minSize
}
