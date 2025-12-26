// Package finger 提供 Web 指纹识别功能
package finger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// saveResults 根据文件扩展名自动选择保存格式
// 支持的格式：
//   - .json: JSON 格式，便于程序解析
//   - .xlsx: Excel 格式，便于人工查看和分析
//
// 参数：
//   - filename: 输出文件路径，根据扩展名决定保存格式
//   - results: 指纹识别结果切片
func saveResults(filename string, results []Result) {
	// 获取文件扩展名并转为小写进行匹配
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".json":
		saveJSON(filename, results)
	case ".xlsx":
		saveXLSX(filename, results)
	}
}

// saveJSON 将结果保存为 JSON 格式文件
// JSON 格式带有缩进，便于阅读和后续程序处理
//
// 参数：
//   - filename: 输出文件路径
//   - results: 指纹识别结果切片
func saveJSON(filename string, results []Result) {
	// 使用缩进格式化 JSON，提高可读性
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Println("[!] JSON error:", err)
		return
	}

	// 创建输出文件
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println("[!] Create error:", err)
		return
	}
	defer f.Close()

	// 写入 JSON 数据
	f.Write(data)
}

// saveXLSX 将结果保存为 Excel 格式文件
// Excel 格式包含以下列：
//   - url: 目标 URL
//   - cms: 识别到的 CMS/框架名称
//   - server: 服务器类型
//   - status_code: HTTP 状态码
//   - length: 响应内容长度
//   - title: 页面标题
//
// 参数：
//   - filename: 输出文件路径
//   - results: 指纹识别结果切片
func saveXLSX(filename string, results []Result) {
	// 创建新的 Excel 文件
	xlsx := excelize.NewFile()

	// 设置表头
	xlsx.SetCellValue("Sheet1", "A1", "url")
	xlsx.SetCellValue("Sheet1", "B1", "cms")
	xlsx.SetCellValue("Sheet1", "C1", "server")
	xlsx.SetCellValue("Sheet1", "D1", "status_code")
	xlsx.SetCellValue("Sheet1", "E1", "length")
	xlsx.SetCellValue("Sheet1", "F1", "title")

	// 遍历结果，逐行写入数据
	// 从第 2 行开始写入（第 1 行是表头）
	for i, r := range results {
		row := strconv.Itoa(i + 2)
		xlsx.SetCellValue("Sheet1", "A"+row, r.URL)
		xlsx.SetCellValue("Sheet1", "B"+row, r.CMS)
		xlsx.SetCellValue("Sheet1", "C"+row, r.Server)
		xlsx.SetCellValue("Sheet1", "D"+row, r.StatusCode)
		xlsx.SetCellValue("Sheet1", "E"+row, r.Length)
		xlsx.SetCellValue("Sheet1", "F"+row, r.Title)
	}

	// 保存 Excel 文件
	if err := xlsx.SaveAs(filename); err != nil {
		fmt.Println("[!] Save error:", err)
	}
}
