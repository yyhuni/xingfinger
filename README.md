# XingFinger

![Author](https://img.shields.io/badge/Author-yyhuni-green) ![language](https://img.shields.io/badge/language-Golang-green)

```
  __  ___                _____ _                       
  \ \/ (_)___  ____ _   / ____(_)___  ____ ____  _____ 
   \  /| / _ \/ __ `/  / /_  / / __ \/ __ `/ _ \/ ___/ 
   /  \| |  __/ /_/ /  / __/ / / / / / /_/ /  __/ /     
  /_/\_\_|\___/\__, /  /_/   /_/_/ /_/\__, /\___/_/      
              /____/                 /____/   By:yyhuni
```

XingFinger 是一款 Web 指纹识别工具，基于 [chainreactors/fingers](https://github.com/chainreactors/fingers) 多指纹库聚合引擎，帮助红队人员快速识别目标系统的技术栈。

## 特性

- 🔍 **多指纹库聚合** - 集成 fingers、wappalyzer、fingerprinthub、ehole、goby 等指纹库，2888+ 指纹规则
- 🚀 **高性能并发** - 支持自定义线程数，快速扫描大量目标
- 🔄 **指纹自动更新** - 支持从 GitHub 下载最新指纹库
- 🎯 **Favicon 识别** - 主动获取 favicon 进行 hash 匹配
- 📝 **多种输出格式** - 支持 JSON 导出和静默模式

## 安装

```bash
go install github.com/yyhuni/xingfinger@latest
```

或从源码编译：

```bash
git clone https://github.com/yyhuni/xingfinger.git
cd xingfinger
go build -o xingfinger
```

## 使用

```bash
# 单目标扫描
xingfinger -u https://example.com

# 批量扫描
xingfinger -l urls.txt

# 输出到 JSON 文件
xingfinger -l urls.txt -o result.json

# 设置并发线程数
xingfinger -l urls.txt -t 50

# 使用代理
xingfinger -l urls.txt -p http://127.0.0.1:8080

# 静默模式（只输出命中结果）
xingfinger -l urls.txt --silent

# 使用自定义指纹
xingfinger -u https://example.com --ehole my_ehole.json
```

## 参数说明

| 参数 | 说明 |
|------|------|
| `-u, --url` | 单个目标 URL |
| `-l, --list` | URL 列表文件 |
| `-o, --output` | 输出文件路径（JSON 格式） |
| `-t, --thread` | 并发线程数（默认 100） |
| `-p, --proxy` | 代理地址 |
| `--timeout` | 请求超时时间（秒，默认 10） |
| `--silent` | 静默模式 |
| `--ehole` | 自定义 EHole 格式指纹文件 |
| `--goby` | 自定义 Goby 格式指纹文件 |
| `--wappalyzer` | 自定义 Wappalyzer 格式指纹文件 |
| `--fingers` | 自定义 Fingers 原生格式指纹文件 |
| `--fingerprinthub` | 自定义 FingerPrintHub 格式指纹文件 |

## 自定义指纹

支持加载自定义指纹文件，格式与对应的指纹库一致：

```bash
# 使用自定义 EHole 格式指纹
xingfinger -u https://example.com --ehole fingerprints/custom_ehole.json

# 同时使用多个自定义指纹
xingfinger -u https://example.com --ehole fingerprints/custom_ehole.json --goby fingerprints/custom_goby.json
```

自定义指纹文件放在 `fingerprints/` 目录下，详见 [fingerprints/README.md](fingerprints/README.md)。

EHole 格式示例：
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

支持的 method: `keyword`、`regular`、`faviconhash`
支持的 location: `body`、`header`、`title`

## 参考项目

本项目参考或使用了以下优秀的开源项目：

- [chainreactors/fingers](https://github.com/chainreactors/fingers) - 多指纹库聚合识别引擎，提供核心指纹识别能力
- [chainreactors/spray](https://github.com/chainreactors/spray) - 目录爆破工具，参考了指纹更新机制
- [EdgeSecurityTeam/EHole](https://github.com/EdgeSecurityTeam/EHole) - 红队重点攻击系统指纹探测工具，参考了项目结构和 JS 跳转检测逻辑

## 指纹库说明

XingFinger 使用 fingers 引擎聚合了多个指纹库：

| 指纹库 | 说明 |
|--------|------|
| fingers | chainreactors 自有指纹库 |
| wappalyzer | Web 技术检测 |
| fingerprinthub | 指纹中心 |
| ehole | 棱洞指纹 |
| goby | Goby 指纹库 |

## License

MIT License
