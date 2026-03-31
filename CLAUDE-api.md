# CLAUDE.md — FlagBridge API (Go)

## What is FlagBridge

FlagBridge is an open-core Feature Flag Management platform with product intelligence. Boolean, string, number, and JSON flags with targeting rules, environments, percentage rollouts, audit logs, a plugin ecosystem, and managed integrations.

- **Website:** https://flagbridge.io
- **Admin UI:** https://admin.flagbridge.io
- **Docs:** Starlight (Astro) on Cloudflare Pages
- **Org:** https://github.com/flagbridge
- **License:** Apache 2.0 (CE), Commercial (Pro)

## Architecture

```
Client Apps (React, Node, Go, Python via SDKs)
        │
        ▼
┌─────────────────────────────────┐
│   FlagBridge API Server (Go)    │  Fly.io (region gru)
│                                 │
│  Flag Eval Engine    (CE)       │
│  Testing API         (CE+Pro)   │
│  Plugin Runtime      (CE+Pro)   │
│  Product Context     (Pro)      │
│  Dashboard & Metrics (Pro)      │
│  Integrations Layer  (Pro)      │
│         │                       │
│  PostgreSQL (Supabase)          │
└─────────────────────────────────┘
```

### Infrastructure

| Component         | Provider                              | Cost MVP |
|-------------------|---------------------------------------|----------|
| Go API            | Fly.io (region gru)                   | ~$5/mo   |
| Admin UI          | Cloudflare Pages via Wrangler CLI     | $5/mo    |
| Plugin assets     | Cloudflare R2                         | $0       |
| Database          | Supabase (PostgreSQL)                 | $0       |
| Auth (SaaS)       | Supabase Auth (GitHub, Google, email, magic link) | $0 |
| Auth (self-hosted)| Local bcrypt + JWT                    | $0       |
| Cache             | In-memory Go → Fly.io Redis at scale  | $0       |
| Errors            | Sentry                                | $0       |
| Metrics/Logs      | Fly.io native                         | $0       |
| Uptime/Status     | Betterstack                           | $0       |
| CI/CD             | GitHub Actions                        | $0       |
| Payments          | Stripe Connect (80/20 marketplace)    | tx fees  |
| Email             | Resend                                | $0       |
| **TOTAL**         |                                       | **~$16/mo** |

### Editions

| Edition   | License     | What's included                                           |
|-----------|------------|-----------------------------------------------------------|
| Community | Apache 2.0 | Core flags, envs, audit, testing API basic, webhooks, plugins |
| Pro       | License key | Analytics, product cards, lifecycle, integrations, advanced testing, marketplace |
| SaaS      | Hosted     | Managed cloud (future)                                    |

### Pro Gating — ADR-001 Dogfooding

FlagBridge gates its own Pro features using itself:
- Internal project: `_flagbridge`
- Flags: `pro.*` namespace (`pro.analytics`, `pro.marketplace`, etc.)
- UI: `<ProGate flag="pro.analytics">` renders Pro CTA if off
- 14-day trial for new installations

## API — 54 Endpoints (28 CE + 26 Pro)

### Key Scopes

```
fb_sk_eval_...  → evaluation only (SDKs in prod)
fb_sk_mgmt_...  → management (admin UI, CI/CD)
fb_sk_test_...  → testing (E2E, QA)
fb_sk_full_...  → all (local dev)
```

### Endpoints

```
# FLAG MANAGEMENT (CE)
GET    /v1/projects/:project/flags
POST   /v1/projects/:project/flags
GET    /v1/projects/:project/flags/:key
PATCH  /v1/projects/:project/flags/:key
DELETE /v1/projects/:project/flags/:key

# FLAG STATE PER ENV (CE)
GET    /v1/projects/:project/flags/:key/states/:env
PUT    /v1/projects/:project/flags/:key/states/:env

# EVALUATION (CE)
POST   /v1/evaluate
POST   /v1/evaluate/batch
# Resolution: session override > targeting > rollout > env default > flag default

# TESTING API (CE basic / Pro full)
POST   /v1/testing/sessions
DELETE /v1/testing/sessions/:id
PUT    /v1/testing/sessions/:id/overrides
GET    /v1/testing/sessions/:id/overrides
DELETE /v1/testing/sessions/:id/overrides/:flag
GET    /v1/testing/sessions                        # Pro
GET    /v1/testing/sessions/:id/metrics            # Pro
POST   /v1/testing/sessions/:id/snapshot           # Pro
POST   /v1/testing/sessions/:id/restore            # Pro

# WEBHOOKS (CE)
POST   /v1/projects/:project/webhooks
GET    /v1/projects/:project/webhooks
PATCH  /v1/projects/:project/webhooks/:id
DELETE /v1/projects/:project/webhooks/:id
GET    /v1/projects/:project/webhooks/:id/logs
POST   /v1/projects/:project/webhooks/:id/test
# 9 event types, HMAC-SHA256, 5 retries exponential backoff

# PRODUCT CARDS (Pro)
GET    /v1/projects/:project/flags/:key/product-card
PUT    /v1/projects/:project/flags/:key/product-card

# METRICS & DASHBOARD (Pro)
GET    /v1/projects/:project/flags/:key/metrics
GET    /v1/projects/:project/dashboard/overview

# LIFECYCLE (Pro)
GET    /v1/projects/:project/lifecycle/stale
GET    /v1/projects/:project/lifecycle/cleanup-suggestions

# PLUGIN SYSTEM (CE+Pro)
GET    /v1/plugins
POST   /v1/plugins/install
DELETE /v1/plugins/:slug
PATCH  /v1/plugins/:slug/config
GET    /v1/plugins/:slug/status

# MARKETPLACE (Pro)
GET    /v1/marketplace/listings
GET    /v1/marketplace/listings/:slug
POST   /v1/marketplace/listings
PATCH  /v1/marketplace/listings/:slug
POST   /v1/marketplace/listings/:slug/review
GET    /v1/marketplace/developer/earnings
POST   /v1/marketplace/purchase/:slug

# INTEGRATIONS (Pro)
POST   /v1/projects/:project/integrations
GET    /v1/projects/:project/integrations
PATCH  /v1/projects/:project/integrations/:provider
DELETE /v1/projects/:project/integrations/:provider
GET    /v1/projects/:project/integrations/:provider/status
# MVP: Mixpanel, Customer.io, Amplitude, Segment, Datadog, Slack

# ADMIN (CE)
GET    /v1/audit-log
GET    /v1/health
```

## Go Code Structure

```
internal/
├── flag/           # Flag CRUD + evaluation engine
├── testing/        # Sessions, overrides, cleanup goroutine
├── webhook/        # Registration + delivery + retry
├── project/        # Project domain
├── auth/           # AuthProvider interface
├── plugin/         # Plugin runtime engine
├── integration/    # Managed integrations (Pro)
└── api/            # HTTP handlers (Chi router)
pkg/
├── database/       # pgx PostgreSQL client
├── config/         # Env config
├── cache/          # CacheProvider interface
└── middleware/      # HTTP middleware
migrations/         # goose SQL migrations
```

### Key Interfaces

```go
type CacheProvider interface {
    Get(ctx context.Context, key string) ([]byte, error)
    Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
}

type AuthProvider interface {
    ValidateToken(ctx context.Context, token string) (*User, error)
    CreateUser(ctx context.Context, input CreateUserInput) (*User, error)
}
```

## Conventions

- Go 1.22+, Chi router, pgx, sqlc, goose, slog
- snake_case DB, camelCase JSON, PascalCase Go
- Table-driven tests with testify
- Errors: `fmt.Errorf("context: %w", err)`

## Priorities

1. Stabilization & test coverage (eval engine, targeting, rollouts, testing sessions)
2. OpenAPI spec for all 54 endpoints
3. Engineering architecture doc

## Do NOT

- Change Pro gating (ADR-001) without approval
- Modify applied migrations — create new ones
- Add Redis — use in-memory CacheProvider
- Expose `_flagbridge` internal project publicly
- Add Go deps without justification
