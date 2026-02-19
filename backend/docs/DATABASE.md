# 数据库配置指南

Genvid 支持本地 PostgreSQL 和多种云端数据库服务。

## 选项 1: 本地 PostgreSQL (开发推荐)

### 启动本地数据库

```bash
cd backend

# 启动 PostgreSQL 和 Redis
docker-compose up -d

# 检查服务状态
docker-compose ps

# 运行数据库迁移
./scripts/db.sh setup
```

### 连接信息

| 参数 | 值 |
|------|-----|
| Host | localhost |
| Port | 5432 |
| User | postgres |
| Password | postgres |
| Database | genvid |

### 数据库管理

访问 Adminer UI: http://localhost:8081

- 系统: PostgreSQL
- 服务器: postgres
- 用户名: postgres
- 密码: postgres
- 数据库: genvid

---

## 选项 2: Supabase (生产推荐)

### 创建 Supabase 项目

1. 访问 https://supabase.com 并登录
2. 点击 "New Project"
3. 填写项目信息:
   - 名称: genvid
   - 数据库密码: (记录下来)
   - 区域: 选择最近的区域
4. 等待项目创建完成 (约 2 分钟)

### 获取连接信息

1. 进入项目 Dashboard
2. 点击 Settings > Database
3. 找到 "Connection string" > 选择 "URI" 格式
4. 复制连接字符串

### 配置环境变量

编辑 `backend/.env`:

```bash
# Supabase Connection (Transaction 模式 - 推荐)
DATABASE_URL=postgresql://postgres.[project-ref]:[password]@aws-0-[region].pooler.supabase.com:6543/postgres

# Session 模式 (用于迁移)
DATABASE_URL_SESSION=postgresql://postgres.[project-ref]:[password]@aws-0-[region].pooler.supabase.com:5432/postgres

DATABASE_HOST=aws-0-[region].pooler.supabase.com
DATABASE_PORT=6543
DATABASE_USER=postgres.[project-ref]
DATABASE_PASSWORD=[your-password]
DATABASE_NAME=postgres
```

### 运行迁移

```bash
cd backend
./scripts/db.sh setup
```

### 启用 RLS

Supabase 默认启用 RLS。迁移脚本已包含 RLS 策略。

---

## 选项 3: Railway

### 创建数据库

```bash
# 安装 Railway CLI
npm install -g @railway/cli

# 登录
railway login

# 创建 PostgreSQL
railway add --plugin postgresql
```

### 获取连接信息

```bash
railway variables
```

### 配置环境变量

```bash
DATABASE_URL=postgresql://postgres:[password]@[host].railway.app:[port]/railway
```

---

## 选项 4: Neon (Serverless)

### 创建项目

1. 访问 https://neon.tech
2. 创建新项目
3. 复制连接字符串

### 配置环境变量

```bash
DATABASE_URL=postgresql://[user]:[password]@[endpoint].neon.tech/[database]?sslmode=require
```

---

## 数据库迁移

### 运行迁移

```bash
cd backend

# 检查连接
./scripts/db.sh check

# 运行迁移
./scripts/db.sh setup

# 重置数据库 (删除所有数据)
./scripts/db.sh reset
```

### 手动迁移

```bash
# 使用 psql
psql $DATABASE_URL -f migrations/001_initial_schema.sql
psql $DATABASE_URL -f migrations/002_rls_policies.sql
psql $DATABASE_URL -f migrations/003_seed_data.sql
```

---

## 生产环境建议

1. **连接池**: 使用 PgBouncer 或 Supabase Pooler
2. **SSL**: 强制 SSL 连接 (`?sslmode=require`)
3. **备份**: 启用自动备份
4. **监控**: 配置数据库监控告警
5. **迁移**: 在部署前先运行迁移

## 故障排查

### 连接失败

```bash
# 检查网络
ping [database-host]

# 测试连接
psql $DATABASE_URL -c "SELECT 1"
```

### 迁移失败

```bash
# 检查当前版本
psql $DATABASE_URL -c "SELECT * FROM schema_migrations;"

# 手动运行单个迁移
psql $DATABASE_URL -f migrations/001_initial_schema.sql
```

### SSL 错误

添加 SSL 模式到连接字符串:
```bash
DATABASE_URL=postgresql://...?sslmode=require
```
