// Package cmd 提供命令行接口功能
package cmd

import (
	"fmt"
	"os"

	"github.com/yyhuni/xingfinger/module/finger"
	"github.com/yyhuni/xingfinger/module/finger/source"

	"github.com/spf13/cobra"
)

// Banner 程序横幅
const Banner = `
  __  ___                _____ _                       
  \ \/ (_)___  ____ _   / ____(_)___  ____ ____  _____ 
   \  /| / _ \/ __ ` + "`" + `/  / /_  / / __ \/ __ ` + "`" + `/ _ \/ ___/ 
   /  \| |  __/ /_/ /  / __/ / / / / / /_/ /  __/ /     
  /_/\_\_|\___/\__, /  /_/   /_/_/ /_/\__, /\___/_/      
              /____/                 /____/   By:yyhuni
`

// 命令行参数变量
var (
	inputFile  string // 输入文件路径
	targetURL  string // 单个目标 URL
	threadNum  int    // 并发线程数
	outputFile string // 输出文件路径
	proxyAddr  string // 代理服务器地址
	timeout    int    // 请求超时时间（秒）
	silent     bool   // 安静模式
)

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "xingfinger",
	Short: "Web fingerprint scanner",
	Run:   runScan,
}

// Execute 执行根命令
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

// init 初始化命令行参数
func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.Flags().StringVarP(&inputFile, "list", "l", "", "Input file with URLs")
	rootCmd.Flags().StringVarP(&targetURL, "url", "u", "", "Single target URL")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (json)")
	rootCmd.Flags().IntVarP(&threadNum, "thread", "t", 100, "Thread count")
	rootCmd.Flags().StringVarP(&proxyAddr, "proxy", "p", "", "Proxy address")
	rootCmd.Flags().IntVar(&timeout, "timeout", 10, "Request timeout (seconds)")
	rootCmd.Flags().BoolVar(&silent, "silent", false, "Silent mode, only output matched results")
}

// runScan 执行扫描任务
func runScan(cmd *cobra.Command, args []string) {
	// 显示横幅
	if !silent {
		fmt.Print(Banner)
	}

	var urls []string

	switch {
	case inputFile != "":
		urls = deduplicate(source.LoadFromFile(inputFile))
	case targetURL != "":
		urls = []string{targetURL}
	default:
		cmd.Help()
		return
	}

	scanner := finger.NewScanner(urls, threadNum, outputFile, proxyAddr, timeout, silent)
	scanner.Run()
	os.Exit(0)
}

// deduplicate 去重
func deduplicate(arr []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0)
	for _, v := range arr {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}
