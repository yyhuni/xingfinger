# 自定义指纹文件

此目录用于存放自定义指纹文件。

## 文件格式

支持以下格式的指纹文件：

- `*_ehole.json` - EHole 格式
- `*_goby.json` - Goby 格式
- `*_wappalyzer.json` - Wappalyzer 格式
- `*_fingers.json` - Fingers 原生格式
- `*_fingerprinthub.json` - FingerPrintHub 格式

## 使用方法

```bash
# 使用单个自定义指纹
xingfinger -u https://example.com --ehole fingerprints/custom_ehole.json

# 使用多个自定义指纹
xingfinger -u https://example.com --ehole fingerprints/custom_ehole.json --goby fingerprints/custom_goby.json
```

## EHole 格式示例

```json
{
  "fingerprint": [
    {
      "cms": "系统名称",
      "method": "keyword",
      "location": "body",
      "keyword": ["特征字符串1", "特征字符串2"]
    },
    {
      "cms": "另一个系统",
      "method": "regular",
      "location": "title",
      "keyword": ["正则表达式"]
    },
    {
      "cms": "第三个系统",
      "method": "faviconhash",
      "location": "body",
      "keyword": ["favicon_hash_value"]
    }
  ]
}
```

## Goby 格式示例

```json
[
  {
    "name": "系统名称",
    "logic": "a",
    "rule": [
      {
        "label": "a",
        "feature": "特征字符串",
        "is_equal": true
      }
    ]
  },
  {
    "name": "另一个系统",
    "logic": "a||b",
    "rule": [
      {
        "label": "a",
        "feature": "特征1",
        "is_equal": true
      },
      {
        "label": "b",
        "feature": "特征2",
        "is_equal": true
      }
    ]
  }
]
```

### Goby 格式说明

- `name`: 系统/应用名称
- `logic`: 逻辑表达式，使用标签组合（支持 `||`、`&&`、括号）
- `rule`: 规则数组
  - `label`: 规则标签（在 logic 中引用）
  - `feature`: 要匹配的特征字符串
  - `is_equal`: `true` 表示包含该特征，`false` 表示不包含该特征

## 支持的匹配方法

| 方法 | 说明 | 示例 |
|------|------|------|
| `keyword` | 关键字匹配（所有关键字都要匹配） | `["admin", "login"]` |
| `regular` | 正则表达式匹配 | `["admin.*panel"]` |
| `faviconhash` | Favicon hash 匹配 | `["1234567890"]` |

## 支持的检测位置

| 位置 | 说明 |
|------|------|
| `body` | 响应体内容 |
| `header` | HTTP 响应头 |
| `title` | 页面标题 |

## 提示

- 自定义指纹会与内置指纹一起使用
- 如果有重复的指纹，结果会自动去重
- 支持同时加载多个自定义指纹文件
