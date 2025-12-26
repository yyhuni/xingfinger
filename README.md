# XingFinger

![Author](https://img.shields.io/badge/Author-yyhuni-green) ![language](https://img.shields.io/badge/language-Golang-green) ![Go Version](https://img.shields.io/badge/Go-1.15+-blue)

```
  __  ___                _____ _                       
  \ \/ (_)___  ____ _   / ____(_)___  ____ ____  _____ 
   \  /| / _ \/ __ `/  / /_  / / __ \/ __ `/ _ \/ ___/ 
   /  \| |  __/ /_/ /  / __/ / / / / / /_/ /  __/ /     
  /_/\_\_|\___/\__, /  /_/   /_/_/ /_/\__, /\___/_/      
              /____/                 /____/   By:yyhuni
```

XingFinger 是一款高效的 Web 指纹识别工具，基于 [chainreactors/fingers](https://github.com/chainreactors/fingers) 多指纹库聚合引擎，帮助安全人员快速识别目标系统的技术栈。

## 特性

- 🔍 **多指纹库聚合** - 集成 fingers、wappalyzer、fingerprinthub、ehole、goby 等指纹库
- 🚀 **高性能并发** - 支持自定义线程数，快速扫描大量目标
- 🎯 **Favicon 识别** - 主动获取 favicon 进行 hash 匹配
- 📝 **多种输出格式** - 支持 JSON 导出和静默模式
- � **自定义指纹*** - 支持加载自定义指纹文件

## 安装

**方式一：go install（推荐）**

```bash
go install github.com/yyhuni/xingfinger@latest
```

**方式二：源码编译**

```bash
git clone https://github.com/yyhuni/xingfinger.git
cd xingfinger
go build -o xingfinger .
```

**方式三：下载二进制**

从 [Releases](https://github.com/yyhuni/xingfinger/releases) 页面下载对应平台的二进制文件。

## 使用

```bash
# 单目标扫描
xingfinger -u https://example.com

# 批量扫描
xingfinger -f urls.txt

# 输出到 JSON 文件
xingfinger -f urls.txt -o result.json

# 设置并发线程数
xingfinger -f urls.txt -t 100

# 使用代理
xingfinger -f urls.txt -p http://127.0.0.1:8080

# 静默模式（只输出命中结果）
xingfinger -f urls.txt -s

# 使用自定义指纹
xingfinger -u https://example.com --ehole my_ehole.json
```

## 参数说明

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-u, --url` | 目标 URL | - |
| `-f, --file` | URL 列表文件 | - |
| `-t, --thread` | 并发线程数 | 50 |
| `--timeout` | 请求超时时间（秒） | 10 |
| `-o, --output` | 输出文件路径（JSON 格式） | - |
| `-p, --proxy` | 代理地址 | - |
| `-s, --silent` | 静默模式，只输出命中结果 | false |
| `--ehole` | 自定义 EHole 指纹文件 | - |
| `--goby` | 自定义 Goby 指纹文件 | - |
| `--wappalyzer` | 自定义 Wappalyzer 指纹文件 | - |
| `--fingers` | 自定义 Fingers 指纹文件 | - |
| `--fingerprint` | 自定义 FingerPrintHub 指纹文件 | - |

## 自定义指纹

支持加载自定义指纹文件，格式与对应的指纹库一致。指纹文件示例见 `fingerprints/` 目录。

**EHole 格式示例**：
```json
{
  "fingerprint": [
    {
      "cms": "系统名称",
      "method": "keyword",
      "location": "body",
      "keyword": ["特征字符串1", "特征字符串2"]
    }
  ]
}
```

- method: `keyword`（关键词）、`regular`（正则）、`faviconhash`（图标哈希）
- location: `body`、`header`、`title`
- keyword 数组中多个关键词为 AND 关系

## 指纹库说明

| 指纹库 | 说明 |
|--------|------|
| fingers | chainreactors 原生指纹库 |
| wappalyzer | Web 技术栈检测 |
| fingerprinthub | 指纹中心 |
| ehole | 棱洞指纹库 |
| goby | Goby 指纹库 |

## 参考项目

- [chainreactors/fingers](https://github.com/chainreactors/fingers) - 多指纹库聚合识别引擎
- [EdgeSecurityTeam/EHole](https://github.com/EdgeSecurityTeam/EHole) - 红队重点攻击系统指纹探测工具

## License

MIT License
