# Stdout 静默初始化实现计划

## 任务 1：补充资源管理回归测试

**文件**
- 新增：`pkg/scanner_stdout_test.go`

**步骤**
1. 先写一个针对 stdout 静默辅助函数的测试
2. 验证测试在修复前失败，失败原因应为未关闭临时文件
3. 再进行最小修复

## 任务 2：最小修复句柄泄露

**文件**
- 修改：`pkg/scanner.go`

**步骤**
1. 提取一个最小辅助函数管理 stdout 临时重定向
2. 在默认引擎和自定义引擎初始化路径复用该辅助函数
3. 不改变错误处理和引擎启用逻辑

## 任务 3：验证

**命令**
- `go test ./pkg -run 'TestWithSilentStdout' -count=1`
- `go test ./...`
- `go build ./...`
- `go vet ./...`
