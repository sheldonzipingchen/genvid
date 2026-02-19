export interface User {
  id: string
  email: string
  full_name: string | null
  avatar_url: string | null
  company_name: string | null
  credits_remaining: number
  credits_used_total: number
  subscription_tier: 'free' | 'starter' | 'pro' | 'business' | 'enterprise'
  subscription_status: 'active' | 'inactive' | 'canceled' | 'past_due' | 'trialing'
  preferred_language: string
  created_at: string
  updated_at: string
}

export interface Project {
  id: string
  user_id: string
  avatar_id: string | null
  title: string | null
  product_name: string | null
  product_description: string | null
  product_url: string | null
  product_image_url: string | null
  script: string | null
  language: string
  format: '9:16' | '1:1' | '16:9'
  status: 'draft' | 'queued' | 'processing' | 'completed' | 'failed' | 'canceled'
  progress_percent: number
  error_message: string | null
  video_url: string | null
  thumbnail_url: string | null
  created_at: string
  updated_at: string
  completed_at: string | null
}

export interface Avatar {
  id: string
  name: string
  display_name: string | null
  gender: 'male' | 'female' | 'other' | null
  age_range: string | null
  ethnicity: string | null
  style: string
  languages: string[]
  preview_video_url: string | null
  thumbnail_url: string | null
  is_premium: boolean
  usage_count: number
}

export interface ScriptTemplate {
  id: string
  name: string
  category: string
  template: string
  language: string
  is_premium: boolean
}

export interface AuthResponse {
  access_token: string
  refresh_token: string
  token_type: string
  expires_in: number
  user: User
}

export interface APIResponse<T = unknown> {
  success: boolean
  data?: T
  error?: {
    code: string
    message: string
    details?: Record<string, string>
  }
  meta?: {
    page: number
    limit: number
    total: number
    total_pages: number
  }
}
