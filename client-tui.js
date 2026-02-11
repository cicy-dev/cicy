#!/usr/bin/env node
process.env.LANG = 'en_US.UTF-8';

const blessed = require('blessed');
const readline = require('readline');
const axios = require('axios');

const API_URL = process.env.API_URL || 'http://localhost:13001';

const screen = blessed.screen({
    smartCSR: true,
    title: 'CICY',
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

async function sendToServer(text) {
    isLoading = true;
    render();

    try {
        const response = await axios.post(`${API_URL}/message`, { message: text }, { timeout: 5000 });
        lastA = response.data.message || 'Message received!';
    } catch (error) {
        lastA = 'Error: ' + (error.code === 'ECONNREFUSED' ? 'Server not running' : error.message);
    }

    isLoading = false;
    render();
}

function render() {
    const qLines = lastQ ? lastQ.split('\n') : [];
    const aLines = lastA ? lastA.split('\n') : [];

    let content = '';

    content += '{bold}{fg #7aa2f7}◇ CICY{/fg}{/bold}\n\n';

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
        content += '{fg #e0af68}Sending...{/fg}\n';
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
    if (input.trim()) {
        lastQ = input;
        lastA = '';
        await sendToServer(input);
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

render();
process.stdout.write('> ');
