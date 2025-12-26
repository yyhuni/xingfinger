// Package finger 提供 Web 指纹识别核心功能
package finger

import (
	"encoding/json"
	"io/ioutil"
)

// Packjson 指纹规则库结构体
// 包含所有指纹规则的集合
type Packjson struct {
	Fingerprint []Fingerprint // 指纹规则列表
}

// Fingerprint 单条指纹规则结构体
// 定义了如何识别特定的 CMS/框架
type Fingerprint struct {
	Cms      string   // CMS/框架名称
	Method   string   // 匹配方法：keyword（关键字）、regular（正则）、faviconhash（图标hash）
	Location string   // 匹配位置：body（响应体）、header（响应头）、title（页面标题）
	Keyword  []string // 匹配关键字/正则表达式列表
}

// Webfingerprint 全局指纹规则库实例
// 在程序启动时加载，供所有扫描线程共享
var (
	Webfingerprint *Packjson
)

// LoadWebfingerprint 从文件加载指纹规则库
// 指纹文件为 JSON 格式，包含所有指纹规则
//
// 参数：
//   - path: 指纹文件路径（通常为 finger.json）
//
// 返回：
//   - error: 加载失败时返回错误信息
func LoadWebfingerprint(path string) error {
	// 读取文件内容
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// 解析 JSON
	var config Packjson
	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	// 保存到全局变量
	Webfingerprint = &config
	return nil
}

// GetWebfingerprint 获取全局指纹规则库实例
// 返回已加载的指纹规则库，供扫描器使用
func GetWebfingerprint() *Packjson {
	return Webfingerprint
}
