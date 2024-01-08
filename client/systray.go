package client

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"fyne.io/systray"
	hook "github.com/robotn/gohook"
	"go.uber.org/zap"
	"math"
	"orange-clipboard/client/conf"
	"orange-clipboard/client/db"
	"orange-clipboard/client/ui"
	"orange-clipboard/common/resource"
	"os"
	"time"
)

var (
	GlobalWindow           fyne.Window
	SysTrayConnectStatusCh = make(chan bool)
	list                   *widget.List
)

const (
	SuccessTitle = "断开连接"
	FailureTitle = "重新连接"
)

type clipboardTheme struct {
	fyne.Theme
	fontResource fyne.Resource
}

func (c clipboardTheme) Font(style fyne.TextStyle) fyne.Resource {
	return c.fontResource
}

func NewClipboardTheme(parent fyne.Theme, fontResource fyne.Resource) clipboardTheme {
	return clipboardTheme{
		Theme:        parent,
		fontResource: fontResource,
	}
}

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
	resource.Debug("退出本系统")
	os.Exit(0)
}

func InitUI() {
	a := app.New()
	a.Lifecycle().SetOnEnteredForeground(func() {
		resource.Debug("集中焦点")
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		resource.Debug("退出焦点")
		GlobalWindow.Hide()
	})
	loadFont(a)
	drv := fyne.CurrentApp().Driver()
	if drv, ok := drv.(desktop.Driver); ok {
		GlobalWindow = drv.CreateSplashWindow()
		AddShortcuts(GlobalWindow)
		GlobalWindow.Resize(fyne.NewSize(400, 400))
		GlobalWindow.SetContent(makeListTab(GlobalWindow))
		GlobalWindow.ShowAndRun()
	}
}

func loadFont(app fyne.App) {
	font, err := fyne.LoadResourceFromPath(conf.GlobalConfig.FontUrl)
	if err != nil {
		resource.Logger.Error("load font resource failure", zap.Error(err))
		return
	}
	app.Settings().SetTheme(NewClipboardTheme(theme.DefaultTheme(), font))
}

func makeListTab(window fyne.Window) fyne.CanvasObject {
	clipboardModels := db.Query(0, 20)
	list = widget.NewList(
		func() int {
			return len(clipboardModels)
		},
		func() fyne.CanvasObject {
			vBox := container.NewVBox(createHeaderItem(), createContentItem(), createBottomItem())
			//tappableContainer := &TappableContainer{
			//	Container: vBox,
			//	popUpMenu: createMenu(window.Canvas()),
			//}
			return vBox
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			model := clipboardModels[id]
			//container := item.(*TappableContainer).Container
			container := item.(*fyne.Container)
			UpdateHeaderItem(model, container.Objects[0])
			UpdateContentItem(model, container.Objects[1])
			UpdateBottomItem(model, container.Objects[2])
			item.Refresh()
			list.SetItemHeight(id, container.MinSize().Height)
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		model := clipboardModels[id]
		WriteClipboard([]byte(model.Msg))
		list.Unselect(id)
		window.Hide()
	}
	list.OnUnselected = func(id widget.ListItemID) {
	}
	return list
}

func AddShortcuts(window fyne.Window) {
	go func() {
		hook.Register(hook.KeyDown, []string{"`", "ctrl"}, func(e hook.Event) {
			resource.Debug("唤醒剪贴板")
			window.SetMaster()
			window.RequestFocus()
			window.Show()

		})
		s := hook.Start()
		<-hook.Process(s)
	}()
}

func createHeaderItem() fyne.CanvasObject {
	timeTextWidget := ui.CreateDefaultText("2分钟前")
	deviceTextWidget := ui.CreateDefaultText("本机")
	vBox := container.NewVBox(timeTextWidget, deviceTextWidget)
	hBox := container.NewHBox(vBox, layout.NewSpacer(), ui.CreateDefaultText("文本"))
	background := canvas.NewRectangle(ui.Blue)
	headerContainer := container.NewStack(background, hBox)
	return headerContainer
}

func UpdateHeaderItem(dataModel db.ClipboardModel, object fyne.CanvasObject) {
	headerContainer := object.(*fyne.Container)
	hBox := headerContainer.Objects[1].(*fyne.Container)
	leftBox := hBox.Objects[0].(*fyne.Container)
	timeTextWidget := leftBox.Objects[0].(*canvas.Text)
	deviceTextWidget := leftBox.Objects[1].(*canvas.Text)

	// set widget data
	timeTextWidget.Text = convertTimeTextStr(dataModel.CreateTime)
	deviceTextWidget.Text = "本机"

}

func createContentItem() fyne.CanvasObject {
	label := widget.NewLabel("Hello")
	label.Wrapping = fyne.TextWrapBreak
	return label
}

func UpdateContentItem(dataModel db.ClipboardModel, object fyne.CanvasObject) {
	label := object.(*widget.Label)
	label.SetText(calculateTextByMaxRow(label, []rune(dataModel.Msg), 5))
}

func createBottomItem() fyne.CanvasObject {
	charNumTextWidget := ui.CreateDefaultText("共100个字符")
	charNumTextWidget.Alignment = fyne.TextAlignCenter
	return charNumTextWidget
}

func UpdateBottomItem(dataModel db.ClipboardModel, object fyne.CanvasObject) {
	charNumTextWidget := object.(*canvas.Text)
	charNumTextWidget.Text = fmt.Sprintf("共%d个字符", len(dataModel.Msg))
}

func createMenu(canvas fyne.Canvas) *widget.PopUpMenu {
	return widget.NewPopUpMenu(fyne.NewMenu("",
		fyne.NewMenuItem("同步", func() {

		}),
		fyne.NewMenuItem("更改", func() {

		}),
		fyne.NewMenuItem("删除", func() {

		})), canvas)
}

type TappableContainer struct {
	*fyne.Container
	popUpMenu *widget.PopUpMenu
}

func (c *TappableContainer) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(c.Container)
}

func (c *TappableContainer) TappedSecondary(event *fyne.PointEvent) {
	resource.Debug("右键点击")
}

/*
工具方法
*/
func calculateTextByMaxRow(label *widget.Label, text []rune, maxRow int) string {
	var (
		innerPadding = theme.InnerPadding()
		fitSize      = label.Size()
		maxWidth     = fitSize.Width - 2*innerPadding
		wrapWidth    = maxWidth
		max          = fyne.NewSize(maxWidth, fitSize.Height)
		low          = 0
		high         = len(text) - 1
		useRowNum    = 0
	)
	var yPos float32 = 0
	measureWidth := float32(math.Min(float64(wrapWidth), float64(max.Width)))
	widthChecker := func(low int, high int) bool {
		return measurer(text[low:high]).Width <= measureWidth
	}
	for low < high {
		if useRowNum >= maxRow {
			return string(text[0:low])
		}
		measured := measurer(text[low:high])
		//fmt.Printf("fitSize:%v\tlow:%d\thigh:%d\tmeasured:%f\n", fitSize, low, high, measured)
		if measured.Width <= measureWidth {
			useRowNum++
			low = high
			high = len(text) - 1
			measureWidth = max.Width
			yPos += measured.Height
			//fmt.Printf("找到新的一行low:%d\thigh:%d\tmeasureWidth:%f\n", low, high, measureWidth)
		} else {
			newHigh := binarySearch(widthChecker, low, high)
			if newHigh <= low {
				useRowNum++
				low++
				yPos += measured.Height
			} else {
				high = newHigh
			}
		}
	}
	return string(text)
}

func measurer(text []rune) fyne.Size {
	return fyne.MeasureText(string(text), theme.TextSize(), fyne.TextStyle{})
}

func binarySearch(lessMaxWidth func(int, int) bool, low int, maxHigh int) int {
	if low >= maxHigh {
		return low
	}
	if lessMaxWidth(low, maxHigh) {
		return maxHigh
	}
	high := low
	delta := maxHigh - low
	for delta > 0 {
		delta /= 2
		if lessMaxWidth(low, high+delta) {
			high += delta
		}
	}
	for (high < maxHigh) && lessMaxWidth(low, high+1) {
		high++
	}
	return high
}

func convertTimeTextStr(unix int64) string {
	return time.Unix(unix, 0).Format(conf.DateTime)
}
