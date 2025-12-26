// Package finger 提供 Web 指纹识别核心功能
package finger

import (
	"regexp"
	"strings"
)

// matchKeyword 关键字匹配
// 检查内容中是否包含所有指定的关键字（AND 逻辑）
// 所有关键字都必须存在才返回 true
//
// 参数：
//   - content: 待匹配的内容（响应体/响应头/标题）
//   - keywords: 关键字列表
//
// 返回：
//   - true: 所有关键字都存在
//   - false: 有任意关键字不存在
func matchKeyword(content string, keywords []string) bool {
	for _, kw := range keywords {
		if !strings.Contains(content, kw) {
			return false
		}
	}
	return true
}

// matchRegex 正则表达式匹配
// 检查内容是否匹配所有指定的正则表达式（AND 逻辑）
// 所有正则都必须匹配才返回 true
//
// 参数：
//   - content: 待匹配的内容
//   - patterns: 正则表达式列表
//
// 返回：
//   - true: 所有正则都匹配
//   - false: 有任意正则不匹配
func matchRegex(content string, patterns []string) bool {
	for _, p := range patterns {
		re := regexp.MustCompile(p)
		if !re.MatchString(content) {
			return false
		}
	}
	return true
}
