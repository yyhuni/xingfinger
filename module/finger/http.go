// Package finger 提供 Web 指纹识别核心功能
package finger

import (
	"crypto/tls"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// resps HTTP 响应结构体
// 存储解析后的 HTTP 响应信息，用于后续指纹匹配
type resps struct {
	url        string              // 请求的 URL
	body       string              // 响应体内容（已转换为 UTF-8）
	header     map[string][]string // 响应头
	server     string              // 服务器类型
	statuscode int                 // HTTP 状态码
	length     int                 // 响应体长度
	title      string              // 页面标题
	jsurl      []string            // JS 跳转发现的 URL 列表
	favhash    string              // Favicon 的 Murmur3 hash 值
}

// userAgents 常用浏览器 User-Agent 列表
// 随机选择以模拟真实浏览器访问，降低被识别为爬虫的风险
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/96.0.4664.110 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 Chrome/97.0.4692.71 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/97.0.4692.71 Safari/537.36",
}

// randomUA 随机返回一个 User-Agent
func randomUA() string {
	return userAgents[rand.Intn(len(userAgents))]
}

// extractTitle 从 HTML 内容中提取页面标题
// 使用 goquery 解析 HTML，查找 <title> 标签内容
//
// 参数：
//   - body: HTML 内容
//
// 返回：
//   - 页面标题，去除换行和首尾空格
func extractTitle(body string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return ""
	}
	title := doc.Find("title").Text()
	title = strings.ReplaceAll(title, "\n", "")
	return strings.TrimSpace(title)
}

// extractFavicon 提取并计算 Favicon 的 hash 值
// 首先尝试从 HTML 中解析 favicon 路径，如果没有则使用默认路径 /favicon.ico
//
// 参数：
//   - body: HTML 内容
//   - targetURL: 目标 URL，用于构建完整的 favicon URL
//
// 返回：
//   - Favicon 的 Murmur3 hash 值
func extractFavicon(body, targetURL string) string {
	// 从 HTML 中提取 favicon 路径
	paths := extractRegex(`href="(.*?favicon....)"`, body)
	u, _ := url.Parse(targetURL)
	baseURL := u.Scheme + "://" + u.Host

	var faviconURL string
	if len(paths) > 0 {
		fav := paths[0][1]
		// 处理不同格式的 favicon 路径
		if strings.HasPrefix(fav, "//") {
			// 协议相对路径
			faviconURL = "http:" + fav
		} else if strings.HasPrefix(fav, "http") {
			// 完整 URL
			faviconURL = fav
		} else {
			// 相对路径
			faviconURL = baseURL + "/" + fav
		}
	} else {
		// 使用默认 favicon 路径
		faviconURL = baseURL + "/favicon.ico"
	}
	return calcFaviconHash(faviconURL)
}

// doHTTPRequest 发送 HTTP 请求并解析响应
// 支持 HTTPS（跳过证书验证）和代理设置
//
// 参数：
//   - urlData: [url, flag] 数组，flag 用于标识是否为 JS 跳转发现的 URL
//   - proxy: 代理服务器地址，为空则不使用代理
//
// 返回：
//   - resps: 解析后的响应结构体
//   - error: 错误信息
func doHTTPRequest(urlData []string, proxy string) (*resps, error) {
	// 配置 HTTP 传输层
	transport := &http.Transport{
		// 跳过 TLS 证书验证，允许访问自签名证书的站点
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// 设置代理
	if proxy != "" {
		proxyURL, _ := url.Parse(proxy)
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	// 创建 HTTP 客户端
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}

	// 构建请求
	req, err := http.NewRequest("GET", urlData[0], nil)
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

	// 读取响应体
	body, _ := ioutil.ReadAll(resp.Body)

	// 转换编码为 UTF-8
	contentType := strings.ToLower(resp.Header.Get("Content-Type"))
	bodyStr := toUtf8(string(body), contentType)

	// 提取服务器信息
	server := "None"
	if s := resp.Header.Get("Server"); s != "" {
		server = s
	} else if p := resp.Header.Get("X-Powered-By"); p != "" {
		server = p
	}

	// 解析 JS 跳转（仅对初始 URL 进行解析，避免无限循环）
	var jsURLs []string
	if urlData[1] == "0" {
		jsURLs = parseJSRedirect(bodyStr, urlData[0])
	}

	return &resps{
		url:        urlData[0],
		body:       bodyStr,
		header:     resp.Header,
		server:     server,
		statuscode: resp.StatusCode,
		length:     len(bodyStr),
		title:      extractTitle(bodyStr),
		jsurl:      jsURLs,
		favhash:    extractFavicon(bodyStr, urlData[0]),
	}, nil
}
