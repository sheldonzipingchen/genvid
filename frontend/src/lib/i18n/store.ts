import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export type Language = 'zh' | 'en'

interface LanguageState {
  language: Language
  setLanguage: (lang: Language) => void
}

export const useLanguageStore = create<LanguageState>()(
  persist(
    (set) => ({
      language: 'zh',
      setLanguage: (lang) => set({ language: lang }),
    }),
    {
      name: 'genvid-language',
    }
  )
)
