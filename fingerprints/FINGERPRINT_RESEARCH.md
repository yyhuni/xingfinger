# Web 指纹识别格式研究报告

## 研究概述

本报告基于网络搜索和官方文档，详细说明了 5 种 Web 指纹识别格式的来源、特点和应用场景。

---

## 1. EHole 格式

### 来源和背景

**EHole（棱洞）** 是由红队安全团队开发的红队重点攻击系统指纹探测工具。

- **GitHub**：https://github.com/EdgeSecurityTeam/EHole
- **开发者**：EdgeSecurityTeam
- **主要用途**：红队重点攻击系统指纹探测
- **特色**：专注于常见的 CMS、OA 系统等红队关注的目标

### 格式特点

1. **简洁易用**
   - JSON 格式，结构简单
   - 易于理解和维护
   - 适合快速创建指纹

2. **多种匹配方式**
   - keyword：精确字符串匹配
   - regular：正则表达式匹配
   - faviconhash：favicon 图标 hash 匹配

3. **多个检测位置**
   - body：响应体
   - header：HTTP 响应头
   - title：页面标题

### 应用场景

- 红队渗透测试
- 资产发现和识别
- 快速指纹库构建

### 示例

```json
{
  "fingerprint": [
    {
      "cms": "WordPress",
      "method": "keyword",
      "location": "body",
      "keyword": ["wp-content", "wp-includes"]
    }
  ]
}
```

---

## 2. Goby 格式

### 来源和背景

**Goby** 是由赵武（Pangolin、JSky、FOFA 作者）打造的网络安全测试工具。

- **官网**：https://www.gobysec.com/
- **开发者**：赵武（Zwell）
- **主要用途**：网络安全测试、漏洞扫描、资产发现
- **特色**：实战性强、体系性完整、高效率

### 格式特点

1. **灵活的逻辑组合**
   - 支持 AND（&）、OR（|）逻辑
   - 支持复杂的条件组合
   - 适合复杂的检测场景

2. **JSON 数组格式**
   - 每个元素代表一个指纹
   - 支持多个规则组合
   - 易于扩展

3. **丰富的特征匹配**
   - 支持多种特征类型
   - 支持精确匹配和模糊匹配
   - 支持版本提取

### 应用场景

- 网络安全测试
- 漏洞扫描
- 资产发现和分类

### 示例

```json
[
  {
    "name": "WordPress",
    "logic": "a|b",
    "rule": [
      {
        "label": "a",
        "feature": "wp-content",
        "is_equal": false
      }
    ]
  }
]
```

---

## 3. Wappalyzer 格式

### 来源和背景

**Wappalyzer** 是一个浏览器扩展和 Web 应用程序，用于识别网站上使用的技术。

- **GitHub**：https://github.com/AliasIO/wappalyzer
- **开发者**：Alias
- **主要用途**：Web 技术识别和分类
- **特色**：支持 1500+ 种技术识别

### 格式特点

1. **灵活的检测方式**
   - HTML 内容匹配
   - HTTP 响应头匹配
   - JavaScript 脚本匹配
   - Cookie 匹配
   - Meta 标签匹配

2. **版本提取**
   - 支持正则表达式提取版本号
   - 支持多个版本模式

3. **技术依赖关系**
   - 支持 implies 字段表示技术依赖
   - 自动推导相关技术

### 应用场景

- Web 技术识别
- 浏览器扩展集成
- 通用技术栈检测

### 示例

```json
{
  "WordPress": {
    "cats": [1, 6],
    "headers": {
      "X-Powered-By": "WordPress"
    },
    "html": ["/wp-content/", "/wp-includes/"],
    "implies": "PHP"
  }
}
```

---

## 4. Fingers 格式

### 来源和背景

**Fingers** 是由 chainreactors 开发的多指纹库聚合识别引擎。

- **GitHub**：https://github.com/chainreactors/fingers
- **开发者**：chainreactors
- **主要用途**：多指纹库聚合、Web 指纹识别
- **特色**：聚合多个指纹库、功能完整

### 格式特点

1. **功能完整**
   - 支持多种匹配方式
   - 支持复杂的检测逻辑
   - 支持版本提取

2. **多个检测位置**
   - HTML 内容
   - HTTP 响应头
   - JavaScript 脚本
   - Cookie
   - Meta 标签

3. **灵活的扩展**
   - 支持自定义指纹
   - 支持指纹聚合
   - 支持多个指纹库

### 应用场景

- 复杂的指纹检测
- 多指纹库聚合
- 高级指纹识别

### 示例

```json
[
  {
    "name": "WordPress",
    "version": "5.0",
    "category": "CMS",
    "html": ["wp-content", "wp-includes"],
    "scripts": ["/wp-includes/js/"]
  }
]
```

---

## 5. FingerPrintHub 格式

### 来源和背景

**FingerPrintHub** 是一个指纹中心项目，使用 **Nuclei 模板格式**。

- **Nuclei GitHub**：https://github.com/projectdiscovery/nuclei
- **Nuclei 文档**：https://docs.projectdiscovery.io/templates/introduction
- **开发者**：ProjectDiscovery
- **主要用途**：现代化漏洞扫描、指纹识别

### Nuclei 背景

**Nuclei** 是由 ProjectDiscovery 开发的现代化漏洞扫描工具。

- **特点**：
  - 基于 YAML 的简单模板
  - 支持多种协议（TCP、DNS、HTTP 等）
  - 零误报设计
  - 社区驱动

- **模板格式**：
  - 原生格式：YAML
  - 支持格式：JSON（通过转换）
  - 通用语言：用于表示可利用的漏洞

### FingerPrintHub 格式特点

1. **最灵活和强大**
   - 基于 Nuclei 模板格式
   - 支持复杂的匹配器
   - 支持提取器
   - 支持条件逻辑

2. **多种 Matcher 类型**
   - word：字符串匹配
   - regex：正则表达式匹配
   - status-code：HTTP 状态码
   - favicon：Favicon hash
   - 等等...

3. **高级功能**
   - 支持提取器提取信息
   - 支持条件逻辑（AND、OR）
   - 支持多个请求
   - 支持动态变量

### 应用场景

- 高级漏洞检测
- 复杂的指纹识别
- 需要提取信息的场景
- 需要条件逻辑的场景

### 示例

```json
[
  {
    "id": "wordpress-detect",
    "info": {
      "name": "WordPress",
      "author": "test",
      "tags": "detect,tech,wordpress",
      "severity": "info",
      "metadata": {
        "product": "WordPress",
        "vendor": "WordPress"
      }
    },
    "http": [
      {
        "method": "GET",
        "path": ["{{BaseURL}}/"],
        "matchers": [
          {
            "type": "word",
            "words": ["wp-content"],
            "case-insensitive": true
          }
        ]
      }
    ]
  }
]
```

---

## 格式演进历程

### 第一代：简单格式（EHole）
- 时间：2019-2020
- 特点：简洁、易用
- 适用：简单指纹识别

### 第二代：增强格式（Goby、Wappalyzer）
- 时间：2020-2021
- 特点：支持逻辑组合、多种检测方式
- 适用：中等复杂度检测

### 第三代：聚合格式（Fingers）
- 时间：2021-2022
- 特点：多指纹库聚合、功能完整
- 适用：复杂指纹识别

### 第四代：通用格式（FingerPrintHub/Nuclei）
- 时间：2022-现在
- 特点：最灵活、最强大、通用
- 适用：所有场景

---

## 格式对比分析

| 维度 | EHole | Goby | Wappalyzer | Fingers | FingerPrintHub |
|------|-------|------|-----------|---------|----------------|
| **发布时间** | 2019 | 2020 | 2015 | 2021 | 2022 |
| **格式** | JSON | JSON | JSON | JSON | JSON |
| **复杂度** | 低 | 中 | 中 | 高 | 高 |
| **学习难度** | 低 | 中 | 中 | 高 | 高 |
| **功能完整度** | 低 | 中 | 中 | 高 | 高 |
| **社区活跃度** | 中 | 高 | 高 | 中 | 高 |
| **应用广泛度** | 中 | 高 | 高 | 中 | 中 |
| **扩展性** | 低 | 中 | 中 | 高 | 高 |

---

## 选择建议

### 快速开始
- 推荐：**EHole 格式**
- 原因：简单易用，学习成本低

### 中等复杂度
- 推荐：**Goby 格式** 或 **Wappalyzer 格式**
- 原因：功能足够，易于维护

### 复杂场景
- 推荐：**Fingers 格式** 或 **FingerPrintHub 格式**
- 原因：功能完整，支持复杂逻辑

### 生产环境
- 推荐：**FingerPrintHub 格式**
- 原因：最灵活、最强大、社区支持好

---

## 技术趋势

1. **格式统一化**
   - 从多种格式向通用格式发展
   - Nuclei 模板成为事实标准

2. **功能增强**
   - 从简单匹配向复杂逻辑发展
   - 支持提取器、条件逻辑等高级功能

3. **社区驱动**
   - 从单一工具向社区聚合发展
   - 多个工具共享指纹库

4. **自动化**
   - 从手工编写向自动生成发展
   - AI 辅助指纹生成

---

## 参考资源

### 官方文档
- [EHole GitHub](https://github.com/EdgeSecurityTeam/EHole)
- [Goby 官网](https://www.gobysec.com/)
- [Wappalyzer GitHub](https://github.com/AliasIO/wappalyzer)
- [Fingers GitHub](https://github.com/chainreactors/fingers)
- [Nuclei 文档](https://docs.projectdiscovery.io/templates/introduction)

### 相关文章
- [Wappalyzer 指纹配置](https://neptliang.github.io/2021/01/11/Wappalyzer-Fingerprint-Configuration/)
- [Goby 自定义编写 EXP 入门篇](https://www.cnblogs.com/gobybot/p/18619435)
- [Nuclei v2.5.0 Release](https://projectdiscovery.io/blog/nuclei-v2-5-release)

---

## 总结

Web 指纹识别格式经历了从简单到复杂、从单一到聚合、从专用到通用的演进过程。

- **EHole** 代表了第一代简单格式
- **Goby** 和 **Wappalyzer** 代表了第二代增强格式
- **Fingers** 代表了第三代聚合格式
- **FingerPrintHub/Nuclei** 代表了第四代通用格式

未来的发展方向是：
1. 格式进一步统一
2. 功能进一步增强
3. 自动化程度进一步提高
4. 社区合作进一步深化
