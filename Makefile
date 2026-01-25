.PHONY: build-web build-server build clean run dev help docker-build docker-up docker-down docker-logs docker-restart

# 默认目标
all: build

# 帮助信息
help:
	@echo "Housekeeper 构建命令："
	@echo "  make build-web    - 构建前端资源"
	@echo "  make build-server - 构建后端服务"
	@echo "  make build        - 构建前端+后端（完整构建）"
	@echo "  make clean        - 清理构建文件"
	@echo "  make run          - 运行服务"
	@echo "  make dev          - 开发模式（仅启动后端，前端需单独启动）"
	@echo ""
	@echo "Docker 命令："
	@echo "  make docker-build   - 构建 Docker 镜像"
	@echo "  make docker-up      - 启动所有 Docker 服务"
	@echo "  make docker-down    - 停止所有 Docker 服务"
	@echo "  make docker-logs    - 查看 Docker 服务日志"
	@echo "  make docker-restart - 重启 Housekeeper 服务"

# 构建前端
build-web:
	@echo ">>> 构建前端资源..."
	cd web && npm install && npm run build
	@echo ">>> 前端构建完成: web/dist/"

# 构建后端
build-server:
	@echo ">>> 构建后端服务..."
	go build -o housekeeper
	@echo ">>> 后端构建完成: housekeeper"

# 完整构建
build: build-web build-server
	@echo ">>> 构建完成！"
	@echo ">>> 运行: ./housekeeper -c ./internal/conf/config.toml"

# 清理
clean:
	@echo ">>> 清理构建文件..."
	rm -rf web/dist
	rm -f housekeeper
	@echo ">>> 清理完成"

# 运行服务
run:
	@echo ">>> 启动服务..."
	./housekeeper -c ./internal/conf/config.toml

# 开发模式（仅后端）
dev:
	@echo ">>> 开发模式启动后端..."
	@echo ">>> 前端开发服务器请在web目录运行: npm run dev"
	go run . -c ./internal/conf/config.toml

# ==========================================
# Docker 命令
# ==========================================

# 构建 Docker 镜像
docker-build:
	@echo ">>> 构建 Docker 镜像..."
	cd deploy && docker-compose build
	@echo ">>> Docker 镜像构建完成"

# 启动所有 Docker 服务
docker-up:
	@echo ">>> 启动 Docker 服务..."
	cd deploy && docker-compose up -d
	@echo ">>> Docker 服务已启动"
	@echo ">>> 访问: http://localhost:19090"

# 停止所有 Docker 服务
docker-down:
	@echo ">>> 停止 Docker 服务..."
	cd deploy && docker-compose down
	@echo ">>> Docker 服务已停止"

# 查看 Docker 服务日志
docker-logs:
	@echo ">>> 查看 Docker 服务日志..."
	cd deploy && docker-compose logs -f

# 重启 Housekeeper 服务
docker-restart:
	@echo ">>> 重启 Housekeeper 服务..."
	cd deploy && docker-compose restart housekeeper
	@echo ">>> 服务已重启"

# 查看 Docker 服务状态
docker-status:
	@echo ">>> Docker 服务状态："
	cd deploy && docker-compose ps
