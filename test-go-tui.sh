#!/bin/bash
# CICY Go TUI 完整测试脚本

echo "🚀 CICY Go TUI 测试"
echo "===================="
echo ""

# 1. 停止旧进程
echo "1️⃣ 停止旧进程..."
pkill -9 -f cicy-go
sleep 1

# 2. 启动 cicy-go
echo "2️⃣ 启动 cicy-go..."
cd /Users/ton/Desktop/skills/cicy/server-go
./cicy-go &
CICY_PID=$!
echo "   进程 PID: $CICY_PID"
sleep 3

# 3. 检查服务器
echo ""
echo "3️⃣ 检查服务器状态..."
curl -s http://localhost:13001/health | jq .
echo ""

# 4. 获取 token
echo "4️⃣ 读取 token..."
TOKEN=$(cat ~/data/cicy-server.txt)
echo "   Token: ${TOKEN:0:20}..."
echo ""

# 5. 发送文本消息
echo "5️⃣ 发送文本消息..."
curl -s -X POST http://localhost:13001/api/message \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content": [{"type": "text", "text": "测试文本消息"}]}' | jq .
sleep 1
echo ""

# 6. 发送图片消息
echo "6️⃣ 发送图片消息..."
IMAGE_BASE64=$(base64 -i ~/Desktop/avatr.png | tr -d '\n')
curl -s -X POST http://localhost:13001/api/message \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"content\": [{\"type\": \"image\", \"data\": \"$IMAGE_BASE64\"}]}" | jq .
sleep 1
echo ""

# 7. 检查图片保存
echo "7️⃣ 检查图片保存位置..."
ls -lh ~/Desktop/images/ | tail -5
echo ""

echo "===================="
echo "✅ 测试完成！"
echo ""
echo "📝 现在你可以："
echo "   1. 查看 TUI 界面，应该看到文本和图片消息"
echo "   2. 按 'o' 打开图片"
echo "   3. 检查 ~/Desktop/images/ 目录"
echo ""
echo "🛑 停止服务器: kill $CICY_PID"
