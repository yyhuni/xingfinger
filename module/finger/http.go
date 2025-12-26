// Package finger 提供 Web 指纹识别核心功能
package finger

import (
	"crypto/tls"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Timeout 请求超时时间（秒）
var Timeout = 10

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
func extractTitle(body string) string {
	// 先尝试用正则提取，更可靠
	re := regexp.MustCompile(`(?is)<title[^>]*>(.*?)</title>`)
	match := re.FindStringSubmatch(body)
	if len(match) > 1 {
		title := strings.TrimSpace(match[1])
		title = strings.ReplaceAll(title, "\n", "")
		title = strings.ReplaceAll(title, "\r", "")
		title = strings.ReplaceAll(title, "\t", "")
		return title
	}

	// 备用方案：使用 goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return ""
	}
	title := doc.Find("title").Text()
	title = strings.ReplaceAll(title, "\n", "")
	return strings.TrimSpace(title)
}

// extractFavicon 提取并计算 Favicon 的 hash 值
func extractFavicon(body, targetURL string) string {
	paths := extractRegex(`href="(.*?favicon....)"`, body)
	u, _ := url.Parse(targetURL)
	baseURL := u.Scheme + "://" + u.Host

	var faviconURL string
	if len(paths) > 0 {
		fav := paths[0][1]
		if strings.HasPrefix(fav, "//") {
			faviconURL = "http:" + fav
		} else if strings.HasPrefix(fav, "http") {
			faviconURL = fav
		} else {
			faviconURL = baseURL + "/" + fav
		}
	} else {
		faviconURL = baseURL + "/favicon.ico"
	}
	return calcFaviconHash(faviconURL)
}

// doHTTPRequest 发送 HTTP 请求并解析响应
func doHTTPRequest(urlData []string, proxy string) (*resps, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if proxy != "" {
		proxyURL, _ := url.Parse(proxy)
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	client := &http.Client{
		Timeout:   time.Duration(Timeout) * time.Second,
		Transport: transport,
	}

	req, err := http.NewRequest("GET", urlData[0], nil)
	if err != nil {
		return nil, err
	}

	req.AddCookie(&http.Cookie{Name: "rememberMe", Value: "me"})
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "close")
	req.Header.Set("User-Agent", randomUA())

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	contentType := strings.ToLower(resp.Header.Get("Content-Type"))
	bodyStr := toUtf8(string(body), contentType)

	server := "None"
	if s := resp.Header.Get("Server"); s != "" {
		server = s
	} else if p := resp.Header.Get("X-Powered-By"); p != "" {
		server = p
	}

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
