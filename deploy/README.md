# Housekeeper Docker 部署指南

本目录包含使用 Docker 和 Docker Compose 部署 Housekeeper 的所有配置文件。

## 快速开始

### 1. 构建并启动所有服务

```bash
cd deploy
docker-compose up -d
```

### 2. 仅构建镜像（不启动）

```bash
cd deploy
docker-compose build
```

### 3. 查看服务状态

```bash
docker-compose ps
```

### 4. 查看服务日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f housekeeper
docker-compose logs -f redis
docker-compose logs -f mongo
```

### 5. 停止服务

```bash
docker-compose down
```

### 6. 停止并清理数据

```bash
docker-compose down -v
```

## 服务说明

### Housekeeper 应用

- **容器名**: `housekeeper-app`
- **端口**: `19090`
- **配置文件**: `../internal/conf/config.docker.toml`
- **日志目录**: `./containers/housekeeper/logs`

### Redis

- **容器名**: `housekeeper-redis`
- **端口**: `6379`
- **密码**: `1433223qQ`
- **数据目录**: `./containers/redis/data`

### MongoDB

- **容器名**: `housekeeper-mongo`
- **端口**: `27017`
- **数据库**: `housekeeper`
- **数据目录**: `./containers/mongo/data`

## 目录结构

```
deploy/
├── docker-compose.yml           # Docker Compose 配置
├── README.md                   # 本文件
└── containers/                 # 持久化数据目录（自动创建）
    ├── housekeeper/
    │   └── logs/              # 应用日志
    ├── redis/
    │   └── data/              # Redis 数据
    └── mongo/
        ├── data/              # MongoDB 数据
        └── configdb/          # MongoDB 配置
```

## 配置说明

### 修改配置

如需修改应用配置，编辑 `../internal/conf/config.docker.toml`：

```toml
Port=19090

[RdbConfig.redis]
host="redis"              # Docker 网络内的服务名
port=6379
password="1433223qQ"

[RdbConfig.mongo]
host="mongo"              # Docker 网络内的服务名
port=27017

[Jwt]
secret_key="your-secret"
expire=256000

[Admin]
Username="admin"
Password="your-password"
```

**重要提示**：
- 修改配置后需要重启容器：`docker-compose restart housekeeper`
- Redis 和 MongoDB 的 host 必须使用容器服务名（redis、mongo），不能使用 localhost
- 生产环境请务必修改默认密码和密钥

### 端口映射

默认端口映射：

| 服务 | 容器端口 | 主机端口 |
|------|---------|---------|
| Housekeeper | 19090 | 19090 |
| Redis | 6379 | 6379 |
| MongoDB | 27017 | 27017 |

如需修改主机端口，编辑 `docker-compose.yml` 中的 `ports` 配置：

```yaml
ports:
  - "宿主机端口:容器端口"
```

## 网络

所有服务运行在同一个 Docker 网络 `housekeeper-network` 中，可以通过服务名互相访问。

## 数据持久化

数据卷挂载：

- **Redis**: `./containers/redis/data` → `/data`
- **MongoDB**: `./containers/mongo/data` → `/data/db`
- **应用日志**: `./containers/housekeeper/logs` → `/app/logs`

## 健康检查

服务健康检查配置：

- **Housekeeper**: HTTP 健康检查 `/v1/health`
- **Redis**: Redis CLI ping
- **MongoDB**: mongosh ping

查看健康状态：

```bash
docker-compose ps
```

## 故障排查

### 1. 服务无法启动

检查日志：

```bash
docker-compose logs housekeeper
```

常见问题：
- 端口被占用：修改 `docker-compose.yml` 中的端口映射
- 配置文件错误：检查 `config.docker.toml` 语法
- 数据库连接失败：确认 Redis 和 MongoDB 已启动

### 2. 前端无法访问

检查：
- 浏览器访问 `http://localhost:19090/`
- 检查防火墙是否开放 19090 端口
- 查看容器日志是否有错误

### 3. 数据库连接失败

确认：
- Redis 和 MongoDB 容器状态：`docker-compose ps`
- 配置文件中的 host 使用服务名（redis、mongo）
- 网络连接正常：`docker network inspect deploy_housekeeper-network`

### 4. 重建容器

完全重建（不删除数据）：

```bash
docker-compose down
docker-compose up -d --build
```

完全重建（删除所有数据）：

```bash
docker-compose down -v
docker-compose up -d --build
```

## 生产部署建议

### 1. 安全配置

- 修改所有默认密码
- 修改 JWT 密钥
- 使用环境变量管理敏感信息
- 限制端口访问范围

### 2. 性能优化

```yaml
# 添加资源限制
deploy:
  resources:
    limits:
      cpus: '2'
      memory: 2G
    reservations:
      cpus: '1'
      memory: 1G
```

### 3. 日志管理

```yaml
# 添加日志驱动
logging:
  driver: "json-file"
  options:
    max-size: "10m"
    max-file: "3"
```

### 4. 备份策略

定期备份数据：

```bash
# 备份 MongoDB
docker exec housekeeper-mongo mongodump --out /data/backup

# 备份 Redis
docker exec housekeeper-redis redis-cli --rdb /data/dump.rdb
```

### 5. 使用 Nginx 反向代理

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:19090;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## 更新应用

### 1. 拉取最新代码

```bash
cd /path/to/housekeeper
git pull
```

### 2. 重新构建镜像

```bash
cd deploy
docker-compose build housekeeper
```

### 3. 重启服务

```bash
docker-compose up -d housekeeper
```

### 4. 验证更新

```bash
docker-compose logs -f housekeeper
```

## 监控

### 查看资源使用

```bash
docker stats
```

### 查看容器详情

```bash
docker inspect housekeeper-app
```

### 进入容器调试

```bash
docker exec -it housekeeper-app sh
```

## 常用命令

```bash
# 启动所有服务
docker-compose up -d

# 停止所有服务
docker-compose down

# 重启特定服务
docker-compose restart housekeeper

# 查看实时日志
docker-compose logs -f

# 清理未使用的镜像
docker image prune -a

# 清理未使用的卷
docker volume prune

# 查看网络
docker network ls
```

## 支持

如有问题，请查看：
- [主 README](../README.md)
- [部署指南](../DEPLOYMENT.md)
- 或提交 Issue

---

**Housekeeper** - Docker 容器化部署 🐳
