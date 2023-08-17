package client

import (
	"context"
	"fmt"
	"golang.design/x/clipboard"
)

func InitClipboard() {
	ListenClipboardText()
	ListenClipboardImage()
}

func ListenClipboardText() {
	ctx := context.Background()
	textCh := clipboard.Watch(ctx, clipboard.FmtText)
	go func() {
		for messageBytes := range textCh {
			message := string(messageBytes)
			fmt.Println("剪贴板文本信息:", message)
			// TODO 发送信息websocket
		}
	}()
}

func ListenClipboardImage() {
	ctx := context.Background()
	imgCh := clipboard.Watch(ctx, clipboard.FmtImage)
	go func() {
		for messageBytes := range imgCh {
			fmt.Println("剪贴板图像信息:", len(messageBytes))
			// TODO 发送信息websocket
		}
	}()
}
