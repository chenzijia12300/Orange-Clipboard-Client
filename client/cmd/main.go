package main

import (
	"orange-clipboard/client"
	"orange-clipboard/client/conf"
	"orange-clipboard/client/db"
	"orange-clipboard/common/resource"
)

func main() {
	resource.InitLog()
	conf.InitConf()
	db.InitDB()
	client.InitClipboard()
	//client.InitConnectServer(ctx)
	client.InitUI()
}
