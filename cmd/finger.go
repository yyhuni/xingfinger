// Package cmd 提供命令行接口功能
package cmd

import (
	"os"

	"github.com/yyhuni/xingfinger/module/finger"
	"github.com/yyhuni/xingfinger/module/finger/source"

	"github.com/spf13/cobra"
)

// scanCmd 指纹扫描子命令
// 用于执行 Web 指纹识别扫描任务
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Fingerprint scan module",
	Run:   runScan,
}

// 命令行参数变量
var (
	inputFile  string // 输入文件路径（包含 URL 列表）
	targetURL  string // 单个目标 URL
	threadNum  int    // 并发线程数
	outputFile string // 输出文件路径
	proxyAddr  string // 代理服务器地址
)

// init 初始化扫描命令
// 注册子命令和命令行参数
func init() {
	// 将 scan 命令添加到根命令
	rootCmd.AddCommand(scanCmd)

	// 注册命令行参数
	scanCmd.Flags().StringVarP(&inputFile, "file", "f", "", "Input file with URLs")
	scanCmd.Flags().StringVarP(&targetURL, "url", "u", "", "Single target URL")
	scanCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (json/xlsx)")
	scanCmd.Flags().IntVarP(&threadNum, "thread", "t", 100, "Thread count")
	scanCmd.Flags().StringVarP(&proxyAddr, "proxy", "p", "", "Proxy address")
}

// runScan 执行扫描任务
// 根据不同的输入源获取目标 URL 列表，然后启动扫描器
//
// 支持的输入源：
//  1. -f/--file: 从文件读取 URL 列表
//  2. -u/--url: 扫描单个 URL
func runScan(cmd *cobra.Command, args []string) {
	var urls []string

	// 根据输入参数选择数据源
	switch {
	case inputFile != "":
		// 从本地文件加载 URL 列表并去重
		urls = deduplicate(source.LoadFromFile(inputFile))
	case targetURL != "":
		// 单个 URL 直接使用
		urls = []string{targetURL}
	default:
		// 未提供任何输入，显示帮助信息
		cmd.Help()
		return
	}

	// 创建扫描器并执行扫描
	scanner := finger.NewScanner(urls, threadNum, outputFile, proxyAddr)
	scanner.Run()
	os.Exit(0)
}

// deduplicate 对字符串切片进行去重
// 使用 map 实现 O(n) 时间复杂度的去重
//
// 参数：
//   - arr: 待去重的字符串切片
//
// 返回：
//   - 去重后的字符串切片，保持原有顺序
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
