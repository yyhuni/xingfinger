// Package source 提供数据源功能
// 支持从本地文件获取目标 URL
package source

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// GetExePath 获取可执行文件所在目录
// 用于定位指纹文件
//
// 返回：
//   - 可执行文件所在目录的绝对路径
func GetExePath() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	// 解析符号链接，获取真实路径
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

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
		if strings.Contains(line, "http") {
			urls = append(urls, line)
		} else {
			// 默认添加 https:// 前缀
			urls = append(urls, "https://"+line)
		}
	}
	return
}
