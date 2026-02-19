# Genvid - 系统架构设计 (System Architecture)

**版本**: 1.0  
**日期**: 2026-02-19

---

## 1. 系统架构概览

### 1.1 整体架构图

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                                    Genvid 系统架构                                        │
├─────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                         │
│    ┌─────────────────────────────────────────────────────────────────────────────┐     │
│    │                              用户层 (Client Layer)                            │     │
│    │                                                                             │     │
│    │   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐   │     │
│    │   │   Browser   │   │   Mobile    │   │  Desktop    │   │   API       │   │     │
│    │   │  (Web App)  │   │   (PWA)     │   │   (Electron)│   │  Clients    │   │     │
│    │   └──────┬──────┘   └──────┬──────┘   └──────┬──────┘   └──────┬──────┘   │     │
│    └──────────┼─────────────────┼─────────────────┼─────────────────┼──────────┘     │
│               │                 │                 │                 │                 │
│               └─────────────────┴─────────────────┴─────────────────┘                 │
│                                          │                                             │
│                                          ▼                                             │
│    ┌─────────────────────────────────────────────────────────────────────────────┐   │
│    │                              CDN / 静态资源层                                  │   │
│    │                                                                             │   │
│    │   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐                       │   │
│    │   │ Cloudflare  │   │    Vercel   │   │   AWS S3    │                       │   │
│    │   │     CDN     │   │   Edge      │   │  (Videos)   │                       │   │
│    │   └──────┬──────┘   └──────┬──────┘   └──────┬──────┘                       │   │
│    └──────────┼─────────────────┼─────────────────┼──────────────────────────────┘   │
│               │                 │                 │                                   │
│               └─────────────────┴─────────────────┘                                   │
│                                          │                                             │
│  ════════════════════════════════════════╪════════════════════════════════════════   │
│                                          │                                             │
│    ┌─────────────────────────────────────────────────────────────────────────────┐   │
│    │                              前端层 (Frontend Layer)                          │   │
│    │                                                                             │   │
│    │   ┌─────────────────────────────────────────────────────────────────────┐   │   │
│    │   │                     Next.js 15 App Router                            │   │   │
│    │   │                                                                     │   │   │
│    │   │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  │   │   │
│    │   │  │ Landing │  │Dashboard│  │ Create  │  │ Pricing │  │  Auth   │  │   │   │
│    │   │  │  Page   │  │   Page  │  │  Flow   │  │  Page   │  │  Pages  │  │   │   │
│    │   │  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘  │   │   │
│    │   │       │            │            │            │            │        │   │   │
│    │   │       └────────────┴────────────┴────────────┴────────────┘        │   │   │
│    │   │                                │                                   │   │   │
│    │   │  ┌─────────────────────────────┴─────────────────────────────┐    │   │   │
│    │   │  │              Server Actions / API Routes                  │    │   │   │
│    │   │  └─────────────────────────────┬─────────────────────────────┘    │   │   │
│    │   └────────────────────────────────┼──────────────────────────────────┘   │   │
│    └────────────────────────────────────┼──────────────────────────────────────┘   │
│                                         │                                              │
│  ═══════════════════════════════════════╪════════════════════════════════════════   │
│                                         │                                              │
│    ┌────────────────────────────────────┴──────────────────────────────────────┐   │
│    │                              后端层 (Backend Layer)                         │   │
│    │                                                                            │   │
│    │   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐                   │   │
│    │   │   Golang    │    │   Supabase  │    │   Redis     │                   │   │
│    │   │   Auth      │◄───│   Postgres  │◄───│   Queue     │                   │   │
│    │   │   Service   │    │   DB        │    │   (BullMQ)  │                   │   │
│    │   └──────┬──────┘    └──────┬──────┘    └──────┬──────┘                   │   │
│    │          │                  │                  │                          │   │
│    │          │    ┌─────────────┴─────────────┐    │                          │   │
│    │          │    │       Video Workers       │    │                          │   │
│    │          │    │  ┌─────────────────────┐  │    │                          │   │
│    │          │    │  │ Worker 1 │ Worker 2 │  │    │                          │   │
│    │          │    │  │(Python)  │ (Python) │  │    │                          │   │
│    │          │    │  └─────────────────────┘  │    │                          │   │
│    │          │    └─────────────┬─────────────┘    │                          │   │
│    │          │                  │                  │                          │   │
│    └──────────┼──────────────────┼──────────────────┼──────────────────────────┘   │
│               │                  │                  │                               │
│  ═════════════╪══════════════════╪══════════════════╪═══════════════════════════   │
│               │                  │                  │                               │
│    ┌──────────┴──────────────────┴──────────────────┴──────────────────────────┐   │
│    │                              外部服务层 (External Services)                 │   │
│    │                                                                            │   │
│    │   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐                   │   │
│    │   │   HeyGen    │    │   Stripe    │    │   Resend    │                   │   │
│    │   │   API       │    │   Payments  │    │   Email     │                   │   │
│    │   └─────────────┘    └─────────────┘    └─────────────┘                   │   │
│    │                                                                            │   │
│    │   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐                   │   │
│    │   │  ElevenLabs │    │   OpenAI    │    │   AWS S3    │                   │   │
│    │   │   Voice     │    │   GPT-4     │    │   Storage   │                   │   │
│    │   └─────────────┘    └─────────────┘    └─────────────┘                   │   │
│    │                                                                            │   │
│    └───────────────────────────────────────────────────────────────────────────┘   │
│                                                                                      │
└──────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. 技术栈详情

### 2.1 前端技术栈

| 组件 | 技术 | 版本 | 用途 |
|------|------|------|------|
| **框架** | Next.js | 15.x | App Router, SSR, API Routes |
| **UI 库** | React | 19.x | 组件化 UI |
| **样式** | Tailwind CSS | 3.x | 原子化 CSS |
| **类型** | TypeScript | 5.x | 类型安全 |
| **状态管理** | Zustand | 4.x | 客户端状态 |
| **数据获取** | TanStack Query | 5.x | 服务端状态缓存 |
| **表单** | React Hook Form | 7.x | 表单处理 |
| **验证** | Zod | 3.x | 数据验证 |
| **动画** | Framer Motion | 11.x | 交互动画 |
| **组件库** | Radix UI | 1.x | 无障碍基础组件 |
| **图标** | Lucide React | 0.x | SVG 图标 |

### 2.2 后端技术栈

| 组件 | 技术 | 版本 | 用途 |
|------|------|------|------|
| **认证服务** | Golang | 1.22+ | JWT 认证、OAuth |
| **数据库** | PostgreSQL | 16.x | 主数据存储 |
| **BaaS** | Supabase | Latest | 数据库托管、认证、实时 |
| **缓存/队列** | Redis | 7.x | 任务队列、缓存 |
| **任务队列** | BullMQ | 5.x | 视频生成任务 |
| **视频处理** | Python Workers | 3.11+ | 视频生成调用 |

### 2.3 外部服务

| 服务 | 用途 | 定价模型 |
|------|------|---------|
| **HeyGen API** | AI Avatar 视频生成 | $0.10-0.50/秒 |
| **Kling AI** | 备选视频生成 | 按次计费 |
| **Stripe** | 支付处理 | 2.9% + $0.30/笔 |
| **Resend** | 邮件发送 | $0.10/100封 |
| **ElevenLabs** | 语音合成 | $5/月起 |
| **OpenAI** | 脚本生成 | GPT-4 API |
| **AWS S3** | 视频存储 | $0.023/GB |
| **Cloudflare** | CDN + 图片 | 免费层 + Pro |

---

## 3. 认证架构

### 3.1 认证流程

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              认证流程图                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐   │
│  │  User   │────▶│  Frontend   │────▶│   Golang    │────▶│  Supabase   │   │
│  │         │     │  (Next.js)  │     │   Auth      │     │   Auth      │   │
│  └─────────┘     └─────────────┘     └─────────────┘     └─────────────┘   │
│       │                │                    │                    │          │
│       │                │                    │                    │          │
│       │   1. Login     │                    │                    │          │
│       │   Request      │                    │                    │          │
│       │───────────────▶│                    │                    │          │
│       │                │                    │                    │          │
│       │                │  2. Forward to     │                    │          │
│       │                │     Golang Auth    │                    │          │
│       │                │───────────────────▶│                    │          │
│       │                │                    │                    │          │
│       │                │                    │  3. Validate &     │          │
│       │                │                    │     Create Session │          │
│       │                │                    │───────────────────▶│          │
│       │                │                    │                    │          │
│       │                │                    │  4. Return JWT     │          │
│       │                │                    │◀───────────────────│          │
│       │                │                    │                    │          │
│       │                │  5. Return Token   │                    │          │
│       │                │◀───────────────────│                    │          │
│       │                │                    │                    │          │
│       │  6. Set Cookie │                    │                    │          │
│       │◀───────────────│                    │                    │          │
│       │                │                    │                    │          │
│       │                │                    │                    │          │
│       │  7. API Calls  │                    │                    │          │
│       │  with Token    │                    │                    │          │
│       │───────────────▶│───────────────────▶│───────────────────▶│          │
│       │                │                    │                    │          │
│       │  8. Response   │                    │                    │          │
│       │◀───────────────│◀───────────────────│◀───────────────────│          │
│       │                │                    │                    │          │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 JWT Token 结构

```json
{
  "header": {
    "alg": "RS256",
    "typ": "JWT",
    "kid": "key-id"
  },
  "payload": {
    "sub": "user-uuid",
    "email": "user@example.com",
    "role": "authenticated",
    "app_metadata": {
      "subscription_tier": "pro",
      "credits_remaining": 45
    },
    "iat": 1708300000,
    "exp": 1708303600
  }
}
```

### 3.3 认证 API 端点

| 端点 | 方法 | 描述 |
|------|------|------|
| `/api/auth/register` | POST | 邮箱注册 |
| `/api/auth/login` | POST | 邮箱登录 |
| `/api/auth/google` | GET | Google OAuth |
| `/api/auth/callback` | GET | OAuth 回调 |
| `/api/auth/logout` | POST | 登出 |
| `/api/auth/refresh` | POST | 刷新 Token |
| `/api/auth/reset-password` | POST | 重置密码 |

---

## 4. 视频生成架构

### 4.1 视频生成流程

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              视频生成流程                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────┐     ┌─────────┐     ┌─────────┐     ┌─────────┐               │
│  │  User   │────▶│  Next.js│────▶│  Queue  │────▶│ Worker  │               │
│  │         │     │  API    │     │  Redis  │     │ Python  │               │
│  └─────────┘     └─────────┘     └─────────┘     └─────────┘               │
│       │               │               │               │                     │
│       │ 1. Submit     │               │               │                     │
│       │    Request    │               │               │                     │
│       │──────────────▶│               │               │                     │
│       │               │               │               │                     │
│       │               │ 2. Create Job │               │                     │
│       │               │──────────────▶│               │                     │
│       │               │               │               │                     │
│       │               │ 3. Return     │               │                     │
│       │               │    Job ID     │               │                     │
│       │◀──────────────│               │               │                     │
│       │               │               │               │                     │
│       │               │               │ 4. Dequeue    │                     │
│       │               │               │    Job        │                     │
│       │               │               │──────────────▶│                     │
│       │               │               │               │                     │
│       │               │               │               │ 5. Call HeyGen      │
│       │               │               │               │    / Kling API      │
│       │               │               │               │                     │
│       │               │               │               │    ┌─────────┐      │
│       │               │               │               │    │ HeyGen  │      │
│       │               │               │               │    │   API   │      │
│       │               │               │               │    └────┬────┘      │
│       │               │               │               │         │           │
│       │               │               │               │ 6. Poll │           │
│       │               │               │               │    Status          │
│       │               │               │               │◀────────┘           │
│       │               │               │               │                     │
│       │               │               │               │ 7. Download Video   │
│       │               │               │               │    Upload to S3     │
│       │               │               │               │                     │
│       │               │               │               │    ┌─────────┐      │
│       │               │               │               │    │  AWS    │      │
│       │               │               │               │    │   S3    │      │
│       │               │               │               │    └────┬────┘      │
│       │               │               │               │         │           │
│       │               │               │               │ 8. Update│          │
│       │               │               │               │    DB     │         │
│       │               │               │               │         │           │
│       │               │               │               │    ┌────▼────┐      │
│       │               │               │               │    │Supabase │      │
│       │               │               │               │    │   DB    │      │
│       │               │               │               │    └─────────┘      │
│       │               │               │               │                     │
│       │ 9. Webhook/   │               │               │                     │
│       │    Poll for   │               │               │                     │
│       │    Status     │               │               │                     │
│       │──────────────▶│               │               │                     │
│       │               │               │               │                     │
│       │ 10. Return    │               │               │                     │
│       │     Video URL │               │               │                     │
│       │◀──────────────│               │               │                     │
│       │               │               │               │                     │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 4.2 队列任务结构

```typescript
// BullMQ Job 数据结构
interface VideoGenerationJob {
  id: string;                    // 任务 ID
  name: 'video-generation';      // 任务类型
  
  data: {
    projectId: string;           // 项目 UUID
    userId: string;              // 用户 UUID
    
    // 视频配置
    avatarId: string;            // Avatar ID
    script: string;              // 脚本文本
    language: string;            // 语言代码
    format: '9:16' | '1:1';      // 视频格式
    
    // 产品信息
    productName: string;
    productDescription?: string;
    productImages: string[];     // 图片 URL 数组
    
    // 回调配置
    webhookUrl?: string;         // 完成回调 URL
    notificationEmail?: string;  // 通知邮箱
  };
  
  options: {
    attempts: 3;                 // 重试次数
    backoff: {
      type: 'exponential';
      delay: 5000;               // 5 秒初始延迟
    },
    removeOnComplete: 100;       // 保留最近 100 条完成记录
    removeOnFail: 500;           // 保留最近 500 条失败记录
  };
}
```

### 4.3 Worker 处理逻辑

```python
# Python Worker 伪代码
async def process_video_generation(job: VideoGenerationJob):
    """处理视频生成任务"""
    
    # 1. 更新状态为 processing
    await update_project_status(job.project_id, 'processing')
    
    try:
        # 2. 选择视频生成服务
        provider = select_provider(job.data)
        
        # 3. 调用视频生成 API
        if provider == 'heygen':
            result = await heygen_client.generate_video(
                avatar_id=job.data['avatarId'],
                script=job.data['script'],
                voice=job.data['language']
            )
        elif provider == 'kling':
            result = await kling_client.generate_video(
                image=job.data['productImages'][0],
                prompt=job.data['script']
            )
        
        # 4. 轮询等待完成
        video_url = await poll_for_completion(
            task_id=result.task_id,
            provider=provider,
            timeout=300  # 5 分钟超时
        )
        
        # 5. 下载视频并上传到 S3
        s3_url = await download_and_upload_to_s3(
            video_url,
            f"videos/{job.projectId}/output.mp4"
        )
        
        # 6. 生成缩略图
        thumbnail_url = await generate_thumbnail(s3_url)
        
        # 7. 更新数据库
        await update_project(
            project_id=job.project_id,
            status='completed',
            video_url=s3_url,
            thumbnail_url=thumbnail_url,
            completed_at=datetime.utcnow()
        )
        
        # 8. 发送通知
        await send_completion_notification(job.userId, job.projectId)
        
    except Exception as e:
        # 9. 错误处理
        await update_project(
            project_id=job.project_id,
            status='failed',
            error_message=str(e)
        )
        
        # 10. 退还额度
        await refund_credit(job.userId, job.projectId)
        
        raise  # 重新抛出以触发重试
```

---

## 5. API 设计

### 5.1 RESTful API 端点

#### 认证 API
```
POST   /api/auth/register          # 注册
POST   /api/auth/login             # 登录
POST   /api/auth/logout            # 登出
POST   /api/auth/refresh           # 刷新 Token
POST   /api/auth/reset-password    # 重置密码
GET    /api/auth/google            # Google OAuth
GET    /api/auth/callback          # OAuth 回调
```

#### 用户 API
```
GET    /api/user/profile           # 获取用户信息
PATCH  /api/user/profile           # 更新用户信息
GET    /api/user/credits           # 获取额度信息
```

#### 项目 API
```
GET    /api/projects               # 获取项目列表
POST   /api/projects               # 创建项目
GET    /api/projects/:id           # 获取项目详情
PATCH  /api/projects/:id           # 更新项目
DELETE /api/projects/:id           # 删除项目
POST   /api/projects/:id/generate  # 开始生成视频
GET    /api/projects/:id/status    # 获取生成状态
```

#### Avatar API
```
GET    /api/avatars                # 获取 Avatar 列表
GET    /api/avatars/:id            # 获取 Avatar 详情
```

#### 脚本 API
```
POST   /api/scripts/generate       # AI 生成脚本
GET    /api/scripts/templates      # 获取脚本模板
```

#### 支付 API
```
POST   /api/payments/checkout      # 创建支付会话
POST   /api/payments/webhook       # Stripe Webhook
GET    /api/payments/subscription  # 获取订阅状态
```

### 5.2 API 响应格式

```typescript
// 成功响应
interface SuccessResponse<T> {
  success: true;
  data: T;
  meta?: {
    page?: number;
    limit?: number;
    total?: number;
  };
}

// 错误响应
interface ErrorResponse {
  success: false;
  error: {
    code: string;
    message: string;
    details?: Record<string, string>;
  };
}

// 分页响应
interface PaginatedResponse<T> {
  success: true;
  data: T[];
  meta: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}
```

### 5.3 Server Actions (Next.js)

```typescript
// app/actions/projects.ts
'use server'

import { revalidatePath } from 'next/cache'
import { auth } from '@/lib/auth'
import { createProject, generateVideo } from '@/lib/projects'

export async function createProjectAction(data: CreateProjectInput) {
  const session = await auth()
  if (!session) throw new Error('Unauthorized')
  
  const project = await createProject(session.user.id, data)
  revalidatePath('/dashboard')
  
  return { success: true, data: project }
}

export async function generateVideoAction(projectId: string) {
  const session = await auth()
  if (!session) throw new Error('Unauthorized')
  
  const result = await generateVideo(projectId, session.user.id)
  
  return { success: true, data: result }
}
```

---

## 6. 部署架构

### 6.1 生产环境部署

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              生产环境部署架构                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                          Vercel (Frontend)                           │   │
│  │                                                                     │   │
│  │   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐              │   │
│  │   │ Edge        │   │ Serverless  │   │ Static      │              │   │
│  │   │ Functions   │   │ Functions   │   │ Assets      │              │   │
│  │   └─────────────┘   └─────────────┘   └─────────────┘              │   │
│  │                                                                     │   │
│  │   Regions: Washington DC, Singapore, Frankfurt                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                        │                                    │
│                                        ▼                                    │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      Supabase (Database + Auth)                      │   │
│  │                                                                     │   │
│  │   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐              │   │
│  │   │ PostgreSQL  │   │    Auth     │   │   Storage   │              │   │
│  │   │   (Primary) │   │   Service   │   │   (S3-like) │              │   │
│  │   └─────────────┘   └─────────────┘   └─────────────┘              │   │
│  │                                                                     │   │
│  │   Region: AWS us-east-1                                             │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                        │                                    │
│                                        ▼                                    │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    AWS / DigitalOcean (Workers)                      │   │
│  │                                                                     │   │
│  │   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐              │   │
│  │   │   Redis     │   │   BullMQ    │   │   Python    │              │   │
│  │   │   (Upstash) │   │   Workers   │   │   Workers   │              │   │
│  │   └─────────────┘   └─────────────┘   └─────────────┘              │   │
│  │                                                                     │   │
│  │   Auto-scaling: 2-10 worker instances                               │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                        │                                    │
│                                        ▼                                    │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                          External Services                           │   │
│  │                                                                     │   │
│  │   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐              │   │
│  │   │  Cloudflare │   │    Stripe   │   │   Resend    │              │   │
│  │   │  (CDN/R2)   │   │  (Payments) │   │   (Email)   │              │   │
│  │   └─────────────┘   └─────────────┘   └─────────────┘              │   │
│  │                                                                     │   │
│  │   ┌─────────────┐   ┌─────────────┐                                 │   │
│  │   │   HeyGen    │   │   OpenAI    │                                 │   │
│  │   │   (Video)   │   │   (GPT-4)   │                                 │   │
│  │   └─────────────┘   └─────────────┘                                 │   │
│  │                                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.2 环境配置

```yaml
# 环境变量配置
environments:
  development:
    NEXT_PUBLIC_API_URL: "http://localhost:3000"
    DATABASE_URL: "postgresql://localhost:5432/genvid_dev"
    REDIS_URL: "redis://localhost:6379"
    
  staging:
    NEXT_PUBLIC_API_URL: "https://staging.genvid.app"
    DATABASE_URL: "${{ secrets.STAGING_DATABASE_URL }}"
    REDIS_URL: "${{ secrets.STAGING_REDIS_URL }}"
    
  production:
    NEXT_PUBLIC_API_URL: "https://genvid.app"
    DATABASE_URL: "${{ secrets.PROD_DATABASE_URL }}"
    REDIS_URL: "${{ secrets.PROD_REDIS_URL }}"
```

---

## 7. 安全架构

### 7.1 安全措施

| 层级 | 安全措施 |
|------|---------|
| **传输层** | HTTPS 强制，TLS 1.3 |
| **认证层** | JWT Token + Refresh Token，OAuth 2.0 |
| **授权层** | RLS (Row Level Security)，RBAC |
| **数据层** | 数据库加密，敏感数据脱敏 |
| **API 层** | Rate Limiting，请求签名验证 |
| **应用层** | CSP，XSS 防护，CSRF Token |

### 7.2 Rate Limiting 配置

```typescript
// Rate Limiting 配置
const rateLimits = {
  // 公开端点
  'POST:/api/auth/login': { limit: 5, window: '15m' },
  'POST:/api/auth/register': { limit: 3, window: '1h' },
  
  // 认证端点
  'POST:/api/projects': { limit: 30, window: '1h' },
  'POST:/api/projects/:id/generate': { limit: 10, window: '1h' },
  
  // API Key 端点
  'API:/api/v1/*': { limit: 100, window: '1h' },
}
```

---

## 8. 监控与日志

### 8.1 监控指标

| 指标类型 | 指标名称 | 告警阈值 |
|---------|---------|---------|
| **可用性** | Uptime | < 99.9% |
| **性能** | API 响应时间 | > 500ms (P95) |
| **性能** | 视频生成时间 | > 10 分钟 |
| **错误率** | 5xx 错误率 | > 1% |
| **业务** | 视频生成失败率 | > 5% |
| **业务** | 队列积压 | > 100 任务 |

### 8.2 日志收集

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│ Application │────▶│   Log       │────▶│  Dashboard  │
│    Logs     │     │ Aggregator  │     │  (Grafana)  │
└─────────────┘     └─────────────┘     └─────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │   Alerting  │
                    │  (PagerDuty)│
                    └─────────────┘
```

---

## 9. 成本估算

### 9.1 基础设施成本 (月度)

| 服务 | 配置 | 预估成本 |
|------|------|---------|
| Vercel Pro | Team | $20/月 |
| Supabase Pro | Medium | $25/月 |
| Upstash Redis | Pay-as-you-go | $10/月 |
| Cloudflare Pro | - | $20/月 |
| AWS S3 | 100GB | $3/月 |
| **基础设施小计** | | **~$78/月** |

### 9.2 API 服务成本 (按量)

| 服务 | 使用量 | 预估成本/月 |
|------|--------|------------|
| HeyGen API | 1000 视频 @ 30秒 | $3,000-5,000 |
| OpenAI GPT-4 | 10K 脚本 | $50-100 |
| Resend | 10K 邮件 | $10 |
| Stripe | $5K GMV | $145 |
| **API 服务小计** | | **~$3,200-5,300/月** |

### 9.3 单视频成本分析

| 项目 | 成本 |
|------|------|
| HeyGen API | $1.50-3.00 |
| 计算/存储 | $0.05 |
| **总成本** | **$1.55-3.05/视频** |

---

*文档维护: 架构团队*  
*最后更新: 2026-02-19*
