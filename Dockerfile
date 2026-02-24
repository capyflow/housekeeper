# ==========================================
# 阶段1: 构建前端
# ==========================================
FROM node:18-alpine AS frontend-builder

WORKDIR /app/web

# 复制前端依赖文件
COPY web/package*.json ./

# 安装前端依赖
RUN npm install

# 复制前端源代码
COPY web/ ./

# 构建前端
RUN npm run build

# ==========================================
# 阶段2: 构建后端
# ==========================================
FROM golang:1.25-alpine AS backend-builder

WORKDIR /app

# 安装必要的编译工具
RUN apk add --no-cache git wget vim

# 复制 Go 模块依赖文件
COPY go.mod go.sum ./

# 复制 vendor 目录
COPY vendor vendor

# 复制后端源代码
COPY . .

# 构建后端可执行文件
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o housekeeper .

# ==========================================
# 阶段3: 运行时镜像
# ==========================================
FROM alpine:3.19

WORKDIR /app

# 安装运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 从构建阶段复制可执行文件
COPY --from=backend-builder /app/housekeeper .

# 从前端构建阶段复制静态文件
COPY --from=frontend-builder /app/web/dist ./web/dist

# 复制配置文件（注意：生产环境建议使用挂载卷）
COPY internal/conf/config.toml ./internal/conf/config.toml

# 复制locale国际化文件（如果存在）
COPY locale ./locale

# 暴露端口
EXPOSE 19090

# 设置可执行权限
RUN chmod +x ./housekeeper

# 启动应用
CMD ["./housekeeper", "-c", "./internal/conf/config.toml", "-s", "./web/dist"]
