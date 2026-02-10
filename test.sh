#!/bin/bash
# CICY 系统测试脚本

echo "🧪 CICY 系统测试"
echo "================="

# 测试 1：服务器健康检查
echo ""
echo "测试 1: 服务器健康检查"
HEALTH=$(curl -s http://localhost:13001/health)
if echo "$HEALTH" | grep -q '"status":"ok"'; then
    echo "✅ 服务器正常"
else
    echo "❌ 服务器异常"
    exit 1
fi

# 测试 2：发送消息
echo ""
echo "测试 2: 发送消息"
SEND=$(curl -s -X POST http://localhost:13001/message -H "Content-Type: application/json" -d '{"message":"测试消息"}')
if echo "$SEND" | grep -q '"success":true'; then
    echo "✅ 消息发送成功"
else
    echo "❌ 消息发送失败"
    exit 1
fi

# 测试 3：获取消息
echo ""
echo "测试 3: 获取消息"
MESSAGES=$(curl -s http://localhost:13001/messages)
if echo "$MESSAGES" | grep -q "测试消息"; then
    echo "✅ 消息获取成功"
else
    echo "❌ 消息获取失败"
    exit 1
fi

# 测试 4：MCP 协议
echo ""
echo "测试 4: MCP 协议"
MCP=$(curl -s -X POST http://localhost:13001/mcp -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}')
if echo "$MCP" | grep -q "send_message"; then
    echo "✅ MCP 协议正常"
else
    echo "❌ MCP 协议异常"
    exit 1
fi

# 测试 5：清理消息
echo ""
echo "测试 5: 清理消息"
CLEAR=$(curl -s -X POST http://localhost:13001/mcp -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"clear_messages","arguments":{}}}')
if echo "$CLEAR" | grep -q "cleared"; then
    echo "✅ 消息清理成功"
else
    echo "❌ 消息清理失败"
    exit 1
fi

echo ""
echo "================="
echo "🎉 所有测试通过！"
