// Package finger 提供 Web 指纹识别核心功能
package finger

import (
	"regexp"
	"strings"
)

// extractRegex 使用正则表达式提取匹配内容
// 返回所有匹配结果及其捕获组
//
// 参数：
//   - pattern: 正则表达式模式
//   - content: 待匹配的内容
//
// 返回：
//   - 所有匹配结果的二维数组，每个元素包含完整匹配和捕获组
func extractRegex(pattern, content string) [][]string {
	re := regexp.MustCompile(pattern)
	return re.FindAllStringSubmatch(content, -1)
}

// parseJSRedirect 解析页面中的 JS 跳转
// 检测常见的 JavaScript 重定向模式，提取跳转目标 URL
//
// 支持的跳转模式：
//  1. window.location.href = "url" 或 top.location.href = "url"
//  2. redirectUrl = "url"
//  3. <meta http-equiv="refresh" content="0;url=xxx">
//
// 参数：
//   - body: HTML 页面内容
//   - baseURL: 当前页面 URL，用于构建完整的跳转 URL
//
// 返回：
//   - 发现的跳转目标 URL 列表
func parseJSRedirect(body, baseURL string) []string {
	// 定义常见的 JS 跳转正则模式
	patterns := []string{
		`(window|top)\.location\.href = ['"](.*?)['"]`, // window.location.href 跳转
		`redirectUrl = ['"](.*?)['"]`,                  // redirectUrl 变量赋值
		`<meta.*?http-equiv=.*?refresh.*?url=(.*?)>`,   // meta refresh 跳转
	}

	var results []string
	for _, p := range patterns {
		matches := extractRegex(p, body)
		for _, m := range matches {
			if len(m) == 0 {
				continue
			}
			// 获取最后一个捕获组（即 URL 部分）
			path := m[len(m)-1]

			// 跳过完整 URL 和空路径
			if strings.Contains(path, "http") || len(path) == 0 {
				continue
			}

			// 清理路径
			path = strings.Trim(path, "/")
			path = strings.ReplaceAll(path, "../", "/")

			// 构建完整 URL
			if !strings.HasSuffix(baseURL, "/") {
				baseURL += "/"
			}
			results = append(results, baseURL+path)
		}
	}
	return results
}
