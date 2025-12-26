// Package finger 提供 Web 指纹识别核心功能
// 包含扫描器、指纹匹配、结果处理等核心逻辑
package finger

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/gookit/color"
	"github.com/yyhuni/xingfinger/module/finger/source"
	"github.com/yyhuni/xingfinger/module/queue"
)

// Result 扫描结果结构体
// 存储单个 URL 的指纹识别结果
type Result struct {
	URL        string `json:"url"`         // 目标 URL
	CMS        string `json:"cms"`         // 识别到的 CMS/框架名称，多个用逗号分隔
	Server     string `json:"server"`      // 服务器类型（来自 Server 或 X-Powered-By 头）
	StatusCode int    `json:"status_code"` // HTTP 响应状态码
	Length     int    `json:"length"`      // 响应内容长度
	Title      string `json:"title"`       // 页面标题
}

// Scanner 指纹扫描器
// 管理扫描任务的执行，包括并发控制、结果收集等
type Scanner struct {
	queue        *queue.Queue // 任务队列，存储待扫描的 URL
	wg           sync.WaitGroup
	thread       int           // 并发线程数
	output       string        // 输出文件路径
	proxy        string        // 代理服务器地址
	allResults   []Result      // 所有扫描结果
	hitResults   []Result      // 命中指纹的结果
	fingerprints *Packjson     // 指纹规则库
}

// NewScanner 创建新的扫描器实例
//
// 参数：
//   - urls: 待扫描的 URL 列表
//   - thread: 并发线程数
//   - output: 输出文件路径（可为空）
//   - proxy: 代理服务器地址（可为空）
//
// 返回：
//   - 初始化完成的扫描器实例
func NewScanner(urls []string, thread int, output, proxy string) *Scanner {
	s := &Scanner{
		queue:      queue.NewQueue(),
		thread:     thread,
		output:     output,
		proxy:      proxy,
		allResults: []Result{},
		hitResults: []Result{},
	}

	// 加载指纹规则库
	err := LoadWebfingerprint(source.GetExePath() + "/finger.json")
	if err != nil {
		fmt.Println("[!] Fingerprint file error")
		os.Exit(1)
	}
	s.fingerprints = GetWebfingerprint()

	// 将所有 URL 加入任务队列
	// 队列元素格式：[url, flag]，flag=0 表示初始 URL，flag=1 表示 JS 跳转发现的 URL
	for _, url := range urls {
		s.queue.Push([]string{url, "0"})
	}
	return s
}

// Run 启动扫描任务
// 创建多个 goroutine 并发执行扫描，等待所有任务完成后输出结果
func (s *Scanner) Run() {
	// 启动工作线程
	for i := 0; i <= s.thread; i++ {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.scan()
		}()
	}

	// 等待所有线程完成
	s.wg.Wait()

	// 输出命中指纹的结果汇总
	color.RGBStyleFromString("244,211,49").Println("\n[+] Matched:")
	for _, r := range s.hitResults {
		fmt.Printf("[ %s | ", r.URL)
		color.RGBStyleFromString("237,64,35").Printf("%s", r.CMS)
		fmt.Printf(" | %s | %d | %d | %s ]\n", r.Server, r.StatusCode, r.Length, r.Title)
	}

	// 如果指定了输出文件，保存结果
	if s.output != "" {
		saveResults(s.output, s.allResults)
	}
}

// scan 执行单个工作线程的扫描逻辑
// 从队列中获取 URL，发送 HTTP 请求，进行指纹匹配
func (s *Scanner) scan() {
	for s.queue.Len() != 0 {
		// 从队列获取任务
		item := s.queue.Pop()
		urlData, ok := item.([]string)
		if !ok {
			continue
		}

		// 发送 HTTP 请求获取响应
		resp, err := doHTTPRequest(urlData, s.proxy)
		if err != nil {
			// HTTPS 失败时尝试 HTTP
			urlData[0] = strings.ReplaceAll(urlData[0], "https://", "http://")
			resp, err = doHTTPRequest(urlData, s.proxy)
			if err != nil {
				continue
			}
		}

		// 将 JS 跳转发现的 URL 加入队列
		for _, jurl := range resp.jsurl {
			if jurl != "" {
				s.queue.Push([]string{jurl, "1"})
			}
		}

		// 将响应头转为 JSON 字符串，便于指纹匹配
		headers := toJSON(resp.header)
		var matched []string

		// 遍历指纹规则进行匹配
		for _, fp := range s.fingerprints.Fingerprint {
			// 根据指纹位置选择匹配目标
			var target string
			switch fp.Location {
			case "body":
				target = resp.body
			case "header":
				target = headers
			case "title":
				target = resp.title
			default:
				continue
			}

			// 根据匹配方法执行匹配
			switch fp.Method {
			case "keyword":
				// 关键字匹配：所有关键字都必须存在
				if matchKeyword(target, fp.Keyword) {
					matched = append(matched, fp.Cms)
				}
			case "faviconhash":
				// Favicon hash 匹配
				if resp.favhash == fp.Keyword[0] {
					matched = append(matched, fp.Cms)
				}
			case "regular":
				// 正则表达式匹配
				if matchRegex(target, fp.Keyword) {
					matched = append(matched, fp.Cms)
				}
			}
		}

		// 去重并构建结果
		matched = unique(matched)
		result := Result{
			URL:        resp.url,
			CMS:        strings.Join(matched, ","),
			Server:     resp.server,
			StatusCode: resp.statuscode,
			Length:     resp.length,
			Title:      resp.title,
		}

		s.allResults = append(s.allResults, result)

		// 格式化输出
		line := fmt.Sprintf("[ %s | %s | %s | %d | %d | %s ]",
			result.URL, result.CMS, result.Server, result.StatusCode, result.Length, result.Title)

		// 命中指纹的结果用红色高亮显示
		if result.CMS != "" {
			color.RGBStyleFromString("237,64,35").Println(line)
			s.hitResults = append(s.hitResults, result)
		} else {
			fmt.Println(line)
		}
	}
}

// toJSON 将 map 转换为 JSON 字符串
// 用于将 HTTP 响应头转换为可搜索的字符串格式
func toJSON(data map[string][]string) string {
	b, _ := json.Marshal(data)
	return string(b)
}

// unique 对字符串切片去重
// 同时过滤空字符串
func unique(arr []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, v := range arr {
		if v != "" && !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}
