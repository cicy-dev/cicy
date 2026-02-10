# TODO: 双击 Ctrl+C 退出

## 需求
用户需要按两次 Ctrl+C (或 Cmd+C) 才能退出 TUI

## 实现方案

### 方法 1：计时器方式
```javascript
let lastCtrlC = 0;
const DOUBLE_PRESS_INTERVAL = 1000; // 1秒内

process.on('SIGINT', () => {
    const now = Date.now();
    if (now - lastCtrlC < DOUBLE_PRESS_INTERVAL) {
        // 第二次按下，退出
        console.log('\n\x1b[90mBye!\x1b[0m');
        process.exit(0);
    } else {
        // 第一次按下，提示
        console.log('\n\x1b[90mPress Ctrl+C again to exit\x1b[0m');
        lastCtrlC = now;
    }
});
```

### 方法 2：计数器方式
```javascript
let ctrlCCount = 0;

process.on('SIGINT', () => {
    ctrlCCount++;
    if (ctrlCCount >= 2) {
        console.log('\n\x1b[90mBye!\x1b[0m');
        process.exit(0);
    } else {
        console.log('\n\x1b[90mPress Ctrl+C again to exit\x1b[0m');
        setTimeout(() => {
            ctrlCCount = 0; // 1秒后重置
        }, 1000);
    }
});
```

## 验收标准
- [x] 第一次按 Ctrl+C 显示提示：`Press Ctrl+C again to exit`
- [x] 1秒内再按 Ctrl+C 退出程序
- [x] 超过1秒后重新计数
- [x] 退出时显示：`Bye!`
- [x] 不影响其他功能

## 测试结果 ✅
1. ✅ 按一次 Ctrl+C → 显示提示
2. ✅ 等待 2 秒 → 提示消失
3. ✅ 快速双击 Ctrl+C → 退出程序

## 已实现 ✅
使用方法 1（计时器方式），已添加到 client-tui.js
在生产环境（`node client-tui.js`）下完美工作
