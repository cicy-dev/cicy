# 发布到 npm 指南

## 准备工作

### 1. 注册 npm 账号
```bash
# 如果还没有账号，访问 https://www.npmjs.com/signup 注册

# 登录 npm
npm login
```

### 2. 检查包名是否可用
```bash
npm search cicy
```

如果名称已被占用，需要修改 `package.json` 中的 `name` 字段。

## 发布步骤

### 1. 更新版本号
```bash
# 补丁版本 (1.0.0 -> 1.0.1)
npm version patch

# 次要版本 (1.0.0 -> 1.1.0)
npm version minor

# 主要版本 (1.0.0 -> 2.0.0)
npm version major
```

### 2. 测试包
```bash
# 运行测试
npm test

# 检查将要发布的文件
npm pack --dry-run

# 本地测试安装
npm pack
npm install -g cicy-1.0.0.tgz
```

### 3. 发布到 npm
```bash
# 发布
npm publish

# 如果是 scoped package
npm publish --access public
```

### 4. 验证发布
```bash
# 查看包信息
npm info cicy

# 测试安装
npm install -g cicy
```

## 发布检查清单

- [ ] 代码已提交到 Git
- [ ] 版本号已更新
- [ ] README.md 文档完整
- [ ] 测试通过 (`npm test`)
- [ ] 依赖版本正确
- [ ] .npmignore 配置正确
- [ ] package.json 中的 files 字段正确
- [ ] 已登录 npm (`npm whoami`)

## 更新已发布的包

```bash
# 1. 修改代码
# 2. 更新版本号
npm version patch

# 3. 发布新版本
npm publish
```

## 撤销发布

```bash
# 撤销特定版本（发布后 72 小时内）
npm unpublish cicy@1.0.0

# 撤销整个包（慎用！）
npm unpublish cicy --force
```

## 发布 CLI 包

如果要发布 CLI 包 (cicy-cli)：

```bash
# 1. 准备 CLI 包
cd /path/to/cicy

# 2. 使用 package-cli.json
cp package-cli.json package.json

# 3. 发布
npm publish

# 4. 恢复原 package.json
git checkout package.json
```

## 常见问题

### 1. 包名已被占用
修改 `package.json` 中的 `name`，例如：
- `@yourusername/cicy`
- `cicy-mcp`
- `cicy-server`

### 2. 权限错误
```bash
# 确保已登录
npm whoami

# 重新登录
npm logout
npm login
```

### 3. 发布失败
```bash
# 检查网络
npm config get registry

# 使用官方源
npm config set registry https://registry.npmjs.org/
```

## 最佳实践

1. **语义化版本**：遵循 semver 规范
   - MAJOR: 不兼容的 API 变更
   - MINOR: 向后兼容的功能新增
   - PATCH: 向后兼容的问题修复

2. **变更日志**：维护 CHANGELOG.md

3. **标签**：为每个版本打 Git 标签
   ```bash
   git tag v1.0.0
   git push --tags
   ```

4. **文档**：保持 README.md 更新

5. **测试**：发布前运行完整测试

## 自动化发布

可以使用 GitHub Actions 自动发布：

```yaml
# .github/workflows/publish.yml
name: Publish to npm

on:
  release:
    types: [created]

jobs:
  publish:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-node@v2
        with:
          node-version: '14'
          registry-url: 'https://registry.npmjs.org'
      - run: npm ci
      - run: npm test
      - run: npm publish
        env:
          NODE_AUTH_TOKEN: ${{ secrets.NPM_TOKEN }}
```

## 相关链接

- [npm 文档](https://docs.npmjs.com/)
- [语义化版本](https://semver.org/lang/zh-CN/)
- [npm 包发布指南](https://docs.npmjs.com/packages-and-modules/contributing-packages-to-the-registry)
