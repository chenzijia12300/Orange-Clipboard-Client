package main

import (
	"context"
	"github.com/spf13/cobra"
	"orangeadd.com/clipboard-client/client"
	"orangeadd.com/clipboard-client/client/conf"
	"orangeadd.com/clipboard-client/common/resource"
	"os"
)

func main() {
	execute()
}

func execute() {
	ctx := context.Background()
	var clientCmd = &cobra.Command{
		Use:   "client",
		Short: "运行客户端",
		Run: func(cmd *cobra.Command, args []string) {
			resource.InitLog()
			conf.InitConf()
			client.InitClipboard()
			client.InitConnectServer(ctx)
			client.InitSystemTray()
		},
	}

	clientCmd.PersistentFlags().BoolVarP(&resource.DebugFlag, "debug", "d", true, "启动debug级别日志")
	// 解析命令行参数并执行相应的子命令
	if err := clientCmd.Execute(); err != nil {
		resource.Logger.Error("启动应用失败")
		os.Exit(1)
	}
}
