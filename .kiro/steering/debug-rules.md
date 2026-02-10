---
inclusion: always
---

# CICY 调试规范

## Tmux 窗口规范

### 固定窗口
- **cicy:tui.0** - TUI 客户端（热重载）
- **cicy:api.0** - API 服务器（热重载）

### 严格规则
1. ❌ **禁止**打开新的 Terminal 窗口
2. ❌ **禁止**创建新的 tmux session
3. ✅ **只能**使用 `tmux send-keys` 向现有窗口发送命令
4. ✅ **只能**在这两个固定窗口中操作

## 启动命令

### TUI 客户端 (cicy:tui.0)
```bash
tmux send-keys -t cicy:tui.0 "cd /Users/ton/Desktop/skills/cicy && npm run dev:tui" C-m
```

### API 服务器 (cicy:api.0)
```bash
tmux send-keys -t cicy:api.0 "cd /Users/ton/Desktop/skills/cicy && npm run dev:server" C-m
```

## 重启流程

### 重启 TUI
```bash
# 1. 停止当前进程
tmux send-keys -t cicy:tui.0 C-c

# 2. 等待 1 秒
sleep 1

# 3. 重新启动
tmux send-keys -t cicy:tui.0 "npm run dev:tui" C-m
```

### 重启 API
```bash
# 1. 停止当前进程
tmux send-keys -t cicy:api.0 C-c

# 2. 等待 1 秒
sleep 1

# 3. 重新启动
tmux send-keys -t cicy:api.0 "npm run dev:server" C-m
```

## 热重载说明

### 自动重载
- TUI: 修改 `client-tui.js` 自动重启
- API: 修改 `server.js` 自动重启

### 手动重启
只在以下情况需要手动重启：
- 修改 `package.json`
- 安装新依赖
- nodemon 失效

## 查看日志

### 查看 TUI 输出
```bash
tmux select-window -t cicy:tui.0
```

### 查看 API 输出
```bash
tmux select-window -t cicy:api.0
```

## 常见操作

### 清屏
```bash
tmux send-keys -t cicy:tui.0 C-l
tmux send-keys -t cicy:api.0 C-l
```

### 停止所有
```bash
tmux send-keys -t cicy:tui.0 C-c
tmux send-keys -t cicy:api.0 C-c
```

### 启动所有
```bash
tmux send-keys -t cicy:api.0 "cd /Users/ton/Desktop/skills/cicy && npm run dev:server" C-m
sleep 2
tmux send-keys -t cicy:tui.0 "cd /Users/ton/Desktop/skills/cicy && npm run dev:tui" C-m
```

## AI 助手规则

当用户要求启动、重启或调试时：
1. ✅ 使用 `tmux send-keys` 发送命令
2. ✅ 只操作 `cicy:tui.0` 和 `cicy:api.0`
3. ❌ 不使用 `osascript` 打开新窗口
4. ❌ 不创建新的 tmux session
5. ❌ 不使用 `open -a Terminal`

## 示例

### 错误做法 ❌
```bash
# 错误：打开新窗口
osascript -e 'tell application "Terminal" to do script "npm start"'

# 错误：创建新 session
tmux new-session -d -s cicy-new

# 错误：在后台运行
npm start &
```

### 正确做法 ✅
```bash
# 正确：向现有窗口发送命令
tmux send-keys -t cicy:api.0 C-c
sleep 1
tmux send-keys -t cicy:api.0 "npm run dev:server" C-m
```
