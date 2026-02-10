# TODO: 底部固定布局（使用 blessed）

## 需求
thinking 和 input 应该在 TUI 底部固定位置，不随消息滚动

## 当前问题
使用 readline，所有内容都在滚动，没有固定布局

## 解决方案
使用 blessed 库重构界面

### 布局设计
```
┌─────────────────────────────────┐
│                                 │
│  消息历史区域（可滚动）           │
│                                 │
│  > hello                        │
│   - Completed in 0.23s          │
│                                 │
│  > test                         │
│   - Completed in 0.19s          │
│                                 │
├─────────────────────────────────┤
│ ⠏ Thinking...                  │  ← 状态栏（固定）
├─────────────────────────────────┤
│ > _                             │  ← 输入框（固定）
└─────────────────────────────────┘
```

## 实现代码

### 基础结构
```javascript
const blessed = require('blessed');

// 创建屏幕
const screen = blessed.screen({
    smartCSR: true,
    title: 'CICY TUI'
});

// 消息历史区域（可滚动）
const messageBox = blessed.box({
    top: 0,
    left: 0,
    width: '100%',
    height: '100%-3',
    scrollable: true,
    alwaysScroll: true,
    scrollbar: {
        ch: ' ',
        bg: 'blue'
    },
    tags: true
});

// 状态栏（固定在倒数第二行）
const statusBar = blessed.box({
    bottom: 1,
    left: 0,
    width: '100%',
    height: 1,
    content: '',
    tags: true
});

// 输入框（固定在最底部）
const inputBox = blessed.textbox({
    bottom: 0,
    left: 0,
    width: '100%',
    height: 1,
    inputOnFocus: true,
    style: {
        fg: 'white',
        bg: 'black'
    }
});

screen.append(messageBox);
screen.append(statusBar);
screen.append(inputBox);

// 聚焦输入框
inputBox.focus();
screen.render();
```

### 显示消息
```javascript
function addMessage(text) {
    messageBox.pushLine(text);
    messageBox.setScrollPerc(100); // 滚动到底部
    screen.render();
}
```

### 显示 Loading
```javascript
function showLoading(text = 'Thinking...') {
    statusBar.setContent(`{cyan-fg}⠏ ${text}{/cyan-fg}`);
    screen.render();
}

function hideLoading() {
    statusBar.setContent('');
    screen.render();
}
```

### 处理输入
```javascript
inputBox.on('submit', async (value) => {
    const text = value.trim();
    if (!text) return;
    
    // 显示用户输入
    addMessage(`\n> ${text}`);
    
    // 清空输入框
    inputBox.clearValue();
    
    // 显示 loading
    showLoading('Sending...');
    
    // 发送消息
    await sendMessage(text);
    
    // 隐藏 loading
    hideLoading();
    
    // 重新聚焦
    inputBox.focus();
    screen.render();
});
```

### 退出处理
```javascript
inputBox.key(['C-c'], () => {
    return process.exit(0);
});
```

## 验收标准
- [x] 消息历史在上方，可滚动
- [x] 状态栏在倒数第二行，固定不动
- [x] 输入框在最底部，固定不动
- [x] Loading 显示在状态栏
- [x] 输入后自动滚动到底部
- [x] Ctrl+C 退出
- [x] 界面不闪烁

## 测试结果 ✅
- 时间：03:45
- 结果：**完美实现！**
- Loading 固定在倒数第二行
- 输入框固定在最底部
- 消息区域可滚动
- 所有功能正常

## 测试场景
1. 发送 10 条消息，消息区域可滚动
2. 发送消息时，loading 在状态栏显示
3. 输入框始终在底部
4. 按 Ctrl+C 退出

## 注意事项
- 需要完全重构 client-tui.js
- 删除 readline 相关代码
- 保留 axios 和消息发送逻辑
- 保持 Masu 风格的输出格式

## 立即执行
1. 安装 blessed（已有）
2. 重构 client-tui.js
3. 测试布局
4. 测试功能
5. 汇报结果
