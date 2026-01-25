#!/bin/bash

# Housekeeper Docker 快速启动脚本

set -e

echo "================================================"
echo "   Housekeeper Docker 部署脚本"
echo "================================================"
echo ""

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo "错误: 未检测到 Docker，请先安装 Docker"
    echo "安装指南: https://docs.docker.com/get-docker/"
    exit 1
fi

# 检查 Docker Compose 是否安装
if ! command -v docker-compose &> /dev/null; then
    echo "错误: 未检测到 Docker Compose，请先安装 Docker Compose"
    echo "安装指南: https://docs.docker.com/compose/install/"
    exit 1
fi

echo "✓ Docker 环境检查通过"
echo ""

# 创建必要的目录
echo ">>> 创建数据目录..."
mkdir -p containers/housekeeper/logs
mkdir -p containers/redis/data
mkdir -p containers/mongo/data
mkdir -p containers/mongo/configdb
echo "✓ 数据目录创建完成"
echo ""

# 检查配置文件
if [ ! -f "../internal/conf/config.docker.toml" ]; then
    echo "警告: config.docker.toml 不存在，将使用默认配置"
    echo "生产环境请务必修改配置文件中的密码和密钥！"
fi
echo ""

# 构建镜像
echo ">>> 构建 Docker 镜像..."
docker-compose build
echo "✓ 镜像构建完成"
echo ""

# 启动服务
echo ">>> 启动服务..."
docker-compose up -d
echo "✓ 服务启动完成"
echo ""

# 等待服务启动
echo ">>> 等待服务启动..."
sleep 5

# 检查服务状态
echo ""
echo "================================================"
echo "   服务状态"
echo "================================================"
docker-compose ps
echo ""

# 显示访问信息
echo "================================================"
echo "   部署成功！"
echo "================================================"
echo ""
echo "访问地址:"
echo "  - 前端页面: http://localhost:19090/"
echo "  - API 接口: http://localhost:19090/v1/"
echo ""
echo "默认账号:"
echo "  - 用户名: admin"
echo "  - 密码: change-this-password"
echo ""
echo "常用命令:"
echo "  - 查看日志: docker-compose logs -f"
echo "  - 停止服务: docker-compose down"
echo "  - 重启服务: docker-compose restart"
echo ""
echo "⚠️  生产环境部署提醒:"
echo "  1. 修改 config.docker.toml 中的所有默认密码"
echo "  2. 修改 JWT 密钥"
echo "  3. 配置防火墙规则"
echo "  4. 配置 HTTPS (建议使用 Nginx 反向代理)"
echo ""
echo "详细文档: deploy/README.md"
echo "================================================"
