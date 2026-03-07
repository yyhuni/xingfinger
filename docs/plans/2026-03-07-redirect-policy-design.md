# Redirect Policy 设计

## 背景

当前 `xingfinger` 在扫描主页面后，会额外解析页面中的 JavaScript 跳转和 `meta refresh` 跳转，并把解析出的目标重新入队继续扫描。这会导致用户只输入一个目标时，终端可能打印多行结果。

以 `www.baidu.com` 为例，输入单个目标后除了主页面结果，还会继续扫描 `http://www.baidu.com/baidu.html?from=noscript` 这类内容层派生页。对于指纹识别工具来说，这种行为容易让用户误以为自己输入了多个目标，也会把派生页指纹噪音混入主目标识别结果。

## 目标

- 让跳转处理行为可配置，而不是默认静默扩展扫描范围
- 默认情况下，一个输入目标只对应一个输出结果
- 区分 HTTP 3xx 跳转与内容层跳转（JS / `meta refresh`）
- 保留需要时继续跟随跳转的能力
- 尽量以最小改动完成行为收敛

## 参考调研

### httpx

`httpx` 提供 `-fr`、`-fhr`、`-maxr`、`-location` 等参数，核心关注点是 HTTP 层重定向。它没有把 JS / `meta refresh` 作为默认重定向语义的一部分。

### WhatWeb

`WhatWeb` 提供 `--follow-redirect=never|http-only|meta-only|same-site|always`，明确把 HTTP 跳转和内容层跳转分开控制。这说明同类工具通常会把“是否跟随跳转”做成显式策略，而不是固定行为。

### EHole

`EHole` 会尝试解析 JS / `meta refresh`，但实现较保守，遇到包含 `http` 的内容跳转会直接跳过。它避免了一部分噪音，但也会漏掉真实内容跳转目标。

## 方案对比

### 方案 A：统一策略枚举（推荐）

新增 `--redirect-policy=never|http|all`：

- `never`：不跟任何跳转
- `http`：仅跟 HTTP 3xx
- `all`：跟 HTTP 3xx + JS / `meta refresh`

优点：
- 用户心智清晰，一个参数即可表达全部行为
- 后续可自然扩展为 `same-site`、`content-only` 等更细策略
- 便于帮助文案、测试和结果语义统一

缺点：
- 相比单布尔参数，初次阅读要理解三个枚举值

### 方案 B：两个布尔开关

新增：
- `--follow-http-redirects`
- `--follow-content-redirects`

优点：
- 控制最直观

缺点：
- 参数组合更多，帮助文案更啰嗦
- 后续如果加入更多策略，参数会继续膨胀

### 方案 C：只做一个总开关

新增 `--follow-redirects`，开启后所有跳转都跟，关闭后都不跟。

优点：
- 实现最简单

缺点：
- 无法只跟 HTTP 而不跟内容层跳转
- 不能很好覆盖当前场景中的真实需求

## 设计决策

采用方案 A：新增 `--redirect-policy=never|http|all`。

### 参数语义

- `never`
  - 不跟随任何跳转
  - 不跟 HTTP 3xx
  - 不跟 JS / `meta refresh`
- `http`
  - 仅跟随 HTTP 3xx
  - 不跟 JS / `meta refresh`
- `all`
  - 跟随 HTTP 3xx
  - 跟随 JS / `meta refresh`

默认值为 `never`。

## 输出语义

无论策略如何，默认都保持“一个输入目标只输出一条结果”的终端语义。

- `never`
  - 输出输入目标这一跳的结果
- `http`
  - 输出跟随 HTTP 3xx 后的最终结果
- `all`
  - 输出跟随 HTTP 3xx 与内容层跳转后的最终结果

这样可以避免派生页被单独打印成第二行，减少用户困惑与指纹噪音。

## 本次实现范围

### 要做

- 新增 `--redirect-policy` 参数与非法值校验
- 调整 HTTP 客户端重定向策略
- 调整 JS / `meta refresh` 跳转继续扫描逻辑
- 调整结果收集与输出逻辑，避免一个输入目标产生多条终端结果
- 补充 E2E 与单元测试

### 不做

- 不扩展为 `same-site`、`content-only` 等更细粒度策略
- 不新增复杂跳转链展示
- 不顺手处理默认指纹误报问题
- 不改 JSON 输出结构，除非实现中发现必须同步调整

## 测试策略

至少覆盖以下场景：

- `never`：HTTP 302 不跟随，只输出起始目标结果
- `http`：HTTP 302 跟随到最终 URL，但只输出一条结果
- `never`：JS / `meta refresh` 不跟随，只输出起始页
- `all`：JS / `meta refresh` 跟随到最终页，但仍只输出一条结果
- 参数非法值时报错
- 默认值为 `never`
