package main

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	deviceInfoList = make([]DeviceInfo, 0)
)

type DeviceInfo struct {
	UserId     string          // 用户id UUID
	IPV4       string          // IPv4网络
	Platform   string          // 平台 Windows/MacOS/Android
	Connection *websocket.Conn // WebSocket连接
}

type ReadMessageHandler func([]byte) bool

func InitServer() {
	http.HandleFunc("/ws", handler)
	http.ListenAndServe("localhost:8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		Logger.Error("建立websocket连接失败", zap.Error(err))
		return
	}
	deviceInfo := DeviceInfo{}
	platform := r.Header.Get("platform")
	userId := r.Header.Get("userId")
	deviceInfo.Platform = platform
	deviceInfo.UserId = userId
	deviceInfo.IPV4 = r.RemoteAddr
	deviceInfo.Connection = conn
	deviceInfoList = append(deviceInfoList, deviceInfo)
	ReadMessage(deviceInfo, WriteClipboard)
}

func WriteMessage(messageType int, data []byte) {
	go func() {
		for i, deviceInfo := range deviceInfoList {
			conn := deviceInfo.Connection
			err := conn.WriteMessage(messageType, data)
			if err != nil {
				Logger.Info("传递信息失败",
					zap.String("userId", deviceInfo.UserId),
					zap.String("IPv4", deviceInfo.IPV4))
				deviceInfoList = append(deviceInfoList[:i], deviceInfoList[i+1:]...)
			}
		}
	}()
}

func ReadMessage(deviceInfo DeviceInfo, readHandler ReadMessageHandler) {
	go func() {
		defer func() {
			deviceInfo.Connection.Close()
		}()
		for {
			messageType, message, err := deviceInfo.Connection.ReadMessage()
			if err != nil {
				Logger.Debug("读取消息失败",
					zap.String("userId", deviceInfo.UserId),
					zap.String("IPv4", deviceInfo.IPV4),
					zap.Error(err))
				return
			}
			if messageType == CANCEL {
				Logger.Info("断开连接",
					zap.String("userId", deviceInfo.UserId),
					zap.String("IPv4", deviceInfo.IPV4))
				return
			}
			if !readHandler(message) {
				return
			}
		}
	}()
}
