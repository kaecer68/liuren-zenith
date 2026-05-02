#!/bin/bash
# check-version-consistency.sh
# 檢查 liuren-zenith 版本號在各文件中是否一致

set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VERSION_FILE="$PROJECT_ROOT/VERSION"

# 1. 讀取 VERSION 文件（唯一來源）
if [[ ! -f "$VERSION_FILE" ]]; then
    echo "❌ VERSION 文件不存在: $VERSION_FILE"
    exit 1
fi

VERSION=$(cat "$VERSION_FILE" | tr -d '[:space:]')
echo "📄 VERSION 文件版本: $VERSION"

# 2. 檢查 go.mod 模塊路徑
GO_MOD_VERSION=$(sed -n 's/^module .*\/v\([0-9]*\).*/\1/p' "$PROJECT_ROOT/go.mod" || true)
MAJOR_VERSION="${VERSION#v}"
MAJOR_VERSION="${MAJOR_VERSION%%.*}"
if [[ -n "$GO_MOD_VERSION" && "$GO_MOD_VERSION" != "$MAJOR_VERSION" ]]; then
    echo "⚠️  go.mod 主版本不一致: v$GO_MOD_VERSION (期望: v$MAJOR_VERSION)"
else
    echo "✅ go.mod 主版本一致: v${GO_MOD_VERSION:-未設置}"
fi

# 3. 檢查 /health endpoint 代碼中是否有 version 欄位
if grep -q '"version"\|serviceVersion' "$PROJECT_ROOT/internal/httpapi/liuren.go"; then
    echo "✅ /health endpoint 包含 version 欄位"
else
    echo "❌ /health endpoint 缺少 version 欄位"
    exit 1
fi

# 4. 檢查 AGENTS.md 版本聲明是否一致
AGENTS_VERSION=$(grep -o 'Version: [^ ]*' "$PROJECT_ROOT/AGENTS.md" | awk '{print $2}' || true)
if [[ -n "$AGENTS_VERSION" && "$AGENTS_VERSION" != "$VERSION" ]]; then
    echo "⚠️  AGENTS.md 版本不一致: $AGENTS_VERSION (期望: $VERSION)"
else
    echo "✅ AGENTS.md 版本一致: ${AGENTS_VERSION:-未設置}"
fi

echo ""
echo "🎉 liuren-zenith 版本號檢查通過！"
