package client

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	deviceInfoList = make([]DeviceInfo, 0)
)

const (
	CANCEL = -1
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
		fmt.Println(err)
		return
	}
	deviceInfo := DeviceInfo{}
	err = json.NewDecoder(r.Body).Decode(&deviceInfo)
	if err != nil {
		fmt.Println("解析连接用户信息出错:", err)
		return
	}
	deviceInfo.IPV4 = r.RemoteAddr
	deviceInfo.Connection = conn
	deviceInfoList = append(deviceInfoList, deviceInfo)
}

func WriteMessage(messageType int, data []byte) {
	go func() {
		for i, deviceInfo := range deviceInfoList {
			conn := deviceInfo.Connection
			err := conn.WriteMessage(messageType, data)
			if err != nil {
				fmt.Printf("传递消息失败:userId:%s\tipv4:%s\terr:%v", deviceInfo.UserId, deviceInfo.IPV4, err)
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
				fmt.Printf("读取消息失败:userId:%s\tipv4:%s\terr:%v", deviceInfo.UserId, deviceInfo.IPV4, err)
				return
			}
			if messageType == CANCEL {
				fmt.Printf("断开连接 userId:%s\tipv4:%s\n", deviceInfo.UserId, deviceInfo.IPV4)
				return
			}
			if !readHandler(message) {
				return
			}
		}
	}()
}
