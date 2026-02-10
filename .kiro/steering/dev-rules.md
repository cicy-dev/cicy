---
inclusion: always
---

# CICY 项目开发规范

## 项目简介
基于 MCP (Model Context Protocol) 的消息通信系统，支持文本和图片传输。

## 核心功能
1. **本地 MCP 服务器** (端口 13001)
   - 文本消息发送/接收
   - Base64 图片传输
   - 5 个 MCP 工具

2. **远程浏览器自动化** (通过 curl-rpc)
   - 连接到 https://gcp-docs.cicy.de5.net
   - 支持 electron-mcp 所有工具

## 技术栈
- **服务器**: Express.js + MCP 协议
- **客户端**: blessed (TUI) + axios
- **协议**: JSON-RPC 2.0

## 开发规则

### 1. 代码风格
- 最小化实现，只写必要代码
- 使用中文注释
- 保持函数简洁

### 2. 文件结构
```
cicy/
├── server.js           # MCP 服务器
├── client-tui.js       # TUI 客户端（主要）
├── client.js           # 基础 TUI
├── client-mcp.js       # MCP 命令行
├── client-simple.js    # 最简客户端
├── package.json
└── README.md
```

### 3. 开发流程
- 使用 `npm run dev:tui` 启动热重载
- 修改代码自动重启
- 测试后更新 README

### 4. 命令规范
客户端支持的命令：
- `/quit` 或 `/q` - 退出
- `/list` - 查看消息
- `/clear` - 清空屏幕
- `/test-image` - 测试图片
- `curl-rpc <command>` - 远程自动化

### 5. 功能扩展原则
- 新功能先在 client-tui.js 实现
- 保持向后兼容
- 更新文档和帮助信息
- 添加到命令建议列表

## 常用命令

### 开发
```bash
npm run dev:server    # 服务器热重载
npm run dev:tui       # 客户端热重载
```

### 测试
```bash
npm start             # 启动服务器
npm run client:tui    # 启动客户端
```

## 远程服务器
- URL: https://gcp-docs.cicy.de5.net
- Token: ~/electron-mcp-token.txt
- 工具: electron-mcp 所有工具

## 注意事项
1. 不要修改 MCP 协议实现
2. 保持客户端简洁易用
3. 图片使用 base64 编码
4. 错误处理要友好
