-- Genvid Database Schema
-- Initial migration for Supabase PostgreSQL

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create ENUMs
CREATE TYPE project_status AS ENUM ('draft', 'queued', 'processing', 'completed', 'failed', 'canceled');
CREATE TYPE video_format AS ENUM ('9:16', '1:1', '16:9');
CREATE TYPE avatar_gender AS ENUM ('male', 'female', 'other');
CREATE TYPE avatar_style AS ENUM ('casual', 'professional', 'energetic', 'friendly', 'elegant', 'trendy');
CREATE TYPE script_category AS ENUM ('product_review', 'unboxing', 'tutorial', 'comparison', 'testimonial', 'before_after', 'storytelling', 'custom');
CREATE TYPE subscription_status AS ENUM ('active', 'inactive', 'canceled', 'past_due', 'trialing');
CREATE TYPE billing_period AS ENUM ('monthly', 'yearly');
CREATE TYPE asset_type AS ENUM ('image', 'video', 'audio', 'document');
CREATE TYPE asset_purpose AS ENUM ('product_image', 'product_video', 'generated_video', 'thumbnail', 'background_music', 'voiceover', 'other');
CREATE TYPE log_event AS ENUM ('video_created', 'video_completed', 'video_failed', 'video_downloaded', 'credit_used', 'credit_refunded', 'subscription_started', 'subscription_canceled', 'user_login', 'api_request');

-- Create update timestamp function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Profiles table
CREATE TABLE profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    full_name VARCHAR(100),
    avatar_url TEXT,
    company_name VARCHAR(200),
    
    credits_remaining INTEGER DEFAULT 3 CHECK (credits_remaining >= 0),
    credits_used_total INTEGER DEFAULT 0 CHECK (credits_used_total >= 0),
    
    subscription_tier VARCHAR(20) DEFAULT 'free' CHECK (subscription_tier IN ('free', 'starter', 'pro', 'business', 'enterprise')),
    subscription_status subscription_status DEFAULT 'inactive',
    
    brand_logo_url TEXT,
    brand_primary_color VARCHAR(7),
    brand_font_family VARCHAR(100),
    
    preferred_language VARCHAR(10) DEFAULT 'en',
    email_notifications BOOLEAN DEFAULT true,
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    last_login_at TIMESTAMPTZ
);

CREATE INDEX idx_profiles_email ON profiles(email);
CREATE INDEX idx_profiles_subscription_tier ON profiles(subscription_tier);
CREATE INDEX idx_profiles_created_at ON profiles(created_at DESC);

CREATE TRIGGER update_profiles_updated_at
    BEFORE UPDATE ON profiles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Avatars table
CREATE TABLE avatars (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    display_name VARCHAR(100),
    description TEXT,
    
    gender avatar_gender,
    age_range VARCHAR(20),
    ethnicity VARCHAR(50),
    
    style avatar_style DEFAULT 'casual',
    language VARCHAR(10)[] DEFAULT ARRAY['en'],
    
    preview_video_url TEXT,
    thumbnail_url TEXT,
    
    is_premium BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    is_custom BOOLEAN DEFAULT false,
    
    usage_count INTEGER DEFAULT 0,
    
    sort_order INTEGER DEFAULT 0,
    category VARCHAR(50),
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_avatars_gender ON avatars(gender);
CREATE INDEX idx_avatars_style ON avatars(style);
CREATE INDEX idx_avatars_active ON avatars(is_active) WHERE is_active = true;
CREATE INDEX idx_avatars_usage ON avatars(usage_count DESC);

CREATE TRIGGER update_avatars_updated_at
    BEFORE UPDATE ON avatars
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Script templates table
CREATE TABLE script_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    
    category script_category NOT NULL,
    language VARCHAR(10) DEFAULT 'en',
    
    template_text TEXT NOT NULL,
    placeholders TEXT[],
    
    tips TEXT,
    best_for TEXT,
    
    is_premium BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    is_user_created BOOLEAN DEFAULT false,
    
    created_by UUID REFERENCES profiles(id) ON DELETE SET NULL,
    
    usage_count INTEGER DEFAULT 0,
    success_rate DECIMAL(5, 2),
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_script_templates_category ON script_templates(category);
CREATE INDEX idx_script_templates_language ON script_templates(language);
CREATE INDEX idx_script_templates_active ON script_templates(is_active) WHERE is_active = true;

CREATE TRIGGER update_script_templates_updated_at
    BEFORE UPDATE ON script_templates
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Projects table
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    avatar_id UUID REFERENCES avatars(id) ON DELETE SET NULL,
    script_template_id UUID REFERENCES script_templates(id) ON DELETE SET NULL,
    
    title VARCHAR(200),
    
    product_name VARCHAR(200),
    product_description TEXT,
    product_url TEXT,
    product_price DECIMAL(10, 2),
    product_currency VARCHAR(3) DEFAULT 'USD',
    
    script TEXT,
    language VARCHAR(10) DEFAULT 'en',
    format video_format DEFAULT '9:16',
    video_duration_seconds INTEGER,
    
    status project_status DEFAULT 'draft',
    progress_percent INTEGER DEFAULT 0 CHECK (progress_percent >= 0 AND progress_percent <= 100),
    error_message TEXT,
    
    external_task_id VARCHAR(100),
    external_provider VARCHAR(50),
    
    video_url TEXT,
    thumbnail_url TEXT,
    file_size_bytes BIGINT,
    
    is_favorite BOOLEAN DEFAULT false,
    tags TEXT[],
    view_count INTEGER DEFAULT 0,
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ
);

CREATE INDEX idx_projects_user_id ON projects(user_id);
CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_projects_created_at ON projects(created_at DESC);
CREATE INDEX idx_projects_user_status ON projects(user_id, status);

CREATE TRIGGER update_projects_updated_at
    BEFORE UPDATE ON projects
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Assets table
CREATE TABLE assets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    
    type asset_type NOT NULL,
    purpose asset_purpose NOT NULL,
    
    filename VARCHAR(255) NOT NULL,
    original_filename VARCHAR(255),
    url TEXT NOT NULL,
    thumbnail_url TEXT,
    
    file_size_bytes BIGINT,
    mime_type VARCHAR(100),
    width INTEGER,
    height INTEGER,
    duration_seconds DECIMAL(10, 2),
    
    alt_text VARCHAR(500),
    is_primary BOOLEAN DEFAULT false,
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_assets_project_id ON assets(project_id);
CREATE INDEX idx_assets_user_id ON assets(user_id);

CREATE TRIGGER update_assets_updated_at
    BEFORE UPDATE ON assets
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Subscriptions table
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    
    plan_id VARCHAR(50) NOT NULL,
    plan_name VARCHAR(100) NOT NULL,
    
    status subscription_status DEFAULT 'inactive',
    
    credits_per_period INTEGER NOT NULL,
    credits_used_this_period INTEGER DEFAULT 0,
    credits_rollover INTEGER DEFAULT 0,
    
    billing_period billing_period DEFAULT 'monthly',
    current_period_start TIMESTAMPTZ,
    current_period_end TIMESTAMPTZ,
    
    amount_cents INTEGER NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    
    stripe_subscription_id VARCHAR(100),
    stripe_customer_id VARCHAR(100),
    stripe_price_id VARCHAR(100),
    
    canceled_at TIMESTAMPTZ,
    cancel_at_period_end BOOLEAN DEFAULT false,
    cancellation_reason TEXT,
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_status ON subscriptions(status);

CREATE TRIGGER update_subscriptions_updated_at
    BEFORE UPDATE ON subscriptions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Usage logs table
CREATE TABLE usage_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES profiles(id) ON DELETE SET NULL,
    project_id UUID REFERENCES projects(id) ON DELETE SET NULL,
    
    event log_event NOT NULL,
    description TEXT,
    
    metadata JSONB DEFAULT '{}',
    
    ip_address INET,
    user_agent TEXT,
    
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_usage_logs_user_id ON usage_logs(user_id);
CREATE INDEX idx_usage_logs_event ON usage_logs(event);
CREATE INDEX idx_usage_logs_created_at ON usage_logs(created_at DESC);

-- API keys table
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    
    name VARCHAR(100) NOT NULL,
    key_hash VARCHAR(255) NOT NULL,
    key_prefix VARCHAR(8),
    
    scopes TEXT[] DEFAULT ARRAY['read', 'write'],
    rate_limit INTEGER DEFAULT 100,
    
    is_active BOOLEAN DEFAULT true,
    last_used_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX idx_api_keys_prefix ON api_keys(key_prefix);
