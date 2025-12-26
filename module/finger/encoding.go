// Package finger 提供 Web 指纹识别核心功能
// 本文件负责 HTTP 响应的字符编码检测和转换
// 支持 GBK、GB2312、GB18030、Big5、UTF-8 等常见编码
package finger

import (
	"regexp"
	"strings"

	"github.com/yinheli/mahonia"
	"golang.org/x/net/html/charset"
)

// decodeToUTF8 将 HTTP 响应内容转换为 UTF-8 编码
// 按以下优先级检测编码：
// 1. Content-Type 响应头中的 charset
// 2. HTML meta 标签中的 charset
// 3. 通过 title 标签内容自动检测
//
// 参数：
//   - content: 原始响应内容
//   - contentType: Content-Type 响应头值
//
// 返回：
//   - 转换为 UTF-8 的内容
func decodeToUTF8(content, contentType string) string {
	// 从 Content-Type 检测编码
	encoding := detectEncoding(contentType)

	// 尝试从 meta 标签获取编码
	if enc := extractMetaCharset(content); enc != "" {
		encoding = enc
	}

	// 如果仍为 UTF-8，尝试通过 title 内容检测
	if enc := detectTitleEncoding(content); enc != "" && encoding == "utf-8" {
		encoding = enc
	}

	// 如果不是 UTF-8，进行编码转换
	if encoding != "" && encoding != "utf-8" {
		return convertEncoding(content, encoding, "utf-8")
	}
	return content
}

// detectEncoding 从 Content-Type 字符串检测编码
// 将各种编码名称标准化为统一格式
//
// 参数：
//   - contentType: Content-Type 响应头值或编码名称
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
		// 中文 GBK 系列编码统一使用 gb18030（兼容性最好）
		return "gb18030"
	case strings.Contains(contentType, "big5"):
		// 繁体中文 Big5 编码
		return "big5"
	case strings.Contains(contentType, "utf-8"):
		return "utf-8"
	}
	// 默认假设为 gb18030（中文网站常见）
	return "gb18030"
}

// extractMetaCharset 从 HTML meta 标签提取 charset 声明
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

// detectTitleEncoding 通过 title 标签内容自动检测编码
// 使用 golang.org/x/net/html/charset 库进行编码检测
//
// 参数：
//   - content: HTML 内容
//
// 返回：
//   - 检测到的编码名称，未找到返回空字符串
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
	// 编码相同，无需转换
	if from == to {
		return src
	}

	// 先解码为中间格式
	decoder := mahonia.NewDecoder(from)
	result := decoder.ConvertString(src)

	// 再编码为目标格式
	encoder := mahonia.NewDecoder(to)
	_, data, _ := encoder.Translate([]byte(result), true)
	return string(data)
}
