# TODO: 修复 Loading 显示问题

## 测试结果 ❌

### 发现的问题
1. ❌ Loading 图标在消息中间，不在底部
2. ❌ 显示了 "Loading..." 文字，应该只显示图标
3. ❌ ANSI 转义码不生效，光标控制失败

### 测试截图
```
⠙ another test          ← 错误：图标在消息里
another test
[03:27:21] → another test
⠋ test loading          ← 错误：图标在消息里
test loading
[03:28:25] → test loading
⠙ Loading...            ← 错误：有文字
```

## 正确的实现方案

### 使用 readline 的方式（推荐）
```javascript
function startLoading() {
    spinnerIndex = 0;
    loadingInterval = setInterval(() => {
        // 清除当前行
        readline.clearLine(process.stdout, 0);
        readline.cursorTo(process.stdout, 0);
        // 只显示绿色图标和提示符
        process.stdout.write(`\x1b[32m${spinners[spinnerIndex]}\x1b[0m > `);
        spinnerIndex = (spinnerIndex + 1) % spinners.length;
    }, 80);
}

function stopLoading() {
    if (loadingInterval) {
        clearInterval(loadingInterval);
        loadingInterval = null;
        // 清除 loading，恢复正常提示符
        readline.clearLine(process.stdout, 0);
        readline.cursorTo(process.stdout, 0);
        rl.prompt(true);
    }
}
```

### 关键点
1. 在提示符 `>` 前面显示图标
2. 只显示图标，不要文字
3. 使用 readline 的 clearLine 和 cursorTo
4. 停止时恢复正常提示符

## 验收标准
- [x] 图标在提示符前：`⠋ > `
- [x] 只有图标，没有 "Loading..." 文字
- [x] 绿色显示
- [x] 动画流畅
- [x] 不影响日志输出
- [x] 响应后恢复正常提示符 `> `

## 立即执行
1. ✅ 修改 startLoading() 和 stopLoading()
2. ✅ 删除 "Loading..." 文字
3. ✅ 测试 3 次，确认位置正确
4. ✅ 汇报结果

## 测试结果 ✅
- 时间：03:30
- 方案：readline 方式
- 结果：**完美！**

### 实际效果
```
> final test
final test
[03:30:08] → final test
⠙ >                    ← 正确：图标在提示符前
```

### 验证通过
- ✅ 图标位置正确：`⠙ >`
- ✅ 只显示图标，无文字
- ✅ 绿色显示
- ✅ 旋转流畅
- ✅ 完成后恢复 `>`
- ✅ 不影响日志
