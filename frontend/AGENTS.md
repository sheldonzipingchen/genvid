# FRONTEND KNOWLEDGE BASE

**Generated:** 2026-02-19

## OVERVIEW

Next.js 16 App Router frontend with Zustand state, centralized API client, and shadcn-style UI components.

## STRUCTURE

```
src/
├── app/               # App Router routes
│   ├── page.tsx       # Landing page
│   ├── layout.tsx     # Root layout
│   ├── login/         # Auth routes
│   ├── register/
│   ├── dashboard/     # Protected area
│   │   └── payment/
│   ├── create/        # Video creation wizard
│   └── pricing/
├── components/
│   ├── ui/            # Primitives (Button, Card, Input)
│   ├── layout/        # Header, MarketingLayout
│   └── features/      # ScriptEditor, AvatarSelector
├── hooks/             # use-stripe (custom hooks)
├── stores/            # Zustand stores (auth)
├── lib/
│   ├── api/           # APIClient singleton
│   └── utils.ts       # cn() helper (clsx + tailwind-merge)
└── types/             # TypeScript interfaces
```

## WHERE TO LOOK

| Task | Location |
|------|----------|
| Add new route | `src/app/<route>/page.tsx` |
| Add UI component | `src/components/ui/` (follow CVA pattern) |
| Add feature component | `src/components/features/` |
| Auth state | `src/stores/auth.ts` |
| API endpoints | `src/lib/api/client.ts` |
| TypeScript types | `src/types/index.ts` |

## CONVENTIONS

- **Path alias**: Import with `@/` prefix (e.g., `@/components/ui/button`)
- **Component pattern**: Forward ref + CVA variants + `cn()` for classes
- **State**: Zustand stores with `create()` + `persist()` middleware
- **API client**: Singleton class with token management, returns `APIResponse<T>`
- **Forms**: react-hook-form + @hookform/resolvers + zod

## KEY PATTERNS

### API Client Usage
```typescript
import { api } from '@/lib/api/client'

const response = await api.login(email, password)
if (response.success) {
  useAuthStore.getState().setTokens(response.data.access_token, response.data.refresh_token)
}
```

### Auth Store
```typescript
const { user, isAuthenticated, logout } = useAuthStore()
```

### UI Components
- Use `Button`, `Card`, `Input` from `@/components/ui/`
- Variants: `default`, `secondary`, `outline`, `ghost`, `destructive`
- Sizes: `sm`, `default`, `lg`, `xl`, `icon`

## DEPENDENCIES

- `zustand` - State management
- `framer-motion` - Animations
- `lucide-react` - Icons
- `react-hook-form` + `zod` - Forms
- `class-variance-authority` - Component variants
- `tailwind-merge` + `clsx` - Class utilities
