import { create } from 'zustand'
import { persist, createJSONStorage } from 'zustand/middleware'
import type { User } from '@/types'

interface AuthState {
  user: User | null
  token: string | null
  refreshToken: string | null
  isAuthenticated: boolean
  _hasHydrated: boolean
  setUser: (user: User | null) => void
  setTokens: (token: string, refreshToken: string) => void
  setAuth: (user: User, token: string, refreshToken: string) => void
  logout: () => void
  setHasHydrated: (state: boolean) => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      refreshToken: null,
      isAuthenticated: false,
      _hasHydrated: false,
      setUser: (user) => set(state => ({ 
        user, 
        isAuthenticated: !!user && !!state.token 
      })),
      setTokens: (token, refreshToken) => set(state => ({ 
        token, 
        refreshToken,
        isAuthenticated: !!token && !!state.user
      })),
      setAuth: (user, token, refreshToken) => set({ 
        user, 
        token, 
        refreshToken,
        isAuthenticated: true
      }),
      logout: () => set({ user: null, token: null, refreshToken: null, isAuthenticated: false, _hasHydrated: true }),
      setHasHydrated: (state) => set({ _hasHydrated: state }),
    }),
    {
      name: 'genvid-auth',
      storage: createJSONStorage(() => localStorage),
      onRehydrateStorage: () => (state) => {
        state?.setHasHydrated(true)
        if (state?.token && state?.user) {
          state.isAuthenticated = true
        }
      },
    }
  )
)
