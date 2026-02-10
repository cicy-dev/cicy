const express = require('express');
const app = express();
const PORT = 13001;

app.use(express.json());

// Mock API with delay
app.post('/message', async (req, res) => {
    const { message } = req.body;
    
    // 随机延迟 2-3 秒
    const delay = 2000 + Math.random() * 1000;
    await new Promise(resolve => setTimeout(resolve, delay));
    
    console.log(`[${new Date().toLocaleTimeString()}] Received: ${message}`);
    
    res.json({
        success: true,
        message: 'Message received',
        timestamp: new Date().toISOString()
    });
});

app.get('/health', (req, res) => {
    res.json({ status: 'ok' });
});

app.listen(PORT, () => {
    console.log(`Mock API server running on http://localhost:${PORT}`);
    console.log('Delay: 2-3 seconds per request');
});
