# 实现计划: fingers 引擎集成

## 概述

将 chainreactors/fingers 指纹识别引擎集成到 xingfinger 项目中，替换现有的自定义指纹匹配逻辑。

## 任务

- [x] 1. 添加 fingers 依赖并清理旧文件
  - [x] 1.1 添加 `github.com/chainreactors/fingers` 到 go.mod
    - 运行 `go get github.com/chainreactors/fingers`
    - _需求: 5.1_
  - [x] 1.2 删除旧的指纹相关文件
    - 删除 `finger.json`
    - 删除 `module/finger/getfingerfile.go`
    - 删除 `module/finger/matchfinger.go`
    - 删除 `module/finger/faviconhash.go`
    - _需求: 4.1, 4.2, 4.3, 4.4_

- [x] 2. 修改 HTTP 响应处理
  - [x] 2.1 更新 `module/finger/http.go`
    - 在 Response 结构体中添加 RawContent 字段
    - 修改 fetch 函数，构建并保存原始 HTTP 响应
    - 添加详细的中文注释
    - _需求: 4.5, 6.1, 6.2, 6.3_

- [x] 3. 集成 fingers 引擎
  - [x] 3.1 重构 `module/finger/finger.go`
    - 导入 fingers 包
    - 在 Scanner 结构体中添加 engine 字段
    - 修改 NewScanner 函数，初始化 fingers 引擎
    - 添加 detectFingerprints 方法
    - 修改 scan 方法，使用 fingers 引擎进行检测
    - 移除旧的指纹匹配逻辑
    - 添加详细的中文注释
    - _需求: 1.1, 1.2, 1.3, 2.1, 2.2, 2.3, 2.4, 6.1, 6.2, 6.3, 6.4_

- [x] 4. 检查点 - 验证基本功能
  - 运行 `go build` 确保编译通过
  - 运行 `go mod tidy` 清理依赖
  - 测试单个 URL 扫描功能
  - _需求: 5.2, 5.3_

- [x] 5. 完善输出和注释
  - [x] 5.1 确保输出格式兼容
    - 验证 httpx 风格输出格式
    - 验证多框架用逗号连接
    - 验证 JSON 输出格式
    - _需求: 3.1, 3.2, 3.3_
  - [x] 5.2 为其他文件添加注释
    - 更新 `cmd/root.go` 注释
    - 更新 `module/finger/output.go` 注释
    - 更新 `module/finger/encoding.go` 注释
    - 更新 `module/finger/jsjump.go` 注释
    - 更新 `module/queue/queue.go` 注释
    - _需求: 6.1, 6.2, 6.3, 6.4, 6.5_

- [x] 6. 最终检查点
  - 确保所有测试通过
  - 验证完整扫描流程
  - 如有问题请询问用户

## 注意事项

- 所有注释使用中文
- 保持输出格式与现有版本兼容
- fingers 引擎实例在所有扫描线程间共享
