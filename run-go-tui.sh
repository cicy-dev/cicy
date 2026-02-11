#!/bin/bash
# 快速启动 Go TUI 客户端

echo "🚀 启动 CICY Go TUI 客户端"

# 检查服务器
if ! curl -s http://localhost:13001/health > /dev/null; then
    echo "❌ 服务器未运行"
    echo "请先启动服务器: npm start"
    exit 1
fi

# 编译（如果需要）
if [ ! -f "tui-go/cicy-tui" ]; then
    echo "📦 首次运行，正在编译..."
    cd tui-go && GOOS=darwin GOARCH=amd64 go build -o cicy-tui && cd ..
fi

# 运行
cd tui-go && ./cicy-tui
