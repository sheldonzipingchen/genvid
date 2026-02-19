'use client'

import { useEffect, useState } from 'react'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import type { Avatar } from '@/types'
import { useAuthStore } from '@/stores/auth'
import { Check, Crown, Loader2 } from 'lucide-react'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

const genderFilters = [
  { value: 'all', label: 'All' },
  { value: 'male', label: 'Male' },
  { value: 'female', label: 'Female' },
]

const styleFilters = [
  { value: 'all', label: 'All Styles' },
  { value: 'casual', label: 'Casual' },
  { value: 'professional', label: 'Professional' },
  { value: 'energetic', label: 'Energetic' },
  { value: 'friendly', label: 'Friendly' },
]

interface AvatarSelectorProps {
  selectedAvatar: Avatar | null
  onSelect: (avatar: Avatar) => void
}

export function AvatarSelector({ selectedAvatar, onSelect }: AvatarSelectorProps) {
  const { token } = useAuthStore()
  const [avatars, setAvatars] = useState<Avatar[]>([])
  const [loading, setLoading] = useState(true)
  const [genderFilter, setGenderFilter] = useState('all')
  const [styleFilter, setStyleFilter] = useState('all')

  useEffect(() => {
    fetchAvatars()
  }, [])

  const fetchAvatars = async () => {
    try {
      const response = await fetch(`${API_URL}/api/avatars`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      const data = await response.json()
      if (response.ok) {
        setAvatars(data.data || data || [])
      } else {
        setAvatars(getMockAvatars())
      }
    } catch (error) {
      setAvatars(getMockAvatars())
    } finally {
      setLoading(false)
    }
  }

  const filteredAvatars = avatars.filter((avatar) => {
    if (genderFilter !== 'all' && avatar.gender !== genderFilter) return false
    if (styleFilter !== 'all' && avatar.style !== styleFilter) return false
    return true
  })

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-xl font-semibold text-gray-900 mb-2">Choose Your Avatar</h2>
        <p className="text-gray-500">Select an AI avatar to present your product</p>
      </div>

      <div className="flex gap-4">
        <div className="flex gap-2">
          {genderFilters.map((filter) => (
            <Button
              key={filter.value}
              variant={genderFilter === filter.value ? 'default' : 'outline'}
              size="sm"
              onClick={() => setGenderFilter(filter.value)}
            >
              {filter.label}
            </Button>
          ))}
        </div>
        <select
          value={styleFilter}
          onChange={(e) => setStyleFilter(e.target.value)}
          className="rounded-lg border border-gray-200 px-3 py-2 text-sm"
        >
          {styleFilters.map((filter) => (
            <option key={filter.value} value={filter.value}>
              {filter.label}
            </option>
          ))}
        </select>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="h-8 w-8 animate-spin text-violet-600" />
        </div>
      ) : (
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          {filteredAvatars.map((avatar) => (
            <Card
              key={avatar.id}
              className={`cursor-pointer transition-all hover:shadow-lg ${
                selectedAvatar?.id === avatar.id
                  ? 'ring-2 ring-violet-600 ring-offset-2'
                  : ''
              }`}
              onClick={() => onSelect(avatar)}
            >
              <CardContent className="p-0">
                <div className="aspect-square bg-gradient-to-br from-violet-100 to-purple-100 rounded-t-lg relative">
                  {avatar.thumbnail_url ? (
                    <img
                      src={avatar.thumbnail_url}
                      alt={avatar.name}
                      className="w-full h-full object-cover rounded-t-lg"
                    />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center">
                      <div className="h-16 w-16 rounded-full bg-gradient-to-br from-violet-400 to-purple-400" />
                    </div>
                  )}
                  {selectedAvatar?.id === avatar.id && (
                    <div className="absolute top-2 right-2 h-6 w-6 rounded-full bg-violet-600 flex items-center justify-center">
                      <Check className="h-4 w-4 text-white" />
                    </div>
                  )}
                  {avatar.is_premium && (
                    <div className="absolute top-2 left-2 px-2 py-0.5 rounded-full bg-yellow-400 text-yellow-900 text-xs font-medium flex items-center gap-1">
                      <Crown className="h-3 w-3" />
                      Pro
                    </div>
                  )}
                </div>
                <div className="p-3">
                  <p className="font-medium text-gray-900">{avatar.display_name || avatar.name}</p>
                  <p className="text-xs text-gray-500 capitalize">{avatar.style}</p>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  )
}

function getMockAvatars(): Avatar[] {
  return [
    { id: '1', name: 'emma', display_name: 'Emma', gender: 'female', style: 'casual', languages: ['en'], is_premium: false, usage_count: 100 },
    { id: '2', name: 'james', display_name: 'James', gender: 'male', style: 'professional', languages: ['en'], is_premium: false, usage_count: 85 },
    { id: '3', name: 'sofia', display_name: 'Sofia', gender: 'female', style: 'energetic', languages: ['en', 'es'], is_premium: true, usage_count: 120 },
    { id: '4', name: 'li', display_name: 'Li', gender: 'male', style: 'friendly', languages: ['en', 'zh'], is_premium: false, usage_count: 75 },
    { id: '5', name: 'maria', display_name: 'Maria', gender: 'female', style: 'casual', languages: ['en', 'pt'], is_premium: true, usage_count: 90 },
    { id: '6', name: 'alex', display_name: 'Alex', gender: 'male', style: 'energetic', languages: ['en'], is_premium: false, usage_count: 65 },
    { id: '7', name: 'yuki', display_name: 'Yuki', gender: 'female', style: 'professional', languages: ['en', 'ja'], is_premium: true, usage_count: 110 },
    { id: '8', name: 'david', display_name: 'David', gender: 'male', style: 'casual', languages: ['en'], is_premium: false, usage_count: 55 },
  ]
}
