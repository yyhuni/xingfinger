// Package pkg 提供 xingfinger 的核心功能
// 本文件负责解析页面中的 JS 跳转
package pkg

import (
	"net/url"
	"regexp"
	"strings"
)

var htmlAttrPattern = regexp.MustCompile(`(?is)([a-zA-Z_:][a-zA-Z0-9_:.\-]*)\s*=\s*(?:"([^"]*)"|'([^']*)'|([^\s>]+))`)

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
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return nil
	}

	var results []string
	seen := make(map[string]bool)
	addRedirect := func(redirectPath string) {
		redirectPath = strings.TrimSpace(redirectPath)
		redirectPath = strings.Trim(redirectPath, `"'`)
		if redirectPath == "" {
			return
		}

		ref, err := url.Parse(redirectPath)
		if err != nil {
			return
		}

		resolved := base.ResolveReference(ref).String()
		if !seen[resolved] {
			seen[resolved] = true
			results = append(results, resolved)
		}
	}

	for _, p := range patterns {
		matches := extractRegex(p, body)
		for _, m := range matches {
			if len(m) == 0 {
				continue
			}

			addRedirect(m[len(m)-1])
		}
	}

	for _, redirectPath := range parseMetaRefreshRedirects(body) {
		addRedirect(redirectPath)
	}

	return results
}

func parseMetaRefreshRedirects(body string) []string {
	matches := extractRegex(`(?is)<meta\b[^>]*>`, body)
	if len(matches) == 0 {
		return nil
	}

	var redirects []string
	for _, match := range matches {
		if len(match) == 0 {
			continue
		}

		attrs := extractHTMLAttributes(match[0])
		if !strings.EqualFold(strings.TrimSpace(attrs["http-equiv"]), "refresh") {
			continue
		}

		redirectPath := extractMetaRefreshURL(attrs["content"])
		if redirectPath == "" {
			continue
		}

		redirects = append(redirects, redirectPath)
	}

	return redirects
}

func extractHTMLAttributes(tag string) map[string]string {
	attrs := make(map[string]string)
	matches := htmlAttrPattern.FindAllStringSubmatch(tag, -1)
	for _, match := range matches {
		if len(match) < 5 {
			continue
		}

		value := match[2]
		if value == "" {
			value = match[3]
		}
		if value == "" {
			value = match[4]
		}

		attrs[strings.ToLower(match[1])] = value
	}
	return attrs
}

func extractMetaRefreshURL(content string) string {
	match := regexp.MustCompile(`(?i)(?:^|;)\s*url\s*=\s*([^;]+)`).FindStringSubmatch(content)
	if len(match) < 2 {
		return ""
	}

	return strings.TrimSpace(strings.Trim(match[1], `"'`))
}
