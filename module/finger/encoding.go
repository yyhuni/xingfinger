// Package finger 提供 Web 指纹识别核心功能
package finger

import (
	"regexp"
	"strings"

	"github.com/yinheli/mahonia"
	"golang.org/x/net/html/charset"
)

// toUtf8 将内容转换为 UTF-8 编码
// 自动检测原始编码并进行转换，支持 GBK、GB2312、GB18030、Big5 等编码
//
// 检测优先级：
//  1. Content-Type 头中的 charset
//  2. HTML meta 标签中的 charset
//  3. 根据 title 内容自动检测
//
// 参数：
//   - content: 原始内容
//   - contentType: HTTP Content-Type 头的值
//
// 返回：
//   - 转换为 UTF-8 后的内容
func toUtf8(content, contentType string) string {
	// 从 Content-Type 检测编码
	encoding := detectEncoding(contentType)

	// 从 meta 标签检测编码
	if enc := extractMetaCharset(content); enc != "" {
		encoding = enc
	}

	// 从 title 内容检测编码（作为补充检测）
	if enc := detectTitleEncoding(content); enc != "" && encoding == "utf-8" {
		encoding = enc
	}

	// 如果不是 UTF-8，进行编码转换
	if encoding != "" && encoding != "utf-8" {
		return convertEncoding(content, encoding, "utf-8")
	}
	return content
}

// detectEncoding 根据 Content-Type 字符串检测编码
// 将各种中文编码统一映射到 gb18030（兼容 GBK 和 GB2312）
//
// 参数：
//   - contentType: Content-Type 头或 charset 值
//
// 返回：
//   - 标准化的编码名称
func detectEncoding(contentType string) string {
	contentType = strings.ToLower(contentType)
	switch {
	case strings.Contains(contentType, "gbk"),
		strings.Contains(contentType, "gb2312"),
		strings.Contains(contentType, "gb18030"),
		strings.Contains(contentType, "windows-1252"):
		return "gb18030"
	case strings.Contains(contentType, "big5"):
		return "big5"
	case strings.Contains(contentType, "utf-8"):
		return "utf-8"
	}
	// 默认使用 gb18030，因为大部分中文网站使用此编码
	return "gb18030"
}

// extractMetaCharset 从 HTML meta 标签中提取 charset
// 匹配格式如：<meta charset="utf-8"> 或 <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
//
// 参数：
//   - content: HTML 内容
//
// 返回：
//   - 检测到的编码名称，未找到返回空字符串
func extractMetaCharset(content string) string {
	re := regexp.MustCompile(`(?is)<meta[^>]*charset\s*=["']?\s*([A-Za-z0-9\-]+)`)
	match := re.FindStringSubmatch(content)
	if len(match) > 1 {
		return detectEncoding(match[1])
	}
	return ""
}

// detectTitleEncoding 通过分析 title 内容检测编码
// 使用 golang.org/x/net/html/charset 包进行自动检测
//
// 参数：
//   - content: HTML 内容
//
// 返回：
//   - 检测到的编码名称
func detectTitleEncoding(content string) string {
	re := regexp.MustCompile(`(?is)<title[^>]*>(.*?)<\/title>`)
	match := re.FindStringSubmatch(content)
	if len(match) > 1 {
		_, enc, _ := charset.DetermineEncoding([]byte(match[1]), "")
		return detectEncoding(enc)
	}
	return ""
}

// convertEncoding 转换字符串编码
// 使用 mahonia 库进行编码转换
//
// 参数：
//   - src: 源字符串
//   - from: 源编码
//   - to: 目标编码
//
// 返回：
//   - 转换后的字符串
func convertEncoding(src, from, to string) string {
	if from == to {
		return src
	}
	// 先解码
	decoder := mahonia.NewDecoder(from)
	result := decoder.ConvertString(src)
	// 再编码
	encoder := mahonia.NewDecoder(to)
	_, data, _ := encoder.Translate([]byte(result), true)
	return string(data)
}
