# Housekeeper 部署指南

## 快速开始

### 方式一：使用 Makefile（推荐）

```bash
# 完整构建（前端+后端）
make build

# 运行服务
make run
```

### 方式二：手动构建

```bash
# 1. 构建前端
cd web
npm install
npm run build
cd ..

# 2. 构建后端
go build -o housekeeper

# 3. 运行
./housekeeper -c ./internal/conf/config.toml
```

## 开发模式

### 后端开发

```bash
# 方式1: 使用Makefile
make dev

# 方式2: 直接运行
go run . -c ./internal/conf/config.toml
```

### 前端开发

```bash
cd web
npm install
npm run dev
```

前端开发服务器会运行在 `http://localhost:3000`，并自动代理API请求到后端 `http://localhost:19090`。

### 同时开发前后端

1. 终端1：启动后端
   ```bash
   make dev
   # 或
   go run . -c ./internal/conf/config.toml
   ```

2. 终端2：启动前端
   ```bash
   cd web
   npm run dev
   ```

访问 `http://localhost:3000` 进行开发。

## 生产部署

### 1. 构建

```bash
make build
```

构建完成后会生成：
- `web/dist/` - 前端静态资源
- `housekeeper` - 后端可执行文件

### 2. 配置

编辑 `internal/conf/config.toml` 配置文件：

```toml
[server]
port = 19090

[database]
# MongoDB配置
# Redis配置

[jwt]
# JWT密钥配置
```

### 3. 运行

```bash
./housekeeper -c ./internal/conf/config.toml
```

### 4. 访问

- 前端页面: `http://your-server:19090/`
- API接口: `http://your-server:19090/v1/`

## 服务管理

### 使用 systemd（Linux）

创建服务文件 `/etc/systemd/system/housekeeper.service`:

```ini
[Unit]
Description=Housekeeper Service
After=network.target

[Service]
Type=simple
User=your-user
WorkingDirectory=/path/to/housekeeper
ExecStart=/path/to/housekeeper/housekeeper -c /path/to/housekeeper/internal/conf/config.toml
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable housekeeper
sudo systemctl start housekeeper
sudo systemctl status housekeeper
```

### 使用 Docker

创建 `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /build
COPY . .
RUN apk add --no-cache nodejs npm make
RUN make build

FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/housekeeper .
COPY --from=builder /build/web/dist ./web/dist
COPY --from=builder /build/internal/conf ./internal/conf
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
EXPOSE 19090
CMD ["./housekeeper", "-c", "./internal/conf/config.toml"]
```

构建和运行：

```bash
docker build -t housekeeper .
docker run -d -p 19090:19090 --name housekeeper housekeeper
```

### 使用 Docker Compose

创建 `docker-compose.yml`:

```yaml
version: '3.8'

services:
  housekeeper:
    build: .
    ports:
      - "19090:19090"
    volumes:
      - ./internal/conf:/app/internal/conf
      - ./logs:/app/logs
    restart: unless-stopped
    depends_on:
      - mongodb
      - redis

  mongodb:
    image: mongo:6
    ports:
      - "27017:27017"
    volumes:
      - mongo_data:/data/db
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    restart: unless-stopped

volumes:
  mongo_data:
  redis_data:
```

启动：

```bash
docker-compose up -d
```

## 反向代理配置

### Nginx

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 前端静态资源
    location / {
        proxy_pass http://localhost:19090;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    # API接口
    location /v1/ {
        proxy_pass http://localhost:19090;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### Caddy

```
your-domain.com {
    reverse_proxy localhost:19090
}
```

## 更新部署

```bash
# 1. 停止服务
systemctl stop housekeeper

# 2. 拉取最新代码
git pull

# 3. 重新构建
make build

# 4. 启动服务
systemctl start housekeeper
```

## 故障排查

### 前端无法访问

检查：
1. `web/dist` 目录是否存在
2. 后端服务是否正常运行
3. 防火墙端口是否开放

### API请求失败

检查：
1. 后端配置文件是否正确
2. 数据库连接是否正常
3. 查看日志文件 `logs/`

### 构建失败

前端构建失败：
```bash
cd web
rm -rf node_modules package-lock.json
npm install
npm run build
```

后端构建失败：
```bash
go mod tidy
go build -o housekeeper
```

## 性能优化

### 前端优化
- 启用 Gzip 压缩（Nginx）
- 配置 CDN
- 静态资源长期缓存

### 后端优化
- 调整数据库连接池
- 启用 Redis 缓存
- 配置日志级别

## 监控

建议监控以下指标：
- 服务运行状态
- API响应时间
- 数据库连接数
- 内存使用率
- 磁盘使用率

可使用 Prometheus + Grafana 进行监控。

## 备份

定期备份：
- MongoDB 数据
- Redis 数据
- 配置文件
- 日志文件

备份脚本示例：
```bash
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
mongodump --out /backup/mongo_$DATE
tar -czf /backup/housekeeper_$DATE.tar.gz ./internal/conf ./logs
```
