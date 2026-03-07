# Redirect Policy Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 为 `xingfinger` 增加可配置的跳转策略，默认不跟随任何跳转，并保证一个输入目标默认只输出一条结果。

**Architecture:** 在 CLI 层新增 `--redirect-policy=never|http|all`，将跳转行为拆分为 HTTP 3xx 跟随策略和内容层跳转策略两个维度，再在扫描结果聚合层收敛输出语义，避免派生页单独打印。实现采用 TDD，优先从 E2E 行为锁定开始，再做最小代码改动。

**Tech Stack:** Go、Cobra、标准库 `net/http`、现有扫描队列与测试框架

---

### Task 1: 为跳转策略写失败测试

**Files:**
- Modify: `cmd/root_e2e_test.go`
- Test: `cmd/root_e2e_test.go`

**Step 1: Write failing E2E tests for default `never` behavior**

新增端到端测试，覆盖：
- HTTP 302 场景默认不跟随
- JS / `meta refresh` 场景默认不跟随
- 断言单个输入只产生一条结果

**Step 2: Run tests to verify they fail**

Run: `go test ./cmd -run 'TestExecute(RedirectPolicyDefaultsToNever|RedirectPolicyNeverSkipsContentRedirects)$' -count=1`
Expected: FAIL

### Task 2: 为参数解析写失败测试

**Files:**
- Modify: `cmd/root_e2e_test.go`
- Test: `cmd/root_e2e_test.go`

**Step 1: Add tests for policy parsing**

新增测试覆盖：
- `--redirect-policy=http` 跟随 HTTP 302
- `--redirect-policy=all` 跟随内容层跳转
- 非法值报错

**Step 2: Run focused tests to verify they fail**

Run: `go test ./cmd -run 'TestExecute(RedirectPolicyHTTPFollowsHTTPRedirect|RedirectPolicyAllFollowsContentRedirect|RedirectPolicyRejectsInvalidValue)$' -count=1`
Expected: FAIL

### Task 3: 实现 CLI 参数与策略类型

**Files:**
- Modify: `cmd/root.go`
- Modify: `pkg/scanner.go`
- Modify: `pkg/http.go`

**Step 1: Add redirect policy option**

在 CLI 层新增：
- 参数定义
- 默认值 `never`
- 非法值校验

**Step 2: Thread policy through scanner/fetch flow**

将跳转策略传入扫描器与请求逻辑，区分：
- HTTP 3xx 是否跟随
- 内容层跳转是否解析/入队

**Step 3: Run focused tests**

Run: `go test ./cmd -run 'TestExecute(RedirectPolicyDefaultsToNever|RedirectPolicyHTTPFollowsHTTPRedirect|RedirectPolicyAllFollowsContentRedirect|RedirectPolicyRejectsInvalidValue|RedirectPolicyNeverSkipsContentRedirects)$' -count=1`
Expected: PASS

### Task 4: 收敛输出为单结果语义

**Files:**
- Modify: `pkg/scanner.go`
- Modify: `cmd/root_e2e_test.go`

**Step 1: Adjust result aggregation/output**

确保一个输入目标默认只输出一条结果，不再把内容层派生页单独打印成第二条终端结果。

**Step 2: Run relevant tests**

Run: `go test ./cmd -run 'TestExecute(.*RedirectPolicy.*)$' -count=1`
Expected: PASS

### Task 5: 回归验证

**Files:**
- Test: `cmd/root_e2e_test.go`
- Test: `pkg/jsjump_test.go`
- Test: `pkg/http_integration_test.go`

**Step 1: Run command package tests**

Run: `go test ./cmd -count=1`
Expected: PASS

**Step 2: Run full test suite**

Run: `go test ./... -count=1`
Expected: PASS

### Task 6: 文档同步

**Files:**
- Modify: `docs/plans/2026-03-07-redirect-policy-design.md`
- Modify: `docs/plans/2026-03-07-redirect-policy.md`

**Step 1: Ensure docs match implementation**

核对实现与设计一致，确保默认值、参数名、输出语义无偏差。

**Step 2: Optional commit (user requested only)**

Run only if requested:
- `git add cmd/root.go pkg/scanner.go pkg/http.go cmd/root_e2e_test.go pkg/jsjump_test.go pkg/http_integration_test.go docs/plans/2026-03-07-redirect-policy-design.md docs/plans/2026-03-07-redirect-policy.md`
- `git commit -m "feat: add redirect policy controls"`
