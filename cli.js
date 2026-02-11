#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');
const os = require('os');
const fs = require('fs');

// 检测平台和架构
const platform = os.platform();
const arch = os.arch();

// 二进制文件映射
const binaries = {
    'darwin-x64': 'cicy-darwin-amd64',
    'darwin-arm64': 'cicy-darwin-arm64',
    'linux-x64': 'cicy-linux-amd64',
    'linux-arm64': 'cicy-linux-arm64',
    'win32-x64': 'cicy-windows-amd64.exe',
    'win32-arm64': 'cicy-windows-arm64.exe'
};

// 获取二进制文件路径
function getBinaryPath() {
    const key = `${platform}-${arch}`;
    const binaryName = binaries[key];
    
    if (!binaryName) {
        console.error(`❌ Unsupported platform: ${platform}-${arch}`);
        console.error('Supported platforms:');
        Object.keys(binaries).forEach(k => console.error(`  - ${k}`));
        process.exit(1);
    }
    
    return path.join(__dirname, 'bin', binaryName);
}

// 检查二进制文件是否存在
function checkBinary(binaryPath) {
    if (!fs.existsSync(binaryPath)) {
        console.error(`❌ Binary not found: ${binaryPath}`);
        console.error('Please run: npm install cicy-cli');
        process.exit(1);
    }
    
    // 确保可执行权限（Unix 系统）
    if (platform !== 'win32') {
        try {
            fs.chmodSync(binaryPath, 0o755);
        } catch (err) {
            console.error(`⚠️  Warning: Could not set executable permission: ${err.message}`);
        }
    }
}

// 主函数
function main() {
    const binaryPath = getBinaryPath();
    checkBinary(binaryPath);
    
    // 启动 Go 客户端
    const child = spawn(binaryPath, process.argv.slice(2), {
        stdio: 'inherit',
        env: process.env
    });
    
    child.on('error', (err) => {
        console.error(`❌ Failed to start CICY: ${err.message}`);
        process.exit(1);
    });
    
    child.on('exit', (code) => {
        process.exit(code || 0);
    });
}

// 运行
main();
