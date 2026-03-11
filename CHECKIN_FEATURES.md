# 打卡任务模块 - 实现文档

## 功能概述

这是一个完整的多任务打卡系统，支持：
- 📋 **创建多个打卡任务**（如健身、阅读、学英语等）
- 🎯 **每日打卡** - 每天只能打卡一次
- 📊 **热力图统计** - GitHub 风格的打卡热力图展示
- ⏰ **时间段限制** - 可为每个任务设置允许的打卡时间窗口
- 📈 **统计分析** - 总打卡天数、连续打卡天数统计
- 💾 **补打卡** - 点击热力图上任意历史日期可补打卡

## 技术架构

### 后端 (Go)

#### 新增文件

1. **`internal/model/checkin.go`** - 数据模型定义
   - `CheckInTarget` - 打卡任务核心结构
   - `TimeRange` - 时间段配置
   - `CheckInRequest/Response` - 打卡请求响应

2. **`internal/repository/checkin_repo.go`** - 数据访问层
   - CRUD 操作接口
   - 打卡记录原子更新
   - 时间段校验逻辑

3. **`internal/service/checkin.go`** - 业务逻辑层
   - 权限验证（所有者控制）
   - 业务规则处理

4. **`internal/handler/checkin.go`** - HTTP 处理器
   - RESTful API 端点
   - 请求参数验证
   - 响应格式化

#### 修改文件

1. **`bootstrap.go`** - 添加 CheckInRepo/Service/Handler 初始化
2. **`internal/router/routers.go`** - 注册打卡任务路由
3. **`web/src/types/index.ts`** - 前端类型定义扩展
4. **`web/src/api/checkin.ts`** - 新增前端 API 封装
5. **`web/src/views/CheckInTarget.vue`** - 全新的打卡任务页面
6. **`web/src/router/index.ts`** - 添加检查入口路由
7. **`web/src/layouts/MainLayout.vue`** - 导航菜单更新

### API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/v1/checkin/target/create` | 创建打卡任务 |
| PUT | `/v1/checkin/target/update` | 更新打卡任务 |
| DELETE | `/v1/checkin/target/delete` | 删除打卡任务 |
| POST | `/v1/checkin/target/info` | 查询任务详情 |
| POST | `/v1/checkin/target/list` | 分页查询任务列表 |
| POST | `/v1/checkin/check_in` | 执行打卡 |
| POST | `/v1/checkin/stats` | 获取打卡统计 |

### 数据结构

```typescript
// 打卡任务
interface CheckInTarget {
  id: string
  title: string              // 任务名称
  description: string        // 描述
  create_time: number
  update_time: number
  owner: string              // 所有者
  time_range: TimeRange      // 时间段设置
  total_days: number         // 总打卡天数
  check_in_counts: {         // 每日打卡次数
    [key: string]: number    // key: YYYY-MM-DD, value: 次数
  }
}

// 时间段
interface TimeRange {
  enabled: boolean   // 是否启用限制
  start: number      // 开始分钟 (0-1439, 0=00:00, 1439=23:59)
  end: number        // 结束分钟
}
```

### 数据库设计

**集合名**: `checkin_target`

```javascript
{
  "_id": "ct_xxx",                 // MongoDB ObjectID
  "id": "ct_xxx",                  // 业务 ID
  "title": "晨跑",                 // 任务名称
  "description": "每天早上跑步 5 公里", 
  "owner": "user_xxx",             // 所有者
  "create_time": 1709856000000,    // 创建时间戳
  "update_time": 1709856000000,    // 更新时间戳
  "time_range": {                  // 时间段
    "enabled": true,
    "start": 360,   // 06:00
    "end": 720      // 12:00
  },
  "total_days": 30,                // 总打卡天数
  "check_in_counts": {             // 打卡统计
    "2026-03-01": 1,
    "2026-03-02": 1,
    ...
  }
}
```

## 前端特性

### 可视化界面

1. **Hero 区域**
   - 当前选中的任务标题
   - 今日打卡按钮（状态自动切换）
   - 新建任务按钮

2. **任务列表**
   - Chip 卡片形式展示所有任务
   - 显示总打卡天数和时间段
   - 点击切换当前任务

3. **热力图**
   - 近 365 天打卡情况
   - GitHub 风格的色阶显示
   - 点击格子补打卡
   - 本周标记（周一、三、五）

4. **统计摘要**
   - 总打卡天数
   - 当前连续打卡天数

5. **创建/编辑弹窗**
   - 任务名称（必填）
   - 描述（可选）
   - 时间段设置（开启/关闭限制）
   - 开始/结束时间输入

### 交互逻辑

- ✅ **今日打卡**：点击按钮直接打卡，已打卡显示"✅ 今日已打卡"
- 📅 **补打卡**：点击热力图任意历史格子即可补打卡
- ⚠️ **时间段校验**：如果启用了时间段，非允许时间内打卡会提示
- 🔄 **自动刷新**：打卡成功后本地数据实时更新
- 📱 **响应式设计**：完美适配移动端

## 使用方法

### 1. 启动服务

```bash
# 构建前端
cd web
npm install
npm run build

# 启动后端
cd ..
go run . -c ./internal/conf/config.toml -s ./web/dist
```

### 2. 访问应用

打开浏览器访问：http://localhost:19090/webui

### 3. 操作流程

1. 登录系统
2. 点击左侧菜单"打卡任务"
3. 点击"新建任务"创建第一个打卡任务
4. 填写任务信息，可选择设置时间段
5. 点击"今日打卡"开始坚持
6. 在热力图上查看打卡记录
7. 忘记打卡时，点击对应日期补打卡

## 注意事项

1. **时间格式**: 时间段使用"分钟数"表示，06:00 = 360 分钟
2. **跨天支持**: 支持跨天时间段（如 22:00-06:00）
3. **原子操作**: 打卡是原子操作，避免重复打卡
4. **权限控制**: 只能操作自己的打卡任务
5. **数据隔离**: 不同用户的打卡数据完全独立

## 后续优化方向

- 🌟 添加打卡提醒推送
- 📊 导出打卡统计数据
- 🏆 排行榜功能
- 👥 共享打卡任务
- 🎨 自定义热力图配色
- 🔔 打卡成就系统

---

**作者**: 贾维斯  
**日期**: 2026-03-11
