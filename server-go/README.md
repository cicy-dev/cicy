# CICY Go - 一体化版本

单进程同时运行 TUI 客户端和 MCP 服务器的 Go 实现。

## 特性

- 🚀 单进程运行，无需分别启动服务器和客户端
- ⚡ 高性能 Go 实现
- 💾 低内存占用 (~15MB)
- 🎨 Tokyo Night 配色主题
- ⏱️ 实时显示消息耗时

## 安装

```bash
cd server-go
go mod download
go build -o cicy-go
```

## 使用

```bash
# 默认端口 13001
./cicy-go

# 自定义端口
./cicy-go -p 3000
./cicy-go --port 3000

# 查看帮助
./cicy-go -h
./cicy-go --help

# 查看版本
./cicy-go -v
./cicy-go --version
```

## 快捷键

- `Enter` - 发送消息
- `Ctrl+C` - 退出
- `ESC` - 退出

## 架构

```
┌─────────────────────────────────┐
│     CICY Go (单进程)             │
│                                 │
│  ┌──────────┐   ┌────────────┐ │
│  │   TUI    │   │ HTTP Server│ │
│  │  Client  │◄─►│  (MCP API) │ │
│  └──────────┘   └────────────┘ │
│                                 │
│  • Bubble Tea TUI               │
│  • net/http Server              │
│  • 内存消息存储                  │
└─────────────────────────────────┘
```

## API 端点

- `POST /mcp` - MCP JSON-RPC 接口
- `POST /message` - 发送消息 (Legacy REST)
- `GET /messages` - 获取所有消息
- `GET /health` - 健康检查

## 性能对比

| 指标 | Node.js | Go |
|------|---------|-----|
| 启动时间 | ~200ms | <10ms |
| 内存占用 | ~50MB | ~15MB |
| 二进制大小 | N/A | ~10MB |
| 并发性能 | 中 | 高 |

## 开发

```bash
# 运行
go run main.go

# 编译
go build -o cicy-go

# 编译优化版本
go build -ldflags="-s -w" -o cicy-go
```

## 与 Node.js 版本的区别

1. **单进程**: Go 版本在一个进程中同时运行 TUI 和服务器
2. **性能**: 启动更快，内存占用更低
3. **部署**: 单个二进制文件，无需 Node.js 环境
4. **并发**: 原生 goroutine 支持，并发性能更好

## 兼容性

完全兼容 Node.js 版本的 MCP 协议和 API 接口，可以与现有客户端互操作。
