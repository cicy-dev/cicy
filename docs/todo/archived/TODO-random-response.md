# TODO: API 随机响应

## 需求
server.js 的 API 不要返回固定格式 "Message received: xxx"，改为随机返回不同的句子

## 当前问题
```javascript
// 当前代码
res.json({ 
    success: true, 
    message: 'Message received: ' + message 
});
```

每次都返回相同格式，很无聊。

## 解决方案

### 随机响应列表
```javascript
const responses = [
    '收到！',
    '明白了',
    '好的',
    '了解',
    '知道了',
    '没问题',
    '👍',
    '✓',
    'OK',
    '已记录',
    '已收到',
    '好嘞',
    '得令',
    '遵命',
    '收到消息',
    '已处理',
    '完成',
    '搞定',
    '安排上了',
    '妥了'
];

function getRandomResponse() {
    return responses[Math.floor(Math.random() * responses.length)];
}
```

### 修改 API 端点
```javascript
app.post('/message', (req, res) => {
    const { message } = req.body;
    
    if (!message) {
        return res.status(400).json({ 
            success: false, 
            error: 'Message is required' 
        });
    }
    
    // 存储消息
    messages.push({
        id: Date.now(),
        text: message,
        timestamp: new Date().toISOString(),
        from: req.ip
    });
    
    // 随机响应
    res.json({ 
        success: true, 
        reply: getRandomResponse()
    });
});
```

### 客户端显示响应（可选）
如果需要在 TUI 显示 API 的回复：
```javascript
async function sendMessage(text) {
    const startTime = Date.now();
    showLoading('Sending...');
    
    try {
        const response = await axios.post(`${API_URL}/message`, { message: text });
        const duration = ((Date.now() - startTime) / 1000).toFixed(2);
        
        hideLoading();
        
        // 显示完成
        addMessage(`{gray-fg} - Completed in ${duration}s{/gray-fg}`);
        
        // 显示 API 回复（可选）
        if (response.data.reply) {
            addMessage(`\n> ${response.data.reply}`);
        }
        
    } catch (error) {
        hideLoading();
        addMessage(`{red-fg} ✗ Error: ${error.message}{/red-fg}`);
    }
}
```

## 验收标准
- [x] 每次返回不同的随机句子
- [x] 不再返回 "Message received: xxx"
- [x] 响应列表至少 10 个（已有20个）
- [x] 包含中文和 emoji
- [x] 客户端正常接收
- [x] 不影响其他功能

## 测试结果 ✅
- test 1 → `已收到`
- test 2 → `OK`
- test 3 → `搞定`
- test 4 → `妥了`
- test 5 → `👍`

所有回复都是随机的，包含中文和 emoji！

## 测试场景
1. 发送 5 条消息，观察响应是否随机
2. 检查响应格式是否正确
3. 确认不再有固定格式

## 立即执行
1. 修改 server.js
2. 添加随机响应列表
3. 修改 /message 端点
4. 测试 5 次
5. 汇报结果
