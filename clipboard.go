package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
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
			Logger.Info("剪贴板文本信息:", zap.String("message", message))
			secretData := Encrypt(clipboardConfig.SecretKey, messageBytes)
			WriteMessage(NORMAL, secretData)
		}
	}()
}

func WriteClipboard(secretData []byte) bool {
	data := Decrypt(clipboardConfig.SecretKey, secretData)
	clipboard.Write(clipboard.FmtText, data)
	return true
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
