// Package pkg 提供 xingfinger 的核心功能
// 本文件负责从各种来源加载目标 URL
package pkg

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadFromFile 从本地文件加载 URL 列表
// 支持两种格式：
//  1. 完整 URL（包含 http:// 或 https://）
//  2. 域名/IP（自动添加 https:// 前缀）
//
// 参数：
//   - filename: URL 列表文件路径，每行一个 URL
//
// 返回：
//   - 处理后的 URL 列表
func LoadFromFile(filename string) (urls []string) {
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("[!] File read error")
		os.Exit(1)
	}
	defer file.Close()

	// 逐行读取并处理
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行
		if line == "" {
			continue
		}

		// 检查是否已包含协议前缀
		lowerLine := strings.ToLower(line)
		if strings.HasPrefix(lowerLine, "http://") || strings.HasPrefix(lowerLine, "https://") {
			urls = append(urls, line)
		} else {
			// 默认添加 https:// 前缀
			urls = append(urls, "https://"+line)
		}
	}
	return
}
