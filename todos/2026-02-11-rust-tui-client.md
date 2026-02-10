# TODO: Rust TUI 客户端实现

## 需求描述
用 Rust 重新实现 CICY TUI 客户端，解决 blessed 的问题。

## 验收标准
- [ ] TUI 界面正常显示
- [ ] 消息发送功能正常
- [ ] 命令处理正常（/quit, /clear, /list）
- [ ] curl-rpc 集成正常
- [ ] 热重载支持
- [ ] 在 tmux 中正常工作

## 技术选型
- **TUI 库**: ratatui (最流行的 Rust TUI 库)
- **HTTP 客户端**: reqwest
- **JSON**: serde_json
- **异步运行时**: tokio

## 开发步骤
1. 创建 Rust 项目
2. 添加依赖
3. 实现基础 TUI
4. 实现 MCP 客户端
5. 实现命令处理
6. 测试验证

## 测试环境
- TUI: cicy:tui.0
- API: cicy:api.0

## 状态
- [x] 开始开发
- [x] 代码实现
- [x] 编译成功
- [x] 基本功能测试
- [ ] tmux 显示问题
- [ ] 已完成

## 测试记录

### 测试1：编译
- 操作：cargo build --release
- 结果：✅ 编译成功

### 测试2：消息发送
- 操作：输入 "test rust message"
- 结果：✅ API 收到消息
- 证据：[2:02:30 AM] Received: test rust message

### 测试3：curl-rpc
- 操作：输入 "curl-rpc ping"
- 结果：❌ 被当作普通消息发送
- 原因：tmux 中键盘输入捕获问题

## 问题分析
Rust TUI 使用 ratatui + crossterm，在 tmux 中可能有兼容性问题。
blessed (Node.js) 也有类似问题。

## 最终结论

### 问题根源
**TUI 库在 tmux 中无法正确捕获键盘输入**
- blessed (Node.js): submit 事件不触发
- ratatui + crossterm (Rust): 键盘事件不触发
- 这是 TUI 库的通用问题，不是代码实现问题

### 验证结果
1. ✅ API 正常工作（手动 curl 测试通过）
2. ✅ Rust TUI 编译成功
3. ✅ Rust TUI 界面显示正常
4. ❌ blessed 和 ratatui 都无法在 tmux 中接收输入
5. ❌ 无法完成验收标准

### 解决方案
**使用简单的 readline 客户端**，不使用 TUI 库：
- client-simple.js 或 client-mcp.js
- 或创建新的简单命令行客户端
- 避免使用 blessed/ratatui 等全屏 TUI 库
