# CICY Go TUI Client

基于 Go + Bubble Tea 的 CICY 客户端，提供流畅的终端用户界面。

## 特性

- ✅ 异步消息发送
- ✅ Loading 动画（Spinner）
- ✅ 完成耗时显示
- ✅ Masu 风格界面
- ✅ Tokyo Night 配色
- ✅ 完全兼容 tmux
- ✅ 命令支持（/help, /quit, /clear, /list）
- ✅ 消息历史（显示最近 5 条）
- ✅ 错误提示

## 安装依赖

```bash
go mod tidy
```

## 编译

```bash
# macOS
GOOS=darwin GOARCH=amd64 go build -o cicy-tui

# Linux
GOOS=linux GOARCH=amd64 go build -o cicy-tui

# Windows
GOOS=windows GOARCH=amd64 go build -o cicy-tui.exe
```

## 运行

```bash
./cicy-tui
```

## 使用

### 基本操作

1. 输入消息后按 Enter 发送
2. 发送时显示 Loading 动画
3. 收到回复后显示完成耗时
4. 按 Ctrl+C 或 Esc 退出

### 命令列表

| 命令 | 说明 |
|------|------|
| `/help` 或 `/h` | 显示帮助信息 |
| `/quit` 或 `/q` | 退出程序 |
| `/clear` 或 `/c` | 清空消息历史 |
| `/list` 或 `/l` | 显示所有消息 |
| `Ctrl+C` | 退出程序 |
| `Esc` | 退出帮助/退出程序 |

### 界面效果

```
◇ CICY

test message 1
好的
 - Completed in 2.25s

test message 2
收到了
 - Completed in 1.83s

> _

Type /help for commands • Ctrl+C to quit
```

### 帮助界面

输入 `/help` 查看：

```
◇ CICY - Help

可用命令：

  /help, /h  -  显示此帮助信息
  /quit, /q  -  退出程序
  /clear, /c  -  清空消息历史
  /list, /l  -  显示所有消息
  Ctrl+C  -  退出程序
  Esc  -  退出帮助/退出程序

按 Esc 或 q 返回
```

## 技术栈

- **框架**: Bubble Tea (TUI framework)
- **组件**: Bubbles (textinput, spinner)
- **样式**: Lipgloss (styling)
- **HTTP**: net/http (标准库)

## 配色方案

Tokyo Night 主题：
- 标题：`#7aa2f7` (蓝色)
- 用户消息：`#f7768e` (红色)
- AI 回复：`#9ece6a` (绿色)
- 耗时信息：`#565f89` (灰色)
- 帮助信息：`#bb9af7` (紫色)
- 错误信息：`#f7768e` (红色)

## API 配置

默认连接到 `http://localhost:13001`

修改 API 地址：
```go
const API_URL = "http://your-server:port"
```

## 对比 Node.js 版本

| 特性 | Go 版本 | Node.js 版本 |
|------|---------|--------------|
| 启动速度 | ⚡ 极快 (<10ms) | 🐢 较慢 (~500ms) |
| 内存占用 | 💚 ~10MB | 💛 ~50MB |
| 编译产物 | 单文件 (9.5MB) | 需要 node_modules |
| 跨平台 | ✅ 编译即可 | ✅ 需要 Node.js |
| 依赖管理 | go.mod | package.json |
| 热重载 | ❌ 需重新编译 | ✅ 自动重启 |
| 命令支持 | ✅ 完整 | ✅ 完整 |
| 消息历史 | ✅ 最近 5 条 | ✅ 全部 |

## 功能特性

### 消息历史
- 自动保存所有消息
- 界面显示最近 5 条
- 使用 `/list` 查看全部

### 错误处理
- 服务器未运行提示
- 网络错误提示
- 未知命令提示

### 用户体验
- 异步请求不阻塞界面
- Loading 动画流畅
- 完成耗时精确显示
- 命令提示友好

## 开发

```bash
# 运行
go run main.go

# 格式化
go fmt

# 测试
go test

# 编译并运行
go build -o cicy-tui && ./cicy-tui
```

## 测试

```bash
# 运行测试脚本
./test.sh

# 输出示例：
# 🧪 CICY Go TUI 测试
# ===================
# 
# 检查服务器状态...
# ✅ 服务器正常
# 
# 编译 Go 客户端...
# ✅ 编译成功
# 
# 检查可执行文件...
# ✅ 可执行文件: cicy-tui (9.5M)
# 
# ===================
# 🎉 测试通过！
```

## 部署

### 单文件部署

```bash
# 编译
GOOS=linux GOARCH=amd64 go build -o cicy-tui

# 复制到服务器
scp cicy-tui user@server:/usr/local/bin/

# 运行
ssh user@server "cicy-tui"
```

### Docker 部署

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o cicy-tui

FROM alpine:latest
COPY --from=builder /app/cicy-tui /usr/local/bin/
CMD ["cicy-tui"]
```

## 故障排除

### 服务器未运行

```
Error: Post "http://localhost:13001/message": dial tcp [::1]:13001: connect: connection refused
```

解决：先启动服务器
```bash
cd .. && npm start
```

### 编译错误

```
go: golang.org/x/sys@v0.38.0 requires go >= 1.24.0
```

解决：更新 Go 版本或修改 go.mod
```bash
go mod tidy
```

### 中文显示问题

确保终端支持 UTF-8：
```bash
export LANG=en_US.UTF-8
```

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT

