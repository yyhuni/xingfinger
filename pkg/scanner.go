// Package pkg 提供 xingfinger 的核心功能
// 本文件是指纹扫描器的核心实现，负责：
// 1. 初始化 fingers 指纹识别引擎
// 2. 并发扫描目标 URL
// 3. 收集和输出扫描结果
package pkg

import (
	"encoding/json"
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
	queue        *Queue          // URL 任务队列
	wg           sync.WaitGroup  // 等待组，用于同步所有扫描 goroutine
	mu           sync.Mutex      // 互斥锁，保护结果切片的并发写入
	thread       int             // 并发线程数
	output       string          // 输出文件路径
	proxy        string          // 代理地址
	silent       bool            // 静默模式，只输出命中结果
	jsonOutput   bool            // JSON 格式输出到终端
	allResults   []Result        // 所有扫描结果
	hitResults   []Result        // 命中指纹的结果
	engine       *fingers.Engine // fingers 指纹识别引擎（默认指纹）
	customEngine *fingers.Engine // 自定义指纹引擎
	arlEngine    *ARLEngine      // ARL 指纹匹配引擎
	engines      []string        // 启用的指纹引擎列表
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
//   - jsonOutput: 是否以 JSON 格式输出到终端
//   - customConfig: 自定义指纹配置
//
// 返回：
//   - *Scanner: 扫描器实例
func NewScanner(urls []string, thread int, output, proxy string, timeout int, silent, jsonOutput bool, customConfig *CustomFingerConfig) *Scanner {
	// 检查是否禁用默认指纹
	noDefault := customConfig != nil && customConfig.NoDefault

	// 检查是否有自定义指纹文件
	hasCustomFingers := customConfig != nil && (customConfig.EHole != "" || customConfig.Goby != "" ||
		customConfig.Wappalyzer != "" || customConfig.Fingers != "" || customConfig.FingerPrint != "")

	var engine *fingers.Engine
	var customEngine *fingers.Engine
	var err error

	// 初始化默认指纹引擎（除非禁用）
	if !noDefault {
		if silent || jsonOutput {
			oldStdout := os.Stdout
			os.Stdout, _ = os.Open(os.DevNull)
			engine, err = fingers.NewEngine()
			os.Stdout = oldStdout
		} else {
			engine, err = fingers.NewEngine()
		}
		if err != nil {
			fmt.Printf("[!] 初始化默认指纹引擎失败: %v\n", err)
			os.Exit(1)
		}
	}

	// 初始化自定义指纹引擎（如果有自定义指纹）
	if hasCustomFingers {
		// 加载自定义指纹文件
		if err := LoadCustomFingerprints(customConfig, silent || jsonOutput); err != nil {
			fmt.Printf("[!] 加载自定义指纹失败: %v\n", err)
			os.Exit(1)
		}

		// 确定要启用的自定义引擎
		var customEngines []string
		if customConfig.EHole != "" {
			customEngines = append(customEngines, "ehole")
		}
		if customConfig.Goby != "" {
			customEngines = append(customEngines, "goby")
		}
		if customConfig.Wappalyzer != "" {
			customEngines = append(customEngines, "wappalyzer")
		}
		if customConfig.Fingers != "" {
			customEngines = append(customEngines, "fingers")
		}
		if customConfig.FingerPrint != "" {
			customEngines = append(customEngines, "fingerprinthub")
		}
		customEngines = append(customEngines, "favicon")

		if silent || jsonOutput {
			oldStdout := os.Stdout
			os.Stdout, _ = os.Open(os.DevNull)
			customEngine, err = fingers.NewEngine(customEngines...)
			os.Stdout = oldStdout
		} else {
			customEngine, err = fingers.NewEngine(customEngines...)
		}
		if err != nil {
			fmt.Printf("[!] 初始化自定义指纹引擎失败: %v\n", err)
			os.Exit(1)
		}
	} else if noDefault {
		// 只有 --no-default 但没有自定义指纹，需要清空默认数据
		if err := LoadCustomFingerprints(customConfig, silent || jsonOutput); err != nil {
			fmt.Printf("[!] 处理指纹配置失败: %v\n", err)
			os.Exit(1)
		}
	}

	// 创建扫描器实例
	s := &Scanner{
		queue:        NewQueue(),
		thread:       thread,
		output:       output,
		proxy:        proxy,
		silent:       silent,
		jsonOutput:   jsonOutput,
		allResults:   []Result{},
		hitResults:   []Result{},
		engine:       engine,
		customEngine: customEngine,
	}

	// 初始化 ARL 引擎（如果指定了 ARL 指纹文件）
	if customConfig != nil && customConfig.ARL != "" {
		arlEngine, err := NewARLEngine(customConfig.ARL)
		if err != nil {
			fmt.Printf("[!] 加载 ARL 指纹失败: %v\n", err)
			os.Exit(1)
		}
		s.arlEngine = arlEngine
		if !silent && !jsonOutput {
			fmt.Printf("[*] 已加载 ARL 指纹: %s (%d 条规则)\n", customConfig.ARL, len(arlEngine.fingerprints))
		}
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

	// 输出扫描统计（非静默模式且非 JSON 模式）
	if !s.silent && !s.jsonOutput {
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
	var allNames []string
	seen := make(map[string]bool)

	// 使用默认引擎检测
	if s.engine != nil {
		frameworks, err := s.engine.DetectContent(rawContent)
		if err == nil {
			for _, name := range frameworks.GetNames() {
				if !seen[name] {
					seen[name] = true
					allNames = append(allNames, name)
				}
			}
		}
	}

	// 使用自定义引擎检测
	if s.customEngine != nil {
		frameworks, err := s.customEngine.DetectContent(rawContent)
		if err == nil {
			for _, name := range frameworks.GetNames() {
				if !seen[name] {
					seen[name] = true
					allNames = append(allNames, name)
				}
			}
		}
	}

	return allNames
}

// detectFavicon 使用 fingers 引擎检测 favicon 指纹
// 主动获取 favicon 文件并进行 hash 匹配
//
// 参数：
//   - body: HTML 响应体，用于提取 favicon URL
//   - baseURL: 当前页面 URL
//
// 返回：
//   - []string: 检测到的框架名称列表
func (s *Scanner) detectFavicon(body, baseURL string) []string {
	// 提取 favicon URL
	faviconURL := extractFaviconURL(body, baseURL)
	if faviconURL == "" {
		return nil
	}

	// 获取 favicon 内容
	faviconContent, err := fetchFavicon(faviconURL, s.proxy)
	if err != nil || len(faviconContent) == 0 {
		return nil
	}

	var allNames []string
	seen := make(map[string]bool)

	// 使用默认引擎的 favicon 检测
	if s.engine != nil {
		frameworks := s.engine.MatchFavicon(faviconContent)
		for _, name := range frameworks.GetNames() {
			if !seen[name] {
				seen[name] = true
				allNames = append(allNames, name)
			}
		}
	}

	// 使用自定义引擎的 favicon 检测
	if s.customEngine != nil {
		frameworks := s.customEngine.MatchFavicon(faviconContent)
		for _, name := range frameworks.GetNames() {
			if !seen[name] {
				seen[name] = true
				allNames = append(allNames, name)
			}
		}
	}

	return allNames
}

// getFaviconHash 获取 favicon 的 MMH3 hash
func (s *Scanner) getFaviconHash(body, baseURL string) string {
	faviconURL := extractFaviconURL(body, baseURL)
	if faviconURL == "" {
		return ""
	}

	faviconContent, err := fetchFavicon(faviconURL, s.proxy)
	if err != nil || len(faviconContent) == 0 {
		return ""
	}

	return calcFaviconHash(faviconContent)
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
			continue
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

		// 使用 ARL 引擎进行指纹检测（如果启用）
		if s.arlEngine != nil {
			// 计算 favicon hash（如果需要）
			faviconHash := ""
			if task[1] == "0" {
				faviconHash = s.getFaviconHash(resp.Body, resp.URL)
			}
			// ARL 匹配
			arlMatched := s.arlEngine.Match(resp.Body, resp.Header, resp.Title, faviconHash)
			for _, m := range arlMatched {
				// 避免重复
				exists := false
				for _, existing := range matched {
					if existing == m {
						exists = true
						break
					}
				}
				if !exists {
					matched = append(matched, m)
				}
			}
		}

		// 主动获取 favicon 进行指纹检测（仅对主页面）
		if task[1] == "0" {
			faviconMatches := s.detectFavicon(resp.Body, resp.URL)
			for _, faviconMatch := range faviconMatches {
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
	// JSON 输出模式
	if s.jsonOutput {
		data, _ := json.Marshal(result)
		fmt.Println(string(data))
		return
	}

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
