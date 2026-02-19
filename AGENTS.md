# GENVID PROJECT KNOWLEDGE BASE

**Generated:** 2026-02-19
**Commit:** 3df033f
**Branch:** main

## OVERVIEW

AI UGC video generation platform for e-commerce sellers. Creates TikTok/Reels/Shorts-style product videos via HeyGen API integration.

**Stack:** Next.js 16 (App Router) + Go 1.22 (chi router) + PostgreSQL + Redis

## STRUCTURE

```
genvid/
├── backend/           # Go HTTP API server
│   ├── cmd/server/    # Entry point
│   ├── internal/      # Private packages (handler/service/repository pattern)
│   ├── pkg/auth/      # JWT utilities (shared)
│   └── migrations/    # SQL migrations
├── frontend/          # Next.js 16 App Router
│   └── src/
│       ├── app/       # Routes (page.tsx per route)
│       ├── components/
│       │   ├── ui/        # shadcn-style primitives
│       │   ├── layout/    # Header, MarketingLayout
│       │   └── features/  # ScriptEditor, AvatarSelector
│       ├── hooks/     # use-stripe
│       ├── stores/    # Zustand (auth)
│       └── lib/api/   # Centralized API client
└── .lane/plans/       # Architecture docs
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Add new API endpoint | `backend/internal/handler/` | Create handler, wire in `cmd/server/main.go` |
| Add business logic | `backend/internal/service/` | Service layer orchestrates repos |
| Database queries | `backend/internal/repository/` | sqlx with `db:` tags |
| New frontend route | `frontend/src/app/<route>/page.tsx` | App Router convention |
| UI components | `frontend/src/components/ui/` | shadcn-style with CVA |
| Auth state | `frontend/src/stores/auth.ts` | Zustand with persist |
| API calls | `frontend/src/lib/api/client.ts` | Single APIClient class |
| Video generation | `backend/internal/zhipu/` | Zhipu CogVideoX client |

## CODE MAP

### Backend Entry Points

| Symbol | Location | Role |
|--------|----------|------|
| `main` | `backend/cmd/server/main.go` | Server bootstrap, DI, routing |
| `AuthHandler` | `backend/internal/handler/handler.go` | Auth HTTP handlers |
| `ProjectHandler` | `backend/internal/handler/handler.go` | Project HTTP handlers |
| `AuthService` | `backend/internal/service/service.go` | Auth business logic |
| `ProjectService` | `backend/internal/service/service.go` | Project/video logic |

### Frontend Entry Points

| Symbol | Location | Role |
|--------|----------|------|
| `RootLayout` | `frontend/src/app/layout.tsx` | Global layout wrapper |
| `page` | `frontend/src/app/page.tsx` | Landing page |
| `useAuthStore` | `frontend/src/stores/auth.ts` | Auth state (Zustand) |
| `api` | `frontend/src/lib/api/client.ts` | HTTP client singleton |

## CONVENTIONS

### Backend (Go)
- **Layered architecture**: handler → service → repository (DI via constructors)
- **DB mapping**: sqlx with `db:` struct tags; pointers for nullable fields
- **Errors**: domain errors (`ErrNotFound`, `ErrUnauthorized`) in repository
- **Context keys**: typed `contextKey` in middleware (NOTE: handlers use string key - needs fix)
- **Config**: nested struct in `internal/config/config.go`, loaded via godotenv

### Frontend (TypeScript)
- **Path alias**: `@/*` maps to `./src/*`
- **State**: Zustand stores with persist middleware for auth
- **UI**: shadcn-style primitives (CVA for variants) in `components/ui/`
- **API**: Single `APIClient` class with token management
- **Forms**: react-hook-form + zod validation
- **Strict TypeScript**: `strict: true`, `noEmit: true`

## ANTI-PATTERNS (THIS PROJECT)

1. **Destructive DB script**: `backend/scripts/db.sh reset` is interactive. Avoid in CI/automation.
2. **Context key mismatch**: Middleware uses typed `contextKey`, handlers read string `"user_id"`. Fix by unifying.
3. **Mock avatars**: `avatar_handler.go` returns hardcoded data. Needs DB integration.
4. **No test config**: No jest/vitest/pytest config found. `make test` runs Go tests only.

## COMMANDS

```bash
# Backend
cd backend
make run              # Start API server (go run ./cmd/server)
make test             # Run Go tests with coverage
make docker-up        # Start PostgreSQL + Redis via compose
./scripts/db.sh setup # Run migrations

# Frontend
cd frontend
npm run dev           # Start Next.js dev server
npm run build         # Production build
npm run lint          # ESLint
```

## NOTES

- **Version mismatch**: README says Next.js 15, package.json has 16.1.6
- **No CI workflows**: Missing `.github/workflows/` - add for GitHub Actions
- **Zhipu API key**: Set `ZHIPU_API_KEY` env var for video generation
- **CogVideoX model**: Default is `cogvideox-3`, configurable via `ZHIPU_MODEL`
- **Stripe webhook**: Production needs proper signature verification
