// Package finger 提供 Web 指纹识别核心功能
package finger

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/twmb/murmur3"
)

// calcMurmur3 计算数据的 Murmur3 hash 值
// Murmur3 是一种非加密哈希函数，速度快且分布均匀
// FOFA 使用此算法计算 favicon hash
//
// 参数：
//   - data: 待计算 hash 的字节数据
//
// 返回：
//   - 32 位有符号整数形式的 hash 字符串
func calcMurmur3(data []byte) string {
	h := murmur3.New32()
	h.Write(data)
	// 转换为有符号 32 位整数，与 FOFA 保持一致
	return fmt.Sprintf("%d", int32(h.Sum32()))
}

// encodeBase64 将数据进行 Base64 编码
// 按照 FOFA 的格式要求，每 76 个字符换行
// 这是计算 favicon hash 的标准格式
//
// 参数：
//   - data: 待编码的字节数据
//
// 返回：
//   - 格式化后的 Base64 编码字节数据
func encodeBase64(data []byte) []byte {
	encoded := base64.StdEncoding.EncodeToString(data)
	var buf bytes.Buffer
	// 每 76 个字符添加换行符
	for i, ch := range encoded {
		buf.WriteByte(byte(ch))
		if (i+1)%76 == 0 {
			buf.WriteByte('\n')
		}
	}
	// 末尾添加换行符
	buf.WriteByte('\n')
	return buf.Bytes()
}

// calcFaviconHash 获取 favicon 并计算其 hash 值
// 用于与 FOFA 的 icon_hash 进行匹配
//
// 计算流程：
//  1. 下载 favicon 文件
//  2. 对内容进行 Base64 编码（每 76 字符换行）
//  3. 计算 Murmur3 hash
//
// 参数：
//   - url: favicon 的 URL 地址
//
// 返回：
//   - favicon 的 hash 值，失败返回 "0"
func calcFaviconHash(url string) string {
	// 创建 HTTP 客户端
	client := http.Client{
		Timeout: 8 * time.Second,
		Transport: &http.Transport{
			// 跳过 TLS 证书验证
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		// 不跟随重定向，使用最后一个响应
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// 发送 GET 请求
	resp, err := client.Get(url)
	if err != nil {
		return "0"
	}
	defer resp.Body.Close()

	// 只处理成功响应
	if resp.StatusCode != 200 {
		return "0"
	}

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "0"
	}

	// 计算 hash：先 Base64 编码，再计算 Murmur3
	return calcMurmur3(encodeBase64(body))
}
