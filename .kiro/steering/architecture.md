---
inclusion: always
---

# CICY 项目架构

## 系统架构

```
┌─────────────────────────────────────────────────┐
│              CICY TUI Client                    │
│  (client-tui.js - blessed + axios)              │
│                                                 │
│  ┌──────────────┐      ┌──────────────────┐    │
│  │ 本地消息通信  │      │ 远程浏览器自动化  │    │
│  │ (MCP)        │      │ (curl-rpc)       │    │
│  └──────┬───────┘      └────────┬─────────┘    │
└─────────┼──────────────────────┼───────────────┘
          │                      │
          ▼                      ▼
┌─────────────────┐    ┌──────────────────────┐
│ Local MCP Server│    │ Remote electron-mcp  │
│ (server.js)     │    │ gcp-docs.cicy.de5.net│
│ Port: 13001     │    │ HTTPS                │
└─────────────────┘    └──────────────────────┘
```

## 数据流

### 本地消息流
```
User Input → client-tui.js → MCP Request → server.js
                                              ↓
                                         Store Message
                                              ↓
User ← Display ← MCP Response ← server.js ←─┘
```

### 远程自动化流
```
User: "curl-rpc ping"
  ↓
client-tui.js (检测 curl-rpc 前缀)
  ↓
exec("curl-rpc ping")
  ↓
curl-rpc 脚本
  ↓
HTTPS POST → gcp-docs.cicy.de5.net/rpc/tools/call
  ↓
Response ← electron-mcp server
  ↓
Display in TUI
```

## 模块职责

### server.js
- MCP 协议实现
- JSON-RPC 2.0 接口
- 消息存储（内存）
- 图片处理（base64）

### client-tui.js
- TUI 界面（blessed）
- 本地 MCP 客户端
- curl-rpc 命令转发
- 命令建议弹窗

### curl-rpc
- 远程 RPC 调用
- Token 认证
- JSON 响应解析

## 扩展点

1. **新增本地工具**: 在 server.js 的 tools 数组添加
2. **新增 TUI 功能**: 在 client-tui.js 添加命令处理
3. **新增远程工具**: 使用 curl-rpc 调用 electron-mcp
