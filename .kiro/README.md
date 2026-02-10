# CICY 项目 Kiro 配置

这个目录包含 Kiro CLI 的项目特定配置。

## 文件说明

- `steering/dev-rules.md` - 开发规范（始终加载）
- `steering/lang.md` - 语言设置（始终加载）
- `steering/architecture.md` - 系统架构（始终加载）
- `steering/debug-rules.md` - 调试规范（始终加载）⭐
- `steering/testing-rules.md` - 测试与验收规范（始终加载）🚨
- `steering/work-mode.md` - AI 工作模式规范（始终加载）⚡

## 使用方式

当在此项目目录使用 Kiro CLI 时，这些文档会自动加载到上下文中，
确保 AI 助手了解项目规范和架构。

## 添加新规则

创建新的 `.md` 文件并在文件头添加：

```yaml
---
inclusion: always
---
```

这样文件会自动加载到每次对话中。
