// Package pkg 提供 xingfinger 的核心功能
// 本文件是指纹扫描器的核心实现，负责：
// 1. 初始化 fingers 指纹识别引擎
// 2. 并发扫描目标 URL
// 3. 收集和输出扫描结果
package pkg

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/chainreactors/fingers"
	"github.com/gookit/color"
)

// Result 扫描结果结构体
// 保存单个 URL 的扫描结果，用于输出和 JSON 导出
type Result struct {
	URL        string `json:"url"`         // 目标 URL
	CMS        string `json:"cms"`         // 检测到的 CMS/框架，多个用逗号分隔
	Server     string `json:"server"`      // 服务器信息
	StatusCode int    `json:"status_code"` // HTTP 状态码
	Length     int    `json:"length"`      // 响应体长度
	Title      string `json:"title"`       // 页面标题
}

// Scanner 指纹扫描器
// 负责管理扫描任务队列、并发控制和结果收集
type Scanner struct {
	queue      *Queue          // URL 任务队列
	wg         sync.WaitGroup  // 等待组，用于同步所有扫描 goroutine
	mu         sync.Mutex      // 互斥锁，保护结果切片的并发写入
	thread     int             // 并发线程数
	output     string          // 输出文件路径
	proxy      string          // 代理地址
	silent     bool            // 静默模式，只输出命中结果
	allResults []Result        // 所有扫描结果
	hitResults []Result        // 命中指纹的结果
	engine     *fingers.Engine // fingers 指纹识别引擎
	engines    []string        // 启用的指纹引擎列表
}

// NewScanner 创建扫描器实例
// 初始化 fingers 引擎和任务队列
//
// 参数：
//   - urls: 待扫描的 URL 列表
//   - thread: 并发线程数
//   - output: 输出文件路径，为空则不保存
//   - proxy: 代理地址，为空则不使用代理
//   - timeout: HTTP 请求超时时间（秒）
//   - silent: 是否启用静默模式
//   - customConfig: 自定义指纹配置
//
// 返回：
//   - *Scanner: 扫描器实例
func NewScanner(urls []string, thread int, output, proxy string, timeout int, silent bool, customConfig *CustomFingerConfig) *Scanner {
	// 确定要启用的指纹引擎
	// 如果有自定义指纹，只启用对应的引擎
	var enableEngines []string

	if customConfig != nil && (customConfig.EHole != "" || customConfig.Goby != "" ||
		customConfig.Wappalyzer != "" || customConfig.Fingers != "" || customConfig.FingerPrint != "") {
		// 有自定义指纹，只启用对应的引擎
		if customConfig.EHole != "" {
			enableEngines = append(enableEngines, "ehole")
		}
		if customConfig.Goby != "" {
			enableEngines = append(enableEngines, "goby")
		}
		if customConfig.Wappalyzer != "" {
			enableEngines = append(enableEngines, "wappalyzer")
		}
		if customConfig.Fingers != "" {
			enableEngines = append(enableEngines, "fingers")
		}
		if customConfig.FingerPrint != "" {
			enableEngines = append(enableEngines, "fingerprinthub")
		}
		// 始终启用 favicon 引擎
		enableEngines = append(enableEngines, "favicon")
	}
	// 如果没有自定义指纹，enableEngines 为空，NewEngine 会使用默认引擎

	// 初始化 fingers 指纹识别引擎
	// fingers 引擎聚合了多个指纹库：fingers、wappalyzer、fingerprinthub、ehole、goby
	// 静默模式下抑制库的加载信息输出
	var engine *fingers.Engine
	var err error
	if silent {
		// 临时重定向标准输出以抑制库的打印信息
		oldStdout := os.Stdout
		os.Stdout, _ = os.Open(os.DevNull)
		if len(enableEngines) > 0 {
			engine, err = fingers.NewEngine(enableEngines...)
		} else {
			engine, err = fingers.NewEngine()
		}
		os.Stdout = oldStdout
	} else {
		if len(enableEngines) > 0 {
			engine, err = fingers.NewEngine(enableEngines...)
		} else {
			engine, err = fingers.NewEngine()
		}
	}
	if err != nil {
		fmt.Printf("[!] 初始化指纹引擎失败: %v\n", err)
		os.Exit(1)
	}

	// 创建扫描器实例
	s := &Scanner{
		queue:      NewQueue(),
		thread:     thread,
		output:     output,
		proxy:      proxy,
		silent:     silent,
		allResults: []Result{},
		hitResults: []Result{},
		engine:     engine,
		engines:    enableEngines,
	}

	// 设置 HTTP 请求超时时间
	Timeout = timeout

	// 将 URL 添加到任务队列
	// task[0] 为 URL，task[1] 为任务类型（"0" 表示主页面）
	for _, url := range urls {
		s.queue.Push([]string{url, "0"})
	}

	return s
}

// Run 启动扫描
// 创建多个 goroutine 并发执行扫描任务
func (s *Scanner) Run() {
	// 启动工作 goroutine
	for i := 0; i < s.thread; i++ {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.scan()
		}()
	}

	// 等待所有 goroutine 完成
	s.wg.Wait()

	// 输出扫描统计（非静默模式）
	if !s.silent {
		color.RGBStyleFromString("244,211,49").Printf("\n[+] Scanned: %d, Matched: %d\n", len(s.allResults), len(s.hitResults))
	}

	// 保存结果到文件
	if s.output != "" {
		saveResults(s.output, s.allResults)
	}
}

// detectFingerprints 使用 fingers 引擎检测指纹
// 将原始 HTTP 响应传递给 fingers 引擎进行多指纹库匹配
//
// 参数：
//   - rawContent: 原始 HTTP 响应内容（包含 header 和 body）
//
// 返回：
//   - []string: 检测到的框架名称列表
func (s *Scanner) detectFingerprints(rawContent []byte) []string {
	// 使用 DetectContent 进行指纹检测
	// 该方法会调用所有启用的指纹引擎进行匹配
	frameworks, err := s.engine.DetectContent(rawContent)
	if err != nil {
		return nil
	}

	// 提取框架名称
	// GetNames() 返回所有非猜测（non-guess）框架的名称
	return frameworks.GetNames()
}

// detectFavicon 使用 fingers 引擎检测 favicon 指纹
// 主动获取 favicon 文件并进行 hash 匹配
//
// 参数：
//   - body: HTML 响应体，用于提取 favicon URL
//   - baseURL: 当前页面 URL
//
// 返回：
//   - string: 检测到的框架名称，未匹配返回空字符串
func (s *Scanner) detectFavicon(body, baseURL string) string {
	// 提取 favicon URL
	faviconURL := extractFaviconURL(body, baseURL)
	if faviconURL == "" {
		return ""
	}

	// 获取 favicon 内容
	faviconContent, err := fetchFavicon(faviconURL, s.proxy)
	if err != nil || len(faviconContent) == 0 {
		return ""
	}

	// 使用 fingers 引擎的 favicon 检测
	// MatchFavicon 会计算 MD5 和 MMH3 hash 并匹配指纹库
	frameworks := s.engine.MatchFavicon(faviconContent)
	if len(frameworks) == 0 {
		return ""
	}

	// 返回第一个匹配的框架名称
	names := frameworks.GetNames()
	if len(names) > 0 {
		return names[0]
	}
	return ""
}

// scan 执行扫描任务
// 从队列中获取 URL，发送请求，进行指纹检测，输出结果
func (s *Scanner) scan() {
	for {
		// 从队列获取任务
		item := s.queue.Pop()
		if item == nil {
			return
		}

		task, ok := item.([]string)
		if !ok {
			continue
		}

		// 发送 HTTP 请求
		resp, err := fetch(task, s.proxy)
		if err != nil {
			// 如果 HTTPS 失败，尝试 HTTP
			task[0] = strings.ReplaceAll(task[0], "https://", "http://")
			resp, err = fetch(task, s.proxy)
			if err != nil {
				continue
			}
		}

		// 处理 JS 跳转
		// 将 JS 跳转的 URL 添加到队列继续扫描
		for _, jsURL := range resp.JsURLs {
			if jsURL != "" {
				s.queue.Push([]string{jsURL, "1"})
			}
		}

		// 使用 fingers 引擎进行指纹检测
		matched := s.detectFingerprints(resp.RawContent)

		// 主动获取 favicon 进行指纹检测（仅对主页面）
		if task[1] == "0" {
			if faviconMatch := s.detectFavicon(resp.Body, resp.URL); faviconMatch != "" {
				// 检查是否已存在，避免重复
				exists := false
				for _, m := range matched {
					if m == faviconMatch {
						exists = true
						break
					}
				}
				if !exists {
					matched = append(matched, faviconMatch)
				}
			}
		}

		// 构建扫描结果
		result := Result{
			URL:        resp.URL,
			CMS:        strings.Join(matched, ","),
			Server:     resp.Server,
			StatusCode: resp.StatusCode,
			Length:     resp.Length,
			Title:      resp.Title,
		}

		// 保存结果（线程安全）
		s.mu.Lock()
		s.allResults = append(s.allResults, result)
		if result.CMS != "" {
			s.hitResults = append(s.hitResults, result)
		}
		s.mu.Unlock()

		// 输出结果
		s.printResult(result)
	}
}

// printResult 输出扫描结果
// 根据模式选择不同的输出格式
//
// 参数：
//   - result: 扫描结果
func (s *Scanner) printResult(result Result) {
	// 静默模式：只输出命中指纹的结果
	if s.silent {
		if result.CMS != "" {
			fmt.Printf("%s [%s]\n", result.URL, result.CMS)
		}
		return
	}

	// 正常模式：httpx 风格输出
	var parts []string
	parts = append(parts, result.URL)
	parts = append(parts, fmt.Sprintf("[%d]", result.StatusCode))
	parts = append(parts, fmt.Sprintf("[%d]", result.Length))
	if result.Server != "" {
		parts = append(parts, fmt.Sprintf("[%s]", result.Server))
	}
	if result.Title != "" {
		parts = append(parts, fmt.Sprintf("[%s]", result.Title))
	}
	if result.CMS != "" {
		parts = append(parts, fmt.Sprintf("[%s]", result.CMS))
	}

	line := strings.Join(parts, " ")

	// 命中指纹的结果用红色高亮显示
	if result.CMS != "" {
		color.RGBStyleFromString("237,64,35").Println(line)
	} else {
		fmt.Println(line)
	}
}
