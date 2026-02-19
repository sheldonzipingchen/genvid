-- Enable Row Level Security
ALTER TABLE profiles ENABLE ROW LEVEL SECURITY;
ALTER TABLE projects ENABLE ROW LEVEL SECURITY;
ALTER TABLE assets ENABLE ROW LEVEL SECURITY;
ALTER TABLE subscriptions ENABLE ROW LEVEL SECURITY;
ALTER TABLE usage_logs ENABLE ROW LEVEL SECURITY;
ALTER TABLE api_keys ENABLE ROW LEVEL SECURITY;
ALTER TABLE avatars ENABLE ROW LEVEL SECURITY;
ALTER TABLE script_templates ENABLE ROW LEVEL SECURITY;

-- Profiles: Users can only access their own profile
CREATE POLICY "Users can view own profile"
    ON profiles FOR SELECT
    USING (auth.uid()::text = id::text);

CREATE POLICY "Users can update own profile"
    ON profiles FOR UPDATE
    USING (auth.uid()::text = id::text);

CREATE POLICY "Users can insert own profile"
    ON profiles FOR INSERT
    WITH CHECK (auth.uid()::text = id::text);

-- Projects: Users can only access their own projects
CREATE POLICY "Users can view own projects"
    ON projects FOR SELECT
    USING (auth.uid()::text = user_id::text);

CREATE POLICY "Users can insert own projects"
    ON projects FOR INSERT
    WITH CHECK (auth.uid()::text = user_id::text);

CREATE POLICY "Users can update own projects"
    ON projects FOR UPDATE
    USING (auth.uid()::text = user_id::text);

CREATE POLICY "Users can delete own projects"
    ON projects FOR DELETE
    USING (auth.uid()::text = user_id::text);

-- Assets: Users can only access assets from their projects
CREATE POLICY "Users can view own assets"
    ON assets FOR SELECT
    USING (auth.uid()::text = user_id::text);

CREATE POLICY "Users can insert own assets"
    ON assets FOR INSERT
    WITH CHECK (auth.uid()::text = user_id::text);

CREATE POLICY "Users can delete own assets"
    ON assets FOR DELETE
    USING (auth.uid()::text = user_id::text);

-- Subscriptions: Users can only access their own subscriptions
CREATE POLICY "Users can view own subscriptions"
    ON subscriptions FOR SELECT
    USING (auth.uid()::text = user_id::text);

-- Avatars: All authenticated users can view active avatars
CREATE POLICY "Authenticated users can view active avatars"
    ON avatars FOR SELECT
    USING (is_active = true);

-- Script templates: Users can view public templates and their own
CREATE POLICY "Users can view available templates"
    ON script_templates FOR SELECT
    USING (
        (is_active = true AND is_user_created = false)
        OR (created_by::text = auth.uid()::text)
    );

-- API keys: Users can only access their own keys
CREATE POLICY "Users can view own api keys"
    ON api_keys FOR SELECT
    USING (auth.uid()::text = user_id::text);

CREATE POLICY "Users can insert own api keys"
    ON api_keys FOR INSERT
    WITH CHECK (auth.uid()::text = user_id::text);

CREATE POLICY "Users can delete own api keys"
    ON api_keys FOR DELETE
    USING (auth.uid()::text = user_id::text);

-- Usage logs: Users can view their own logs
CREATE POLICY "Users can view own usage logs"
    ON usage_logs FOR SELECT
    USING (auth.uid()::text = user_id::text);
