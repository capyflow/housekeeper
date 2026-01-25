#!/bin/bash

# Housekeeper 数据恢复脚本

set -e

if [ -z "$1" ]; then
    echo "用法: bash restore.sh <备份文件.tar.gz>"
    echo ""
    echo "可用的备份文件:"
    ls -lh backup/*.tar.gz 2>/dev/null || echo "  未找到备份文件"
    exit 1
fi

BACKUP_FILE="$1"
BACKUP_DIR="./backup"
RESTORE_DIR="./restore_tmp"

if [ ! -f "$BACKUP_FILE" ] && [ ! -f "$BACKUP_DIR/$BACKUP_FILE" ]; then
    echo "错误: 备份文件不存在: $BACKUP_FILE"
    exit 1
fi

# 如果只提供了文件名，添加路径
if [ ! -f "$BACKUP_FILE" ]; then
    BACKUP_FILE="$BACKUP_DIR/$BACKUP_FILE"
fi

echo "================================================"
echo "   Housekeeper 数据恢复"
echo "================================================"
echo ""
echo "⚠️  警告: 此操作将覆盖现有数据！"
echo "备份文件: $BACKUP_FILE"
echo ""
echo "是否继续？(yes/NO)"
read -r response
if [ "$response" != "yes" ]; then
    echo "操作已取消"
    exit 0
fi

echo ""
echo ">>> 解压备份文件..."
mkdir -p "$RESTORE_DIR"
tar -xzf "$BACKUP_FILE" -C "$RESTORE_DIR"
echo "✓ 备份文件解压完成"
echo ""

# 查找解压后的目录
MONGO_DIR=$(find "$RESTORE_DIR" -name "mongo_*" -type d | head -n 1)
REDIS_FILE=$(find "$RESTORE_DIR" -name "redis_*.rdb" | head -n 1)

if [ -n "$MONGO_DIR" ]; then
    echo ">>> 恢复 MongoDB 数据..."
    docker cp "$MONGO_DIR" housekeeper-mongo:/data/restore_tmp
    docker exec housekeeper-mongo mongorestore --drop /data/restore_tmp --quiet
    docker exec housekeeper-mongo rm -rf /data/restore_tmp
    echo "✓ MongoDB 数据恢复完成"
    echo ""
fi

if [ -n "$REDIS_FILE" ]; then
    echo ">>> 恢复 Redis 数据..."
    # 停止 Redis 写入
    docker exec housekeeper-redis redis-cli save > /dev/null 2>&1 || true
    docker cp "$REDIS_FILE" housekeeper-redis:/data/dump.rdb
    docker-compose restart redis
    echo "✓ Redis 数据恢复完成"
    echo ""
fi

# 清理临时文件
echo ">>> 清理临时文件..."
rm -rf "$RESTORE_DIR"
echo "✓ 清理完成"
echo ""

echo "================================================"
echo "   恢复完成"
echo "================================================"
echo ""
echo "请重启服务以确保所有更改生效:"
echo "  docker-compose restart"
echo ""
