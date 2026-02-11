#!/bin/bash
# Go TUI 客户端测试脚本

echo "🧪 CICY Go TUI 测试"
echo "==================="

# 检查服务器
echo ""
echo "检查服务器状态..."
if ! curl -s http://localhost:13001/health > /dev/null; then
    echo "❌ 服务器未运行"
    echo "请先启动服务器: npm start"
    exit 1
fi
echo "✅ 服务器正常"

# 编译
echo ""
echo "编译 Go 客户端..."
if GOOS=darwin GOARCH=amd64 go build -o cicy-tui; then
    echo "✅ 编译成功"
else
    echo "❌ 编译失败"
    exit 1
fi

# 检查可执行文件
echo ""
echo "检查可执行文件..."
if [ -f "./cicy-tui" ]; then
    SIZE=$(ls -lh cicy-tui | awk '{print $5}')
    echo "✅ 可执行文件: cicy-tui ($SIZE)"
else
    echo "❌ 可执行文件不存在"
    exit 1
fi

echo ""
echo "==================="
echo "🎉 测试通过！"
echo ""
echo "运行客户端:"
echo "  ./cicy-tui"
echo ""
echo "或在 tmux 中:"
echo "  tmux send-keys -t cicy:tui.0 'cd /Users/ton/Desktop/skills/cicy/tui-go && ./cicy-tui' C-m"
