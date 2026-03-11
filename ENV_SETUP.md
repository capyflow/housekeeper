# Housekeeper - 环境配置与使用指南

## ✅ 环境已配置完成

### 1. Go 环境配置

**状态：** ✅ 已配置

- **Go 版本:** go1.25.5 linux/amd64
- **安装位置:** `/usr/local/go/bin/go`
- **GOROOT:** `/usr/local/go`
- **GOPATH:** `/home/aaron/go`

**配置文件更新:** `~/.zshrc`

```bash
# Go environment configuration
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$GOROOT/bin:$GOPATH/bin:$PATH
export GO111MODULE=on
```

> 💡 **注意**: 如果您使用的是新打开的终端窗口，需要重新加载配置或重启终端：
> ```bash
> source ~/.zshrc
> ```

### 2. Node.js 环境配置

**状态：** ✅ 已配置

- **Node.js:** v22.22.1
- **npm:** 10.9.4

### 3. 项目依赖

**后端依赖:** ✅ 已初始化（go.mod 存在）
**前端依赖:** ✅ 已安装（node_modules 已创建）

---

## 🚀 启动方式

### 方式一：开发模式（推荐）

分别启动前后端：

#### 1. 启动后端

```bash
cd /home/aaron/.openclaw/workspace/housekeeper
source ~/.zshrc  # 确保 Go 环境生效
go run . -c ./internal/conf/config.toml
```

后端将在 `http://localhost:19090` 启动

#### 2. 启动前端（另一个终端）

```bash
cd /home/aaron/.openclaw/workspace/housekeeper/web
npm run dev
```

前端将在 `http://localhost:3000` 启动

访问：http://localhost:3000/webui

---

### 方式二：生产模式

先构建前端，再启动后端：

```bash
# 1. 构建前端
cd /home/aaron/.openclaw/workspace/housekeeper/web
npm run build

# 2. 启动后端（自动 serving 静态文件）
cd ..
source ~/.zshrc
./housekeeper -c ./internal/conf/config.toml -s ./web/dist
```

访问：http://localhost:19090/webui

---

## 🔧 重新编译

### 修改后端代码后

```bash
cd /home/aaron/.openclaw/workspace/housekeeper
source ~/.zshrc
go build -o housekeeper .
```

### 修改前端代码后

```bash
cd /home/aaron/.openclaw/workspace/housekeeper/web
npm run build
```

---

## 📋 验证环境

运行以下命令检查所有环境：

```bash
# 检查 Go
source ~/.zshrc
go version

# 检查 Node.js
node --version
npm --version

# 检查项目
cd /home/aaron/.openclaw/workspace/housekeeper
ls -la housekeeper  # 应能看到编译好的二进制文件
ls -rf web/dist     # 应能看到构建的前端文件
```

---

## 🐛 常见问题

### Q1: `command not found: go`

**原因:** 终端没有加载 `.zshrc` 配置

**解决:** 
```bash
source ~/.zshrc
```

或者重启终端窗口。

---

### Q2: 前端构建失败

**原因:** node_modules 未安装或损坏

**解决:**
```bash
cd /home/aaron/.openclaw/workspace/housekeeper/web
rm -rf node_modules package-lock.json
npm install
```

---

### Q3: 端口被占用

**原因:** 19090 或 3000 端口已被其他进程占用

**解决:** 
1. 查看占用端口的进程：`lsof -i :19090`
2. 杀死进程：`kill -9 <PID>`
3. 修改配置文件中的端口号

---

### Q4: MongoDB/Redis 连接失败

**原因:** 数据库服务未启动

**解决:** 
使用 Docker Compose 一键启动所有服务：
```bash
cd /home/aaron/.openclaw/workspace/housekeeper
make docker-up
```

---

## 📦 Docker 部署（可选）

如果您的机器上有 Docker，可以使用 Docker Compose 一键启动所有服务：

```bash
# 启动所有服务（Housekeeper + Redis + MongoDB）
cd /home/aaron/.openclaw/workspace/housekeeper
docker-compose -f deploy/docker-compose.yml up -d

# 查看服务状态
docker-compose -f deploy/docker-compose.yml ps

# 查看日志
docker-compose -f deploy/docker-compose.yml logs -f

# 停止服务
docker-compose -f deploy/docker-compose.yml down
```

---

## 🎯 下一步

环境已准备就绪！您可以：

1. **启动应用** - 按照上面的"启动方式"章节操作
2. **登录系统** - 使用您的账号登录
3. **创建打卡任务** - 进入"打卡任务"页面开始使用新功能！
4. **查看文档** - 阅读 `CHECKIN_FEATURES.md` 了解更多功能细节

---

**最后更新时间:** 2026-03-11 15:34
**作者:** 贾维斯 ⚡
