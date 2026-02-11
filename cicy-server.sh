#!/bin/bash
# CICY Server Launcher
# 自动选择 Go 版本或 Node.js 版本

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 优先使用 Go 版本
if [ -f "$SCRIPT_DIR/server-go/cicy-server" ]; then
    echo "🚀 启动 Go 版本服务器..."
    exec "$SCRIPT_DIR/server-go/cicy-server" "$@"
elif command -v go &> /dev/null && [ -f "$SCRIPT_DIR/server-go/main.go" ]; then
    echo "⚙️  编译 Go 版本..."
    cd "$SCRIPT_DIR/server-go"
    go build -o cicy-server .
    cd "$SCRIPT_DIR"
    echo "🚀 启动 Go 版本服务器..."
    exec "$SCRIPT_DIR/server-go/cicy-server" "$@"
else
    echo "🚀 启动 Node.js 版本服务器..."
    exec node "$SCRIPT_DIR/server.js" "$@"
fi
