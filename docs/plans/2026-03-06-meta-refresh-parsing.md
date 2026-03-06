# Meta Refresh Parsing Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 修复 `meta refresh` 跳转提取把标签闭合尾巴误当成 URL 内容的问题，并保留绝对/相对跳转跟随能力。

**Architecture:** 保持 `parseJSRedirect` 作为统一入口，继续复用现有 JS 跳转解析逻辑，仅把 `meta refresh` 分支改成“属性级提取 + content 级 URL 解析”。实现完成后通过回归测试验证旧缺陷消失且原有行为保持稳定。

**Tech Stack:** Go、标准库 `regexp`、`net/url`、`strings`、Go `testing`

---

### Task 1: 写出失败回归测试

**Files:**
- Modify: `pkg/jsjump_test.go`
- Test: `pkg/jsjump_test.go`

**Step 1: Write the failing test**

新增一个用例，覆盖：
- `<meta http-equiv="refresh" content="0;url=http://www.baidu.com/baidu.html?from=noscript" />`
- `baseURL` 为 `https://www.baidu.com`
- 期望结果只包含干净的 `http://www.baidu.com/baidu.html?from=noscript`

**Step 2: Run test to verify it fails**

Run: `go test ./pkg -run TestParseJSRedirectSupportsSelfClosingMetaRefreshWithAbsoluteURL -count=1`
Expected: FAIL，旧实现会返回脏 URL

### Task 2: 实现最小修复

**Files:**
- Modify: `pkg/jsjump.go`

**Step 1: Add helpers for meta tag parsing**

新增最小辅助逻辑：
- 提取候选 `meta` 标签
- 解析标签属性
- 从 `content` 中取出 `url=` 值

**Step 2: Keep existing URL normalization**

对解析出的跳转目标继续走：
- `strings.TrimSpace`
- 去除包裹引号
- `url.Parse`
- `base.ResolveReference`

**Step 3: Run test to verify it passes**

Run: `go test ./pkg -run TestParseJSRedirectSupportsSelfClosingMetaRefreshWithAbsoluteURL -count=1`
Expected: PASS

### Task 3: 做相关回归验证

**Files:**
- Test: `pkg/jsjump_test.go`
- Test: `pkg/http_integration_test.go`

**Step 1: Run focused tests**

Run: `go test ./pkg -run 'TestParseJSRedirect|TestFetchFollowsMultipleRedirectsAndResolvesJSRedirectFromFinalURL' -count=1`
Expected: PASS

**Step 2: Run broader pkg tests if needed**

Run: `go test ./pkg -count=1`
Expected: PASS

### Task 4: 收尾

**Files:**
- Modify: `docs/plans/2026-03-06-meta-refresh-parsing-design.md`
- Modify: `docs/plans/2026-03-06-meta-refresh-parsing.md`

**Step 1: Ensure docs match implementation**

核对设计与实现范围一致，不扩散到无关逻辑。

**Step 2: Optional commit (user requested only)**

Run only if requested:
- `git add pkg/jsjump.go pkg/jsjump_test.go docs/plans/2026-03-06-meta-refresh-parsing-design.md docs/plans/2026-03-06-meta-refresh-parsing.md`
- `git commit -m "fix: harden meta refresh redirect parsing"`
