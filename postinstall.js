#!/usr/bin/env node

const https = require('https');
const fs = require('fs');
const path = require('path');
const os = require('os');
const { execSync } = require('child_process');

const VERSION = '1.0.0';
const GITHUB_REPO = 'yourusername/cicy';

// 平台映射
const platformMap = {
    'darwin-x64': { os: 'darwin', arch: 'amd64' },
    'darwin-arm64': { os: 'darwin', arch: 'arm64' },
    'linux-x64': { os: 'linux', arch: 'amd64' },
    'linux-arm64': { os: 'linux', arch: 'arm64' },
    'win32-x64': { os: 'windows', arch: 'amd64' },
    'win32-arm64': { os: 'windows', arch: 'arm64' }
};

function getPlatformKey() {
    const platform = os.platform();
    const arch = os.arch();
    return `${platform}-${arch}`;
}

function getBinaryName(platformKey) {
    const info = platformMap[platformKey];
    if (!info) return null;
    
    const ext = info.os === 'windows' ? '.exe' : '';
    return `cicy-${info.os}-${info.arch}${ext}`;
}

function downloadBinary(url, dest) {
    return new Promise((resolve, reject) => {
        console.log(`📦 Downloading: ${url}`);
        
        const file = fs.createWriteStream(dest);
        
        https.get(url, (response) => {
            if (response.statusCode === 302 || response.statusCode === 301) {
                // 处理重定向
                return downloadBinary(response.headers.location, dest)
                    .then(resolve)
                    .catch(reject);
            }
            
            if (response.statusCode !== 200) {
                reject(new Error(`Download failed: ${response.statusCode}`));
                return;
            }
            
            response.pipe(file);
            
            file.on('finish', () => {
                file.close();
                console.log('✅ Download complete');
                resolve();
            });
        }).on('error', (err) => {
            fs.unlink(dest, () => {});
            reject(err);
        });
    });
}

async function install() {
    console.log('🚀 Installing CICY CLI...');
    
    const platformKey = getPlatformKey();
    const binaryName = getBinaryName(platformKey);
    
    if (!binaryName) {
        console.error(`❌ Unsupported platform: ${platformKey}`);
        process.exit(1);
    }
    
    // 创建 bin 目录
    const binDir = path.join(__dirname, 'bin');
    if (!fs.existsSync(binDir)) {
        fs.mkdirSync(binDir, { recursive: true });
    }
    
    const binaryPath = path.join(binDir, binaryName);
    
    // 如果已存在，跳过下载
    if (fs.existsSync(binaryPath)) {
        console.log('✅ Binary already exists');
        return;
    }
    
    // 下载 URL
    const downloadUrl = `https://github.com/${GITHUB_REPO}/releases/download/v${VERSION}/${binaryName}`;
    
    try {
        await downloadBinary(downloadUrl, binaryPath);
        
        // 设置可执行权限（Unix 系统）
        if (os.platform() !== 'win32') {
            fs.chmodSync(binaryPath, 0o755);
        }
        
        console.log('✅ CICY CLI installed successfully!');
        console.log('');
        console.log('Run: npx cicy');
        console.log('Or:  cicy');
        
    } catch (err) {
        console.error(`❌ Installation failed: ${err.message}`);
        console.error('');
        console.error('Please try:');
        console.error('  1. Check your internet connection');
        console.error('  2. Download manually from GitHub releases');
        console.error(`  3. Place binary in: ${binDir}`);
        process.exit(1);
    }
}

// 运行安装
install().catch(err => {
    console.error(err);
    process.exit(1);
});
