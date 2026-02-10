#!/usr/bin/env node
process.env.LANG = 'en_US.UTF-8';

const blessed = require('blessed');
const axios = require('axios');
const { exec } = require('child_process');
const { promisify } = require('util');

const execAsync = promisify(exec);
const API_URL = 'http://localhost:13001';

// Loading 状态
const spinners = ['⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'];
let spinnerIndex = 0;
let loadingInterval = null;

// 统计信息
let totalMessages = 0;
let totalTime = 0;

// 双击 Ctrl+C 退出
let lastCtrlC = 0;
const DOUBLE_PRESS_INTERVAL = 1000;

// 命令列表
const COMMANDS = ['/quit', '/q', '/clear', '/help', '/test-image', '/list'];

// 输入历史
const MAX_HISTORY_SIZE = 50;
const commandHistory = [];
let historyIndex = -1;
let currentInput = '';

// 创建屏幕
const screen = blessed.screen({
    smartCSR: true,
    title: 'CICY TUI',
    fullUnicode: true
});

// 消息历史区域（可滚动）
const messageBox = blessed.box({
    top: 0,
    left: 0,
    width: '100%',
    height: '100%-5',
    scrollable: true,
    alwaysScroll: true,
    scrollbar: {
        ch: ' ',
        bg: 'blue'
    },
    tags: true,
    content: '{bold}{blue-fg}=== opencode-message ==={/blue-fg}{/bold}\n{gray-fg}Server: ' + API_URL + '{/gray-fg}\n'
});

// 统计窗口（左下角，与 statusBar 同行）
const statsBox = blessed.box({
    bottom: 1,
    left: 0,
    width: 15,
    height: 1,
    tags: true,
    style: {
        fg: 'cyan'
    },
    content: 'Msgs:0'
});

// 状态栏（固定在倒数第二行，右侧）
const statusBar = blessed.box({
    bottom: 1,
    left: 15,
    width: '100%-15',
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
screen.append(statsBox);
screen.append(statusBar);
screen.append(inputBox);

// 显示消息
function addMessage(text) {
    messageBox.pushLine(text);
    messageBox.setScrollPerc(100);
    screen.render();
}

// 显示 Loading
function showLoading(text = 'Thinking...') {
    spinnerIndex = 0;
    loadingInterval = setInterval(() => {
        statusBar.setContent(`{cyan-fg}${spinners[spinnerIndex]} ${text}{/cyan-fg} {gray-fg}| 📊 ${totalMessages} msgs, ${totalTime.toFixed(2)}s{/gray-fg}`);
        screen.render();
        spinnerIndex = (spinnerIndex + 1) % spinners.length;
    }, 80);
}

// 隐藏 Loading
function hideLoading() {
    if (loadingInterval) {
        clearInterval(loadingInterval);
        loadingInterval = null;
        updateStats();
    }
}

// 更新统计窗口
function updateStats() {
    // 更新 statsBox
    statsBox.setContent(`{cyan-fg}Msgs:${totalMessages}{/cyan-fg}`);
    // 在 statusBar 也显示
    if (!loadingInterval) {
        statusBar.setContent(`{gray-fg}Time:${totalTime.toFixed(2)}s{/gray-fg}`);
    }
    screen.render();
}

// 发送消息
async function sendMessage(text) {
    const startTime = Date.now();
    const minLoadingTime = 500; // 最少显示 500ms
    
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
        
        // 显示服务器回复
        if (response.data && response.data.message) {
            addMessage(`\n{gray-fg}> ${response.data.message}{/gray-fg}`);
        }
        
        addMessage(`{gray-fg} - Completed in ${duration}s{/gray-fg}\n`);
        
        // 更新统计
        totalMessages++;
        totalTime += parseFloat(duration);
        updateStats();
        
    } catch (error) {
        hideLoading();
        if (error.code === 'ECONNREFUSED') {
            addMessage(`{red-fg} ✗ Error: Server not running{/red-fg}\n`);
        } else {
            addMessage(`{red-fg} ✗ Error: ${error.message}{/red-fg}\n`);
        }
    }
    
    inputBox.clearValue();
    inputBox.focus();
    screen.render();
}

// 执行 curl-rpc
async function executeCurlRpc(command) {
    const startTime = Date.now();
    
    showLoading('Thinking...');
    
    try {
        const { stdout } = await execAsync(`export ELECTRON_MCP_URL=https://gcp-docs.cicy.de5.net && curl-rpc ${command}`);
        const duration = ((Date.now() - startTime) / 1000).toFixed(2);
        
        hideLoading();
        
        if (stdout) {
            stdout.split('\n').forEach(line => {
                if (line && !line.includes('https://gcp-docs.cicy.de5.net') && !line.includes('---')) {
                    addMessage(line);
                }
            });
            addMessage('');
        }
        
        addMessage(`{gray-fg} - Completed in ${duration}s{/gray-fg}\n`);
        
        // 更新统计
        totalMessages++;
        totalTime += parseFloat(duration);
        updateStats();
        
    } catch (error) {
        hideLoading();
        addMessage(`{red-fg} ✗ Error: ${error.message}{/red-fg}\n`);
    }
    
    inputBox.clearValue();
    inputBox.focus();
    screen.render();
}

// 处理输入
inputBox.on('submit', async (value) => {
    const text = value.trim();
    if (!text) {
        inputBox.clearValue();
        inputBox.focus();
        return;
    }

    if (commandHistory.length === 0 || commandHistory[commandHistory.length - 1] !== text) {
        commandHistory.push(text);
        if (commandHistory.length > MAX_HISTORY_SIZE) {
            commandHistory.shift();
        }
    }
    historyIndex = -1;

    addMessage(`\n> ${text}`);

    if (text === '/quit' || text === '/q') {
        process.exit(0);
    } else if (text === '/clear') {
        messageBox.setContent('{bold}{blue-fg}=== opencode-message ==={/blue-fg}{/bold}\n{gray-fg}Server: ' + API_URL + '{/gray-fg}\n');
        inputBox.clearValue();
        inputBox.focus();
        screen.render();
    } else if (text === '/help') {
        addMessage(`\n{cyan-fg}Commands:{/cyan-fg}`);
        addMessage(`  {yellow-fg}/quit, /q{/yellow-fg}   - 退出程序`);
        addMessage(`  {yellow-fg}/clear{/yellow-fg}       - 清屏`);
        addMessage(`  {yellow-fg}/help{/yellow-fg}       - 显示帮助`);
        addMessage(`  {yellow-fg}/list{/yellow-fg}       - 查看所有消息`);
        addMessage(`  {yellow-fg}/test-image{/yellow-fg} - 发送测试图片`);
        addMessage(`  {yellow-fg}curl-rpc <cmd>{/yellow-fg} - 执行 curl-rpc 命令`);
        addMessage('');
        inputBox.clearValue();
        inputBox.focus();
        screen.render();
    } else if (text.startsWith('curl-rpc ')) {
        const command = text.substring(9);
        await executeCurlRpc(command);
    } else {
        await sendMessage(text);
    }
});

// 快捷键处理
inputBox.key(['C-c'], () => {
    const now = Date.now();
    if (now - lastCtrlC < DOUBLE_PRESS_INTERVAL) {
        addMessage('\n{gray-fg}Bye!{/gray-fg}');
        screen.render();
        setTimeout(() => process.exit(0), 100);
    } else {
        statusBar.setContent('{gray-fg}Press Ctrl+C again to exit{/gray-fg}');
        screen.render();
        lastCtrlC = now;
        setTimeout(() => {
            if (Date.now() - lastCtrlC >= DOUBLE_PRESS_INTERVAL) {
                statusBar.setContent('');
                screen.render();
            }
        }, DOUBLE_PRESS_INTERVAL);
    }
});

inputBox.key(['C-l'], () => {
    messageBox.setContent('{bold}{blue-fg}=== opencode-message ==={/blue-fg}{/bold}\n{gray-fg}Server: ' + API_URL + '{/gray-fg}\n');
    inputBox.clearValue();
    inputBox.focus();
    screen.render();
});

inputBox.key(['C-u'], () => {
    inputBox.clearValue();
    inputBox.focus();
    screen.render();
});

inputBox.key(['C-a'], () => {
    const value = inputBox.getValue();
    inputBox.setValue(value);
    inputBox.focus();
    screen.render();
});

inputBox.key(['C-e'], () => {
    const value = inputBox.getValue();
    inputBox.setValue(value);
    inputBox.focus();
    screen.render();
});

inputBox.key(['up'], () => {
    if (commandHistory.length > 0 && historyIndex < commandHistory.length - 1) {
        if (historyIndex === -1) {
            currentInput = inputBox.getValue();
        }
        historyIndex++;
        inputBox.setValue(commandHistory[commandHistory.length - 1 - historyIndex]);
        inputBox.focus();
        screen.render();
    }
});

inputBox.key(['down'], () => {
    if (historyIndex > 0) {
        historyIndex--;
        inputBox.setValue(commandHistory[commandHistory.length - 1 - historyIndex]);
        inputBox.focus();
        screen.render();
    } else if (historyIndex === 0) {
        historyIndex = -1;
        inputBox.setValue(currentInput);
        inputBox.focus();
        screen.render();
    }
});

inputBox.key(['tab'], () => {
    const value = inputBox.getValue().trim();
    if (value.startsWith('/')) {
        const matches = COMMANDS.filter(cmd => cmd.startsWith(value));
        if (matches.length === 1) {
            inputBox.setValue(matches[0] + ' ');
            inputBox.focus();
            screen.render();
        } else if (matches.length > 1) {
            addMessage(`{yellow-fg}Candidates: ${matches.join(', ')}{/yellow-fg}`);
            inputBox.focus();
            screen.render();
        }
    }
});

// 聚焦输入框
inputBox.focus();
updateStats(); // 初始化统计显示
screen.render();

// 确保输入框始终可用
screen.key(['enter'], () => {
    inputBox.submit();
});
