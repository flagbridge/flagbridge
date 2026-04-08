# CLAUDE.md — FlagBridge API (Go)

> Copiar pra: flagbridge/flagbridge/CLAUDE.md

## O que é

FlagBridge — open-core feature flag management com product intelligence.
Este repo é a API Go central.

## Stack

Go 1.22+, Chi router, pgx (PostgreSQL via Supabase), sqlc, goose, slog.
Cache: in-memory (CacheProvider interface). Auth: Supabase Auth SaaS / bcrypt+JWT self-hosted (AuthProvider interface).
Deploy: Fly.io (region gru). CI: GitHub Actions.

## Estrutura

```
internal/
├── flag/        # CRUD + evaluation engine
├── testing/     # Sessions, overrides, cleanup
├── webhook/     # Registration, delivery, retry
├── project/     # Project domain
├── auth/        # AuthProvider interface
├── plugin/      # Plugin runtime
├── integration/ # Managed integrations (Pro)
└── api/         # HTTP handlers (Chi)
pkg/
├── database/    # pgx client
├── config/      # Env config
├── cache/       # CacheProvider interface
└── middleware/   # HTTP middleware
migrations/      # goose SQL
```

## API — 72 endpoints (40 CE + 32 Pro)

Keys: fb_sk_{eval|test|mgmt|full}_...
Resolution: session override > targeting > rollout (MurmurHash3) > env default > flag default.
Pro gating: ADR-001 dogfooding via _flagbridge internal project.

## Convenções

- snake_case DB, camelCase JSON, PascalCase Go
- Table-driven tests com testify
- Errors: `fmt.Errorf("context: %w", err)`
- Structured logging via slog

## NÃO faça

- Não mude Pro gating (ADR-001) sem aprovação
- Não modifique migrations aplicadas — crie novas
- Não adicione Redis — use CacheProvider in-memory
- Não exponha _flagbridge no público
- Não adicione deps Go sem justificativa
