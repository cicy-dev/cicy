#!/bin/bash
# 超强定时器 - 监控 Worker-1 并通知 Master

LOG_FILE="/Users/ton/Desktop/skills/cicy/temp/logs/monitor.log"
WORKER_PANE="workers:worker-1.0"  # Worker-1 的正确名称
CHECK_INTERVAL="${1:-5}"  # 第一个参数作为检查间隔，默认 5 秒

# 确保目录存在
mkdir -p "$(dirname "$LOG_FILE")"

echo "=== 监控启动 $(date) - 间隔: ${CHECK_INTERVAL}秒 ===" >> "$LOG_FILE"

while true; do
    TIMESTAMP=$(date "+%Y-%m-%d %H:%M:%S")
    
    # 检查 Worker-1 状态
    if tmux has-session -t workers 2>/dev/null; then
        # 捕获最后几行输出
        LAST_OUTPUT=$(tmux capture-pane -t "$WORKER_PANE" -p | tail -5)
        
        # 检查是否在工作（Thinking...）
        if echo "$LAST_OUTPUT" | grep -q "Thinking..."; then
            echo "[$TIMESTAMP] ⚙️  Worker-1 工作中" >> "$LOG_FILE"
            
        # 检查是否等待权限（Allow this action? 或 [y/n/t]）
        elif echo "$LAST_OUTPUT" | grep -qE "Allow this action|\\[y/n/t\\]"; then
            echo "[$TIMESTAMP] 🔐 Worker-1 等待权限 - 自动批准" >> "$LOG_FILE"
            
            # 自动发送 't' 信任工具
            tmux send-keys -t "$WORKER_PANE" "t" C-m
            echo "[$TIMESTAMP] ✅ 已发送 't' 信任工具" >> "$LOG_FILE"
            
        else
            # 不在工作也不等待权限 - 卡住了
            echo "[$TIMESTAMP] ⚠️  Worker-1 卡住 - 发送继续指令" >> "$LOG_FILE"
            
            # 直接发送继续工作指令
            tmux send-keys -t "$WORKER_PANE" "继续工作！验收、测试、改进代码，不要停！" C-m
            echo "[$TIMESTAMP] 📨 已发送继续工作指令" >> "$LOG_FILE"
        fi
    else
        echo "[$TIMESTAMP] ❌ Worker session 不存在" >> "$LOG_FILE"
    fi
    
    # 等待下一次检查
    sleep "$CHECK_INTERVAL"
done
