---
inclusion: always
---

# 文件管理规范

## 核心原则
❌ **禁止**在项目根目录乱写文件
✅ **必须**按分类存放到指定目录

## 目录结构

### 1. TODO 文件
```
docs/todo/
├── active/           # 进行中的任务
├── completed/        # 已完成的任务
└── archived/         # 已归档的任务
```

**命名规范**：
- `TODO-功能名称.md`
- 例如：`TODO-masu-style.md`

### 2. 测试记录
```
docs/tests/
├── manual/           # 手动测试记录
├── automated/        # 自动化测试
└── reports/          # 测试报告
```

### 3. 开发文档
```
docs/dev/
├── architecture.md   # 架构文档
├── api.md           # API 文档
└── guides/          # 开发指南
```

### 4. 临时文件
```
temp/
├── logs/            # 临时日志
├── cache/           # 缓存文件
└── scratch/         # 草稿文件
```

## 文件操作规则

### ✅ 允许的操作
1. 在指定目录创建文件
2. 更新现有文件
3. 移动文件到正确分类
4. 删除过期临时文件

### ❌ 禁止的操作
1. 在根目录创建 TODO 文件
2. 在根目录创建测试文件
3. 在根目录创建临时文件
4. 创建未分类的文件

## 工作流程

### 创建 TODO
```bash
# ❌ 错误
/Users/ton/Desktop/skills/cicy/TODO-new-feature.md

# ✅ 正确
/Users/ton/Desktop/skills/cicy/docs/todo/active/TODO-new-feature.md
```

### 完成任务后
```bash
# 移动到 completed
mv docs/todo/active/TODO-feature.md docs/todo/completed/
```

### 归档旧任务
```bash
# 移动到 archived
mv docs/todo/completed/TODO-old.md docs/todo/archived/
```

## 当前清理

### 需要整理的文件
```
根目录的 TODO 文件：
- TODO-blessed-layout.md
- TODO-loading-fix.md
- TODO-loading.md
- TODO-masu-style.md
- TODO-random-response.md
- TODO-MASU-FINAL.md
- TODO-double-ctrlc.md
- TODO-loading-fix2.md
- TODO-SUMMARY.md
```

### 整理方案
1. 创建目录结构
2. 按状态分类移动
3. 更新引用路径
4. 删除重复文件

## AI 助手规则

### 创建文件前必须：
1. ✅ 确定文件类型
2. ✅ 选择正确目录
3. ✅ 使用规范命名
4. ✅ 检查是否已存在

### 禁止行为：
1. ❌ 在根目录创建 TODO
2. ❌ 创建未分类文件
3. ❌ 使用随意命名
4. ❌ 创建重复文件

## 检查清单

创建文件时问自己：
- [ ] 这个文件属于什么类型？
- [ ] 应该放在哪个目录？
- [ ] 文件名是否符合规范？
- [ ] 是否已经存在类似文件？

## 记住

> **分类存放，井然有序**
> 
> **根目录只放核心文件**
> 
> **文档归档，便于查找**
