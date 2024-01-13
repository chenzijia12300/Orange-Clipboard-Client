package main

import (
	"orange-clipboard/client"
	"orange-clipboard/client/conf"
	"orange-clipboard/client/db"
	"orange-clipboard/common/resource"
	"time"
)

func main() {
	resource.InitLog()
	conf.InitConf()
	db.InitDB()
	client.InitClipboard()
	//client.InitConnectServer(ctx)
	client.AddMessageListener(func(messageContainer client.MessageContainer) bool {
		db.InsertOrUpdate(db.ClipboardModel{
			Msg:        string(messageContainer.Data),
			MsgType:    db.MsgTextType,
			CreateTime: time.Now().Unix(),
		})
		return true
	})
	client.InitUI()

}
