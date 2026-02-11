# 快速发布指南

## 当前状态
✅ 所有测试通过
✅ 已登录 npm (cicybot)
✅ 包名: cicy
✅ 版本: 1.0.0
✅ 文件准备完毕

## 立即发布

### 方法 1: 直接发布
```bash
npm publish
```

### 方法 2: 先本地测试
```bash
# 1. 打包
npm pack

# 2. 全局安装测试
npm install -g cicy-1.0.0.tgz

# 3. 测试命令
cicy-server

# 4. 卸载测试包
npm uninstall -g cicy

# 5. 发布
npm publish
```

## 发布后验证

```bash
# 1. 查看包信息
npm info cicy

# 2. 安装测试
npm install -g cicy

# 3. 运行服务器
cicy-server

# 4. 查看版本
npm view cicy version
```

## 包含的文件

发布包将包含以下文件：
- server.js (MCP 服务器)
- client-tui.js (TUI 客户端)
- client-remote.js (远程客户端)
- README.md (文档)
- AGENTS.md (开发指南)
- test.sh (测试脚本)
- package.json (配置)

## 使用方式

用户安装后可以：

```bash
# 安装
npm install -g cicy

# 启动服务器
cicy-server

# 或者
npx cicy
```

## 注意事项

1. 发布后无法撤销（72小时后）
2. 确保 Git 已提交所有更改
3. 版本号遵循语义化版本
4. 发布后记得打 Git 标签

## Git 标签

```bash
git add .
git commit -m "chore: prepare for v1.0.0 release"
git tag v1.0.0
git push origin main --tags
```

## 下次更新

```bash
# 1. 修改代码
# 2. 更新版本
npm version patch  # 1.0.0 -> 1.0.1

# 3. 发布
npm publish

# 4. 推送标签
git push --tags
```

## 紧急撤销

如果发布后发现问题（72小时内）：

```bash
npm unpublish cicy@1.0.0
```

## 帮助

如有问题，查看：
- PUBLISHING.md - 完整发布指南
- CHANGELOG.md - 版本历史
- README.md - 使用文档
