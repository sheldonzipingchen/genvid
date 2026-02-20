'use client'

import { useEffect, useState, useRef } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useAuthStore } from '@/stores/auth'
import { AvatarSelector } from '@/components/features/avatar-selector'
import { ScriptEditor } from '@/components/features/script-editor'
import type { Avatar, Project } from '@/types'
import { 
  Video, 
  ArrowLeft, 
  ArrowRight, 
  Loader2, 
  Upload,
  Link as LinkIcon,
  Sparkles,
  X,
  Check
} from 'lucide-react'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

const steps = [
  { id: 1, name: 'Product', description: 'Add your product' },
  { id: 2, name: 'Avatar', description: 'Choose your avatar' },
  { id: 3, name: 'Script', description: 'Create your script' },
  { id: 4, name: 'Generate', description: 'Generate video' },
]

export default function CreatePage() {
  const router = useRouter()
  const { user, token, isAuthenticated, _hasHydrated } = useAuthStore()
  const [currentStep, setCurrentStep] = useState(1)
  const [loading, setLoading] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)
  const [productImages, setProductImages] = useState<string[]>([])
  const [productImageURL, setProductImageURL] = useState<string | null>(null)
  const [uploadingImage, setUploadingImage] = useState(false)
  const [isDragging, setIsDragging] = useState(false)
  
  const [project, setProject] = useState<Partial<Project>>({
    product_name: '',
    product_description: '',
    product_url: '',
    language: 'en',
    format: '9:16',
  })
  const [selectedAvatar, setSelectedAvatar] = useState<Avatar | null>(null)
  const [script, setScript] = useState('')
  const [videoDuration, setVideoDuration] = useState(5)
  const [projectId, setProjectId] = useState<string | null>(null)

  useEffect(() => {
    if (!_hasHydrated) return
    if (!isAuthenticated) {
      router.push('/login')
    }
  }, [isAuthenticated, _hasHydrated, router])

  const uploadImage = async (file: File): Promise<string | null> => {
    const formData = new FormData()
    formData.append('file', file)
    
    try {
      const response = await fetch(`${API_URL}/api/upload`, {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: formData,
      })
      
      const data = await response.json()
      if (response.ok && data.data?.url) {
        return data.data.url
      }
    } catch (error) {
      console.error('Failed to upload image:', error)
    }
    return null
  }

  const handleFileSelect = async (files: FileList | null) => {
    if (!files) return
    
    const file = files[0]
    if (file && file.type.startsWith('image/')) {
      const preview = URL.createObjectURL(file)
      setProductImages([preview])
      
      setUploadingImage(true)
      const url = await uploadImage(file)
      setUploadingImage(false)
      
      if (url) {
        setProductImageURL(url)
      }
    }
  }

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(false)
    handleFileSelect(e.dataTransfer.files)
  }

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(true)
  }

  const handleDragLeave = () => {
    setIsDragging(false)
  }

  const removeImage = () => {
    setProductImages([])
    setProductImageURL(null)
  }

  const handleProductSubmit = async () => {
    if (!project.product_name) return
    
    setLoading(true)
    try {
      const response = await fetch(`${API_URL}/api/projects`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          product_name: project.product_name,
          product_description: project.product_description,
          product_url: project.product_url,
          product_image_url: productImageURL,
        }),
      })
      
      const data = await response.json()
      if (response.ok && (data.data || data.id)) {
        setProjectId(data.data?.id || data.id)
        setCurrentStep(2)
      }
    } catch (error) {
      console.error('Failed to create project:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleGenerateVideo = async () => {
    if (!projectId || !selectedAvatar || !script) return
    
    setLoading(true)
    try {
      const response = await fetch(`${API_URL}/api/projects/${projectId}/generate`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          avatar_id: selectedAvatar.id,
          script: script,
          language: project.language,
          format: project.format,
          video_duration: videoDuration,
        }),
      })
      
      const data = await response.json()
      if (response.ok) {
        router.push('/dashboard')
      }
    } catch (error) {
      console.error('Failed to generate video:', error)
    } finally {
      setLoading(false)
    }
  }

  const canProceed = () => {
    switch (currentStep) {
      case 1:
        return !!project.product_name
      case 2:
        return !!selectedAvatar
      case 3:
        return !!script && script.length >= 10
      default:
        return false
    }
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white border-b border-gray-200">
        <div className="mx-auto max-w-5xl px-4 sm:px-6 lg:px-8">
          <div className="flex h-16 items-center justify-between">
            <Link href="/dashboard" className="flex items-center gap-2 text-gray-600 hover:text-gray-900">
              <ArrowLeft className="h-4 w-4" />
              Back to Dashboard
            </Link>
            <div className="flex items-center gap-2">
              <span className="text-sm text-gray-500">{user?.credits_remaining} credits</span>
            </div>
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-5xl px-4 sm:px-6 lg:px-8 py-8">
        <nav className="mb-8">
          <ol className="flex items-center justify-between">
            {steps.map((step, index) => (
              <li key={step.id} className="flex-1">
                <div className={`flex items-center ${index !== steps.length - 1 ? 'pr-8' : ''}`}>
                  <div className={`flex h-10 w-10 items-center justify-center rounded-full border-2 ${
                    currentStep >= step.id
                      ? 'border-violet-600 bg-violet-600 text-white'
                      : 'border-gray-300 text-gray-400'
                  }`}>
                    {currentStep > step.id ? (
                      <svg className="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
                        <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                      </svg>
                    ) : (
                      step.id
                    )}
                  </div>
                  <div className="ml-4 flex-1">
                    <p className={`text-sm font-medium ${currentStep >= step.id ? 'text-gray-900' : 'text-gray-400'}`}>
                      {step.name}
                    </p>
                    <p className="text-xs text-gray-500">{step.description}</p>
                  </div>
                  {index !== steps.length - 1 && (
                    <div className={`h-0.5 flex-1 ${currentStep > step.id ? 'bg-violet-600' : 'bg-gray-200'}`} />
                  )}
                </div>
              </li>
            ))}
          </ol>
        </nav>

        <Card>
          <CardContent className="p-8">
            {currentStep === 1 && (
              <div className="space-y-6">
                <div>
                  <h2 className="text-xl font-semibold text-gray-900 mb-2">Add Your Product</h2>
                  <p className="text-gray-500">Enter your product details to generate a UGC video</p>
                </div>

                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Product Name <span className="text-red-500">*</span>
                    </label>
                    <Input
                      value={project.product_name || ''}
                      onChange={(e) => setProject({ ...project, product_name: e.target.value })}
                      placeholder="e.g., Premium Wireless Earbuds"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Product URL
                    </label>
                    <div className="relative">
                      <LinkIcon className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
                      <Input
                        value={project.product_url || ''}
                        onChange={(e) => setProject({ ...project, product_url: e.target.value })}
                        placeholder="https://yourstore.com/product"
                        className="pl-10"
                      />
                    </div>
                    <p className="text-xs text-gray-500 mt-1">We'll extract product details automatically</p>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Product Description
                    </label>
                    <textarea
                      value={project.product_description || ''}
                      onChange={(e) => setProject({ ...project, product_description: e.target.value })}
                      placeholder="Describe your product's key features and benefits..."
                      rows={4}
                      className="w-full rounded-lg border border-gray-200 p-3 text-sm focus:outline-none focus:ring-2 focus:ring-violet-500"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Product Images
                    </label>
                    <input
                      ref={fileInputRef}
                      type="file"
                      accept="image/*"
                      multiple
                      className="hidden"
                      onChange={(e) => handleFileSelect(e.target.files)}
                    />
                    <div
                      onClick={() => !uploadingImage && fileInputRef.current?.click()}
                      onDrop={handleDrop}
                      onDragOver={handleDragOver}
                      onDragLeave={handleDragLeave}
                      className={`border-2 border-dashed rounded-lg p-8 text-center transition-colors ${
                        uploadingImage ? 'cursor-not-allowed opacity-50' : 'cursor-pointer'
                      } ${
                        isDragging ? 'border-violet-500 bg-violet-50' : 'border-gray-200 hover:border-violet-400'
                      }`}
                    >
                      {uploadingImage ? (
                        <Loader2 className="h-8 w-8 mx-auto text-violet-600 mb-2 animate-spin" />
                      ) : (
                        <Upload className="h-8 w-8 mx-auto text-gray-400 mb-2" />
                      )}
                      <p className="text-sm text-gray-600">
                        {uploadingImage ? 'Uploading...' : 'Click or drag images to upload'}
                      </p>
                      <p className="text-xs text-gray-400 mt-1">PNG, JPG up to 10MB</p>
                    </div>
                    {productImages.length > 0 && (
                      <div className="flex flex-wrap gap-3 mt-4">
                        {productImages.map((img, index) => (
                          <div key={index} className="relative group">
                            <img
                              src={img}
                              alt={`Product ${index + 1}`}
                              className="h-20 w-20 object-cover rounded-lg border border-gray-200"
                            />
                            {productImageURL && (
                              <div className="absolute -bottom-1 -right-1 h-5 w-5 bg-green-500 text-white rounded-full flex items-center justify-center">
                                <Check className="h-3 w-3" />
                              </div>
                            )}
                            <button
                              onClick={(e) => { e.stopPropagation(); removeImage(); }}
                              className="absolute -top-2 -right-2 h-5 w-5 bg-red-500 text-white rounded-full opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center"
                            >
                              <X className="h-3 w-3" />
                            </button>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                </div>
              </div>
            )}

            {currentStep === 2 && (
              <AvatarSelector
                selectedAvatar={selectedAvatar}
                onSelect={setSelectedAvatar}
              />
            )}

            {currentStep === 3 && (
              <ScriptEditor
                productName={project.product_name || ''}
                productDescription={project.product_description || ''}
                script={script}
                videoDuration={videoDuration}
                onChange={setScript}
                onDurationChange={setVideoDuration}
              />
            )}

            {currentStep === 4 && (
              <div className="text-center py-8">
                <div className="h-20 w-20 mx-auto mb-6 rounded-full bg-gradient-to-br from-violet-100 to-purple-100 flex items-center justify-center">
                  <Sparkles className="h-10 w-10 text-violet-600" />
                </div>
                <h2 className="text-xl font-semibold text-gray-900 mb-2">Ready to Generate!</h2>
                <p className="text-gray-500 mb-6 max-w-md mx-auto">
                  Your video will be generated in 2-5 minutes. We'll notify you when it's ready.
                </p>
                
                <div className="bg-gray-50 rounded-lg p-6 max-w-md mx-auto text-left">
                  <h3 className="font-medium text-gray-900 mb-4">Summary</h3>
                  <dl className="space-y-2 text-sm">
                    <div className="flex justify-between">
                      <dt className="text-gray-500">Product</dt>
                      <dd className="text-gray-900 font-medium">{project.product_name}</dd>
                    </div>
                    <div className="flex justify-between">
                      <dt className="text-gray-500">Avatar</dt>
                      <dd className="text-gray-900 font-medium">{selectedAvatar?.display_name || selectedAvatar?.name}</dd>
                    </div>
                    <div className="flex justify-between">
                      <dt className="text-gray-500">Duration</dt>
                      <dd className="text-gray-900 font-medium">{videoDuration}s</dd>
                    </div>
                    <div className="flex justify-between">
                      <dt className="text-gray-500">Format</dt>
                      <dd className="text-gray-900 font-medium">{project.format}</dd>
                    </div>
                    <div className="flex justify-between">
                      <dt className="text-gray-500">Language</dt>
                      <dd className="text-gray-900 font-medium uppercase">{project.language}</dd>
                    </div>
                    <div className="flex justify-between pt-2 border-t border-gray-200">
                      <dt className="text-gray-500">Cost</dt>
                      <dd className="text-violet-600 font-medium">1 credit</dd>
                    </div>
                  </dl>
                </div>
              </div>
            )}

            <div className="flex justify-between mt-8 pt-6 border-t border-gray-100">
              <Button
                variant="outline"
                onClick={() => setCurrentStep(currentStep - 1)}
                disabled={currentStep === 1}
              >
                <ArrowLeft className="h-4 w-4 mr-2" />
                Previous
              </Button>
              
              {currentStep < 4 ? (
                <Button
                  onClick={() => {
                    if (currentStep === 1) handleProductSubmit()
                    else setCurrentStep(currentStep + 1)
                  }}
                  disabled={!canProceed() || loading}
                >
                  {loading ? (
                    <Loader2 className="h-4 w-4 animate-spin" />
                  ) : (
                    <>
                      Next
                      <ArrowRight className="h-4 w-4 ml-2" />
                    </>
                  )}
                </Button>
              ) : (
                <Button
                  onClick={handleGenerateVideo}
                  disabled={loading}
                >
                  {loading ? (
                    <>
                      <Loader2 className="h-4 w-4 animate-spin mr-2" />
                      Generating...
                    </>
                  ) : (
                    <>
                      <Video className="h-4 w-4 mr-2" />
                      Generate Video
                    </>
                  )}
                </Button>
              )}
            </div>
          </CardContent>
        </Card>
      </main>
    </div>
  )
}
