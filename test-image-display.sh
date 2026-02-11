#!/bin/bash
# 测试图片显示 - 直接在终端运行（不在 tmux 中）

TOKEN=$(cat ~/data/cicy-server.txt)
IMAGE_BASE64=$(base64 -i /Users/ton/Desktop/avatr.png)

echo "📸 发送图片到 CICY 服务器..."
echo ""

curl -X POST http://localhost:13001/api/message \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"content\": [
      {
        \"type\": \"text\",
        \"text\": \"这是一张测试图片：\"
      },
      {
        \"type\": \"image\",
        \"data\": \"$IMAGE_BASE64\"
      }
    ]
  }"

echo ""
echo ""
echo "✅ 图片已发送"
echo "💡 提示：如果你在 iTerm2 中运行，应该能看到图片显示在服务器日志中"
echo "💡 如果在 tmux 中，需要退出 tmux 后直接运行 cicy-go 才能看到图片"
