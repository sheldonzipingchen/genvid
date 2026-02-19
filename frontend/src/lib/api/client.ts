import type { User, Project, Avatar, AuthResponse, APIResponse } from '@/types'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

class APIClient {
  private token: string | null = null

  setToken(token: string) {
    this.token = token
  }

  clearToken() {
    this.token = null
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<APIResponse<T>> {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    }

    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`
    }

    const response = await fetch(`${API_URL}${endpoint}`, {
      ...options,
      headers,
    })

    return response.json()
  }

  async register(email: string, password: string, fullName?: string): Promise<APIResponse<AuthResponse>> {
    return this.request('/api/auth/register', {
      method: 'POST',
      body: JSON.stringify({ email, password, full_name: fullName }),
    })
  }

  async login(email: string, password: string): Promise<APIResponse<AuthResponse>> {
    return this.request('/api/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    })
  }

  async refreshToken(refreshToken: string): Promise<APIResponse<AuthResponse>> {
    return this.request('/api/auth/refresh', {
      method: 'POST',
      body: JSON.stringify({ refresh_token: refreshToken }),
    })
  }

  async getProfile(): Promise<APIResponse<User>> {
    return this.request('/api/user/profile')
  }

  async updateProfile(data: Partial<User>): Promise<APIResponse<User>> {
    return this.request('/api/user/profile', {
      method: 'PATCH',
      body: JSON.stringify(data),
    })
  }

  async getProjects(page = 1, limit = 20): Promise<APIResponse<Project[]>> {
    return this.request(`/api/projects?page=${page}&limit=${limit}`)
  }

  async getProject(id: string): Promise<APIResponse<Project>> {
    return this.request(`/api/projects/${id}`)
  }

  async createProject(data: {
    product_name: string
    product_description?: string
    product_url?: string
  }): Promise<APIResponse<Project>> {
    return this.request('/api/projects', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async deleteProject(id: string): Promise<APIResponse<void>> {
    return this.request(`/api/projects/${id}`, {
      method: 'DELETE',
    })
  }

  async generateVideo(
    projectId: string,
    data: {
      avatar_id: string
      script: string
      language: string
      format: string
    }
  ): Promise<APIResponse<Project>> {
    return this.request(`/api/projects/${projectId}/generate`, {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async getAvatars(): Promise<APIResponse<Avatar[]>> {
    return this.request('/api/avatars')
  }

  async getAvatar(id: string): Promise<APIResponse<Avatar>> {
    return this.request(`/api/avatars/${id}`)
  }
}

export const api = new APIClient()
