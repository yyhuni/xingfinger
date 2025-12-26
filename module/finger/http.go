// Package finger 提供 Web 指纹识别核心功能
// 本文件负责 HTTP 请求发送和响应解析
package finger

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// Timeout 请求超时时间（秒）
// 可通过命令行参数修改
var Timeout = 10

// Response HTTP 响应结构体
// 包含 HTTP 响应的所有关键信息，供指纹识别使用
type Response struct {
	URL        string              // 请求的目标 URL
	RawContent []byte              // 原始 HTTP 响应内容（包含 header 和 body），供 fingers 引擎使用
	Body       string              // 响应体内容（已解码为 UTF-8）
	Header     map[string][]string // 响应头
	Server     string              // 服务器信息（从 Server 或 X-Powered-By 头获取）
	StatusCode int                 // HTTP 状态码
	Length     int                 // 响应体长度
	Title      string              // 页面标题（从 <title> 标签提取）
	JsURLs     []string            // JS 跳转 URL 列表
}

// userAgents 常用浏览器 User-Agent 列表
// 用于模拟真实浏览器请求，避免被目标服务器拦截
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/96.0.4664.110 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 Chrome/97.0.4692.71 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/97.0.4692.71 Safari/537.36",
}

// randomUA 随机返回一个 User-Agent
// 每次请求使用不同的 UA，降低被识别为扫描器的风险
func randomUA() string {
	return userAgents[rand.Intn(len(userAgents))]
}

// extractTitle 从 HTML 提取页面标题
// 使用正则表达式匹配 <title> 标签内容
//
// 参数：
//   - body: HTML 响应体
//
// 返回：
//   - 页面标题，如果未找到则返回空字符串
func extractTitle(body string) string {
	re := regexp.MustCompile(`(?is)<title[^>]*>(.*?)</title>`)
	match := re.FindStringSubmatch(body)
	if len(match) > 1 {
		title := strings.TrimSpace(match[1])
		// 清理标题中的换行和制表符
		title = strings.ReplaceAll(title, "\n", "")
		title = strings.ReplaceAll(title, "\r", "")
		title = strings.ReplaceAll(title, "\t", "")
		return title
	}
	return ""
}

// buildRawResponse 构建原始 HTTP 响应内容
// 将 HTTP 响应转换为原始格式，供 fingers 引擎的 DetectContent 方法使用
//
// 参数：
//   - resp: HTTP 响应对象
//   - body: 响应体内容
//
// 返回：
//   - 原始 HTTP 响应的字节数组
func buildRawResponse(resp *http.Response, body []byte) []byte {
	var buf bytes.Buffer

	// 写入状态行
	buf.WriteString(fmt.Sprintf("HTTP/%d.%d %d %s\r\n",
		resp.ProtoMajor, resp.ProtoMinor, resp.StatusCode, http.StatusText(resp.StatusCode)))

	// 写入响应头
	for key, values := range resp.Header {
		for _, value := range values {
			buf.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
		}
	}

	// 写入空行分隔头和体
	buf.WriteString("\r\n")

	// 写入响应体
	buf.Write(body)

	return buf.Bytes()
}

// fetch 发送 HTTP 请求并解析响应
// 这是核心的 HTTP 请求函数，负责：
// 1. 创建 HTTP 客户端（支持代理和 TLS）
// 2. 发送请求并读取响应
// 3. 解析响应内容（编码转换、标题提取等）
// 4. 构建原始响应供 fingers 引擎使用
//
// 参数：
//   - task: 任务数组，task[0] 为 URL，task[1] 为任务类型（"0" 表示主页面，"1" 表示 JS 跳转页面）
//   - proxy: 代理地址，为空则不使用代理
//
// 返回：
//   - *Response: 解析后的响应结构体
//   - error: 错误信息
func fetch(task []string, proxy string) (*Response, error) {
	// 创建 HTTP 传输层，跳过 TLS 证书验证
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// 配置代理
	if proxy != "" {
		proxyURL, _ := url.Parse(proxy)
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	// 创建 HTTP 客户端
	client := &http.Client{
		Timeout:   time.Duration(Timeout) * time.Second,
		Transport: transport,
	}

	// 创建请求
	req, err := http.NewRequest("GET", task[0], nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	// 添加 rememberMe cookie 用于检测 Shiro 框架
	req.AddCookie(&http.Cookie{Name: "rememberMe", Value: "me"})
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "close")
	req.Header.Set("User-Agent", randomUA())

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取原始响应体
	rawBody, _ := io.ReadAll(resp.Body)

	// 构建原始 HTTP 响应（供 fingers 引擎使用）
	rawContent := buildRawResponse(resp, rawBody)

	// 解码响应体为 UTF-8
	contentType := strings.ToLower(resp.Header.Get("Content-Type"))
	body := decodeToUTF8(string(rawBody), contentType)

	// 提取服务器信息
	server := ""
	if s := resp.Header.Get("Server"); s != "" {
		server = s
	} else if p := resp.Header.Get("X-Powered-By"); p != "" {
		server = p
	}

	// 解析 JS 跳转（仅对主页面进行）
	var jsURLs []string
	if task[1] == "0" {
		jsURLs = parseJSRedirect(body, task[0])
	}

	return &Response{
		URL:        task[0],
		RawContent: rawContent,
		Body:       body,
		Header:     resp.Header,
		Server:     server,
		StatusCode: resp.StatusCode,
		Length:     len(body),
		Title:      extractTitle(body),
		JsURLs:     jsURLs,
	}, nil
}

// extractFaviconURL 从 HTML 中提取 favicon URL
// 查找 <link rel="icon" href="xxx"> 或 <link rel="shortcut icon" href="xxx">
//
// 参数：
//   - body: HTML 响应体
//   - baseURL: 当前页面 URL，用于构建完整的 favicon URL
//
// 返回：
//   - favicon 的完整 URL
func extractFaviconURL(body, baseURL string) string {
	// 解析基础 URL
	u, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}
	base := u.Scheme + "://" + u.Host

	// 尝试从 HTML 中提取 favicon 路径
	// 匹配 <link rel="icon" href="xxx"> 或 <link rel="shortcut icon" href="xxx">
	patterns := []string{
		`<link[^>]*rel=["'](?:shortcut )?icon["'][^>]*href=["']([^"']+)["']`,
		`<link[^>]*href=["']([^"']+)["'][^>]*rel=["'](?:shortcut )?icon["']`,
		`href=["']([^"']*favicon[^"']*)["']`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		match := re.FindStringSubmatch(body)
		if len(match) > 1 {
			faviconPath := match[1]
			// 处理不同格式的路径
			if strings.HasPrefix(faviconPath, "//") {
				return "http:" + faviconPath
			} else if strings.HasPrefix(faviconPath, "http") {
				return faviconPath
			} else if strings.HasPrefix(faviconPath, "/") {
				return base + faviconPath
			} else {
				return base + "/" + faviconPath
			}
		}
	}

	// 默认使用 /favicon.ico
	return base + "/favicon.ico"
}

// fetchFavicon 获取 favicon 内容
// 发送 HTTP 请求获取 favicon 文件的原始字节内容
//
// 参数：
//   - faviconURL: favicon 的完整 URL
//   - proxy: 代理地址，为空则不使用代理
//
// 返回：
//   - []byte: favicon 文件内容
//   - error: 错误信息
func fetchFavicon(faviconURL, proxy string) ([]byte, error) {
	// 创建 HTTP 传输层
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// 配置代理
	if proxy != "" {
		proxyURL, _ := url.Parse(proxy)
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	// 创建 HTTP 客户端，使用较短的超时时间
	client := &http.Client{
		Timeout:   time.Duration(5) * time.Second,
		Transport: transport,
	}

	// 创建请求
	req, err := http.NewRequest("GET", faviconURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", randomUA())

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("favicon request failed: %d", resp.StatusCode)
	}

	// 读取响应体
	return io.ReadAll(resp.Body)
}
