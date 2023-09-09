package cmd

import (
	"github.com/spf13/cobra"
	"orangeadd.com/clipboard-client/client"
	"orangeadd.com/clipboard-client/conf"
	"orangeadd.com/clipboard-client/resource"
	"orangeadd.com/clipboard-client/server"
	"os"
)

func Init() {
	resource.InitLog()
	conf.InitConf()
}

func Execute() {
	var rootCmd = &cobra.Command{Use: "clipboard", Run: func(cmd *cobra.Command, args []string) {
		server.InitServer()
	}}
	var clientCmd = &cobra.Command{
		Use:   "client",
		Short: "运行客户端",
		Run: func(cmd *cobra.Command, args []string) {
			resource.Logger.Info("执行客户端初始化操作")
			client.InitClipboard()
			client.ConnectServer()
			client.InitSystemTray()

		},
	}
	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "运行服务端",
		Run: func(cmd *cobra.Command, args []string) {
			resource.Logger.Info("执行服务端初始化操作")
			server.InitServer()
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&resource.DebugFlag, "debug", "d", true, "启动debug级别日志")
	rootCmd.AddCommand(clientCmd)
	rootCmd.AddCommand(serverCmd)
	Init()
	// 解析命令行参数并执行相应的子命令
	if err := rootCmd.Execute(); err != nil {
		resource.Logger.Error("启动应用失败")
		os.Exit(1)
	}
}
