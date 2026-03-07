# XingFinger

![Author](https://img.shields.io/badge/Author-yyhuni-green) ![language](https://img.shields.io/badge/language-Golang-green) ![Go Version](https://img.shields.io/badge/Go-1.26+-blue)

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
- 📝 **多种输出格式** - 支持终端 JSON 输出、文件导出和静默模式
- 🔧 **自定义指纹** - 支持加载自定义指纹文件，默认与内置指纹叠加使用
- 🌐 **ARL 指纹支持** - 支持灯塔 ARL YAML 格式指纹（9000+ 条规则）

## 安装

> 需要 Go 1.26 或更高版本。

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
xingfinger -l urls.txt

# 终端输出 JSON 格式（方便管道处理）
xingfinger -l urls.txt -j

# 保存结果到 JSON 文件
xingfinger -l urls.txt -o result.json

# 设置并发线程数
xingfinger -l urls.txt -t 100

# 使用代理
xingfinger -l urls.txt -p http://127.0.0.1:8080

# 静默模式（只输出命中结果）
xingfinger -l urls.txt -s

# 默认不跟跳转；只看原始响应
xingfinger -u www.baidu.com

# 只跟随 HTTP 3xx 跳转
xingfinger -u www.baidu.com --redirect-policy http

# 跟随 HTTP 3xx + JS / meta refresh 跳转
xingfinger -u www.baidu.com --redirect-policy all

# 使用自定义指纹（与默认指纹叠加）
xingfinger -u https://example.com --ehole my_ehole.json

# 使用 ARL 指纹（与默认指纹叠加）
xingfinger -u https://example.com --arl fingerprints/ARL.yaml

# 禁用默认指纹，仅使用自定义指纹
xingfinger -u https://example.com --no-default --arl fingerprints/ARL.yaml

# JSON 输出配合 jq 过滤
xingfinger -l urls.txt -j | jq 'select(.cms | contains("shiro"))'
```

## 参数说明

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-u, --url` | 目标 URL | - |
| `-l, --list` | URL 列表文件 | - |
| `-t, --thread` | 并发线程数 | 50 |
| `--timeout` | 请求超时时间（秒） | 10 |
| `-o, --output` | 输出文件路径（JSON 格式） | - |
| `-p, --proxy` | 代理地址 | - |
| `--redirect-policy` | 跳转策略：`never` 不跟跳转，`http` 只跟 HTTP 3xx，`all` 跟 HTTP 3xx + JS / `meta refresh` | `never` |
| `-s, --silent` | 静默模式，只输出命中结果 | false |
| `-j, --json` | 终端输出 JSON 格式 | false |
| `--no-default` | 禁用默认指纹，仅使用自定义指纹 | false |
| `--ehole` | 自定义 EHole 指纹文件 | - |
| `--goby` | 自定义 Goby 指纹文件 | - |
| `--wappalyzer` | 自定义 Wappalyzer 指纹文件 | - |
| `--fingers` | 自定义 Fingers 指纹文件 | - |
| `--fingerprint` | 自定义 FingerPrintHub 指纹文件 | - |
| `--arl` | 自定义 ARL YAML 指纹文件 | - |

## 自定义指纹

支持加载自定义指纹文件，格式与对应的指纹库一致。自定义指纹默认与内置指纹**叠加使用**，如需禁用内置指纹，请使用 `--no-default` 参数。

默认 `--redirect-policy` 为 `never`，即不跟随任何跳转，因此一个输入目标默认只输出原始页面结果；如需跟随 HTTP 3xx 或内容跳转，请显式指定 `http` 或 `all`。

指纹文件示例见 `fingerprints/` 目录。

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

**ARL YAML 格式示例**：
```yaml
- name: nginx
  rule: header="nginx"

- name: apache
  rule: header="Apache" && body="Apache Server"

- name: wordpress
  rule: body="wp-content" && body="wp-includes"

- name: favicon_example
  rule: icon_hash="116323821"
```

- 支持的条件类型：`body`、`header`、`title`、`icon_hash`
- 多个条件使用 `&&` 连接，表示 AND 关系
- 大小写不敏感匹配

## 指纹库说明

| 指纹库 | 说明 | 规则数量 |
|--------|------|----------|
| fingers | chainreactors 原生指纹库 | 内置 |
| wappalyzer | Web 技术栈检测 | 内置 |
| fingerprinthub | 指纹中心 | 内置 |
| ehole | 棱洞指纹库 | 内置 |
| goby | Goby 指纹库 | 内置 |
| ARL | 灯塔 ARL 指纹库 | 9264 条 |

## 输出格式

### 终端 JSON 输出 (`-j`)

每行一个 JSON 对象，方便管道处理：

```json
{"url":"https://example.com","cms":"nginx,php","server":"nginx/1.18.0","status_code":200,"length":12345,"title":"Example"}
```

### 文件输出 (`-o`)

JSON 数组格式：

```json
[
  {
    "url": "https://example.com",
    "cms": "WordPress,PHP",
    "server": "nginx/1.18.0",
    "status_code": 200,
    "length": 12345,
    "title": "Example Site"
  }
]
```

| 字段 | 类型 | 说明 |
|------|------|------|
| `url` | string | 目标 URL |
| `cms` | string | 检测到的指纹，多个用逗号分隔 |
| `server` | string | Server 响应头 |
| `status_code` | int | HTTP 状态码 |
| `length` | int | 响应体长度 |
| `title` | string | 页面标题 |

## 参考项目

- [chainreactors/fingers](https://github.com/chainreactors/fingers) - 多指纹库聚合识别引擎
- [EdgeSecurityTeam/EHole](https://github.com/EdgeSecurityTeam/EHole) - 红队重点攻击系统指纹探测工具
- [TophantTechnology/ARL](https://github.com/TophantTechnology/ARL) - 灯塔资产灯塔系统

## License

MIT License
