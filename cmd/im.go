package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	ConfigPath string
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(
		&ConfigPath,
		"config",
		"im.yaml",
		"config file (default is ./im.yaml)")
}

// 命令行工具本身是一棵树,和路由解析调用对应的handler一个道理
var rootCmd = &cobra.Command{
	Use:   "im",
	Short: "这是一个超牛逼的IM系统",
	Run:   MyIM,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func MyIM(cmd *cobra.Command, args []string) {

}

func initConfig() {

}
