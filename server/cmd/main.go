package main

import (
	"github.com/spf13/cobra"
	"orangeadd.com/clipboard-client/common/resource"
	"orangeadd.com/clipboard-client/server"
	"os"
)

func main() {
	execute()
}
func execute() {
	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "运行服务端",
		Run: func(cmd *cobra.Command, args []string) {
			resource.InitLog()
			server.InitConf()
			server.InitServer()
		},
	}
	serverCmd.PersistentFlags().BoolVarP(&resource.DebugFlag, "debug", "d", true, "启动debug级别日志")
	// 解析命令行参数并执行相应的子命令
	if err := serverCmd.Execute(); err != nil {
		resource.Logger.Error("启动应用失败")
		os.Exit(1)
	}
}
