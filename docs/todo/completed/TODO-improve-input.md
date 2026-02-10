# TODO: 优化 TUI 输入体验

## 问题
当前 TUI 在 tmux 环境下输入框无法正常提交消息。

## 原因
blessed textbox 在 tmux 中的兼容性问题。

## 解决方案

### 方案 1：改进 blessed 配置
```javascript
const inputBox = blessed.textbox({
    bottom: 0,
    left: 0,
    width: '100%',
    height: 1,
    inputOnFocus: true,
    keys: true,  // 添加
    mouse: true, // 添加
    style: {
        fg: 'white',
        bg: 'black'
    }
});
```

### 方案 2：使用 readline 替代
创建混合模式：blessed 显示 + readline 输入

### 方案 3：文档说明
在 README 中说明：
- TUI 最佳使用方式：直接运行，不通过 tmux
- 或使用 CLI 客户端（client-cli.js）

## 验收标准
- [ ] 输入框可以正常提交消息
- [ ] 或提供替代方案
- [ ] 更新文档说明

## 优先级
P2 - 中等（不影响核心功能）

## 状态
- [x] 待处理
- [x] 已完成

## 实施方案
选择方案 3：文档说明

### 已完成
1. ✅ 更新 README.md
   - 添加 TUI 使用说明
   - 更新客户端对比表
   - 说明 tmux 兼容性问题

2. ✅ 提供替代方案
   - 推荐使用 CLI 客户端（client-cli.js）
   - 或在独立终端运行 TUI

### 验收标准
- [x] 更新文档说明
- [x] 提供替代方案
- [x] 用户知道如何正确使用

## 完成时间
2026-02-11 05:09
