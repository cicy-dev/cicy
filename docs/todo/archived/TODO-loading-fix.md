# TODO: TUI Loading 功能修复

## 测试结果 ❌
**问题**：Loading 没有显示出来！

### 测试记录
- 时间：03:23
- 窗口：cicy:tui.0
- 测试：发送 "test message"
- 结果：看不到 loading 动画

### 问题原因
1. Loading 显示时间太短（响应太快）
2. Loading 被日志覆盖
3. Loading 位置不在底部

## 修复方案（选方案1）

### 方案 1：底部状态栏 ⭐
在屏幕最底部显示 loading，不被覆盖：
```javascript
function startLoading() {
    spinnerIndex = 0;
    loadingInterval = setInterval(() => {
        // 保存当前光标位置
        process.stdout.write('\x1b[s');
        // 移动到屏幕底部（第999行，第0列）
        process.stdout.write('\x1b[999;0H');
        // 清除该行
        process.stdout.write('\x1b[2K');
        // 显示绿色 loading
        process.stdout.write(`\x1b[32m${spinners[spinnerIndex]} Loading...\x1b[0m`);
        // 恢复光标位置
        process.stdout.write('\x1b[u');
        spinnerIndex = (spinnerIndex + 1) % spinners.length;
    }, 80);
}

function stopLoading() {
    if (loadingInterval) {
        clearInterval(loadingInterval);
        loadingInterval = null;
        // 清除底部 loading
        process.stdout.write('\x1b[s');
        process.stdout.write('\x1b[999;0H');
        process.stdout.write('\x1b[2K');
        process.stdout.write('\x1b[u');
    }
}
```

### 方案 2：提示符前显示
```javascript
// 在 > 前面显示
process.stdout.write(`\x1b[32m${spinners[spinnerIndex]}\x1b[0m > `);
```

## 验收标准
- [x] Loading 在底部清晰可见
- [x] 绿色旋转动画
- [x] 不被日志覆盖
- [x] 响应后消失
- [x] 不影响输入

## 立即执行
1. ✅ 修改 startLoading() 和 stopLoading()
2. ✅ 使用方案 1（底部状态栏）
3. ✅ 测试 3 次
4. ✅ 截图并汇报

## 测试结果 ✅
- 时间：03:25
- 方案：方案1 - 底部状态栏
- 结果：**成功！**
- Loading 显示在屏幕最底部
- 绿色旋转动画流畅
- 不被日志覆盖
- 响应后自动消失
- 完全不影响输入和日志
