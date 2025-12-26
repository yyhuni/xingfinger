// Package cmd 提供命令行接口功能
// 本文件定义了 xingfinger 的命令行参数和执行入口
package cmd

import (
	"fmt"
	"os"

	"github.com/yyhuni/xingfinger/pkg"

	"github.com/spf13/cobra"
)

// Banner 程序横幅
// 在非静默模式下启动时显示
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
	inputFile  string // 输入文件路径，包含待扫描的 URL 列表
	targetURL  string // 单个目标 URL
	threadNum  int    // 并发线程数，默认 100
	outputFile string // 输出文件路径，支持 JSON 格式
	proxyAddr  string // 代理服务器地址，格式如 http://127.0.0.1:8080
	timeout    int    // HTTP 请求超时时间（秒），默认 10
	silent     bool   // 静默模式，只输出命中指纹的结果

	// 自定义指纹文件
	eholeFile       string // EHole 格式指纹文件
	gobyFile        string // Goby 格式指纹文件
	wappalyzerFile  string // Wappalyzer 格式指纹文件
	fingersFile     string // Fingers 原生格式指纹文件
	fingerprintFile string // FingerPrintHub 格式指纹文件
)

// rootCmd 根命令
// xingfinger 的主命令，直接执行扫描功能
var rootCmd = &cobra.Command{
	Use:   "xingfinger",
	Short: "Web fingerprint scanner",
	Long: `XingFinger - Web 指纹识别工具

基于 chainreactors/fingers 多指纹库聚合引擎，支持：
- fingers 指纹库
- wappalyzer 技术检测
- fingerprinthub 指纹中心
- ehole 棱洞指纹
- goby 指纹库

使用示例：
  xingfinger -u http://example.com
  xingfinger -l urls.txt -o result.json
  xingfinger -u http://example.com --silent
  xingfinger -u http://example.com --ehole my_ehole.json`,
	Run: runScan,
}

// Execute 执行根命令
// 这是程序的入口点，由 main 函数调用
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

// init 初始化命令行参数
// 定义所有支持的命令行标志
func init() {
	// 禁用默认的 completion 命令
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// 定义命令行参数
	rootCmd.Flags().StringVarP(&inputFile, "list", "l", "", "输入文件路径，包含待扫描的 URL 列表")
	rootCmd.Flags().StringVarP(&targetURL, "url", "u", "", "单个目标 URL")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "输出文件路径（JSON 格式）")
	rootCmd.Flags().IntVarP(&threadNum, "thread", "t", 100, "并发线程数")
	rootCmd.Flags().StringVarP(&proxyAddr, "proxy", "p", "", "代理地址（如 http://127.0.0.1:8080）")
	rootCmd.Flags().IntVar(&timeout, "timeout", 10, "请求超时时间（秒）")
	rootCmd.Flags().BoolVar(&silent, "silent", false, "静默模式，只输出命中指纹的结果")

	// 自定义指纹文件参数
	rootCmd.Flags().StringVar(&eholeFile, "ehole", "", "自定义 EHole 格式指纹文件")
	rootCmd.Flags().StringVar(&gobyFile, "goby", "", "自定义 Goby 格式指纹文件")
	rootCmd.Flags().StringVar(&wappalyzerFile, "wappalyzer", "", "自定义 Wappalyzer 格式指纹文件")
	rootCmd.Flags().StringVar(&fingersFile, "fingers", "", "自定义 Fingers 原生格式指纹文件")
	rootCmd.Flags().StringVar(&fingerprintFile, "fingerprinthub", "", "自定义 FingerPrintHub 格式指纹文件")
}

// runScan 执行扫描任务
// 根据命令行参数加载 URL 并启动扫描器
//
// 参数：
//   - cmd: cobra 命令对象
//   - args: 命令行参数
func runScan(cmd *cobra.Command, args []string) {
	// 非静默模式下显示横幅
	if !silent {
		fmt.Print(Banner)
	}

	// 加载自定义指纹文件
	customConfig := &pkg.CustomFingerConfig{
		EHole:       eholeFile,
		Goby:        gobyFile,
		Wappalyzer:  wappalyzerFile,
		Fingers:     fingersFile,
		FingerPrint: fingerprintFile,
	}
	if err := pkg.LoadCustomFingerprints(customConfig); err != nil {
		fmt.Printf("[!] %v\n", err)
		os.Exit(1)
	}

	var urls []string

	// 根据输入方式加载 URL
	switch {
	case inputFile != "":
		// 从文件加载 URL 列表并去重
		urls = deduplicate(pkg.LoadFromFile(inputFile))
	case targetURL != "":
		// 使用单个目标 URL
		urls = []string{targetURL}
	default:
		// 未提供输入，显示帮助信息
		cmd.Help()
		return
	}

	// 创建扫描器并执行扫描
	scanner := pkg.NewScanner(urls, threadNum, outputFile, proxyAddr, timeout, silent)
	scanner.Run()
	os.Exit(0)
}

// deduplicate 对字符串切片去重
// 保持原有顺序，移除重复的 URL
//
// 参数：
//   - arr: 原始字符串切片
//
// 返回：
//   - 去重后的字符串切片
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
