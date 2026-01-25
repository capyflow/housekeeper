#!/bin/bash

# Housekeeper 部署测试脚本

set -e

echo "================================================"
echo "   Housekeeper 部署测试"
echo "================================================"
echo ""

BASE_URL="http://localhost:19090"
API_URL="$BASE_URL/v1"

# 测试前端页面
echo ">>> 测试前端页面..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/")
if [ "$HTTP_CODE" == "200" ]; then
    echo "✓ 前端页面访问正常 (HTTP $HTTP_CODE)"
else
    echo "✗ 前端页面访问失败 (HTTP $HTTP_CODE)"
    exit 1
fi
echo ""

# 测试静态资源
echo ">>> 测试静态资源加载..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/assets/")
if [ "$HTTP_CODE" == "200" ] || [ "$HTTP_CODE" == "404" ]; then
    echo "✓ 静态资源路由正常"
else
    echo "✗ 静态资源路由异常 (HTTP $HTTP_CODE)"
fi
echo ""

# 测试 API 可访问性
echo ">>> 测试 API 可访问性..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$API_URL/user/login")
if [ "$HTTP_CODE" == "200" ] || [ "$HTTP_CODE" == "400" ] || [ "$HTTP_CODE" == "401" ]; then
    echo "✓ API 接口可访问 (HTTP $HTTP_CODE)"
else
    echo "✗ API 接口访问失败 (HTTP $HTTP_CODE)"
    exit 1
fi
echo ""

# 测试登录接口
echo ">>> 测试登录接口..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/user/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"change-this-password"}')

echo "登录响应: $LOGIN_RESPONSE"
echo ""

# 检查 Docker 容器状态
echo ">>> 检查 Docker 容器状态..."
echo ""
docker-compose ps
echo ""

# 检查容器日志是否有错误
echo ">>> 检查容器日志..."
ERROR_COUNT=$(docker-compose logs housekeeper 2>&1 | grep -i "error\|panic\|fatal" | wc -l)
if [ "$ERROR_COUNT" -gt 0 ]; then
    echo "⚠️  检测到 $ERROR_COUNT 个错误日志"
    echo "查看详细日志: docker-compose logs housekeeper"
else
    echo "✓ 未检测到错误日志"
fi
echo ""

# 检查资源使用
echo ">>> 检查资源使用情况..."
docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}"
echo ""

echo "================================================"
echo "   测试完成"
echo "================================================"
echo ""
echo "访问地址:"
echo "  - 前端: $BASE_URL/"
echo "  - API: $API_URL/"
echo ""
echo "如有问题，请检查:"
echo "  1. docker-compose logs -f housekeeper"
echo "  2. docker-compose ps"
echo "  3. netstat -tlnp | grep 19090"
echo ""
