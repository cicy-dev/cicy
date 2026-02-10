#!/bin/bash
# Master 永不停止工作循环

WORKER_SESSION="workers:worker-1"
CHECK_INTERVAL=30

echo "=== Master 工作循环启动 $(date) ==="

while true; do
    echo "[$(date '+%H:%M:%S')] 检查 Worker-1..."
    
    LAST_LINE=$(tmux capture-pane -t "$WORKER_SESSION" -p | tail -1)
    
    if echo "$LAST_LINE" | grep -q "What would you like to do next?"; then
        echo "[$(date '+%H:%M:%S')] Worker-1 停了 - 发送指令"
        tmux send-keys -t "$WORKER_SESSION" "继续工作！检查 docs/todo/active/ 任务，不要停，持续改进代码直到完美！" C-m
    elif echo "$LAST_LINE" | grep -q "Thinking..."; then
        echo "[$(date '+%H:%M:%S')] Worker-1 工作中"
    else
        echo "[$(date '+%H:%M:%S')] Worker-1 状态: $LAST_LINE"
    fi
    
    sleep "$CHECK_INTERVAL"
done
