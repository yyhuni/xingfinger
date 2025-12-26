// Package cmd 提供 xingfinger 的命令行接口
// 使用 cobra 框架实现命令行参数解析
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yyhuni/xingfinger/pkg"
)

var (
	// 命令行参数
	targetURL string // 单个目标 URL
	urlFile   string // URL 列表文件
	thread    int    // 并发线程数
	timeout   int    // 请求超时时间
	output    string // 输出文件路径
	proxy     string // 代理地址
	silent    bool   // 静默模式

	// 自定义指纹文件
	eholeFile       string // EHole 指纹文件
	gobyFile        string // Goby 指纹文件
	wappalyzerFile  string // Wappalyzer 指纹文件
	fingersFile     string // Fingers 指纹文件
	fingerprintFile string // FingerPrintHub 指纹文件
)

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "xingfinger",
	Short: "Web 指纹识别工具",
	Long: `xingfinger - 高效的 Web 指纹识别工具

支持多种指纹库：EHole、Goby、Wappalyzer、Fingers、FingerPrintHub
支持单个 URL 或批量扫描，支持代理和自定义指纹`,
	Run: runScan,
}

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// 禁用参数按字母排序，按定义顺序显示
	rootCmd.Flags().SortFlags = false

	// 目标参数
	rootCmd.Flags().StringVarP(&targetURL, "url", "u", "", "目标 URL")
	rootCmd.Flags().StringVarP(&urlFile, "file", "f", "", "URL 列表文件")

	// 扫描参数
	rootCmd.Flags().IntVarP(&thread, "thread", "t", 50, "并发线程数")
	rootCmd.Flags().IntVar(&timeout, "timeout", 10, "请求超时时间（秒）")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "输出文件路径（JSON 格式）")
	rootCmd.Flags().StringVarP(&proxy, "proxy", "p", "", "代理地址")
	rootCmd.Flags().BoolVarP(&silent, "silent", "s", false, "静默模式，只输出命中结果")

	// 自定义指纹文件
	rootCmd.Flags().StringVar(&eholeFile, "ehole", "", "自定义 EHole 指纹文件")
	rootCmd.Flags().StringVar(&gobyFile, "goby", "", "自定义 Goby 指纹文件")
	rootCmd.Flags().StringVar(&wappalyzerFile, "wappalyzer", "", "自定义 Wappalyzer 指纹文件")
	rootCmd.Flags().StringVar(&fingersFile, "fingers", "", "自定义 Fingers 指纹文件")
	rootCmd.Flags().StringVar(&fingerprintFile, "fingerprint", "", "自定义 FingerPrintHub 指纹文件")
}

// runScan 执行扫描
func runScan(cmd *cobra.Command, args []string) {
	// 收集目标 URL
	var urls []string

	if targetURL != "" {
		// 单个 URL
		if !strings.HasPrefix(targetURL, "http") {
			targetURL = "https://" + targetURL
		}
		urls = append(urls, targetURL)
	}

	if urlFile != "" {
		// 从文件加载
		urls = append(urls, pkg.LoadFromFile(urlFile)...)
	}

	// 检查是否有目标
	if len(urls) == 0 {
		fmt.Println("[!] 请指定目标 URL (-u) 或 URL 文件 (-f)")
		cmd.Help()
		os.Exit(1)
	}

	// 构建自定义指纹配置
	var customConfig *pkg.CustomFingerConfig
	if eholeFile != "" || gobyFile != "" || wappalyzerFile != "" || fingersFile != "" || fingerprintFile != "" {
		customConfig = &pkg.CustomFingerConfig{
			EHole:       eholeFile,
			Goby:        gobyFile,
			Wappalyzer:  wappalyzerFile,
			Fingers:     fingersFile,
			FingerPrint: fingerprintFile,
		}
	}

	// 创建扫描器并运行
	scanner := pkg.NewScanner(urls, thread, output, proxy, timeout, silent, customConfig)
	scanner.Run()
}
