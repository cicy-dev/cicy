---
inclusion: always
---

# CICY 项目长期记忆

## 项目概况

### 基本信息
- **项目名称**：CICY - MCP 消息通信系统
- **技术栈**：Node.js + Express + blessed/readline
- **端口**：13001
- **开发环境**：tmux (cicy:tui.0, cicy:api.0)

### 核心功能
1. 本地 MCP 服务器（文本/图片消息）
2. TUI/CLI 客户端
3. 远程浏览器自动化（curl-rpc）

## 工作经验总结

### ✅ 成功经验

#### 1. 验收流程
- **实际测试**：不能只看代码，必须实际运行测试
- **捕获输出**：用 `tmux capture-pane` 查看真实效果
- **对照标准**：逐项检查验收标准，不能遗漏

#### 2. 问题发现
- **多次测试**：连续发送消息，观察累加效果
- **边界情况**：测试错误处理、空输入等
- **细节检查**：Loading 动画、中文显示、统计信息

#### 3. 文档管理
- **分类存放**：TODO 放 `docs/todo/active/`
- **详细标准**：验收标准要具体可测试
- **及时更新**：验收结果立即记录

#### 4. Worker 通信方式 ⚠️ 重要
**正确方式**：
- ✅ 创建 TODO 文件在 `docs/todo/active/`
- ✅ Worker 自己查看并执行任务
- ✅ Master 负责验收结果

**错误方式**：
- ❌ 用 `tmux send-keys` 发送文本命令
- ❌ 用 `curl` 发送 HTTP 请求
- ❌ 尝试"发送指令"给 Worker

**记住**：
> Worker 是独立的 AI 实例
> 
> TODO 文件就是任务分配
> 
> 不需要"发送命令"，只需要创建任务文档

#### 5. 定时器监控系统 🔔 重要

**工作机制**：
```bash
# 脚本位置
/Users/ton/Desktop/skills/cicy/monitor-worker.sh

# 工作流程
每 5 秒检查一次 Worker-1 状态
↓
检查是否在 "Thinking..." (工作中)
↓
如果不在工作：
  1. 捕获 Worker-1 最后输出
  2. 发送消息到 Master (cicy:tui.0)
  3. sleep 1 秒
  4. 发送回车键 (C-m)
  5. 记录到日志
```

**关键点**：
- ✅ 每 5 秒检查一次
- ✅ 只检查是否有 "Thinking..."
- ✅ 通知 Master，不直接操作 Worker
- ✅ 先发消息，sleep 1秒，再发回车
- ✅ 记录所有操作到日志

**启动/停止**：
```bash
# 启动
nohup /Users/ton/Desktop/skills/cicy/monitor-worker.sh > /dev/null 2>&1 &

# 查看进程
ps aux | grep monitor-worker

# 停止
kill <PID>

# 查看日志
tail -f /Users/ton/Desktop/skills/cicy/temp/logs/monitor.log
```

**Master 收到警报后的职责**：
1. 检查 Worker-1 状态
2. 查看他在做什么
3. 如果卡住：发送指令让他继续
4. 如果完成：验收工作
5. 如果等待权限：发送 't' 信任工具

**记住**：
> 定时器是 Master 的助手
> 
> 定时器只通知，不操作
> 
> Master 负责决策和处理

### ❌ 失败教训

#### 1. 验收不严格
**问题**：
- 第一次验收说"通过"，实际有问题
- 没有仔细检查所有验收标准
- 给了 105 分（A 级），但问题很多

**教训**：
- ✅ 必须逐项对照验收标准
- ✅ 发现问题立即指出，不能放过
- ✅ 不合格就是不合格，不能妥协

#### 2. 代码检查不全面
**问题**：
- 只检查了 `sendMessage` 函数
- 没发现 `executeCurlRpc` 还在调用 `showStats()`
- 导致统计信息还在显示

**教训**：
- ✅ 用 `grep` 搜索所有调用
- ✅ 检查所有相关函数
- ✅ 不能只看一个地方

#### 3. 中文显示问题
**问题**：
- blessed 在 tmux 中文显示有问题
- 尝试了 `fullUnicode: true` 和 `LANG=UTF-8`
- 但问题依然存在

**教训**：
- ✅ blessed 中文支持有限，应该早点建议用 CLI
- ✅ 技术选型要考虑实际问题
- ✅ 遇到编码问题，换方案比修复快

#### 4. Worker 管理不严
**问题**：
- Worker 在根目录创建 9 个 TODO 文件
- 验收不通过还给了高分
- 没有及时纠正错误

**教训**：
- ✅ 发现违规立即扣分
- ✅ 不合格就打回重做
- ✅ KPI 要严格执行

## 指挥者工作规范

### 验收流程（标准）

#### 第一步：准备测试环境
```bash
# 1. 重启服务
tmux send-keys -t cicy:api.0 C-c
sleep 1
tmux send-keys -t cicy:api.0 "npm run dev:server" C-m

# 2. 重启客户端
tmux send-keys -t cicy:tui.0 C-c
sleep 1
tmux send-keys -t cicy:tui.0 "npm run dev:tui" C-m
```

#### 第二步：执行测试场景
```bash
# 场景 1：基础功能
tmux send-keys -t cicy:tui.0 "test1" C-m
sleep 2
tmux capture-pane -t cicy:tui.0 -p

# 场景 2：连续发送
tmux send-keys -t cicy:tui.0 "test2" C-m
sleep 2
tmux send-keys -t cicy:tui.0 "test3" C-m
sleep 2
tmux capture-pane -t cicy:tui.0 -p

# 场景 3：错误处理
# 停止服务器，测试错误提示
```

#### 第三步：检查代码
```bash
# 搜索关键函数调用
grep -n "showStats" client-tui.js
grep -n "Loading" client-tui.js

# 查看具体实现
fs_read 相关函数
```

#### 第四步：对照验收标准
```markdown
- [ ] 标准 1：实际测试 → 通过/失败
- [ ] 标准 2：实际测试 → 通过/失败
- [ ] 标准 3：实际测试 → 通过/失败
```

#### 第五步：评分和反馈
```markdown
## 验收结果
✅ 通过：X/Y 项
❌ 失败：原因

## KPI 评分
之前：X 分
本次：+/- Y 分
总分：Z 分

## 下一步
- 通过：归档 TODO，分配新任务
- 失败：更新 TODO，要求重做
```

### KPI 管理规范

#### 扣分必须严格
- 在根目录乱写文件：-10 分（立即）
- 验收不通过：-15 分（每次）
- 不听指令：-10 分（立即）
- 代码质量差：-5 分

#### 加分要有依据
- 一次性通过所有标准：+10 分
- 代码简洁优雅：+8 分
- 主动发现问题：+10 分

#### 评级要准确
- S 级（120+）：完美完成，超出预期
- A 级（100-119）：完全符合要求
- B 级（80-99）：基本完成，有小问题
- C 级（60-79）：问题较多，需要返工
- D 级（<60）：不合格，停止工作

### 文件管理规范

#### 创建文件前检查
```bash
# 1. 确定文件类型
TODO → docs/todo/active/
测试 → docs/tests/
文档 → docs/dev/

# 2. 检查是否存在
ls docs/todo/active/TODO-*.md

# 3. 使用规范命名
TODO-功能名称.md
```

#### 定期清理
```bash
# 移动已完成
mv docs/todo/active/TODO-done.md docs/todo/completed/

# 归档旧任务
mv docs/todo/completed/TODO-old.md docs/todo/archived/
```

## 常见问题处理

### 问题 1：中文显示为 ???
**原因**：blessed 在 tmux 中文支持有限
**解决**：
1. 建议使用 `client-cli.js`（readline）
2. 或者接受问号显示，测试功能逻辑

### 问题 2：Loading 动画不显示
**原因**：请求太快，动画来不及显示
**解决**：
```javascript
const minLoadingTime = 500;
const elapsed = Date.now() - startTime;
if (elapsed < minLoadingTime) {
    await new Promise(resolve => setTimeout(resolve, minLoadingTime - elapsed));
}
```

### 问题 3：统计信息删不掉
**原因**：多个地方调用 `showStats()`
**解决**：
```bash
# 搜索所有调用
grep -n "showStats" client-tui.js

# 全部删除
# 1. 删除函数定义
# 2. 删除所有调用
```

### 问题 4：热重载不生效
**原因**：nodemon 配置或文件未保存
**解决**：
```bash
# 手动重启
tmux send-keys -t cicy:tui.0 C-c
sleep 1
tmux send-keys -t cicy:tui.0 "npm run dev:tui" C-m
```

## 下次改进

### 验收方面
1. ✅ 第一次就严格验收，不放过问题
2. ✅ 用 grep 检查所有相关代码
3. ✅ 测试至少 3 个场景
4. ✅ 截图保存验收证据

### 管理方面
1. ✅ 发现违规立即扣分
2. ✅ KPI 评分要准确
3. ✅ 不合格就打回重做
4. ✅ 警告要有实际后果

### 技术方面
1. ✅ 遇到编码问题早点换方案
2. ✅ 优先使用简单可靠的技术
3. ✅ blessed 中文问题 → 用 readline
4. ✅ 复杂功能 → 分步实现和测试

## 记住

> **验收必须严格，不能妥协**
> 
> **发现问题立即指出，不能放过**
> 
> **KPI 要准确执行，不能手软**
> 
> **代码要全面检查，不能遗漏**

## 当前状态

### Worker-1 评分
- **当前分数**：70 分（C 级 ❌）
- **警告次数**：1 次
- **待修复**：TODO-fix-tui-issues.md

### 下一步
1. 严格监督 Worker-1 修复
2. 验收必须全部通过
3. 不通过继续扣分
4. 2 次 C 级就降级

### 项目进度
- ✅ 基础 MCP 服务器
- ✅ TUI/CLI 客户端
- ⚠️ Masu 风格界面（部分完成）
- ⏳ 中文显示问题
- ⏳ Loading 动画优化
