package main

import (
	"orangeadd.com/clipboard-client/client"
	"orangeadd.com/clipboard-client/client/conf"
	"orangeadd.com/clipboard-client/client/db"
	"orangeadd.com/clipboard-client/common/resource"
)

func main() {
	resource.InitLog()
	conf.InitConf()
	db.InitDB()
	client.InitClipboard()
	//client.InitConnectServer(ctx)
	client.InitUI()
}
