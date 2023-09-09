package client

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"orangeadd.com/clipboard-client/conf"
	"orangeadd.com/clipboard-client/resource"
	"time"
)

var (
	messageCh      chan messageContainer
	WriteMessageCh chan<- messageContainer
)

type messageContainer struct {
	Type int
	Data []byte
}

const (
	pongWait   = 30 * time.Second
	pingPeriod = pongWait * 9 / 10
)

func ConnectServer() {
	messageCh = make(chan messageContainer)
	WriteMessageCh = messageCh
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
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(s string) error {
		resource.Logger.Debug("服务端响应pong", zap.String("message", s))
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil && messageType == conf.CANCEL {
			resource.Logger.Info("服务器断开连接",
				zap.String("serverUrl", conn.RemoteAddr().String()))
			return
		}
		if err != nil {
			resource.Logger.Debug("读取服务器消息失败",
				zap.Int("type", messageType),
				zap.String("serverUrl", conn.RemoteAddr().String()),
				zap.Error(err))
			return
		}
		if !readHandler(message) {
			return
		}
	}
}

func WriteServerMessage(conn *websocket.Conn, readMessageCh <-chan messageContainer) {
	defer conn.Close()
	ticker := time.NewTicker(pingPeriod)
	for {
		select {
		case <-ticker.C:
			err := conn.WriteMessage(websocket.PingMessage, []byte("ping"))
			if err != nil {
				resource.Logger.Error("服务器心跳连接失败")
				time.AfterFunc(5*time.Second, func() {
					ConnectServer()
				})
				return
			}
		case messageContainer := <-readMessageCh:
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
}
