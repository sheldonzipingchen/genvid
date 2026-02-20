import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { Language } from './translations'

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
