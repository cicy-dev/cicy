#!/bin/bash
# 重启 cicy-go

pkill -9 -f cicy-go
sleep 1
cd /Users/ton/Desktop/skills/cicy/server-go
./cicy-go
