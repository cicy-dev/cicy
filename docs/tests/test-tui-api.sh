#!/bin/bash
# 完整的 TUI + API 测试脚本

echo "🚀 CICY TUI + API 完整测试"
echo "=========================="
echo ""

# 1. 停止旧进程
echo "1️⃣ 清理旧进程..."
pkill -f cicy-go
sleep 1

# 2. 创建 tmux session
echo "2️⃣ 创建 tmux session: cicy-test..."
tmux kill-session -t cicy-test 2>/dev/null
tmux new-session -d -s cicy-test -n main

# 3. 启动 cicy-go
echo "3️⃣ 启动 cicy-go 服务器..."
tmux send-keys -t cicy-test:main "cd /Users/ton/Desktop/skills/cicy/server-go && ./cicy-go" C-m
sleep 3

# 4. 显示初始状态
echo "4️⃣ 显示 TUI 初始状态..."
tmux capture-pane -t cicy-test:main -p | tail -20
echo ""
echo "按回车继续..."
read

# 5. 发送测试消息
echo "5️⃣ 发送测试消息..."
TOKEN=$(cat ~/data/cicy-server.txt)

echo "   📝 发送文本消息 1..."
curl -s -X POST http://localhost:13001/api/message \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content": [{"type": "text", "text": "第一条测试消息"}]}' > /dev/null
sleep 1

echo "   📝 发送文本消息 2..."
curl -s -X POST http://localhost:13001/api/message \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content": [{"type": "text", "text": "第二条测试消息"}]}' > /dev/null
sleep 1

echo "   📝 发送文本消息 3..."
curl -s -X POST http://localhost:13001/api/message \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content": [{"type": "text", "text": "第三条测试消息"}]}' > /dev/null
sleep 1

# 6. 显示更新后的状态
echo ""
echo "6️⃣ 显示 TUI 更新后的状态..."
tmux capture-pane -t cicy-test:main -p | tail -25
echo ""
echo "按回车继续..."
read

# 7. 发送图片消息
echo "7️⃣ 发送图片消息..."
IMAGE_BASE64=$(base64 -i /Users/ton/Desktop/avatr.png)
curl -s -X POST http://localhost:13001/api/message \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"content\": [{\"type\": \"text\", \"text\": \"这是一张图片：\"}, {\"type\": \"image\", \"data\": \"$IMAGE_BASE64\"}]}" > /dev/null
sleep 2

# 8. 显示最终状态
echo ""
echo "8️⃣ 显示 TUI 最终状态..."
tmux capture-pane -t cicy-test:main -p | tail -30
echo ""

# 9. 询问是否保持 session
echo "=========================="
echo "✅ 测试完成！"
echo ""
echo "tmux session 'cicy-test' 仍在运行"
echo "你可以："
echo "  - 运行 'tmux attach -t cicy-test' 查看 TUI"
echo "  - 运行 'tmux kill-session -t cicy-test' 关闭"
echo ""
echo "是否现在关闭 session? (y/N)"
read -r response
if [[ "$response" =~ ^[Yy]$ ]]; then
    tmux kill-session -t cicy-test
    echo "✅ Session 已关闭"
else
    echo "✅ Session 保持运行"
    echo "   运行 'tmux attach -t cicy-test' 查看"
fi
