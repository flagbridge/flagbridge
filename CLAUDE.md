# FlagBridge

Feature flags with product intelligence. Open source.

## Project Structure

- **Monorepo** with Go API + Next.js Admin UI + TypeScript SDK packages
- `apps/api/` — Go API server (chi, pgx, zerolog)
- `apps/web/` — Next.js 15 Admin UI (App Router, Tailwind, Radix UI, next-intl, TanStack Query)
- `packages/` — TypeScript packages (@flagbridge/sdk-node, sdk-react, plugin-sdk, openfeature-provider, cli, create-plugin)
- `docs/` — Documentation site (Docusaurus)

## Tech Stack

### API (Go)
- Router: chi
- Database: PostgreSQL via pgx
- Logging: zerolog
- Migrations: golang-migrate

### Web (TypeScript)
- Framework: Next.js 15 (App Router)
- Styling: Tailwind CSS v4 with design tokens from DESIGN.md
- Components: Radix UI primitives
- i18n: next-intl (en, pt-BR)
- State: Zustand (client), TanStack Query (server)
- Package manager: pnpm

### SDKs
- TypeScript packages use TypeScript project references
- Published under @flagbridge scope on npm

## Conventions

- Language: code and technical terms in English, comments can be bilingual
- Commits: Conventional Commits (feat:, fix:, chore:, docs:, etc.)
- Go: standard project layout, internal/ for private packages, pkg/ for shared
- TypeScript: strict mode, ESM-first
- CSS: follow DESIGN.md design system tokens — no arbitrary colors or spacing
- i18n: all user-facing strings go through next-intl, keys in en.json and pt-BR.json
- API: RESTful, JSON responses, proper HTTP status codes
- Database: all schema changes via numbered migrations in apps/api/migrations/

## Design System

See DESIGN.md for the complete design system specification. Key rules:
- Dark theme, Electric Blue (#3B82F6) primary
- No visible borders for section separation — use tonal surface transitions
- Inter font for everything
- Developer-first, left-aligned, generous whitespace

## Development

```bash
# Start local environment
docker compose up -d

# API
cd apps/api && go run ./cmd/server

# Web
cd apps/web && pnpm dev
```

## Important

- Never use #000000 — darkest color is surface_container_lowest (#060e20)
- Never use opaque borders for layout separation
- Always validate design choices against DESIGN.md tokens
- All new API endpoints need corresponding migration if they touch the schema
