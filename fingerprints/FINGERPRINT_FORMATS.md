# Web 指纹识别格式详解

本文档详细说明了 xingfinger 支持的 5 种指纹格式的特点、用途和格式规范。

## 1. EHole 格式

### 简介

**EHole（棱洞）** 是由红队安全团队开发的红队重点攻击系统指纹探测工具。EHole 格式是一种简洁的 JSON 格式，专门用于 Web 系统指纹识别。

**特点**：
- 格式简单易懂
- 支持多种匹配方式（keyword、regular、faviconhash）
- 支持多个检测位置（body、header、title）
- 轻量级，易于维护

### 格式规范

```json
{
  "fingerprint": [
    {
      "cms": "系统名称",
      "method": "keyword|regular|faviconhash",
      "location": "body|header|title",
      "keyword": ["特征字符串1", "特征字符串2"],
      "version": "版本号（可选）"
    }
  ]
}
```

### 字段说明

| 字段 | 说明 | 示例 |
|------|------|------|
| cms | 系统名称 | "WordPress" |
| method | 匹配方式 | "keyword"、"regular"、"faviconhash" |
| location | 检测位置 | "body"、"header"、"title" |
| keyword | 特征字符串数组 | ["wp-content", "wp-includes"] |
| version | 版本号（可选） | "5.0" |

### 匹配方式

- **keyword**：精确字符串匹配，支持多个关键词（AND 逻辑）
- **regular**：正则表达式匹配
- **faviconhash**：favicon 图标 hash 匹配

### 示例

```json
{
  "fingerprint": [
    {
      "cms": "WordPress",
      "method": "keyword",
      "location": "body",
      "keyword": ["wp-content", "wp-includes"]
    },
    {
      "cms": "Joomla",
      "method": "regular",
      "location": "body",
      "keyword": ["Joomla!\\s+([\\d.]+)"]
    }
  ]
}
```

---

## 2. Goby 格式

### 简介

**Goby** 是由赵武（Pangolin、JSky、FOFA 作者）打造的网络安全测试工具。Goby 格式使用 JSON 数组，支持复杂的逻辑组合。

**特点**：
- 支持复杂的逻辑组合（AND、OR）
- 支持多种特征匹配方式
- 灵活的规则定义
- 适合复杂的指纹检测场景

### 格式规范

```json
[
  {
    "name": "系统名称",
    "logic": "a|b|c|...",
    "rule": [
      {
        "label": "a",
        "feature": "特征字符串",
        "is_equal": true|false
      }
    ]
  }
]
```

### 字段说明

| 字段 | 说明 | 示例 |
|------|------|------|
| name | 系统名称 | "WordPress" |
| logic | 逻辑组合 | "a"、"a\|b"、"a&b" |
| rule | 规则数组 | 见下表 |
| label | 规则标签 | "a"、"b"、"c" |
| feature | 特征字符串 | "wp-content" |
| is_equal | 是否精确匹配 | true、false |

### 逻辑说明

- **a**：单个规则 a 匹配
- **a\|b**：规则 a 或 b 匹配（OR 逻辑）
- **a&b**：规则 a 和 b 都匹配（AND 逻辑）
- **a&(b\|c)**：规则 a 匹配，且 b 或 c 匹配

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
      },
      {
        "label": "b",
        "feature": "wp-includes",
        "is_equal": false
      }
    ]
  }
]
```

---

## 3. Wappalyzer 格式

### 简介

**Wappalyzer** 是一个浏览器扩展和 Web 应用程序，用于识别网站上使用的技术。Wappalyzer 格式是一种灵活的 JSON 对象格式。

**特点**：
- 支持多种检测方式（HTML、headers、scripts、cookies 等）
- 支持正则表达式和版本提取
- 支持技术依赖关系（implies）
- 广泛应用于 Web 技术识别

### 格式规范

```json
{
  "技术名称": {
    "cats": [1, 6],
    "headers": {
      "X-Powered-By": "技术名称"
    },
    "html": "正则表达式",
    "scripts": "脚本路径",
    "cookies": {
      "cookie名": "值"
    },
    "meta": {
      "generator": "正则表达式"
    },
    "implies": "依赖的技术",
    "icon": "图标文件",
    "website": "官网地址"
  }
}
```

### 字段说明

| 字段 | 说明 | 示例 |
|------|------|------|
| cats | 分类 ID | [1, 6] |
| headers | HTTP 响应头匹配 | {"X-Powered-By": "WordPress"} |
| html | HTML 内容匹配 | "wp-content" |
| scripts | 脚本路径匹配 | "/wp-includes/js/" |
| cookies | Cookie 匹配 | {"wordpress_logged_in": ""} |
| meta | Meta 标签匹配 | {"generator": "WordPress"} |
| implies | 依赖的技术 | "PHP" |
| icon | 图标文件 | "WordPress.svg" |
| website | 官网地址 | "https://wordpress.org" |

### 分类 ID

- 1: CMS
- 6: Ecommerce
- 11: JavaScript frameworks
- 18: Web frameworks
- 等等...

### 示例

```json
{
  "WordPress": {
    "cats": [1, 6],
    "headers": {
      "X-Powered-By": "WordPress"
    },
    "html": [
      "/wp-content/",
      "/wp-includes/",
      "<link rel='stylesheet' id='wp-block-library"
    ],
    "scripts": "/wp-includes/js/",
    "implies": "PHP",
    "icon": "WordPress.svg",
    "website": "https://wordpress.org"
  }
}
```

---

## 4. Fingers 格式

### 简介

**Fingers** 是由 chainreactors 开发的多指纹库聚合识别引擎。Fingers 格式是该引擎的原生格式，支持多种匹配方式和复杂的检测逻辑。

**特点**：
- 功能最完整
- 支持多种匹配方式
- 支持复杂的逻辑组合
- 支持版本提取
- 支持多个检测位置

### 格式规范

Fingers 格式是一个 JSON 对象数组，每个对象代表一个指纹：

```json
[
  {
    "name": "系统名称",
    "version": "版本号",
    "category": "分类",
    "website": "官网",
    "icon": "图标",
    "headers": {
      "X-Powered-By": "特征"
    },
    "html": ["特征1", "特征2"],
    "scripts": ["脚本路径"],
    "cookies": {
      "cookie名": "值"
    },
    "meta": {
      "generator": "特征"
    }
  }
]
```

### 示例

```json
[
  {
    "name": "WordPress",
    "version": "5.0",
    "category": "CMS",
    "website": "https://wordpress.org",
    "headers": {
      "X-Powered-By": "WordPress"
    },
    "html": [
      "wp-content",
      "wp-includes"
    ],
    "scripts": [
      "/wp-includes/js/"
    ]
  }
]
```

---

## 5. FingerPrintHub 格式

### 简介

**FingerPrintHub** 是一个指纹中心项目，使用 **Nuclei 模板格式**。Nuclei 是由 ProjectDiscovery 开发的现代化漏洞扫描工具，使用 YAML 或 JSON 格式定义检测模板。

**特点**：
- 基于 Nuclei 模板格式
- 支持复杂的匹配器（matchers）
- 支持多种匹配类型（word、regex、favicon、status-code 等）
- 支持提取器（extractors）
- 支持条件逻辑（AND、OR）
- 最灵活和强大的格式

### 格式规范

```json
[
  {
    "id": "指纹ID",
    "info": {
      "name": "系统名称",
      "author": "作者",
      "tags": "标签",
      "severity": "info|low|medium|high|critical",
      "metadata": {
        "product": "产品名",
        "vendor": "厂商名"
      }
    },
    "http": [
      {
        "method": "GET|POST|PUT|DELETE",
        "path": ["{{BaseURL}}/path"],
        "headers": {
          "User-Agent": "Mozilla/5.0"
        },
        "matchers": [
          {
            "type": "word|regex|status-code|favicon",
            "words": ["特征字符串"],
            "case-insensitive": true
          }
        ]
      }
    ]
  }
]
```

### 字段说明

| 字段 | 说明 | 示例 |
|------|------|------|
| id | 指纹 ID | "wordpress-detect" |
| name | 系统名称 | "WordPress" |
| author | 作者 | "test" |
| tags | 标签 | "detect,tech,wordpress" |
| severity | 严重级别 | "info" |
| product | 产品名 | "WordPress" |
| vendor | 厂商名 | "WordPress" |
| method | HTTP 方法 | "GET"、"POST" |
| path | 请求路径 | ["{{BaseURL}}/"] |
| matchers | 匹配器数组 | 见下表 |

### Matcher 类型

| 类型 | 说明 | 示例 |
|------|------|------|
| word | 字符串匹配 | {"type": "word", "words": ["wp-content"]} |
| regex | 正则表达式匹配 | {"type": "regex", "regex": ["WordPress\\s+([\\d.]+)"]} |
| status-code | HTTP 状态码 | {"type": "status-code", "status": [200, 301]} |
| favicon | Favicon hash | {"type": "favicon", "hash": ["hash值"]} |

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
            "words": ["wp-content", "wp-includes"],
            "case-insensitive": true
          }
        ]
      }
    ]
  }
]
```

---

## 格式对比

| 特性 | EHole | Goby | Wappalyzer | Fingers | FingerPrintHub |
|------|-------|------|-----------|---------|----------------|
| 格式 | JSON | JSON | JSON | JSON | JSON |
| 复杂度 | 低 | 中 | 中 | 高 | 高 |
| 逻辑组合 | 否 | 是 | 否 | 是 | 是 |
| 版本提取 | 否 | 否 | 是 | 是 | 是 |
| 正则支持 | 是 | 是 | 是 | 是 | 是 |
| 多位置检测 | 是 | 是 | 是 | 是 | 是 |
| 学习难度 | 低 | 中 | 中 | 高 | 高 |
| 应用场景 | 简单指纹 | 中等复杂 | Web 技术 | 复杂指纹 | 高级检测 |

---

## 选择建议

### 何时使用各种格式

1. **EHole 格式**
   - 适合快速创建简单指纹
   - 适合初学者
   - 适合轻量级指纹库

2. **Goby 格式**
   - 适合需要逻辑组合的指纹
   - 适合中等复杂度的检测
   - 适合 Goby 用户

3. **Wappalyzer 格式**
   - 适合 Web 技术识别
   - 适合浏览器扩展集成
   - 适合通用技术栈检测

4. **Fingers 格式**
   - 适合复杂的指纹检测
   - 适合需要多种匹配方式的场景
   - 适合 Fingers 引擎用户

5. **FingerPrintHub 格式**
   - 适合最复杂的检测场景
   - 适合需要高级功能的用户
   - 适合 Nuclei 模板用户
   - 适合需要提取器和复杂逻辑的场景

---

## 参考资源

- [EHole GitHub](https://github.com/EdgeSecurityTeam/EHole)
- [Goby 官网](https://www.gobysec.com/)
- [Wappalyzer GitHub](https://github.com/AliasIO/wappalyzer)
- [Fingers GitHub](https://github.com/chainreactors/fingers)
- [Nuclei 文档](https://docs.projectdiscovery.io/templates/introduction)
- [FingerPrintHub GitHub](https://github.com/0x727/FingerprintHub)
