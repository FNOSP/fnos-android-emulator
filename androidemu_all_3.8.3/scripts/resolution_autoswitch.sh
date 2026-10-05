#!/bin/bash
# 分辨率自动切换守护启动脚本（架构自适应）
# 根据 CPU 架构选择对应的 Go 二进制：x86_64 用 amd64，ARM64 用 arm64
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ARCH=$(uname -m)

case "$ARCH" in
    x86_64|amd64)
        BIN="${SCRIPT_DIR}/resolution_autoswitch_amd64"
        ;;
    aarch64|arm64)
        BIN="${SCRIPT_DIR}/resolution_autoswitch_arm64"
        ;;
    *)
        echo "不支持的架构: $ARCH" >&2
        exit 1
        ;;
esac

if [ ! -x "$BIN" ]; then
    chmod +x "$BIN" 2>/dev/null || true
fi

exec "$BIN"
