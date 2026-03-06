# Fetch More Boundary Tests Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 在不修改业务逻辑的前提下，为 `fetch` 增加更多边界集成测试覆盖。

**Architecture:** 使用 `httptest` 构造多跳重定向、错误压缩头、无内容类型、空响应等真实 HTTP 场景，优先编写应当通过的绿色测试。

**Tech Stack:** Go, `net/http`, `net/http/httptest`

---

### Task 1: Redirect 链与最终 URL 解析
1. 写测试覆盖多跳重定向。
2. 写测试覆盖最终 URL 上的相对 JS 跳转。
3. 运行测试确认通过。

### Task 2: 压缩头容错
1. 写测试覆盖错误 `deflate` 头但正文仍是纯 HTML。
2. 运行测试确认通过。

### Task 3: 缺省头与空响应
1. 写测试覆盖缺失 `Content-Type`。
2. 写测试覆盖 `204 No Content`。
3. 运行测试确认通过。

### Task 4: 全量验证
1. 运行 `go test ./...`
2. 运行 `go build ./...`
3. 运行 `go vet ./...`
