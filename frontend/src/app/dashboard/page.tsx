'use client'

import { useEffect, useState, useCallback } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { useAuthStore } from '@/stores/auth'
import { LanguageSelector } from '@/components/layout/language-selector'
import { useTranslation } from '@/lib/i18n'
import type { Project } from '@/types'
import { 
  Video, 
  Plus, 
  Clock, 
  CheckCircle, 
  XCircle, 
  Play,
  Download,
  Trash2,
  Loader2
} from 'lucide-react'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

const statusConfig = {
  draft: { label: 'Draft', color: 'bg-gray-100 text-gray-600', icon: Clock },
  queued: { label: 'Queued', color: 'bg-blue-100 text-blue-600', icon: Clock },
  processing: { label: 'Processing', color: 'bg-yellow-100 text-yellow-600', icon: Loader2 },
  completed: { label: 'Completed', color: 'bg-green-100 text-green-600', icon: CheckCircle },
  failed: { label: 'Failed', color: 'bg-red-100 text-red-600', icon: XCircle },
  canceled: { label: 'Canceled', color: 'bg-gray-100 text-gray-600', icon: XCircle },
}

export default function DashboardPage() {
  const router = useRouter()
  const { user, token, isAuthenticated, logout, _hasHydrated } = useAuthStore()
  const { t } = useTranslation()
  const [projects, setProjects] = useState<Project[]>([])
  const [loading, setLoading] = useState(true)

  const fetchProjects = useCallback(async () => {
    try {
      const response = await fetch(`${API_URL}/api/projects`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      const data = await response.json()
      if (response.ok) {
        setProjects(data.data || data || [])
      }
    } catch (error) {
      console.error('Failed to fetch projects:', error)
    } finally {
      setLoading(false)
    }
  }, [token])

  useEffect(() => {
    if (!_hasHydrated) return
    if (!isAuthenticated) {
      router.push('/login')
      return
    }
    fetchProjects()
  }, [isAuthenticated, _hasHydrated, router, fetchProjects])

  useEffect(() => {
    const hasProcessing = projects.some(p => p.status === 'processing' || p.status === 'queued')
    if (!hasProcessing) return

    const interval = setInterval(fetchProjects, 5000)
    return () => clearInterval(interval)
  }, [projects, fetchProjects])

  const handleDelete = async (projectId: string) => {
    if (!confirm('Are you sure you want to delete this project?')) return
    
    try {
      await fetch(`${API_URL}/api/projects/${projectId}`, {
        method: 'DELETE',
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      setProjects(projects.filter(p => p.id !== projectId))
    } catch (error) {
      console.error('Failed to delete project:', error)
    }
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white border-b border-gray-200">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="flex h-16 items-center justify-between">
            <Link href="/" className="flex items-center space-x-2">
              <div className="h-8 w-8 rounded-lg bg-gradient-to-br from-violet-600 to-purple-600 flex items-center justify-center">
                <Video className="h-5 w-5 text-white" />
              </div>
              <span className="text-xl font-bold bg-gradient-to-r from-violet-600 to-purple-600 bg-clip-text text-transparent">
                Genvid
              </span>
            </Link>
            <div className="flex items-center gap-4">
              <LanguageSelector />
              <div className="text-right">
                <p className="text-sm font-medium text-gray-900">{user?.full_name || user?.email}</p>
                <p className="text-xs text-gray-500">{user?.credits_remaining} {t('dashboard.creditsRemaining').toLowerCase()}</p>
              </div>
              <Button 
                variant="outline" 
                onClick={() => { logout(); router.push('/') }}
                className="text-red-600 border-red-200 hover:bg-red-50 hover:text-red-700"
              >
                {t('nav.logout')}
              </Button>
            </div>
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
          <Card>
            <CardContent className="p-4">
              <p className="text-sm text-gray-500">Credits Remaining</p>
              <p className="text-2xl font-bold text-gray-900">{user?.credits_remaining || 0}</p>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4">
              <p className="text-sm text-gray-500">Videos Created</p>
              <p className="text-2xl font-bold text-gray-900">{user?.credits_used_total || 0}</p>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4">
              <p className="text-sm text-gray-500">Plan</p>
              <p className="text-2xl font-bold text-gray-900 capitalize">{user?.subscription_tier || 'Free'}</p>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4 flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-500">Need more credits?</p>
                <p className="text-sm text-violet-600">Upgrade your plan</p>
              </div>
              <Button variant="outline" size="sm" asChild>
                <Link href="/pricing">Upgrade</Link>
              </Button>
            </CardContent>
          </Card>
        </div>

        <div className="flex items-center justify-between mb-6">
          <h1 className="text-2xl font-bold text-gray-900">Your Videos</h1>
          <Button asChild>
            <Link href="/create">
              <Plus className="h-4 w-4 mr-2" />
              Create Video
            </Link>
          </Button>
        </div>

        {loading ? (
          <div className="flex items-center justify-center py-12">
            <Loader2 className="h-8 w-8 animate-spin text-violet-600" />
          </div>
        ) : projects.length === 0 ? (
          <Card className="text-center py-12">
            <CardContent>
              <div className="h-16 w-16 mx-auto mb-4 rounded-full bg-violet-100 flex items-center justify-center">
                <Video className="h-8 w-8 text-violet-600" />
              </div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">No videos yet</h3>
              <p className="text-gray-500 mb-6">Create your first UGC video to get started</p>
              <Button asChild>
                <Link href="/create">
                  <Plus className="h-4 w-4 mr-2" />
                  Create Your First Video
                </Link>
              </Button>
            </CardContent>
          </Card>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {projects.map((project) => {
              const StatusIcon = statusConfig[project.status].icon
              return (
                <Card key={project.id} className="overflow-hidden hover:shadow-lg transition-shadow">
                  <div className="aspect-video bg-gray-100 relative">
                    {project.thumbnail_url ? (
                      <img 
                        src={project.thumbnail_url} 
                        alt={project.title || 'Video thumbnail'}
                        className="w-full h-full object-cover"
                      />
                    ) : (
                      <div className="w-full h-full flex items-center justify-center">
                        <Video className="h-12 w-12 text-gray-300" />
                      </div>
                    )}
                    {project.status === 'completed' && (
                      <div className="absolute inset-0 flex items-center justify-center bg-black/30 opacity-0 hover:opacity-100 transition-opacity">
                        <Button size="icon" className="h-12 w-12 rounded-full">
                          <Play className="h-6 w-6" />
                        </Button>
                      </div>
                    )}
                    {project.status === 'processing' && (
                      <div className="absolute inset-0 flex items-center justify-center bg-black/50">
                        <div className="text-center text-white">
                          <Loader2 className="h-8 w-8 animate-spin mx-auto mb-2" />
                          <p className="text-sm">{project.progress_percent}%</p>
                        </div>
                      </div>
                    )}
                  </div>
                  <CardContent className="p-4">
                    <div className="flex items-start justify-between mb-2">
                      <h3 className="font-semibold text-gray-900 truncate">
                        {project.title || project.product_name || 'Untitled'}
                      </h3>
                      <span className={`flex items-center gap-1 px-2 py-1 rounded-full text-xs ${statusConfig[project.status].color}`}>
                        <StatusIcon className="h-3 w-3" />
                        {statusConfig[project.status].label}
                      </span>
                    </div>
                    <p className="text-sm text-gray-500 mb-4 truncate">
                      {project.product_description || 'No description'}
                    </p>
                    <div className="flex items-center gap-2">
                      {project.status === 'completed' && project.video_url && (
                        <Button size="sm" variant="outline" asChild>
                          <a href={project.video_url} download>
                            <Download className="h-3 w-3 mr-1" />
                            Download
                          </a>
                        </Button>
                      )}
                      {project.status === 'draft' && (
                        <Button size="sm" asChild>
                          <Link href={`/create/${project.id}`}>
                            Continue
                          </Link>
                        </Button>
                      )}
                      <Button 
                        size="sm" 
                        variant="ghost" 
                        className="ml-auto text-gray-400 hover:text-red-500"
                        onClick={() => handleDelete(project.id)}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              )
            })}
          </div>
        )}
      </main>
    </div>
  )
}
