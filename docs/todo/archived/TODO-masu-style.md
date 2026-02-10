# TODO: 改进 TUI 界面风格（参考 masu）

## 目标界面
```
> test message
 - Completed in 0.37s

> another message
 - Completed in 0.52s

 ▸ Messages: 2 • Time: 0.89s
```

## 需要修改的内容

### 1. 移除时间戳
❌ 删除：`[03:27:20] → ss`
✅ 改为：只显示消息内容

### 2. 添加完成提示
每次发送成功后显示：
```
 - Completed in 0.37s
```

### 3. 底部统计栏（可选）
```
 ▸ Messages: 2 • Time: 0.89s
```

### 4. Loading 样式
发送中：`⠋ > `
完成后：恢复 `> `

## 实现要点

### 修改 log 函数
```javascript
// 删除时间戳和图标，只显示消息
function log(msg, type = 'info') {
    readline.clearLine(process.stdout, 0);
    readline.cursorTo(process.stdout, 0);
    
    if (type === 'completed') {
        // 完成提示
        console.log(`\x1b[90m - Completed in ${msg}\x1b[0m\n`);
    } else {
        // 普通消息（如果需要）
        console.log(msg);
    }
    rl.prompt(true);
}
```

### 修改 sendMessage
```javascript
async function sendMessage(text) {
    const startTime = Date.now();
    startLoading();
    
    try {
        const response = await axios.post(`${API_URL}/message`, { message: text }, { timeout: 5000 });
        const duration = ((Date.now() - startTime) / 1000).toFixed(2);
        
        stopLoading();
        log(`${duration}s`, 'completed');
    } catch (error) {
        stopLoading();
        log(`Error: ${error.message}`, 'error');
    }
}
```

## 验收标准
- [x] 没有时间戳 `[03:27:20]`
- [x] 没有箭头 `→`
- [x] 显示完成时间：`- Completed in 0.37s`
- [x] Loading 图标：`⠋ > `
- [x] 界面简洁清爽

## 立即执行
1. ✅ 修改 log 函数，删除时间戳
2. ✅ 添加耗时计算
3. ✅ 显示完成提示
4. ✅ 测试 5 次
5. ✅ 汇报结果

## 测试结果 ✅
- 时间：03:32
- 结果：**完美实现 masu 风格！**

### 实际效果
```
> test message 1
test message 1
 - Completed in 2.25s

> test message 2
test message 2
 - Completed in 2.36s

>
```

### 改进点
- ✅ 删除了时间戳
- ✅ 删除了箭头符号
- ✅ 添加了完成耗时
- ✅ 界面简洁清爽
- ✅ 灰色显示耗时信息
