# 自定义指纹文档索引

本目录包含了关于 xingfinger 自定义指纹的完整文档和示例。

## 📚 文档导航

### 快速开始
- **[README.md](README.md)** - 自定义指纹使用说明
  - 如何创建自定义指纹
  - 如何加载自定义指纹
  - 常见问题解答

### 详细指南
- **[FINGERPRINT_FORMATS.md](FINGERPRINT_FORMATS.md)** - 指纹格式详解
  - 5 种指纹格式的详细说明
  - 格式规范和字段说明
  - 示例代码
  - 格式对比和选择建议

- **[FINGERPRINT_RESEARCH.md](FINGERPRINT_RESEARCH.md)** - 指纹格式研究报告
  - 5 种格式的来源和背景
  - 格式演进历程
  - 技术趋势分析
  - 参考资源

### 测试结果
- **[TEST_RESULTS.md](TEST_RESULTS.md)** - 自定义指纹测试结果
  - 所有 5 种格式的测试结果
  - 测试命令和输出
  - 关键发现

### 示例对比
- **[EXAMPLES_COMPARISON.md](EXAMPLES_COMPARISON.md)** - 指纹格式示例对比
  - 5 种格式的示例对比分析
  - 每种格式的独特特性
  - 功能对比表
  - 选择建议

## 📋 示例文件

### 改进的示例文件

所有示例都使用真实的 CMS 系统（**WordPress、Joomla、Drupal**）来展示每种格式的独特特性。

详见 [EXAMPLES_COMPARISON.md](EXAMPLES_COMPARISON.md) 了解详细对比。

### EHole 格式
- **[custom_ehole.json](custom_ehole.json)** - EHole 格式示例
  - 展示了 keyword、regular 匹配方式
  - 展示了 body、header 检测位置
  - 适合简单指纹

### Goby 格式
- **[custom_goby.json](custom_goby.json)** - Goby 格式示例
  - 展示了 OR 逻辑（a|b|c）
  - 展示了 AND 逻辑（a&b）
  - 适合中等复杂度指纹

### Wappalyzer 格式
- **[custom_wappalyzer.json](custom_wappalyzer.json)** - Wappalyzer 格式示例
  - 展示了多种检测方式（headers、html、scripts、meta）
  - 展示了 implies 技术依赖
  - 适合 Web 技术识别

### Fingers 格式
- **[custom_fingers.json](custom_fingers.json)** - Fingers 格式示例
  - 展示了完整的检测方式（headers、html、scripts、cookies、meta）
  - 展示了元数据（category、website）
  - 适合复杂指纹

### FingerPrintHub 格式
- **[custom_fingerprinthub.json](custom_fingerprinthub.json)** - FingerPrintHub 格式示例
  - 展示了多种 matcher 类型（word、regex）
  - 展示了 part 字段指定检测位置
  - 适合高级检测

## 🚀 快速使用

### 1. 选择合适的格式

根据你的需求选择合适的格式：

| 需求 | 推荐格式 | 原因 |
|------|---------|------|
| 快速开始 | EHole | 简单易用 |
| 中等复杂度 | Goby 或 Wappalyzer | 功能足够 |
| 复杂场景 | Fingers 或 FingerPrintHub | 功能完整 |

### 2. 参考示例文件

查看对应格式的示例文件，了解格式规范。

### 3. 创建自己的指纹

基于示例文件创建自己的指纹文件。

### 4. 加载自定义指纹

```bash
# 使用单个自定义指纹
./xingfinger -u https://example.com --ehole fingerprints/custom_ehole.json

# 使用多个自定义指纹
./xingfinger -u https://example.com \
  --ehole fingerprints/custom_ehole.json \
  --goby fingerprints/custom_goby.json
```

## 📖 详细文档

### 格式选择指南

**EHole 格式** - 最简单
```json
{
  "fingerprint": [
    {
      "cms": "系统名称",
      "method": "keyword",
      "location": "body",
      "keyword": ["特征字符串"]
    }
  ]
}
```

**Goby 格式** - 支持逻辑
```json
[
  {
    "name": "系统名称",
    "logic": "a|b",
    "rule": [
      {
        "label": "a",
        "feature": "特征字符串",
        "is_equal": false
      }
    ]
  }
]
```

**Wappalyzer 格式** - 灵活多样
```json
{
  "系统名称": {
    "cats": [1, 6],
    "headers": {"X-Powered-By": "系统名称"},
    "html": "特征字符串"
  }
}
```

**Fingers 格式** - 功能完整
```json
[
  {
    "name": "系统名称",
    "category": "CMS",
    "html": ["特征1", "特征2"],
    "scripts": ["/脚本路径/"]
  }
]
```

**FingerPrintHub 格式** - 最强大
```json
[
  {
    "id": "指纹ID",
    "info": {
      "name": "系统名称",
      "metadata": {"product": "产品名"}
    },
    "http": [
      {
        "method": "GET",
        "path": ["{{BaseURL}}/"],
        "matchers": [
          {
            "type": "word",
            "words": ["特征字符串"]
          }
        ]
      }
    ]
  }
]
```

## 🔍 常见问题

### Q: 如何选择指纹格式？
A: 根据复杂度选择：
- 简单指纹 → EHole
- 中等复杂度 → Goby 或 Wappalyzer
- 复杂指纹 → Fingers 或 FingerPrintHub

### Q: 可以同时使用多个格式吗？
A: 可以，使用多个参数即可：
```bash
./xingfinger -u url --ehole file1.json --goby file2.json
```

### Q: 自定义指纹会覆盖默认指纹吗？
A: 不会，自定义指纹会与默认指纹一起使用。

### Q: 如何验证自定义指纹是否正确？
A: 使用测试服务器验证，参考 TEST_RESULTS.md。

### Q: 支持哪些文件格式？
A: 支持 `.json` 和 `.json.gz` 格式。

## 📚 相关资源

### 官方项目
- [EHole](https://github.com/EdgeSecurityTeam/EHole)
- [Goby](https://www.gobysec.com/)
- [Wappalyzer](https://github.com/AliasIO/wappalyzer)
- [Fingers](https://github.com/chainreactors/fingers)
- [Nuclei](https://github.com/projectdiscovery/nuclei)

### 文档
- [Nuclei 文档](https://docs.projectdiscovery.io/templates/introduction)
- [Wappalyzer 指纹配置](https://neptliang.github.io/2021/01/11/Wappalyzer-Fingerprint-Configuration/)

## 📝 文件清单

```
fingerprints/
├── INDEX.md                    # 本文件 - 文档索引
├── README.md                   # 使用说明
├── FINGERPRINT_FORMATS.md      # 格式详解
├── FINGERPRINT_RESEARCH.md     # 研究报告
├── EXAMPLES_COMPARISON.md      # 示例对比分析 ⭐ 新增
├── TEST_RESULTS.md             # 测试结果
├── custom_ehole.json           # EHole 示例（改进）
├── custom_goby.json            # Goby 示例（改进）
├── custom_wappalyzer.json      # Wappalyzer 示例（改进）
├── custom_fingers.json         # Fingers 示例（改进）
└── custom_fingerprinthub.json  # FingerPrintHub 示例（改进）
```

## 🎯 下一步

1. **阅读 README.md** - 了解基本使用方法
2. **查看示例文件** - 选择合适的格式
3. **参考 FINGERPRINT_FORMATS.md** - 了解格式规范
4. **创建自己的指纹** - 基于示例创建
5. **测试验证** - 使用 xingfinger 验证

## 💡 提示

- 从简单格式开始，逐步学习复杂格式
- 参考示例文件，快速上手
- 使用测试服务器验证指纹
- 查看 TEST_RESULTS.md 了解测试方法

---

**最后更新**：2025-12-26

**版本**：1.0

**作者**：xingfinger 项目
