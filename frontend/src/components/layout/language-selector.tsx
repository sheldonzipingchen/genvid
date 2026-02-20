'use client'

import { useLanguageStore } from '@/lib/i18n/store'
import { Globe } from 'lucide-react'
import { Button } from '@/components/ui/button'

export function LanguageSelector() {
  const { language, setLanguage } = useLanguageStore()
  
  return (
    <Button
      variant="ghost"
      size="sm"
      onClick={() => setLanguage(language === 'zh' ? 'en' : 'zh')}
      className="gap-2"
    >
      <Globe className="h-4 w-4" />
      {language === 'zh' ? 'EN' : '中文'}
    </Button>
  )
}
