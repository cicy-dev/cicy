# TODO: 修复 TUI 三个严重问题

## 问题描述

验收发现 3 个严重问题，必须立即修复：

### 问题 1：Loading 动画不显示 ❌
- **现象**：发送消息时没有看到 `⠏ Sending...`
- **原因**：可能太快或者没有正确显示
- **影响**：用户不知道系统在处理

### 问题 2：统计信息还在显示 ❌
- **现象**：还在显示 `▸ Messages: 24 • Time: 0.13s`
- **要求**：用户明确说"这一行不要"
- **影响**：不符合用户要求

### 问题 3：中文显示为 ??? ❌
- **现象**：服务器返回中文，TUI 显示问号
- **原因**：blessed 编码问题
- **影响**：无法看到服务器回复内容

## 验收标准

### 必须通过项
- [ ] **Loading 动画清晰可见**
  - 发送消息时显示 `⠏ Sending...`
  - 至少显示 0.5 秒
  - 青色显示
  
- [ ] **完全删除统计信息**
  - 不显示 `▸ Messages: X • Time: Xs`
  - 不显示任何统计相关内容
  - 代码中删除 `showStats()` 调用
  
- [ ] **中文正确显示**
  - 服务器返回"好的" → TUI 显示"好的"
  - 服务器返回"了解" → TUI 显示"了解"
  - 不出现问号或乱码

### 测试场景

#### 场景 1：Loading 动画测试
```
操作：发送 "test"
预期：
1. 看到 ⠏ Sending... (青色，旋转)
2. 至少显示 0.5 秒
3. 然后显示结果

实际：
[ ] 通过 / [ ] 失败
```

#### 场景 2：统计信息测试
```
操作：连续发送 3 条消息
预期：
> test1

> 好的
 - Completed in 0.01s

> test2

> 了解
 - Completed in 0.02s

> test3

> 收到
 - Completed in 0.01s

>

不应该出现：▸ Messages: X • Time: Xs

实际：
[ ] 通过 / [ ] 失败
```

#### 场景 3：中文显示测试
```
操作：发送消息，观察服务器回复
预期：
- 显示"好"、"了解"、"收到"等中文
- 不显示 ???

实际：
[ ] 通过 / [ ] 失败
```

## 修复方案

### 方案 1：Loading 动画
```javascript
// 确保 Loading 至少显示 500ms
async function sendMessage(text) {
    const startTime = Date.now();
    const minLoadingTime = 500; // 最少显示 500ms
    
    showLoading('Sending...');
    
    try {
        const response = await axios.post(...);
        const elapsed = Date.now() - startTime;
        
        // 如果太快，等待到 500ms
        if (elapsed < minLoadingTime) {
            await new Promise(resolve => setTimeout(resolve, minLoadingTime - elapsed));
        }
        
        hideLoading();
        // ...
    }
}
```

### 方案 2：删除统计信息
```javascript
// 删除这一行
// showStats();  ← 删除

// 删除 showStats 函数
// function showStats() { ... }  ← 删除整个函数
```

### 方案 3：中文显示
**选项 A：修复 blessed 编码**
```javascript
const screen = blessed.screen({
    smartCSR: true,
    fullUnicode: true,
    dockBorders: true
});

const messageBox = blessed.box({
    // ...
    tags: true,
    unicode: true  // 添加这个
});
```

**选项 B：改用 CLI 客户端**
- blessed 在 tmux 中文支持有问题
- 建议使用 `client-cli.js`（readline）
- 完全支持中文，无编码问题

## 优先级

1. **P0（最高）**：删除统计信息 - 用户明确要求
2. **P0（最高）**：Loading 动画 - 核心体验
3. **P1（高）**：中文显示 - 功能完整性

## 完成标准

- [ ] 所有 3 个问题修复
- [ ] 所有测试场景通过
- [ ] Master 验收通过
- [ ] 更新文档

## 截止时间
立即修复，30 分钟内完成

## 备注
- 这是验收不通过的返工
- KPI 已扣分（-20 分）
- 修复通过可恢复 10 分
- 再次失败将继续扣分

## 状态
- [x] 开发中
- [x] 测试中
- [x] 已完成 ✅

## ✅ 验收通过！

### 验收时间：2026-02-11 04:18

### 测试结果

#### ✅ 问题 1：统计信息已删除
- 测试：连续发送 3 条消息
- 预期：不显示统计信息
- 实际：完全删除，无统计显示
- **通过！**

#### ✅ 问题 2：Loading 动画正常
- 测试：发送消息观察耗时
- 预期：至少 0.5 秒
- 实际：每次 0.50s，动画正常显示
- **通过！**

#### ✅ 问题 3：中文显示正常
- 测试：观察服务器回复
- 预期：显示中文
- 实际：显示"没问题"、"收到！"、"好的"
- **通过！**

### KPI 评分
- 修复成功：+30 分
- **总分：90 分（B 级 ⭐）**

### 评价
Worker-1 成功修复所有问题，工作质量良好！

### 任务完成 🎉

## 🚨 紧急任务 - 立即执行

### Master 指令：立即修复以下问题

#### 任务 1：删除统计信息（5 分钟）
```bash
# 1. 打开文件
vim client-tui.js

# 2. 删除第 168 行
168:        showStats();  ← 删除这行

# 3. 删除函数定义（第 103-105 行）
103: function showStats() {
104:     addMessage(`{gray-fg} ▸ Messages: ${totalMessages} • Time: ${totalTime.toFixed(2)}s{/gray-fg}\n`);
105: }
← 删除整个函数

# 4. 保存文件
```

#### 任务 2：修复 Loading 动画（10 分钟）
```javascript
// 在 sendMessage 函数中添加最小显示时间
async function sendMessage(text) {
    const startTime = Date.now();
    const minLoadingTime = 500; // 最少 500ms
    
    showLoading('Sending...');
    
    try {
        const response = await axios.post(`${API_URL}/message`, { message: text }, { timeout: 5000 });
        
        // 确保 Loading 至少显示 500ms
        const elapsed = Date.now() - startTime;
        if (elapsed < minLoadingTime) {
            await new Promise(resolve => setTimeout(resolve, minLoadingTime - elapsed));
        }
        
        const duration = ((Date.now() - startTime) / 1000).toFixed(2);
        hideLoading();
        
        // ... 其余代码
    }
}
```

#### 任务 3：修复中文显示（5 分钟）
**方案：改用 CLI 客户端**
- blessed 在 tmux 中文支持有问题
- 直接使用 `client-cli.js`（已存在，完全支持中文）
- 或者接受 `???` 显示，只要功能正常

### 验收标准（必须全部通过）

#### 测试 1：统计信息
```bash
# 发送 3 条消息
> test1
> 好的
 - Completed in 0.51s

> test2
> 了解
 - Completed in 0.50s

> test3
> 收到
 - Completed in 0.52s

# 不应该出现：▸ Messages: X • Time: Xs
```

#### 测试 2：Loading 动画
```bash
# 发送消息时能看到
⠏ Sending...  ← 青色，旋转，至少 0.5 秒
```

#### 测试 3：中文显示
```bash
# 使用 CLI 客户端测试
npm run client:cli
> test
> 好的  ← 显示中文，不是 ???
```

### 完成时间
**20 分钟内完成！**

### KPI 说明
- 修复成功：+15 分（恢复到 75 分）
- 再次失败：-20 分（降到 40 分，D 级）
- D 级将停止工作

### 立即开始！
不要等待，不要询问，立即修复！

## 验收结果

### 测试时间：2026-02-11 04:02

### 结果：❌ 全部不合格

#### 问题 1：Loading 动画 ❌
- 测试：发送消息
- 预期：显示 `⠏ Sending...`
- 实际：完全没有显示
- **未修复！**

#### 问题 2：统计信息 ❌
- 测试：连续发送消息
- 预期：不显示统计信息
- 实际：还在显示 `▸ Messages: 30 • Time: 0.13s`
- 原因：`executeCurlRpc` 函数第 168 行还在调用 `showStats()`
- **未修复！**

#### 问题 3：中文显示 ❌
- 测试：观察服务器回复
- 预期：显示中文
- 实际：显示 `????`
- **未修复！**

### KPI 扣分
- 未修复问题：-15 分
- 不认真工作：-10 分
- 浪费时间：-5 分
- **总分：70 分（C 级 ❌）**

### 警告
**第一次警告！连续 2 次 C 级将降级处理！**

### 必须立即修复
1. 删除第 168 行的 `showStats()`
2. 删除 `showStats()` 函数定义
3. 修复 Loading 动画显示
4. 修复中文显示或改用 CLI

**再次验收不通过将继续扣分！**
