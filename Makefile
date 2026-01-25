.PHONY: build-web build-server build clean run dev help

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
