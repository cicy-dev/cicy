#!/bin/bash
# 测试 CICY API - MCP 格式

TOKEN=$(cat ~/data/cicy-server.txt)

echo "测试 MCP content 数组格式..."
echo ""

# 测试 1: 单个文本消息
echo "1️⃣ 单个文本消息"
curl -X POST http://localhost:13001/api/message \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": [
      {
        "type": "text",
        "text": "Hello from MCP format!"
      }
    ]
  }'
echo -e "\n"

# 测试 2: 多个文本消息
echo "2️⃣ 多个文本消息"
curl -X POST http://localhost:13001/api/message \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": [
      {
        "type": "text",
        "text": "First message"
      },
      {
        "type": "text",
        "text": "Second message"
      },
      {
        "type": "text",
        "text": "Third message"
      }
    ]
  }'
echo -e "\n"

# 测试 3: 文本 + 图片
echo "3️⃣ 文本 + 图片"
IMAGE_BASE64=$(base64 -i /Users/ton/Desktop/avatr.png)
curl -X POST http://localhost:13001/api/message \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"content\": [
      {
        \"type\": \"text\",
        \"text\": \"Here is an image:\"
      },
      {
        \"type\": \"image\",
        \"data\": \"$IMAGE_BASE64\"
      }
    ]
  }"
echo -e "\n"

echo "✅ 测试完成"
