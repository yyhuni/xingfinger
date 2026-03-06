# CLI 端到端边界覆盖计划

## 任务 1：补 TLS 服务辅助函数

**文件**
- 修改：`cmd/root_e2e_test.go`

**步骤**
1. 增加本地 TLS `httptest` 服务工厂
2. 复用现有 CLI 子进程执行辅助函数

## 任务 2：补边界 E2E

**文件**
- 修改：`cmd/root_e2e_test.go`

**测试**
- HTTPS 自签名扫描成功
- HTTP 302 跳转扫描成功
- JS redirect 扫描成功
- deflate 压缩响应扫描成功
- gzip 压缩响应扫描成功
- HTTP 代理扫描成功

## 任务 3：验证

**命令**
- `GOTOOLCHAIN=go1.26.0 ... go test ./cmd -run 'TestExecute(.*TLS.*|.*Redirect.*|.*Deflate.*|.*Gzip.*|.*Proxy.*)' -count=1`
- `GOTOOLCHAIN=go1.26.0 ... go test ./...`
- `GOTOOLCHAIN=go1.26.0 ... go build ./...`
- `GOTOOLCHAIN=go1.26.0 ... go vet ./...`
