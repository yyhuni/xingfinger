# CLI 端到端覆盖扩展计划

## 任务 1：补测试辅助函数

**文件**
- 修改：`cmd/root_e2e_test.go`

**步骤**
1. 提供本地测试服务创建辅助函数
2. 提供 EHole / ARL 临时指纹文件生成函数
3. 提供 CLI 子进程调用与输出解析辅助函数

## 任务 2：补命令可用性测试

**文件**
- 修改：`cmd/root_e2e_test.go`

**测试**
- `-l` + `-o` 输出文件可写且 JSON 可解析
- `-s` 静默模式只输出命中项

## 任务 3：补指纹有效性测试

**文件**
- 修改：`cmd/root_e2e_test.go`

**测试**
- EHole `body` 命中
- EHole `header` 命中
- EHole `title` 命中
- EHole `faviconhash` 命中
- ARL 命中

## 任务 4：验证

**命令**
- `GOTOOLCHAIN=go1.26.0 ... go test ./cmd -run 'TestExecute' -count=1`
- `GOTOOLCHAIN=go1.26.0 ... go test ./...`
- `GOTOOLCHAIN=go1.26.0 ... go build ./...`
- `GOTOOLCHAIN=go1.26.0 ... go vet ./...`
