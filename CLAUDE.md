# CLAUDE.md — FlagBridge API (Go)

## What is FlagBridge

FlagBridge is an open-core Feature Flag Management platform with product intelligence.
This repo is the Go API server — the core of the self-hosted distribution.

- **Website:** https://flagbridge.io
- **Admin UI:** https://admin.flagbridge.io
- **Docs:** https://docs.flagbridge.io
- **Org:** https://github.com/flagbridge
- **License:** Apache 2.0 (CE), Commercial (Pro)

## Stack

Go 1.22+, Chi router, pgx (PostgreSQL via Supabase), sqlc, goose, slog.
Cache: in-memory (CacheProvider interface). Auth: Supabase Auth SaaS / bcrypt+JWT self-hosted (AuthProvider interface).
Deploy: Fly.io (region gru). CI: GitHub Actions.

## Structure

```
cmd/server/          # Entry point
internal/
├── flag/            # Flag CRUD + evaluation engine
├── project/         # Project domain
├── environment/     # Environment management
├── evaluation/      # Flag evaluation engine
├── targeting/       # Targeting rules
├── testing/         # Sessions, overrides, cleanup
├── webhook/         # Registration, delivery, retry
├── apikey/          # API key management
├── audit/           # Audit logging
├── sse/             # Real-time SSE streaming
├── auth/            # AuthProvider interface (JWT + API keys)
├── middleware/       # HTTP middleware
├── cache/           # CacheProvider interface (in-memory)
├── config/          # Env config
├── database/        # pgx PostgreSQL client
└── testutil/        # Test helpers
migrations/          # goose SQL migrations
openapi.yaml         # OpenAPI 3.1 spec
```

## API — 72 endpoints (40 CE + 32 Pro)

Keys: fb_sk_{eval|test|mgmt|full}_...
Resolution: session override > targeting > rollout (MurmurHash3) > env default > flag default.
Pro gating: ADR-001 dogfooding via _flagbridge internal project.

## Conventions

- snake_case DB, camelCase JSON, PascalCase Go
- Table-driven tests com testify
- Errors: `fmt.Errorf("context: %w", err)`
- Structured logging via slog

## Do NOT

- Change Pro gating (ADR-001) without approval
- Modify applied migrations — create new ones
- Add Redis — use CacheProvider in-memory
- Expose `_flagbridge` internal project publicly
- Add Go deps without justification
