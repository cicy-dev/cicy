#!/bin/bash
# 测试 SSH 选择功能

echo "🧪 测试 SSH 选择功能"
echo "===================="

# 1. 停止旧进程
echo "1️⃣ 停止旧进程..."
pkill -9 -f cicy-go
sleep 1

# 2. 创建 tmux session
echo "2️⃣ 创建 tmux session..."
tmux kill-session -t cicy-test 2>/dev/null
tmux new-session -d -s cicy-test -x 120 -y 40

# 3. 启动 cicy-go
echo "3️⃣ 启动 cicy-go..."
tmux send-keys -t cicy-test "cd /Users/ton/Desktop/skills/cicy/server-go && ./cicy-go" C-m
echo "   等待 5 秒让程序完全启动..."
sleep 5

# 4. 发送 /ssh 命令
echo "4️⃣ 发送 /ssh 命令..."
tmux send-keys -t cicy-test "/ssh" C-m
echo "   等待 2 秒让界面渲染..."
sleep 2

# 5. 捕获输出
echo "5️⃣ 捕获输出..."
tmux capture-pane -t cicy-test -p > /tmp/cicy-ssh-test.txt

# 6. 检查结果
echo ""
echo "📊 测试结果："
echo "===================="
cat /tmp/cicy-ssh-test.txt
echo "===================="
echo ""

# 7. 验证
if grep -q "ssh gcp" /tmp/cicy-ssh-test.txt; then
    echo "✅ SSH 选择框显示正常"
else
    echo "❌ SSH 选择框未显示"
fi

if grep -q "▶" /tmp/cicy-ssh-test.txt; then
    echo "✅ 选中标记显示正常"
else
    echo "❌ 选中标记未显示"
fi

echo ""
echo "💡 提示："
echo "   - 查看完整输出: cat /tmp/cicy-ssh-test.txt"
echo "   - 连接到 session: tmux attach -t cicy-test"
echo "   - 停止测试: tmux kill-session -t cicy-test"
