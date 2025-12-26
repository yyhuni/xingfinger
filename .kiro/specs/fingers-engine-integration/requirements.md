# 需求文档

## 简介

重构 xingfinger 的指纹识别引擎，使用 chainreactors/fingers 库替换现有的自定义指纹匹配逻辑。fingers 是一个多指纹库聚合识别引擎，支持 fingers、wappalyzer、fingerprinthub、ehole、goby 等多个指纹库，性能更强、指纹更全。

## 术语表

- **Fingers_Engine**: chainreactors/fingers 库提供的多指纹库聚合识别引擎
- **Scanner**: xingfinger 的扫描器，负责并发扫描和结果收集
- **Framework**: fingers 库中的指纹识别结果结构体
- **HTTP_Response**: HTTP 响应的原始内容，包含 header 和 body

## 需求

### 需求 1: 集成 fingers 引擎

**用户故事:** 作为开发者，我希望集成 fingers 引擎，以便使用多个指纹库进行更好的检测。

#### 验收标准

1. 当 Scanner 初始化时，Fingers_Engine 应使用 `fingers.NewEngine()` 创建
2. 当 Fingers_Engine 初始化失败时，Scanner 应退出并显示错误信息
3. Scanner 应存储 Fingers_Engine 实例，供所有扫描线程复用

### 需求 2: 替换指纹匹配逻辑

**用户故事:** 作为开发者，我希望用 fingers 引擎替换自定义指纹匹配，以便利用其优化的匹配算法。

#### 验收标准

1. 当收到 HTTP 响应时，Scanner 应使用 `Fingers_Engine.DetectContent()` 进行指纹检测
2. 当指纹检测完成时，Scanner 应从 Frameworks 结果中提取框架名称
3. Scanner 应移除旧的基于 finger.json 的匹配逻辑
4. Scanner 应移除旧的 Fingerprints 结构体及相关加载函数

### 需求 3: 保持输出格式兼容

**用户故事:** 作为用户，我希望输出格式保持不变，以便现有工作流程不受影响。

#### 验收标准

1. 当显示结果时，Scanner 应以相同的 httpx 风格格式显示框架名称
2. 当检测到多个框架时，Scanner 应用逗号连接它们
3. 当保存到 JSON 时，Scanner 应保持相同的 Result 结构体格式

### 需求 4: 清理旧代码

**用户故事:** 作为开发者，我希望移除未使用的代码，以便代码库干净且易于维护。

#### 验收标准

1. 系统应移除 `finger.json` 文件
2. 系统应移除 `getfingerfile.go`（Fingerprints 加载）
3. 系统应移除 `matchfinger.go`（关键字/正则匹配）
4. 系统应移除 `faviconhash.go`（favicon 哈希计算 - fingers 已处理）
5. 系统应更新 `http.go` 以返回原始 HTTP 响应供 fingers 引擎使用

### 需求 5: 更新依赖

**用户故事:** 作为开发者，我希望更新依赖，以便项目正确使用 fingers 库。

#### 验收标准

1. 系统应将 `github.com/chainreactors/fingers` 添加到 go.mod
2. 系统应运行 `go mod tidy` 清理未使用的依赖
3. 系统应验证项目构建成功

### 需求 6: 添加详细代码注释

**用户故事:** 作为开发者，我希望所有代码都有详细的中文注释，以便代码易于理解和维护。

#### 验收标准

1. 所有 Go 源文件应包含文件级别的包注释，说明该文件的用途
2. 所有导出的函数应包含函数注释，说明功能、参数和返回值
3. 所有导出的结构体应包含结构体注释，说明用途和字段含义
4. 复杂的逻辑块应包含行内注释，解释实现细节
5. 注释应使用中文编写，保持风格一致
