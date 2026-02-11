#!/usr/bin/env node
process.env.LANG = 'en_US.UTF-8';

const blessed = require('blessed');
const readline = require('readline');
const axios = require('axios');

// 从命令行参数或环境变量获取远程服务器地址
const REMOTE_URL = process.argv[2] || process.env.CICY_REMOTE_URL || 'http://localhost:13001';

console.log(`🌐 Connecting to: ${REMOTE_URL}`);

const screen = blessed.screen({
    smartCSR: true,
    title: 'CICY Remote',
    fullUnicode: true
});

const mainBox = blessed.box({
    top: 0,
    left: 0,
    width: '100%',
    height: '100%',
    tags: true,
    style: {
        bg: '#1a1b26',
        fg: '#a9b1d6'
    }
});

screen.append(mainBox);

let lastQ = '';
let lastA = '';
let isLoading = false;
let connectionStatus = 'checking';

async function checkConnection() {
    try {
        const response = await axios.get(`${REMOTE_URL}/health`, { timeout: 3000 });
        if (response.data.status === 'ok') {
            connectionStatus = 'connected';
            return true;
        }
    } catch (error) {
        connectionStatus = 'disconnected';
        return false;
    }
    return false;
}

async function sendToServer(text) {
    isLoading = true;
    render();

    try {
        const response = await axios.post(`${REMOTE_URL}/message`, { message: text }, { timeout: 10000 });
        lastA = response.data.message || 'Message received!';
        connectionStatus = 'connected';
    } catch (error) {
        if (error.code === 'ECONNREFUSED') {
            lastA = 'Error: Cannot connect to remote server';
            connectionStatus = 'disconnected';
        } else if (error.code === 'ETIMEDOUT') {
            lastA = 'Error: Connection timeout';
            connectionStatus = 'timeout';
        } else {
            lastA = 'Error: ' + error.message;
            connectionStatus = 'error';
        }
    }

    isLoading = false;
    render();
}

function getStatusColor() {
    switch (connectionStatus) {
        case 'connected': return '#9ece6a';
        case 'disconnected': return '#f7768e';
        case 'timeout': return '#e0af68';
        case 'checking': return '#7aa2f7';
        default: return '#565f89';
    }
}

function getStatusText() {
    switch (connectionStatus) {
        case 'connected': return '● Connected';
        case 'disconnected': return '● Disconnected';
        case 'timeout': return '● Timeout';
        case 'checking': return '● Checking...';
        default: return '● Unknown';
    }
}

function render() {
    const qLines = lastQ ? lastQ.split('\n') : [];
    const aLines = lastA ? lastA.split('\n') : [];

    let content = '';

    // 标题和状态
    content += '{bold}{fg #7aa2f7}◇ CICY Remote{/fg}{/bold}\n';
    content += `{fg ${getStatusColor()}}${getStatusText()}{/fg} `;
    content += `{fg #565f89}${REMOTE_URL}{/fg}\n\n`;

    if (qLines.length > 0) {
        qLines.forEach((line) => {
            content += `{fg #f7768e}> ${line}\n`;
        });
        content += '\n';
    }

    if (aLines.length > 0) {
        aLines.forEach((line) => {
            content += `{fg #9ece6a}  ${line}\n`;
        });
        content += '\n';
    }

    if (qLines.length === 0 && aLines.length === 0) {
        content += '{fg #565f89}No messages yet{/fg}\n\n';
    }

    content += '{fg #414868}────────────────────────────────────{/fg}\n\n';

    if (isLoading) {
        content += '{fg #e0af68}Sending to remote server...{/fg}\n';
    } else {
        content += '{fg #c0caf5}Type a message and press Enter{/fg}\n\n';
        content += lastQ ? `{fg #bb9af7}${lastQ}{/fg}\n` : '{fg #565f89}(empty){/fg}\n';
    }

    mainBox.setContent(content);
    screen.render();
}

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout,
    terminal: true
});

rl.on('line', async (input) => {
    const trimmed = input.trim();
    
    // 命令处理
    if (trimmed === '/quit' || trimmed === '/q') {
        console.log('\n{fg #f7768e}Bye!{/fg}\n');
        process.exit(0);
    }
    
    if (trimmed === '/status') {
        await checkConnection();
        render();
        process.stdout.write('> ');
        return;
    }
    
    if (trimmed === '/help') {
        console.log('\nCommands:');
        console.log('  /status - Check connection status');
        console.log('  /quit   - Exit');
        console.log('  /help   - Show this help');
        process.stdout.write('> ');
        return;
    }
    
    if (trimmed) {
        lastQ = trimmed;
        lastA = '';
        await sendToServer(trimmed);
    }
    process.stdout.write('> ');
});

readline.emitKeypressEvents(process.stdin);

process.stdin.on('keypress', (ch, key) => {
    if (key.ctrl && key.name === 'c') {
        console.log('\n{fg #f7768e}Bye!{/fg}\n');
        process.exit(0);
    }
});

// 初始化
(async () => {
    console.log('🔍 Checking connection...');
    const connected = await checkConnection();
    
    if (connected) {
        console.log('✅ Connected successfully!');
    } else {
        console.log('⚠️  Warning: Cannot connect to server');
        console.log('You can still send messages, but they may fail.');
    }
    
    console.log('\nCommands: /status, /help, /quit');
    console.log('');
    
    render();
    process.stdout.write('> ');
})();
