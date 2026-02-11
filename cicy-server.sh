#!/bin/bash
# CICY Server Launcher
# 自动选择 Go 版本或 Node.js 版本，支持自动下载

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_OWNER="cicy-dev"
REPO_NAME="cicy"
BINARY_NAME="cicy-go"
GITHUB_REPO="$REPO_OWNER/$REPO_NAME"

# 检测平台
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)

    case "$os" in
        linux*) os="linux" ;;
        darwin*) os="darwin" ;;
        mingw*|cygwin*|msys*) os="windows" ;;
        *) echo "Unsupported OS: $os" && exit 1 ;;
    esac

    case "$arch" in
        x86_64|amd64) arch="amd64" ;;
        arm64|aarch64) arch="arm64" ;;
        armv7l) arch="armv7l" ;;
        *) echo "Unsupported arch: $arch" && exit 1 ;;
    esac

    echo "${os}-${arch}"
}

# 下载 Go 二进制
download_binary() {
    local platform=$1
    local download_dir="$SCRIPT_DIR/.bin"
    local binary_path="$download_dir/$BINARY_NAME"
    local version_file="$download_dir/.version"
    local current_version=""
    local latest_version=""

    # 创建下载目录
    mkdir -p "$download_dir"

    # 检查是否需要更新
    if [ -f "$version_file" ]; then
        current_version=$(cat "$version_file")
    fi

    # 获取最新版本
    latest_version=$(curl -s "https://api.github.com/repos/$GITHUB_REPO/releases/latest" | grep '"tag_name"' | cut -d'"' -f4 | sed 's/v//')

    if [ -z "$latest_version" ]; then
        echo "❌ 无法获取最新版本"
        return 1
    fi

    # 如果版本不匹配，下载新版本
    if [ "$current_version" != "$latest_version" ] || [ ! -f "$binary_path" ]; then
        echo "📦 下载 $BINARY_NAME v$latest_version ($platform)..."

        local download_url="https://github.com/$GITHUB_REPO/releases/download/v$latest_version/$BINARY_NAME-$platform.tar.gz"

        # 下载并解压
        curl -fsSL "$download_url" -o /tmp/cicy-go.tar.gz || {
            echo "❌ 下载失败，尝试下载源码..."
            download_source
            return $?
        }

        tar -xzf /tmp/cicy-go.tar.gz -C "$download_dir"
        rm -f /tmp/cicy-go.tar.gz
        chmod +x "$binary_path"

        # 保存版本
        echo "$latest_version" > "$version_file"

        echo "✅ 下载完成: $binary_path"
    else
        echo "✅ 使用缓存版本: $binary_path"
    fi

    echo "🚀 启动 Go 版本服务器..."
    exec "$binary_path" "$@"
}

# 备用：下载源码编译
download_source() {
    echo "📦 下载源码..."
    local temp_dir=$(mktemp -d)
    cd "$temp_dir"

        curl -fsSL "https://github.com/$GITHUB_REPO/archive/refs/heads/main.tar.gz" -o source.tar.gz
    tar -xzf source.tar.gz

    cd "$REPO_NAME-$latest_version/server-go"
    echo "⚙️  编译中..."
    go build -o "$download_dir/$BINARY_NAME" .

    chmod +x "$download_dir/$BINARY_NAME"
    echo "$latest_version" > "$download_dir/.version"

    cd "$SCRIPT_DIR"
    rm -rf "$temp_dir"

    echo "✅ 编译完成"
    echo "🚀 启动 Go 版本服务器..."
    exec "$download_dir/$BINARY_NAME" "$@"
}

# 优先使用本地预编译版本
PLATFORM=$(detect_platform)
if [ -f "$SCRIPT_DIR/server-go/cicy-go-$PLATFORM" ]; then
    echo "🚀 启动本地 Go 版本 ($PLATFORM)..."
    exec "$SCRIPT_DIR/server-go/cicy-go-$PLATFORM" "$@"
elif [ -f "$SCRIPT_DIR/server-go/cicy-go" ]; then
    echo "🚀 启动本地 Go 版本 (通用)..."
    exec "$SCRIPT_DIR/server-go/cicy-go" "$@"
fi

# 尝试使用本地 Go 编译
if command -v go &> /dev/null && [ -f "$SCRIPT_DIR/server-go/main.go" ]; then
    echo "⚙️  使用本地 Go 编译..."
    cd "$SCRIPT_DIR/server-go"
    go build -o "$SCRIPT_DIR/.bin/cicy-go" .
    cd "$SCRIPT_DIR"
    echo "🚀 启动 Go 版本服务器..."
    exec "$SCRIPT_DIR/.bin/cicy-go" "$@"
fi

# 自动下载
echo "📥 自动下载 Go 二进制..."
PLATFORM=$(detect_platform)
download_binary "$PLATFORM"
