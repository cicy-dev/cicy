# TODO: TUI Loading 功能

## 当前状态
✅ 已实现基础 loading 功能：
- 使用旋转符号动画
- 绿色显示
- 在发送消息时显示
- 响应后自动消失

## 需要测试的验收标准
- [x] 输入消息后立即显示绿色 loading
- [x] Loading 动画流畅（旋转符号）
- [x] 服务器响应后 loading 消失
- [x] 不影响日志输出
- [x] 不影响提示符显示
- [x] curl-rpc 命令也显示 loading

## 测试步骤
1. 启动 API 服务器（cicy:api.0）
2. 在 workers:worker-1.0 运行 `npm run dev:tui`
3. 输入普通消息，观察 loading
4. 输入 curl-rpc 命令，观察 loading
5. 测试多次，确认稳定

## 问题修复
如果发现问题：
- Loading 位置不对 → 调整显示位置
- 动画不流畅 → 调整 interval 时间
- 没有消失 → 检查 stopLoading 调用

## 完成标准
✅ 所有验收标准通过
✅ README.md 已更新
✅ 测试结果已汇报

## 测试记录
- 测试时间：2026-02-11 03:18
- 测试环境：cicy:tui.0 + cicy:api.0
- 测试结果：所有验收标准 100% 通过
- 状态：✅ 已完成
