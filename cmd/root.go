// Package cmd 提供命令行接口功能
// 使用 cobra 框架实现命令行参数解析和子命令管理
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cfgFile 配置文件路径，可通过 --config 参数指定
var cfgFile string

// rootCmd 根命令，所有子命令都挂载在此命令下
var rootCmd = &cobra.Command{
	Use:   "xingfinger",
	Short: "Web fingerprint scanner",
}

// Execute 执行根命令
// 这是命令行程序的入口点，由 main 函数调用
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

// init 初始化函数
// 在包加载时自动执行，用于设置命令行参数和初始化配置
func init() {
	// 注册配置加载函数，在命令执行前加载配置
	cobra.OnInitialize(loadConfig)

	// 禁用默认的 completion 子命令
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// 添加全局配置文件参数
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}

// loadConfig 加载配置文件
// 支持两种方式：
//  1. 通过 --config 参数指定配置文件路径
//  2. 自动从用户主目录加载 .xingfinger.yaml 配置文件
func loadConfig() {
	if cfgFile != "" {
		// 使用用户指定的配置文件
		viper.SetConfigFile(cfgFile)
	} else {
		// 获取用户主目录
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// 在主目录下查找配置文件
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".xingfinger")
	}

	// 自动读取环境变量
	viper.AutomaticEnv()

	// 读取配置文件（忽略错误，配置文件可选）
	viper.ReadInConfig()
}
