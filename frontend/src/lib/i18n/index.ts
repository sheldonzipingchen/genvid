import { translations, Language } from './translations'
import { useLanguageStore } from './store'

type NestedKeyOf<ObjectType extends object> = {
  [Key in keyof ObjectType & string]: ObjectType[Key] extends object
    ? `${Key}.${NestedKeyOf<ObjectType[Key]>}`
    : Key
}[keyof ObjectType & string]

type TranslationKeys = NestedKeyOf<typeof translations.zh>

function getNestedValue(obj: Record<string, unknown>, path: string): string {
  const value = path.split('.').reduce<unknown>((acc, key) => {
    if (acc && typeof acc === 'object' && key in acc) {
      return (acc as Record<string, unknown>)[key]
    }
    return undefined
  }, obj)
  
  return typeof value === 'string' ? value : path
}

export function useTranslation() {
  const language = useLanguageStore((state) => state.language)
  
  const t = (key: string): string => {
    const langTranslations = translations[language] || translations.zh
    return getNestedValue(langTranslations as unknown as Record<string, unknown>, key)
  }
  
  return { t, language }
}

export { Language }
