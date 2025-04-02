package cmd

import (
	"github.com/HuanXin-Chen/MyIM/client"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(clientCmd)
}

var clientCmd = &cobra.Command{
	Use: "client",
	Run: ClientHandle,
}

// 调用逻辑层运行
func ClientHandle(cmd *cobra.Command, args []string) {
	client.RunMain()
}
