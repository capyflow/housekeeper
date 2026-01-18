# Housekeeper Web

基于 Vue 3 + TypeScript + Vite 的管理后台。

## 功能特性

- 用户登录认证
- 仪表盘首页
- 共享看板管理（CRUD操作）
- 响应式设计
- 状态管理（Pinia）
- 路由守卫

## 技术栈

- Vue 3
- TypeScript
- Vite
- Vue Router
- Pinia
- Axios

## 开发

### 安装依赖

```bash
npm install
```

### 启动开发服务器

```bash
npm run dev
```

应用将运行在 http://localhost:3000

### 构建生产版本

```bash
npm run build
```

### 预览生产版本

```bash
npm run preview
```

## 项目结构

```
src/
├── api/          # API 请求
├── assets/       # 静态资源
├── components/   # 公共组件
├── layouts/      # 布局组件
├── router/       # 路由配置
├── stores/       # 状态管理
├── styles/       # 全局样式
├── types/        # TypeScript 类型定义
├── views/        # 页面组件
├── App.vue       # 根组件
└── main.ts       # 入口文件
```

## API 接口

应用通过代理连接后端服务：

- 开发环境：`/api` -> `http://localhost:8080`
- API 基础路径：`/api/v1`

### 接口列表

- `POST /api/v1/notes/share_board/create` - 创建看板
- `PUT /api/v1/notes/share_board/update` - 更新看板
- `DELETE /api/v1/notes/share_board/delete` - 删除看板
- `POST /api/v1/notes/share_board/list` - 查询看板列表

## 默认登录

开发环境使用 Mock 登录，输入任意用户名和密码即可登录。
