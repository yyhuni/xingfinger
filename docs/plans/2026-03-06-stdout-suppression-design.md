# Stdout 静默初始化设计

## 背景

`pkg/scanner.go` 在静默初始化 `fingers.NewEngine()` 时，会把 `os.Stdout` 临时重定向到 `os.DevNull`。当前实现打开了 `os.DevNull`，但没有关闭该文件句柄。

## 目标

- 修复 `os.DevNull` 文件描述符泄露
- 不改变扫描结果与指纹匹配语义
- 保持 `silent` / `jsonOutput` 下的静默初始化行为
- 为资源管理行为补充回归测试

## 方案

采用一个很小的辅助函数封装“临时静默 stdout”流程：

1. 打开 `os.DevNull`
2. 保存旧的 `os.Stdout`
3. 临时切换 `os.Stdout`
4. 执行初始化函数
5. 恢复 `os.Stdout`
6. 关闭 `os.DevNull`

## 影响评估

这是资源生命周期修复，不调整扫描策略、请求策略和结果判定，预期无业务语义变化。
