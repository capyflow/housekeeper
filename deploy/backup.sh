#!/bin/bash

# Housekeeper 数据备份脚本

set -e

BACKUP_DIR="./backup"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_NAME="housekeeper_backup_$DATE"

echo "================================================"
echo "   Housekeeper 数据备份"
echo "================================================"
echo ""

# 创建备份目录
mkdir -p "$BACKUP_DIR"

echo ">>> 备份 MongoDB 数据..."
docker exec housekeeper-mongo mongodump --out "/data/backup_$DATE" --quiet
docker cp housekeeper-mongo:/data/backup_$DATE "$BACKUP_DIR/mongo_$DATE"
docker exec housekeeper-mongo rm -rf "/data/backup_$DATE"
echo "✓ MongoDB 备份完成: $BACKUP_DIR/mongo_$DATE"
echo ""

echo ">>> 备份 Redis 数据..."
docker exec housekeeper-redis redis-cli save > /dev/null 2>&1 || true
docker cp housekeeper-redis:/data/dump.rdb "$BACKUP_DIR/redis_$DATE.rdb"
echo "✓ Redis 备份完成: $BACKUP_DIR/redis_$DATE.rdb"
echo ""

echo ">>> 备份配置文件..."
mkdir -p "$BACKUP_DIR/config_$DATE"
cp ../internal/conf/config.docker.toml "$BACKUP_DIR/config_$DATE/" 2>/dev/null || true
echo "✓ 配置文件备份完成"
echo ""

echo ">>> 创建备份压缩包..."
cd "$BACKUP_DIR"
tar -czf "$BACKUP_NAME.tar.gz" mongo_$DATE redis_$DATE.rdb config_$DATE
rm -rf mongo_$DATE redis_$DATE.rdb config_$DATE
cd ..
echo "✓ 备份压缩包创建完成: $BACKUP_DIR/$BACKUP_NAME.tar.gz"
echo ""

# 获取文件大小
BACKUP_SIZE=$(du -h "$BACKUP_DIR/$BACKUP_NAME.tar.gz" | cut -f1)

echo "================================================"
echo "   备份完成"
echo "================================================"
echo ""
echo "备份文件: $BACKUP_DIR/$BACKUP_NAME.tar.gz"
echo "备份大小: $BACKUP_SIZE"
echo ""
echo "恢复备份:"
echo "  bash restore.sh $BACKUP_NAME.tar.gz"
echo ""
echo "清理旧备份（保留最近7天）:"
echo "  find $BACKUP_DIR -name '*.tar.gz' -mtime +7 -delete"
echo ""
