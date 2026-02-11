#!/bin/bash
# 发布前检查脚本

echo "🔍 CICY 发布前检查"
echo "==================="

# 检查 Git 状态
echo ""
echo "1. 检查 Git 状态..."
if [ -n "$(git status --porcelain)" ]; then
    echo "⚠️  警告: 有未提交的更改"
    git status --short
else
    echo "✅ Git 工作区干净"
fi

# 检查 npm 登录
echo ""
echo "2. 检查 npm 登录..."
if npm whoami > /dev/null 2>&1; then
    USER=$(npm whoami)
    echo "✅ 已登录: $USER"
else
    echo "❌ 未登录 npm"
    echo "请运行: npm login"
    exit 1
fi

# 检查包名
echo ""
echo "3. 检查包名..."
PACKAGE_NAME=$(node -p "require('./package.json').name")
echo "包名: $PACKAGE_NAME"

# 检查版本
echo ""
echo "4. 检查版本..."
VERSION=$(node -p "require('./package.json').version")
echo "版本: $VERSION"

# 运行测试
echo ""
echo "5. 运行测试..."
if npm test; then
    echo "✅ 测试通过"
else
    echo "❌ 测试失败"
    exit 1
fi

# 检查文件列表
echo ""
echo "6. 检查将要发布的文件..."
npm pack --dry-run

# 检查依赖
echo ""
echo "7. 检查依赖..."
npm outdated || true

echo ""
echo "==================="
echo "✅ 检查完成！"
echo ""
echo "准备发布:"
echo "  npm publish"
echo ""
echo "或者先本地测试:"
echo "  npm pack"
echo "  npm install -g cicy-$VERSION.tgz"
