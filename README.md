# CICY - MCP 消息通信系统

基于 Model Context Protocol (MCP) 的文本和图片消息传输系统。

## 快速开始

### 安装依赖
```bash
npm install
```

### 启动服务器
```bash
npm start
# 服务器运行在 http://localhost:13001
```

### 启动客户端
```bash
# CLI 客户端（推荐 - 完全兼容 tmux）
node client-cli.js

# 或使用脚本
npm run client:cli

# TUI 客户端（图形界面，建议直接运行）
npm run dev:tui
# 注意：TUI 在 tmux 中可能有输入问题，建议在独立终端运行
```

## 功能特性

- ✅ MCP 协议 (2024-11-05)
- ✅ JSON-RPC 2.0 接口
- ✅ 文本消息发送/接收
- ✅ Base64 图片传输
- ✅ CLI 命令行界面
- ✅ 实时状态显示
- ✅ 完全兼容 tmux

## 客户端选项

| 文件 | 说明 | 界面 | tmux 兼容 | 推荐场景 |
|------|------|------|-----------|----------|
| `tui-go/cicy-tui` | **⚡ Go 版本** 高性能 TUI | Bubble Tea | ✅ | 生产环境 |
| `client-remote.js` | **🌐 远程客户端** 连接远程服务器 | blessed | ✅ | 远程连接 |
| `client-cli.js` | **推荐** 简单命令行 | readline | ✅ | tmux/脚本 |
| `client-tui.js` | 完整 TUI，支持图片 | blessed | ⚠️ | 独立终端 |
| `client-mcp.js` | MCP 命令行客户端 | readline | ✅ | 开发调试 |
| `client-simple.js` | 最简单客户端 | readline | ✅ | 学习示例 |
| `client.js` | 基础 TUI，仅文本 | blessed | ⚠️ | 已废弃 |

**注意**：
- ✅ = 完全兼容
- ⚠️ = 在 tmux 中可能有输入问题，建议直接运行
- ⚡ Go 版本：启动快、内存低、单文件部署
- 🌐 远程客户端：可连接到任何 CICY 服务器
- TUI 客户端功能：实时统计、Loading 动画、统计窗口

### 远程客户端 (连接远程服务器)

```bash
# 连接到本地服务器
npm run client:remote

# 连接到远程服务器
node client-remote.js http://remote-server:13001

# 使用环境变量
export CICY_REMOTE_URL=http://remote-server:13001
npm run client:remote
```

特性：
- 🌐 连接到任何 CICY 服务器
- 📊 实时连接状态显示
- ⚡ 自动重连机制
- 🎨 Tokyo Night 配色
- 📝 命令支持（/status, /help, /quit）

命令：
- `/status` - 检查连接状态
- `/help` - 显示帮助
- `/quit` - 退出

### Go 客户端 (推荐生产环境)

```bash
cd tui-go
./test.sh              # 测试并编译
./cicy-tui             # 运行
```

特性：
- ⚡ 启动速度极快（<10ms）
- 💚 内存占用低（~10MB）
- 📦 单文件部署，无需依赖
- 🎨 Masu 风格界面
- ⏱️ 完成耗时显示
- 🔄 异步 Loading 动画

## 命令

在客户端中输入：

- `/test-image` - 发送测试图片
- `/list` - 查看所有消息
- `/clear` - 清空屏幕
- `/quit` 或 `/q` - 退出
- `Ctrl+C` - 退出

## MCP 工具

服务器提供 5 个 MCP 工具：

1. `send_message` - 发送文本消息
2. `send_image` - 发送 base64 图片
3. `get_messages` - 获取所有消息
4. `get_images` - 获取所有图片
5. `clear_messages` - 清空消息

## API 端点

### MCP 协议
```bash
POST http://localhost:13001/mcp
Content-Type: application/json

{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "send_message",
    "arguments": { "message": "Hello" }
  }
}
```

### REST API (传统)
```bash
# 发送消息
POST http://localhost:13001/message
{"message": "Hello"}

# 获取消息
GET http://localhost:13001/messages

# 健康检查
GET http://localhost:13001/health
```

## 图片传输

图片使用 base64 编码传输：

```javascript
const imageBase64 = fs.readFileSync('image.png').toString('base64');

await mcpRequest('tools/call', {
    name: 'send_image',
    arguments: {
        name: 'image.png',
        data: imageBase64,
        mimeType: 'image/png'
    }
});
```

## 技术栈

- **服务器**: Express.js
- **客户端**: blessed (TUI), axios (HTTP)
- **协议**: MCP, JSON-RPC 2.0

## 开发

```bash
# 开发模式（自动重启）
npm run dev:server  # 服务器
npm run dev:client  # 客户端
npm run dev         # 同时启动
```

## 端口

默认端口: `13001`

修改端口:
```bash
PORT=3000 npm start
```
