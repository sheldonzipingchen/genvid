# Genvid - 数据模型设计 (Data Model)

**版本**: 1.0  
**日期**: 2026-02-19  
**数据库**: PostgreSQL 16 + Supabase

---

## 1. 数据库架构概览

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           数据库架构图                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌──────────────┐       ┌──────────────┐       ┌──────────────┐           │
│  │   profiles   │───┬───│   projects   │───┬───│    assets    │           │
│  │   (用户表)    │   │   │   (项目表)    │   │   │   (素材表)    │           │
│  └──────────────┘   │   └──────────────┘   │   └──────────────┘           │
│         │           │          │           │                              │
│         │           │          │           │                              │
│         ▼           │          ▼           │                              │
│  ┌──────────────┐   │   ┌──────────────┐   │                              │
│  │subscriptions │   │   │    avatars   │◄──┘                              │
│  │  (订阅表)     │   │   │  (Avatar库)  │                                  │
│  └──────────────┘   │   └──────────────┘                                  │
│                     │                                                      │
│                     │   ┌──────────────┐                                   │
│                     └──▶│script_templates│                                │
│                         │ (脚本模板表)   │                                  │
│                         └──────────────┘                                   │
│                                                                             │
│  ┌──────────────┐       ┌──────────────┐                                   │
│  │   api_keys   │       │usage_logs    │                                   │
│  │  (API密钥)   │       │ (使用日志)    │                                   │
│  └──────────────┘       └──────────────┘                                   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. 表结构定义

### 2.1 profiles (用户表)

存储用户基本信息和账户设置。

```sql
CREATE TABLE profiles (
  -- 主键
  id UUID PRIMARY KEY REFERENCES auth.users(id) ON DELETE CASCADE,
  
  -- 基本信息
  email VARCHAR(255) NOT NULL UNIQUE,
  full_name VARCHAR(100),
  avatar_url TEXT,
  company_name VARCHAR(200),
  
  -- 额度管理
  credits_remaining INTEGER DEFAULT 3 CHECK (credits_remaining >= 0),
  credits_used_total INTEGER DEFAULT 0 CHECK (credits_used_total >= 0),
  
  -- 订阅状态
  subscription_tier VARCHAR(20) DEFAULT 'free' 
    CHECK (subscription_tier IN ('free', 'starter', 'pro', 'business', 'enterprise')),
  subscription_status VARCHAR(20) DEFAULT 'inactive'
    CHECK (subscription_status IN ('active', 'inactive', 'canceled', 'past_due')),
  
  -- 品牌设置 (Phase 2)
  brand_logo_url TEXT,
  brand_primary_color VARCHAR(7),  -- Hex color
  brand_font_family VARCHAR(100),
  
  -- 偏好设置
  preferred_language VARCHAR(10) DEFAULT 'en',
  email_notifications BOOLEAN DEFAULT true,
  
  -- 时间戳
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  last_login_at TIMESTAMPTZ,
  
  -- 索引优化
  CONSTRAINT valid_email CHECK (email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$')
);

-- 索引
CREATE INDEX idx_profiles_email ON profiles(email);
CREATE INDEX idx_profiles_subscription_tier ON profiles(subscription_tier);
CREATE INDEX idx_profiles_created_at ON profiles(created_at DESC);

-- 触发器: 自动更新 updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_profiles_updated_at
  BEFORE UPDATE ON profiles
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

-- 注释
COMMENT ON TABLE profiles IS '用户账户信息表';
COMMENT ON COLUMN profiles.credits_remaining IS '当前剩余视频额度';
COMMENT ON COLUMN profiles.subscription_tier IS '订阅等级: free/starter/pro/business/enterprise';
```

### 2.2 projects (视频项目表)

存储视频项目的完整信息，包括产品信息和生成状态。

```sql
CREATE TYPE project_status AS ENUM (
  'draft',        -- 草稿
  'queued',       -- 排队中
  'processing',   -- 生成中
  'completed',    -- 已完成
  'failed',       -- 失败
  'canceled'      -- 已取消
);

CREATE TYPE video_format AS ENUM ('9:16', '1:1', '16:9');

CREATE TABLE projects (
  -- 主键
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  
  -- 外键
  user_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
  avatar_id UUID REFERENCES avatars(id) ON DELETE SET NULL,
  script_template_id UUID REFERENCES script_templates(id) ON DELETE SET NULL,
  
  -- 项目基本信息
  title VARCHAR(200),  -- 自动生成或用户指定
  
  -- 产品信息
  product_name VARCHAR(200),
  product_description TEXT,
  product_url TEXT,
  product_price DECIMAL(10, 2),
  product_currency VARCHAR(3) DEFAULT 'USD',
  
  -- 视频配置
  script TEXT,                    -- 最终使用的脚本
  language VARCHAR(10) DEFAULT 'en',
  format video_format DEFAULT '9:16',
  video_duration_seconds INTEGER, -- 预估或实际时长
  
  -- 生成状态
  status project_status DEFAULT 'draft',
  progress_percent INTEGER DEFAULT 0 CHECK (progress_percent >= 0 AND progress_percent <= 100),
  error_message TEXT,             -- 失败时的错误信息
  
  -- 外部服务信息
  external_task_id VARCHAR(100),  -- 视频生成 API 的任务 ID
  external_provider VARCHAR(50),  -- 使用的视频生成服务
  
  -- 输出结果
  video_url TEXT,                 -- 生成的视频 URL
  thumbnail_url TEXT,             -- 缩略图 URL
  file_size_bytes BIGINT,
  
  -- 元数据
  is_favorite BOOLEAN DEFAULT false,
  tags TEXT[],                    -- 用户标签
  view_count INTEGER DEFAULT 0,   -- 预览次数
  
  -- 时间戳
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  started_at TIMESTAMPTZ,         -- 开始生成时间
  completed_at TIMESTAMPTZ,       -- 完成时间
  expires_at TIMESTAMPTZ,         -- 视频过期时间 (可选)
  
  -- 约束
  CONSTRAINT valid_product CHECK (
    product_name IS NOT NULL OR product_url IS NOT NULL
  )
);

-- 索引
CREATE INDEX idx_projects_user_id ON projects(user_id);
CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_projects_created_at ON projects(created_at DESC);
CREATE INDEX idx_projects_user_status ON projects(user_id, status);
CREATE INDEX idx_projects_external_task ON projects(external_task_id) 
  WHERE external_task_id IS NOT NULL;

-- 全文搜索索引
CREATE INDEX idx_projects_search ON projects 
  USING GIN(to_tsvector('english', coalesce(title, '') || ' ' || coalesce(product_name, '') || ' ' || coalesce(product_description, '')));

-- 触发器
CREATE TRIGGER update_projects_updated_at
  BEFORE UPDATE ON projects
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

-- 注释
COMMENT ON TABLE projects IS '视频项目表，存储所有视频生成任务';
COMMENT ON COLUMN projects.external_task_id IS '视频生成服务的任务标识';
COMMENT ON COLUMN projects.external_provider IS '使用的视频生成服务商: heygen/kling/runwayml';
```

### 2.3 assets (素材表)

存储项目相关的图片、视频、音频等素材。

```sql
CREATE TYPE asset_type AS ENUM ('image', 'video', 'audio', 'document');
CREATE TYPE asset_purpose AS ENUM (
  'product_image',    -- 产品图片
  'product_video',    -- 产品视频
  'generated_video',  -- 生成的视频
  'thumbnail',        -- 缩略图
  'background_music', -- 背景音乐
  'voiceover',        -- 配音
  'other'             -- 其他
);

CREATE TABLE assets (
  -- 主键
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  
  -- 外键
  project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
  
  -- 素材信息
  type asset_type NOT NULL,
  purpose asset_purpose NOT NULL,
  
  -- 文件信息
  filename VARCHAR(255) NOT NULL,
  original_filename VARCHAR(255),
  url TEXT NOT NULL,
  thumbnail_url TEXT,
  
  -- 文件属性
  file_size_bytes BIGINT,
  mime_type VARCHAR(100),
  width INTEGER,      -- 图片/视频宽度
  height INTEGER,     -- 图片/视频高度
  duration_seconds DECIMAL(10, 2),  -- 音频/视频时长
  
  -- 元数据
  alt_text VARCHAR(500),
  is_primary BOOLEAN DEFAULT false,  -- 主要素材标识
  
  -- 时间戳
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_assets_project_id ON assets(project_id);
CREATE INDEX idx_assets_user_id ON assets(user_id);
CREATE INDEX idx_assets_type ON assets(type);
CREATE INDEX idx_assets_purpose ON assets(purpose);

-- 触发器
CREATE TRIGGER update_assets_updated_at
  BEFORE UPDATE ON assets
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

-- 注释
COMMENT ON TABLE assets IS '素材表，存储项目相关的图片、视频、音频等';
```

### 2.4 avatars (Avatar 库表)

存储可供用户选择的 AI Avatar 信息。

```sql
CREATE TYPE avatar_gender AS ENUM ('male', 'female', 'other');
CREATE TYPE avatar_style AS ENUM (
  'casual',       -- 休闲
  'professional', -- 专业
  'energetic',    -- 活力
  'friendly',     -- 亲切
  'elegant',      -- 优雅
  'trendy'        -- 时尚
);

CREATE TABLE avatars (
  -- 主键
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  
  -- 基本信息
  name VARCHAR(100) NOT NULL,
  display_name VARCHAR(100),
  description TEXT,
  
  -- 特征属性
  gender avatar_gender,
  age_range VARCHAR(20),          -- '20s', '30s', '40s', '50s+'
  ethnicity VARCHAR(50),
  
  -- 风格设置
  style avatar_style DEFAULT 'casual',
  language VARCHAR(10)[] DEFAULT ARRAY['en'],
  
  -- 媒体资源
  preview_video_url TEXT,
  thumbnail_url TEXT,
  
  -- 状态
  is_premium BOOLEAN DEFAULT false,
  is_active BOOLEAN DEFAULT true,
  is_custom BOOLEAN DEFAULT false,  -- 用户自定义 Avatar
  
  -- 使用统计
  usage_count INTEGER DEFAULT 0,
  
  -- 排序和分组
  sort_order INTEGER DEFAULT 0,
  category VARCHAR(50),            -- 分类标签
  
  -- 时间戳
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_avatars_gender ON avatars(gender);
CREATE INDEX idx_avatars_style ON avatars(style);
CREATE INDEX idx_avatars_active ON avatars(is_active) WHERE is_active = true;
CREATE INDEX idx_avatars_premium ON avatars(is_premium);
CREATE INDEX idx_avatars_usage ON avatars(usage_count DESC);

-- 触发器
CREATE TRIGGER update_avatars_updated_at
  BEFORE UPDATE ON avatars
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

-- 注释
COMMENT ON TABLE avatars IS 'AI Avatar 库，存储所有可用虚拟人物';
```

### 2.5 script_templates (脚本模板表)

存储视频脚本模板和 AI 生成的脚本。

```sql
CREATE TYPE script_category AS ENUM (
  'product_review',   -- 产品评测
  'unboxing',         -- 开箱体验
  'tutorial',         -- 使用教程
  'comparison',       -- 产品对比
  'testimonial',      -- 用户证言
  'before_after',     -- 前后对比
  'storytelling',     -- 故事叙述
  'custom'            -- 自定义
);

CREATE TABLE script_templates (
  -- 主键
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  
  -- 基本信息
  name VARCHAR(100) NOT NULL,
  description TEXT,
  
  -- 分类
  category script_category NOT NULL,
  language VARCHAR(10) DEFAULT 'en',
  
  -- 模板内容
  template_text TEXT NOT NULL,
  placeholders TEXT[],            -- 可替换变量列表
  
  -- 使用提示
  tips TEXT,
  best_for TEXT,                  -- 最适合的产品类型
  
  -- 状态
  is_premium BOOLEAN DEFAULT false,
  is_active BOOLEAN DEFAULT true,
  is_user_created BOOLEAN DEFAULT false,
  
  -- 关联用户 (用户自定义模板)
  created_by UUID REFERENCES profiles(id) ON DELETE SET NULL,
  
  -- 使用统计
  usage_count INTEGER DEFAULT 0,
  success_rate DECIMAL(5, 2),     -- 使用该模板的视频成功率
  
  -- 时间戳
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_script_templates_category ON script_templates(category);
CREATE INDEX idx_script_templates_language ON script_templates(language);
CREATE INDEX idx_script_templates_active ON script_templates(is_active) WHERE is_active = true;

-- 触发器
CREATE TRIGGER update_script_templates_updated_at
  BEFORE UPDATE ON script_templates
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

-- 注释
COMMENT ON TABLE script_templates IS '视频脚本模板库';
```

### 2.6 subscriptions (订阅表)

存储用户订阅和支付信息。

```sql
CREATE TYPE subscription_status AS ENUM (
  'active',      -- 活跃
  'inactive',    -- 未激活
  'canceled',    -- 已取消
  'past_due',    -- 逾期
  'trialing'     -- 试用中
);

CREATE TYPE billing_period AS ENUM ('monthly', 'yearly');

CREATE TABLE subscriptions (
  -- 主键
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  
  -- 外键
  user_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
  
  -- 订阅计划
  plan_id VARCHAR(50) NOT NULL,   -- 'starter_monthly', 'pro_yearly', etc.
  plan_name VARCHAR(100) NOT NULL,
  
  -- 状态
  status subscription_status DEFAULT 'inactive',
  
  -- 额度设置
  credits_per_period INTEGER NOT NULL,
  credits_used_this_period INTEGER DEFAULT 0,
  credits_rollover INTEGER DEFAULT 0,  -- 可累积额度
  
  -- 计费周期
  billing_period billing_period DEFAULT 'monthly',
  current_period_start TIMESTAMPTZ,
  current_period_end TIMESTAMPTZ,
  
  -- 价格
  amount_cents INTEGER NOT NULL,
  currency VARCHAR(3) DEFAULT 'USD',
  
  -- Stripe 集成
  stripe_subscription_id VARCHAR(100),
  stripe_customer_id VARCHAR(100),
  stripe_price_id VARCHAR(100),
  
  -- 取消信息
  canceled_at TIMESTAMPTZ,
  cancel_at_period_end BOOLEAN DEFAULT false,
  cancellation_reason TEXT,
  
  -- 时间戳
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_status ON subscriptions(status);
CREATE INDEX idx_subscriptions_stripe_customer ON subscriptions(stripe_customer_id);
CREATE INDEX idx_subscriptions_period_end ON subscriptions(current_period_end) 
  WHERE status = 'active';

-- 触发器
CREATE TRIGGER update_subscriptions_updated_at
  BEFORE UPDATE ON subscriptions
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

-- 注释
COMMENT ON TABLE subscriptions IS '用户订阅和计费信息';
```

### 2.7 usage_logs (使用日志表)

记录用户行为和系统事件，用于分析和审计。

```sql
CREATE TYPE log_event AS ENUM (
  'video_created',      -- 视频创建
  'video_completed',    -- 视频完成
  'video_failed',       -- 视频失败
  'video_downloaded',   -- 视频下载
  'credit_used',        -- 额度消耗
  'credit_refunded',    -- 额度退还
  'subscription_started', -- 订阅开始
  'subscription_canceled', -- 订阅取消
  'user_login',         -- 用户登录
  'api_request'         -- API 请求
);

CREATE TABLE usage_logs (
  -- 主键
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  
  -- 外键
  user_id UUID REFERENCES profiles(id) ON DELETE SET NULL,
  project_id UUID REFERENCES projects(id) ON DELETE SET NULL,
  
  -- 事件信息
  event log_event NOT NULL,
  description TEXT,
  
  -- 关联数据
  metadata JSONB DEFAULT '{}',    -- 额外的上下文数据
  
  -- IP 和设备信息
  ip_address INET,
  user_agent TEXT,
  
  -- 时间戳
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_usage_logs_user_id ON usage_logs(user_id);
CREATE INDEX idx_usage_logs_event ON usage_logs(event);
CREATE INDEX idx_usage_logs_created_at ON usage_logs(created_at DESC);
CREATE INDEX idx_usage_logs_project_id ON usage_logs(project_id);

-- 分区 (按月)
-- CREATE TABLE usage_logs_2026_02 PARTITION OF usage_logs
--   FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');

-- 注释
COMMENT ON TABLE usage_logs IS '使用日志表，记录用户行为和系统事件';
```

### 2.8 api_keys (API 密钥表)

存储用户 API 密钥，用于 API 访问。

```sql
CREATE TABLE api_keys (
  -- 主键
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  
  -- 外键
  user_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
  
  -- 密钥信息
  name VARCHAR(100) NOT NULL,
  key_hash VARCHAR(255) NOT NULL,  -- 哈希后的密钥
  key_prefix VARCHAR(8),           -- 密钥前缀 (用于识别)
  
  -- 权限
  scopes TEXT[] DEFAULT ARRAY['read', 'write'],
  
  -- 限制
  rate_limit INTEGER DEFAULT 100,  -- 每小时请求数限制
  
  -- 状态
  is_active BOOLEAN DEFAULT true,
  last_used_at TIMESTAMPTZ,
  expires_at TIMESTAMPTZ,
  
  -- 时间戳
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX idx_api_keys_prefix ON api_keys(key_prefix);

-- 注释
COMMENT ON TABLE api_keys IS '用户 API 密钥表';
```

---

## 3. RLS (Row Level Security) 策略

### 3.1 profiles 表

```sql
ALTER TABLE profiles ENABLE ROW LEVEL SECURITY;

-- 用户只能查看自己的 profile
CREATE POLICY "Users can view own profile"
  ON profiles FOR SELECT
  USING (auth.uid() = id);

-- 用户只能更新自己的 profile
CREATE POLICY "Users can update own profile"
  ON profiles FOR UPDATE
  USING (auth.uid() = id);

-- 用户可以插入自己的 profile (注册时)
CREATE POLICY "Users can insert own profile"
  ON profiles FOR INSERT
  WITH CHECK (auth.uid() = id);
```

### 3.2 projects 表

```sql
ALTER TABLE projects ENABLE ROW LEVEL SECURITY;

-- 用户只能访问自己的项目
CREATE POLICY "Users can view own projects"
  ON projects FOR SELECT
  USING (auth.uid() = user_id);

CREATE POLICY "Users can insert own projects"
  ON projects FOR INSERT
  WITH CHECK (auth.uid() = user_id);

CREATE POLICY "Users can update own projects"
  ON projects FOR UPDATE
  USING (auth.uid() = user_id);

CREATE POLICY "Users can delete own projects"
  ON projects FOR DELETE
  USING (auth.uid() = user_id);
```

### 3.3 avatars 表

```sql
ALTER TABLE avatars ENABLE ROW LEVEL SECURITY;

-- 所有认证用户可以查看活跃的 Avatar
CREATE POLICY "Authenticated users can view active avatars"
  ON avatars FOR SELECT
  USING (auth.role() = 'authenticated' AND is_active = true);
```

### 3.4 script_templates 表

```sql
ALTER TABLE script_templates ENABLE ROW LEVEL SECURITY;

-- 用户可以查看公共模板和自己的模板
CREATE POLICY "Users can view available templates"
  ON script_templates FOR SELECT
  USING (
    (is_active = true AND is_user_created = false)
    OR created_by = auth.uid()
  );
```

---

## 4. 数据库视图

### 4.1 用户项目统计视图

```sql
CREATE VIEW user_project_stats AS
SELECT 
  p.id AS user_id,
  p.email,
  p.credits_remaining,
  p.subscription_tier,
  COUNT(pr.id) AS total_projects,
  COUNT(CASE WHEN pr.status = 'completed' THEN 1 END) AS completed_projects,
  COUNT(CASE WHEN pr.status = 'failed' THEN 1 END) AS failed_projects,
  MAX(pr.created_at) AS last_project_at
FROM profiles p
LEFT JOIN projects pr ON p.id = pr.user_id
GROUP BY p.id, p.email, p.credits_remaining, p.subscription_tier;

COMMENT ON VIEW user_project_stats IS '用户项目统计视图';
```

### 4.2 热门 Avatar 视图

```sql
CREATE VIEW popular_avatars AS
SELECT 
  a.id,
  a.name,
  a.thumbnail_url,
  a.gender,
  a.style,
  COUNT(p.id) AS usage_count,
  COUNT(CASE WHEN p.status = 'completed' THEN 1 END) AS success_count
FROM avatars a
LEFT JOIN projects p ON a.id = p.avatar_id
WHERE a.is_active = true
GROUP BY a.id, a.name, a.thumbnail_url, a.gender, a.style
ORDER BY usage_count DESC;

COMMENT ON VIEW popular_avatars IS '热门 Avatar 统计视图';
```

---

## 5. 数据库函数

### 5.1 消耗额度函数

```sql
CREATE OR REPLACE FUNCTION use_credit(
  p_user_id UUID,
  p_project_id UUID
) RETURNS BOOLEAN AS $$
DECLARE
  v_current_credits INTEGER;
BEGIN
  -- 获取当前额度
  SELECT credits_remaining INTO v_current_credits
  FROM profiles WHERE id = p_user_id;
  
  -- 检查额度是否足够
  IF v_current_credits <= 0 THEN
    RETURN FALSE;
  END IF;
  
  -- 扣减额度
  UPDATE profiles 
  SET 
    credits_remaining = credits_remaining - 1,
    credits_used_total = credits_used_total + 1
  WHERE id = p_user_id;
  
  -- 记录日志
  INSERT INTO usage_logs (user_id, project_id, event, description)
  VALUES (p_user_id, p_project_id, 'credit_used', '1 credit used for video generation');
  
  RETURN TRUE;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

COMMENT ON FUNCTION use_credit IS '消耗用户视频额度';
```

### 5.2 退还额度函数

```sql
CREATE OR REPLACE FUNCTION refund_credit(
  p_user_id UUID,
  p_project_id UUID,
  p_reason TEXT
) RETURNS VOID AS $$
BEGIN
  -- 增加额度
  UPDATE profiles 
  SET credits_remaining = credits_remaining + 1
  WHERE id = p_user_id;
  
  -- 记录日志
  INSERT INTO usage_logs (user_id, project_id, event, description)
  VALUES (p_user_id, p_project_id, 'credit_refunded', p_reason);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

COMMENT ON FUNCTION refund_credit IS '退还用户视频额度';
```

### 5.3 检查用户额度函数

```sql
CREATE OR REPLACE FUNCTION check_user_credits(
  p_user_id UUID
) RETURNS TABLE(
  has_credits BOOLEAN,
  credits_remaining INTEGER,
  subscription_tier VARCHAR,
  can_generate BOOLEAN
) AS $$
BEGIN
  RETURN QUERY
  SELECT 
    credits_remaining > 0 AS has_credits,
    credits_remaining,
    subscription_tier,
    credits_remaining > 0 OR subscription_tier = 'enterprise' AS can_generate
  FROM profiles
  WHERE id = p_user_id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

COMMENT ON FUNCTION check_user_credits IS '检查用户视频额度状态';
```

---

## 6. 初始数据 (Seed Data)

### 6.1 Avatar 初始数据

```sql
-- 示例 Avatar 数据
INSERT INTO avatars (name, display_name, gender, age_range, style, language, thumbnail_url, is_premium) VALUES
('emma_casual', 'Emma', 'female', '20s', 'casual', ARRAY['en', 'es'], '/avatars/emma_thumb.jpg', false),
('james_pro', 'James', 'male', '30s', 'professional', ARRAY['en'], '/avatars/james_thumb.jpg', false),
('sofia_energy', 'Sofia', 'female', '20s', 'energetic', ARRAY['en', 'pt'], '/avatars/sofia_thumb.jpg', true),
('li_friendly', 'Li', 'male', '30s', 'friendly', ARRAY['en', 'zh'], '/avatars/li_thumb.jpg', false),
('maria_elegant', 'Maria', 'female', '40s', 'elegant', ARRAY['en', 'es', 'pt'], '/avatars/maria_thumb.jpg', true);
```

### 6.2 脚本模板初始数据

```sql
-- 示例脚本模板
INSERT INTO script_templates (name, category, template_text, language, is_premium) VALUES
('Product Review', 'product_review', 
'I''ve been using {product_name} for {time_period} now, and I have to say... {review_content}. 
If you''re looking for {benefit}, this is definitely worth checking out!',
'en', false),

('Unboxing Experience', 'unboxing',
'Hey everyone! I just got my {product_name} delivered today. Let''s unbox this together! 
First impressions... {first_impression}. Stay tuned for my full review!',
'en', false),

('Before & After', 'before_after',
'So I''ve been using {product_name} for {time_period}. Here''s what it looked like before... 
And here''s after {time_period}. The difference is {result}!',
'en', true),

('Comparison', 'comparison',
'Today I''m comparing {product_name} with {competitor}. 
Let''s break it down: price, quality, and overall value. Here''s my honest take...',
'en', true);
```

---

## 7. 数据迁移脚本

### 7.1 创建迁移文件结构

```
/supabase/migrations/
├── 20260219000000_initial_schema.sql
├── 20260219000001_rls_policies.sql
├── 20260219000002_seed_data.sql
└── 20260219000003_functions.sql
```

---

*文档维护: 后端团队*  
*最后更新: 2026-02-19*
