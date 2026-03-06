// Package pkg 提供 xingfinger 的核心功能
// 本文件负责解析页面中的 JS 跳转
package pkg

import (
	"net/url"
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
		`(window|top)\.location\.href\s*=\s*['"](.*?)['"]`, // window.location.href 跳转
		`redirectUrl\s*=\s*['"](.*?)['"]`,                  // redirectUrl 变量赋值
		`<meta.*?http-equiv=.*?refresh.*?url=(.*?)>`,       // meta refresh 跳转
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return nil
	}

	var results []string
	seen := make(map[string]bool)
	for _, p := range patterns {
		matches := extractRegex(p, body)
		for _, m := range matches {
			if len(m) == 0 {
				continue
			}

			redirectPath := strings.TrimSpace(m[len(m)-1])
			redirectPath = strings.Trim(redirectPath, `"'`)
			if redirectPath == "" {
				continue
			}

			ref, err := url.Parse(redirectPath)
			if err != nil {
				continue
			}

			resolved := base.ResolveReference(ref).String()
			if !seen[resolved] {
				seen[resolved] = true
				results = append(results, resolved)
			}
		}
	}
	return results
}
