# Genvid - TikTok UGC 视频生成平台

一个专为跨境电商卖家设计的 AI UGC 视频生成平台，帮助用户在几分钟内创建专业、真实的 TikTok/Reels/Shorts 风格产品视频。

## 项目结构

```
Genvid/
├── .lane/plans/              # 架构文档
│   ├── spec.md               # 产品规格说明书
│   ├── data-model.md         # 数据库设计
│   └── architecture.md       # 系统架构设计
├── research.md               # 市场调研报告
├── README.md
├── backend/                  # Golang 后端服务
│   ├── cmd/server/           # 入口文件
│   ├── internal/
│   │   ├── config/           # 配置管理
│   │   ├── handler/          # HTTP 处理器
│   │   ├── zhipu/            # 智谱 CogVideoX API 客户端
│   │   ├── middleware/       # 中间件
│   │   ├── model/            # 数据模型
│   │   ├── repository/       # 数据访问层
│   │   ├── service/          # 业务逻辑层
│   │   └── worker/           # 视频生成 Worker
│   ├── pkg/auth/             # JWT 认证
│   ├── migrations/           # 数据库迁移
│   ├── Dockerfile
│   ├── docker-compose.yml
│   └── Makefile
└── frontend/                 # Next.js 16 前端
    └── src/
        ├── app/              # App Router 页面
        │   ├── create/       # 视频创建流程
        │   ├── dashboard/    # 用户仪表盘
        │   ├── login/        # 登录页
        │   ├── pricing/      # 定价页
        │   └── register/     # 注册页
        ├── components/
        │   ├── features/     # 功能组件
        │   ├── layout/       # 布局组件
        │   └── ui/           # UI 基础组件
        ├── hooks/            # React Hooks
        ├── lib/api/          # API 客户端
        ├── stores/           # Zustand 状态
        └── types/            # TypeScript 类型
```

## 已实现功能

### 后端 API
- ✅ 用户注册/登录 (JWT 认证)
- ✅ 用户信息管理
- ✅ 项目 CRUD
- ✅ 视频生成 API
- ✅ Avatar 列表 API
- ✅ Stripe 支付集成
- ✅ 智谱 CogVideoX 视频生成客户端
- ✅ 视频生成 Worker

### 前端页面
- ✅ Landing Page (Hero, Features, How It Works, Testimonials)
- ✅ 登录/注册页面
- ✅ Dashboard 仪表盘
- ✅ 视频创建流程 (4 步引导)
- ✅ Avatar 选择器
- ✅ AI 脚本生成器
- ✅ 定价页面

### 数据库
- ✅ 8 张核心表
- ✅ RLS 安全策略
- ✅ 种子数据

## 快速开始

### 环境要求

- Go 1.22+
- Node.js 18+
- Docker (本地开发) 或 云端数据库

### 数据库配置

支持 **本地 PostgreSQL** 或 **云端数据库** (Supabase/Railway/Neon)。

#### 选项 1: 本地开发 (推荐)

```bash
cd backend

# 复制配置文件
cp .env.example .env

# 启动 PostgreSQL + Redis
docker-compose up -d

# 运行数据库迁移
./scripts/db.sh setup

# 启动后端服务
make run
```

#### 选项 2: Supabase (生产推荐)

```bash
cd backend

# 1. 在 https://supabase.com 创建项目
# 2. 获取数据库连接字符串 (Settings > Database)
# 3. 编辑 .env 文件:
DATABASE_URL=postgresql://postgres.[ref]:[password]@aws-0-[region].pooler.supabase.com:6543/postgres

# 4. 运行迁移
./scripts/db.sh setup

# 5. 启动服务
make run
```

详细配置: [backend/docs/DATABASE.md](backend/docs/DATABASE.md)

### 前端启动

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

访问 http://localhost:3000

## API 端点

| 端点 | 方法 | 描述 |
|------|------|------|
| `/health` | GET | 健康检查 |
| `/api/auth/register` | POST | 用户注册 |
| `/api/auth/login` | POST | 用户登录 |
| `/api/auth/refresh` | POST | 刷新 Token |
| `/api/user/profile` | GET/PATCH | 用户信息 |
| `/api/projects` | GET/POST | 项目列表/创建 |
| `/api/projects/:id` | GET/DELETE | 项目详情/删除 |
| `/api/projects/:id/generate` | POST | 生成视频 |
| `/api/avatars` | GET | Avatar 列表 |
| `/api/upload` | POST | 上传图片 |
| `/api/payments/checkout` | POST | 创建支付会话 |
| `/api/payments/webhook` | POST | Stripe Webhook |

## 环境变量配置

### 后端 (backend/.env)

```bash
# 数据库 - Docker 本地开发
DATABASE_URL=postgresql://postgres:postgres@localhost:5433/genvid?sslmode=disable

# JWT
JWT_SECRET=your-jwt-secret

# 智谱 AI - 视频生成
ZHIPU_API_KEY=your-zhipu-api-key
ZHIPU_MODEL=cogvideox-3

# Stripe - 支付
STRIPE_SECRET_KEY=sk_test_xxx
STRIPE_WEBHOOK_SECRET=whsec_xxx

# 其他
OPENAI_API_KEY=sk-xxx
RESEND_API_KEY=re_xxx
```

### 前端 (frontend/.env.local)

```bash
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## 注意事项

### 数据库迁移

`./scripts/db.sh setup` 脚本不会自动读取 `.env` 文件，推荐使用以下方式：

```bash
# 方式 1: 加载环境变量后执行
export $(grep -v '^#' .env | xargs) && ./scripts/db.sh setup

# 方式 2: 使用 Docker 执行迁移
docker exec -i genvid-postgres psql -U postgres -d genvid < ./migrations/001_initial_schema.sql
```

### 视频生成

- **文生视频**: 只需要 `prompt` 参数
- **图生视频**: 需要 `prompt` + `image_url` 参数（图片会自动转为 base64）
- 图片格式支持: JPG, PNG, GIF, WebP（最大 10MB）
- 视频生成时间: 约 2-5 分钟

### 端口冲突

Docker PostgreSQL 使用 **5433** 端口（避免与本地 PostgreSQL 5432 冲突）。

## 技术栈

### 前端
- Next.js 16 (App Router)
- React 19
- TypeScript
- Tailwind CSS
- Zustand (状态管理)
- Framer Motion (动画)

### 后端
- Golang 1.22
- PostgreSQL 16 (Supabase)
- Redis 7
- JWT 认证

### 外部服务
- 智谱 CogVideoX (视频生成)
- Stripe (支付)
- Resend (邮件)
- OpenAI (脚本生成)
- AWS S3 (存储)

## 开发命令

### 后端

```bash
make build          # 构建
make run            # 运行
make test           # 测试
make docker-up      # 启动 Docker 服务
```

### 前端

```bash
npm run dev         # 开发模式
npm run build       # 构建
npm run start       # 生产模式
npm run lint        # 代码检查
```

## 部署

### 前端 (Vercel)
```bash
cd frontend
vercel
```

### 后端 (Railway/Fly.io)
```bash
cd backend
fly deploy
```

## 下一步优化

1. **Supabase 集成**: 替换本地 PostgreSQL
2. **真实 CogVideoX API**: 配置生产 API Key
3. **Redis 队列**: 使用 BullMQ 替代内存队列
4. **CDN 配置**: Cloudflare 加速视频下载
5. **监控告警**: 添加 Sentry + Grafana

## License

MIT
