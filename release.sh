#!/bin/bash
# CICY 一键发布脚本
# 用法: ./release.sh <version> [npm]

set -e

VERSION=${1:-$(date +%Y.%m.%d)}
PUBLISH_NPM=${2:-false}

REPO="cicy-dev/cicy"
GITHUB_TOKEN=${GITHUB_TOKEN:-""}
NPM_TOKEN=${NPM_TOKEN:-""}

echo "🚀 CICY 一键发布"
echo "================"
echo "版本: v$VERSION"
echo "发布到 npm: $PUBLISH_NPM"
echo ""

# 检查 GitHub Token
if [ -z "$GITHUB_TOKEN" ]; then
    echo "❌ 请设置 GITHUB_TOKEN 环境变量"
    echo "export GITHUB_TOKEN='your_github_token'"
    exit 1
fi

# 1. 更新版本号
echo "📝 更新 package.json 版本..."
npm version $VERSION --no-git-tag-version 2>/dev/null || {
    # 手动更新版本
    sed -i "s/\"version\": \"[0-9.]*\"/\"version\": \"$VERSION\"/" package.json
}
echo "✅ 版本更新为 $VERSION"

# 2. 构建 Go 平台二进制文件
echo "🔨 构建 Go 二进制文件..."
export GOROOT=/usr/local/go
export PATH=$GOROOT/bin:$PATH

cd server-go
GOOS=linux GOARCH=amd64 go build -o cicy-go-linux-amd64 .
GOOS=darwin GOARCH=amd64 go build -o cicy-go-darwin-amd64 .
GOOS=darwin GOARCH=arm64 go build -o cicy-go-darwin-arm64 .

tar -czf cicy-go-linux-amd64.tar.gz cicy-go-linux-amd64
tar -czf cicy-go-darwin-amd64.tar.gz cicy-go-darwin-amd64
tar -czf cicy-go-darwin-arm64.tar.gz cicy-go-darwin-arm64
cd ..

echo "✅ 构建完成"

# 3. 创建 GitHub Release
echo "📦 创建 GitHub Release v$VERSION..."

# 创建 Release
RELEASE_RESPONSE=$(curl -s -X POST \
    -H "Authorization: token $GITHUB_TOKEN" \
    -H "Accept: application/vnd.github.v3+json" \
    https://api.github.com/repos/$REPO/releases \
    -d "{
        \"tag_name\": \"v$VERSION\",
        \"name\": \"CICY v$VERSION\",
        \"body\": \"## 主要变更\n- 自动发布 Go 二进制文件\n- 支持多平台 (Linux/macOS)\",
        \"draft\": false,
        \"prerelease\": false
    }")

RELEASE_ID=$(echo $RELEASE_RESPONSE | grep -o '"id": [0-9]*' | head -1 | cut -d' ' -f2)

if [ -z "$RELEASE_ID" ]; then
    echo "❌ 创建 Release 失败"
    echo $RELEASE_RESPONSE
    exit 1
fi

echo "✅ Release 创建成功: $RELEASE_ID"

# 4. 上传二进制文件
echo "⬆️  上传二进制文件..."

upload_asset() {
    local file=$1
    local name=$2
    curl -s -X POST \
        -H "Authorization: token $GITHUB_TOKEN" \
        -H "Content-Type: application/gzip" \
        --data-binary "@$file" \
        "https://uploads.github.com/repos/$REPO/releases/$RELEASE_ID/assets?name=$name"
    echo "✅ 上传 $name"
}

cd server-go
upload_asset "cicy-go-linux-amd64.tar.gz" "cicy-go-linux-amd64.tar.gz"
upload_asset "cicy-go-darwin-amd64.tar.gz" "cicy-go-darwin-amd64.tar.gz"
upload_asset "cicy-go-darwin-arm64.tar.gz" "cicy-go-darwin-arm64.tar.gz"
cd ..

# 5. 发布到 npm (可选)
if [ "$PUBLISH_NPM" = "true" ]; then
    if [ -z "$NPM_TOKEN" ]; then
        echo "⚠️  未设置 NPM_TOKEN，跳过 npm 发布"
    else
        echo "📦 发布到 npm..."
        echo "//registry.npmjs.org/:_authToken=$NPM_TOKEN" > .npmrc
        npm publish
        rm -f .npmrc
        echo "✅ npm 发布成功"
    fi
fi

# 6. 提交更改
echo "📝 提交更改..."
git add package.json server-go/.gitignore
git commit -m "chore: release v$VERSION"
git tag "v$VERSION"
git push origin main
git push origin "v$VERSION"

echo ""
echo "🎉 发布完成！"
echo "================"
echo "GitHub Release: https://github.com/$REPO/releases/tag/v$VERSION"
