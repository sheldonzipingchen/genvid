'use client'

import { useState } from 'react'
import { api } from '@/lib/api/client'
import { useAuthStore } from '@/stores/auth'

interface CheckoutOptions {
  planId: string
  successUrl?: string
  cancelUrl?: string
}

export function useStripe() {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const { token } = useAuthStore()

  const createCheckoutSession = async (options: CheckoutOptions) => {
    if (!token) {
      setError('Not authenticated')
      return null
    }

    setLoading(true)
    setError(null)

    api.setToken(token)

    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/payments/checkout`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({
          plan_id: options.planId,
          success_url: options.successUrl || `${window.location.origin}/dashboard?payment=success`,
          cancel_url: options.cancelUrl || `${window.location.origin}/pricing?payment=canceled`,
        }),
      })

      const data = await response.json()

      if (data.success && data.data?.session_url) {
        window.location.href = data.data.session_url
        return data.data
      } else {
        setError(data.error?.message || 'Failed to create checkout session')
        return null
      }
    } catch (err) {
      setError('Failed to connect to payment service')
      return null
    } finally {
      setLoading(false)
    }
  }

  return {
    createCheckoutSession,
    loading,
    error,
  }
}
