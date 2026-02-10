# TODO 汇总 - Worker-1 任务清单

## 当前待完成任务

### 1. ✅ 双击 Ctrl+C 退出
文件：`TODO-double-ctrlc.md`
状态：已下发

### 2. 🔄 Masu 风格界面
文件：`TODO-MASU-FINAL.md`
状态：开发中
- 回显用户输入
- Loading 状态
- 完成提示
- 统计栏

### 3. ⏳ 底部固定布局
**需求**：thinking 和 input 在 TUI 底部固定位置
**状态**：已通知，待创建 TODO

### 4. ⏳ API 随机响应
**需求**：server.js 返回随机句子，不要固定格式
**状态**：已通知，待创建 TODO

## 工作流程

### Master (我) 的职责
1. 接收用户需求
2. 创建 TODO 文档（包含验收标准）
3. 通知 worker-1 执行
4. 测试验收
5. 汇报结果

### Worker-1 的职责
1. 查看 TODO 文档
2. 按要求实现功能
3. 自测通过
4. 等待 master 验收

## 验收流程
1. Master 测试 TUI 功能
2. 对照验收标准逐项检查
3. 发现问题 → 创建修复 TODO
4. 全部通过 → 标记完成 ✅

## 下一步
Master 需要为任务 3 和 4 创建正式 TODO 文档
