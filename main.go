// xingfinger - Web 指纹识别工具
// 用于识别目标网站使用的 CMS、框架、服务器等技术栈信息
//
// 主要功能：
//   - 支持单个 URL 或批量 URL 扫描
//   - 支持 FOFA 搜索引擎集成查询
//   - 支持多种指纹匹配方式（关键字、正则、favicon hash）
//   - 支持结果导出为 JSON 或 Excel 格式
package main

import "github.com/yyhuni/xingfinger/cmd"

// main 程序入口函数
// 调用 cmd.Execute() 启动命令行解析和执行
func main() {
	cmd.Execute()
}
