package server

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"orangeadd.com/clipboard-client/client/conf"
	"orangeadd.com/clipboard-client/common/resource"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	deviceInfoList = make([]DeviceInfo, 0)
	messageCh      = make(chan messageContainer)
)

type messageContainer struct {
	DeviceInfo DeviceInfo
	Data       []byte
	Type       int
}

type DeviceInfo struct {
	DeviceName string          // 设备名称
	IPV4       string          // IPv4网络
	Platform   string          // 平台 Windows/MacOS/Android
	Connection *websocket.Conn // WebSocket连接
}

func InitServer() {
	http.HandleFunc("/ws", handler)
	go WriteClientMessage()

	http.ListenAndServe("localhost:8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		resource.Logger.Error("建立websocket连接失败", zap.Error(err))
		return
	}
	deviceInfo := DeviceInfo{}
	platform := r.Header.Get("platform")
	deviceName := r.Header.Get("deviceName")
	deviceInfo.Platform = platform
	deviceInfo.DeviceName = deviceName
	deviceInfo.IPV4 = r.RemoteAddr
	deviceInfo.Connection = conn
	deviceInfoList = append(deviceInfoList, deviceInfo)
	resource.Logger.Info("建立websocket连接", zap.String("ip", deviceInfo.IPV4), zap.String("deviceName", deviceInfo.DeviceName))
	ReadClientMessage(deviceInfo)
}

func WriteClientMessage() {
	for container := range messageCh {
		data := container.Data
		messageType := container.Type
		sendDeviceInfo := container.DeviceInfo
		for i, deviceInfo := range deviceInfoList {
			if sendDeviceInfo.IPV4 == deviceInfo.IPV4 {
				continue
			}
			conn := deviceInfo.Connection
			err := conn.WriteMessage(messageType, data)
			if err != nil {
				resource.Logger.Info("传递信息失败",
					zap.String("deviceName", deviceInfo.DeviceName),
					zap.String("IPv4", deviceInfo.IPV4))
				deviceInfoList = append(deviceInfoList[:i], deviceInfoList[i+1:]...)
			}
		}
	}

}

func ReadClientMessage(deviceInfo DeviceInfo) {
	defer deviceInfo.Connection.Close()
	for {
		messageType, message, err := deviceInfo.Connection.ReadMessage()
		resource.Logger.Debug("接受客户端信息",
			zap.String("deviceName", deviceInfo.DeviceName),
			zap.String("IPv4", deviceInfo.IPV4),
			zap.String("message", string(message)))
		if err != nil {
			resource.Logger.Debug("读取消息失败",
				zap.String("deviceName", deviceInfo.DeviceName),
				zap.String("IPv4", deviceInfo.IPV4),
				zap.Error(err))
			return
		}
		if messageType == conf.CANCEL {
			resource.Logger.Info("断开连接",
				zap.String("deviceName", deviceInfo.DeviceName),
				zap.String("IPv4", deviceInfo.IPV4))
			return
		}
		messageCh <- messageContainer{
			DeviceInfo: deviceInfo,
			Data:       message,
			Type:       messageType,
		}
	}
}
