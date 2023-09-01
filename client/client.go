package client

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"orangeadd.com/clipboard-client/conf"
	"orangeadd.com/clipboard-client/resource"
)

var (
	messageCh      chan messageContainer   = make(chan messageContainer)
	WriteMessageCh chan<- messageContainer = messageCh
)

type messageContainer struct {
	Type int
	Data []byte
}

func ConnectServer() {
	serverUrl := conf.GlobalConfig.ServerUrl
	header := http.Header{}
	header.Add(conf.SystemName, conf.GlobalConfig.SystemName)
	header.Add(conf.DeviceName, conf.GlobalConfig.DeviceName)
	conn, _, err := websocket.DefaultDialer.Dial(serverUrl, header)
	if err != nil {
		resource.Logger.Error("连接服务器失败", zap.String("serverUrl", conf.GlobalConfig.ServerUrl), zap.Error(err))
		return
	}
	go ReadServerMessage(conn, WriteClipboard)
	go WriteServerMessage(conn, messageCh)
}

func ReadServerMessage(conn *websocket.Conn, readHandler ReadMessageHandler) {
	defer conn.Close()
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			resource.Logger.Debug("读取服务器消息失败",
				zap.String("serverUrl", conn.RemoteAddr().String()),
				zap.Error(err))
			return
		}
		if messageType == conf.CANCEL {
			resource.Logger.Info("服务器断开连接",
				zap.String("serverUrl", conn.RemoteAddr().String()))
			return
		}
		if !readHandler(message) {
			return
		}
	}
}

func WriteServerMessage(conn *websocket.Conn, readMessageCh <-chan messageContainer) {
	defer conn.Close()
	for messageContainer := range readMessageCh {
		messageType := messageContainer.Type
		data := messageContainer.Data
		err := conn.WriteMessage(messageType, data)
		if err != nil {
			resource.Logger.Info("传递信息失败",
				zap.String("serverUrl", conn.RemoteAddr().String()), zap.Error(err))
			return
		}
	}
}
