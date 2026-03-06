# Fetch Boundary Behavior Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 为 `fetch` 补上重定向、压缩响应和非 HTML 边界的集成测试，并修复暴露的问题。

**Architecture:** 使用 `httptest` 构造真实 HTTP 响应，先写失败测试，再在 `pkg/http.go` 做最小改动。重点验证最终 URL、`gzip` / `deflate` 解压、以及非 HTML 内容保真。

**Tech Stack:** Go, `net/http`, `net/http/httptest`, `compress/gzip`, `compress/zlib`

---

### Task 1: Redirect 集成测试

**Files:**
- Modify: `pkg/http_integration_test.go`
- Test: `pkg/http_integration_test.go`

1. 写失败测试，验证 `fetch` 返回最终 URL。
2. 运行测试确认失败。
3. 最小化修改 `pkg/http.go`。
4. 运行测试确认通过。

### Task 2: gzip / deflate 集成测试

**Files:**
- Modify: `pkg/http_integration_test.go`
- Modify: `pkg/http.go`
- Test: `pkg/http_integration_test.go`

1. 写失败测试，覆盖 `gzip` 与 `deflate` HTML 响应。
2. 运行测试确认至少 `deflate` 失败。
3. 最小化修改解压逻辑。
4. 运行测试确认通过。

### Task 3: 非 HTML 边界测试

**Files:**
- Modify: `pkg/http_integration_test.go`
- Test: `pkg/http_integration_test.go`

1. 写测试验证 JSON/纯文本响应不误提取标题。
2. 运行测试确认行为稳定。
3. 如有失败，再做最小修复。

### Task 4: 全量验证

**Files:**
- Test: `pkg/http_integration_test.go`

1. 运行 `go test ./...`
2. 运行 `go build ./...`
3. 运行 `go vet ./...`
