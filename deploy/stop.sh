#!/bin/bash

# Housekeeper Docker 停止脚本

set -e

echo "================================================"
echo "   停止 Housekeeper 服务"
echo "================================================"
echo ""

# 检查是否需要删除数据
if [ "$1" == "--clean" ] || [ "$1" == "-c" ]; then
    echo "⚠️  警告: 即将停止服务并删除所有数据！"
    echo "此操作不可恢复，是否继续？(y/N)"
    read -r response
    if [ "$response" != "y" ] && [ "$response" != "Y" ]; then
        echo "操作已取消"
        exit 0
    fi

    echo ">>> 停止并删除所有服务和数据..."
    docker-compose down -v

    echo ">>> 删除数据目录..."
    rm -rf containers/

    echo "✓ 服务已停止，数据已清理"
else
    echo ">>> 停止服务（保留数据）..."
    docker-compose down
    echo "✓ 服务已停止"
fi

echo ""
echo "================================================"
echo "   停止完成"
echo "================================================"
echo ""
echo "重新启动服务:"
echo "  bash start.sh"
echo ""
echo "完全清理（包括数据）:"
echo "  bash stop.sh --clean"
echo ""
