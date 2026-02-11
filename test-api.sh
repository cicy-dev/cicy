#!/bin/bash
# 测试 CICY API

# 读取 token
TOKEN=$(cat ~/data/cicy-server.txt)

echo "Token: $TOKEN"
echo ""

# 测试发送文本消息
echo "📝 测试发送文本消息..."
curl -X POST http://localhost:13001/api/message \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "text",
    "text": "Hello from API!"
  }'

echo -e "\n"

# 测试发送图片（base64）
echo "🖼️  测试发送图片（base64）..."
IMAGE_BASE64=$(base64 -i /Users/ton/Desktop/avatr.png)
curl -X POST http://localhost:13001/api/message \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"type\": \"image\",
    \"data\": \"$IMAGE_BASE64\"
  }"

echo -e "\n"

# 测试未授权访问
echo "❌ 测试未授权访问..."
curl -X POST http://localhost:13001/api/message \
  -H "Content-Type: application/json" \
  -d '{
    "type": "text",
    "text": "This should fail"
  }'

echo -e "\n"
