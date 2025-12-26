//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"net/http"
)

func main() {
	// EHole 指纹测试页面
	http.HandleFunc("/ehole", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Server", "TestServer/1.0")
		fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head>
    <title>EHole Test Page</title>
</head>
<body>
    <h1>EHole 指纹测试页面</h1>
    <p>这是用于测试 EHole 自定义指纹的页面</p>
</body>
</html>`)
	})

	// Goby 指纹测试页面
	http.HandleFunc("/goby", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Server", "TestServer/1.0")
		fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head>
    <title>Goby Test Page</title>
</head>
<body>
    <h1>Goby 指纹测试页面</h1>
    <p>这是用于测试 Goby 自定义指纹的页面</p>
</body>
</html>`)
	})

	// Wappalyzer 指纹测试页面
	http.HandleFunc("/wappalyzer", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("X-Powered-By", "TestSystem")
		fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head>
    <title>Wappalyzer Test Page</title>
</head>
<body>
    <h1>Wappalyzer 指纹测试页面</h1>
    <script src="example_app.js"></script>
</body>
</html>`)
	})

	// Fingers 指纹测试页面
	http.HandleFunc("/fingers", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head>
    <title>Fingers Test Page</title>
</head>
<body>
    <h1>Fingers 指纹测试页面</h1>
    <p>TestFramework detected</p>
</body>
</html>`)
	})

	// FingerPrintHub 指纹测试页面
	http.HandleFunc("/fingerprinthub", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head>
    <title>FingerPrintHub Test Page</title>
</head>
<body>
    <h1>FingerPrintHub 指纹测试页面</h1>
    <p>TestApp detected</p>
</body>
</html>`)
	})

	// Favicon 测试
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		// 返回一个简单的 ICO 文件内容（最小的有效 ICO 文件）
		icoData := []byte{
			0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x01,
			0x00, 0x00, 0x01, 0x00, 0x18, 0x00, 0x30, 0x00,
			0x00, 0x00, 0x16, 0x00, 0x00, 0x00, 0x28, 0x00,
			0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00,
			0x00, 0x00, 0x01, 0x00, 0x18, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF,
			0xFF, 0x00,
		}
		w.Header().Set("Content-Type", "image/x-icon")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(icoData)))
		w.Write(icoData)
	})

	// 默认页面
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head>
    <title>Test Server</title>
    <link rel="icon" type="image/x-icon" href="/favicon.ico">
</head>
<body>
    <h1>指纹测试服务器</h1>
    <ul>
        <li><a href="/ehole">EHole 指纹测试</a></li>
        <li><a href="/goby">Goby 指纹测试</a></li>
        <li><a href="/wappalyzer">Wappalyzer 指纹测试</a></li>
        <li><a href="/fingers">Fingers 指纹测试</a></li>
        <li><a href="/fingerprinthub">FingerPrintHub 指纹测试</a></li>
        <li><a href="/favicon.ico">Favicon 测试</a></li>
    </ul>
</body>
</html>`)
	})

	fmt.Println("测试服务器启动在 http://localhost:8888")
	http.ListenAndServe(":8888", nil)
}
