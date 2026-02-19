# BACKEND KNOWLEDGE BASE

**Generated:** 2026-02-19

## OVERVIEW

Go 1.22 HTTP API server with layered architecture (handler/service/repository), JWT auth, HeyGen video generation, and Stripe payments.

## STRUCTURE

```
backend/
├── cmd/server/main.go     # Entry point, DI, routing
├── internal/
│   ├── config/            # Config loading (godotenv)
│   ├── handler/           # HTTP handlers
│   ├── service/           # Business logic
│   ├── repository/        # DB access (sqlx)
│   ├── model/             # Data models + API response types
│   ├── middleware/        # Auth, CORS, logging
│   ├── heygen/            # HeyGen API client
│   └── worker/            # Video generation worker
├── pkg/auth/              # JWT + password utilities (shared)
├── migrations/            # SQL migration files
└── scripts/db.sh          # DB management script
```

## WHERE TO LOOK

| Task | Location |
|------|----------|
| Add endpoint | `internal/handler/` + wire in `cmd/server/main.go` |
| Add business logic | `internal/service/` |
| Add DB query | `internal/repository/` |
| Add data model | `internal/model/models.go` |
| Add middleware | `internal/middleware/` |
| HeyGen integration | `internal/heygen/client.go` |
| Video processing | `internal/worker/video_worker.go` |
| Config/env | `internal/config/config.go` |

## CONVENTIONS

### Layered Architecture
```
HTTP Request → Handler → Service → Repository → DB
                  ↓
            External APIs (HeyGen, Stripe)
```

- **Handlers**: Parse requests, call services, format responses
- **Services**: Business logic, orchestrate repos and external clients
- **Repositories**: SQL queries via sqlx, domain errors

### Dependency Injection
```go
// Constructor pattern
func NewAuthHandler(authService *service.AuthService) *AuthHandler
func NewAuthService(profileRepo *repository.ProfileRepository, jwtService *auth.JWTService, cfg *config.Config) *AuthService
```

### DB Patterns
- Use `sqlx` with `db:` struct tags
- Context-aware methods: `GetContext`, `ExecContext`
- Pointers for nullable fields in models
- Domain errors: `ErrNotFound`, `ErrUnauthorized`

### API Response Format
```go
type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   *APIError   `json:"error,omitempty"`
    Meta    *Meta       `json:"meta,omitempty"`
}
```

## ANTI-PATTERNS (FIX NEEDED)

1. **Context key mismatch**: Middleware uses `contextKey("user_id")`, handlers read string `"user_id"`. Unify to typed key.

2. **Mock avatars**: `avatar_handler.go` returns hardcoded list. Wire to repository.

## COMMANDS

```bash
make run          # Start server
make test         # Run tests with coverage
make build        # Build binary to bin/genvid-backend
make docker-up    # Start PostgreSQL + Redis
./scripts/db.sh setup   # Run migrations
```

## EXTERNAL INTEGRATIONS

| Service | Config Key | Location |
|---------|------------|----------|
| HeyGen | `HEYGEN_API_KEY` | `internal/heygen/client.go` |
| Stripe | `STRIPE_SECRET_KEY` | `internal/handler/payment_handler.go` |
| JWT | `JWT_SECRET` | `pkg/auth/jwt.go` |

## MIGRATIONS

Located in `migrations/`. Run with:
```bash
./scripts/db.sh setup   # Apply all
./scripts/db.sh reset   # Drop and recreate (DESTRUCTIVE)
```
