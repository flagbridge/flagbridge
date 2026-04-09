<p  align="center">
  <img width="512" height="140" alt="textual-dark" src="https://github.com/user-attachments/assets/67852e6f-5815-482a-ad7a-a38b74eb5cd3" />
</p>

# FlagBridge

**Feature flags with product intelligence. Open source.**

[![Website](https://img.shields.io/badge/website-flagbridge.io-blue)](https://flagbridge.io)
[![License](https://img.shields.io/badge/license-Apache%202.0-green)](LICENSE)

FlagBridge is an open-core feature flag management platform that combines powerful flag evaluation with product intelligence — helping teams not just toggle features, but understand their impact.

This repository contains the **Go API server**. The admin dashboard, SDKs, docs, and Helm charts live in separate repositories under the [flagbridge GitHub organization](https://github.com/flagbridge).

---

## Features

- **Flag Management** — Create, organize, and manage feature flags across projects and environments
- **Targeting Rules** — Target users by attributes, segments, and percentage rollouts (MurmurHash3)
- **Real-time Streaming** — SSE-based push updates so clients receive flag changes without polling
- **Test Sessions** — Override flag values per session for QA and automated testing workflows
- **Webhook Delivery** — Register webhooks for flag change events with automatic retry
- **Audit Logging** — Full audit trail for every flag and project mutation
- **API Key Scopes** — Fine-grained key scopes: `eval`, `test`, `mgmt`, and `full`
- **OpenFeature Compatible** — Works with the OpenFeature standard via the dedicated provider
- **Multi-language SDKs** — Node.js, Go, Python, and React SDKs available separately
- **Self-hosted** — Run on your own infrastructure with Docker or Kubernetes (Helm)

---

## Quick Start

```bash
# Clone the repository
git clone https://github.com/flagbridge/flagbridge.git
cd flagbridge

# Start the API and database
docker compose up -d

# API is available at http://localhost:8080
```

The compose file starts the Go API server and a PostgreSQL instance. Database migrations run automatically on startup via goose.

---

## Architecture

This repo contains a single application: the Go API server.

```
apps/
└── api/
    ├── cmd/server/main.go       # Entry point
    ├── internal/
    │   ├── flag/                # Flag CRUD + evaluation
    │   ├── project/             # Project management
    │   ├── environment/         # Environment management
    │   ├── evaluation/          # Flag evaluation engine
    │   ├── targeting/           # Targeting rules
    │   ├── testing/             # Test sessions & overrides
    │   ├── webhook/             # Webhook delivery
    │   ├── apikey/              # API key management
    │   ├── audit/               # Audit logging
    │   ├── sse/                 # Real-time SSE streaming
    │   ├── auth/                # Auth (JWT + API keys)
    │   ├── middleware/          # HTTP middleware
    │   ├── cache/               # CacheProvider (in-memory)
    │   ├── config/              # Environment config
    │   └── database/            # pgx PostgreSQL client
    ├── migrations/              # goose SQL migrations
    └── openapi.yaml             # OpenAPI 3.1 spec
```

The API exposes 39 endpoints. Evaluation resolution order: session override > targeting rules > percentage rollout > environment default > flag default.

---

## Tech Stack

| Component | Technology |
|---|---|
| Language | Go 1.22+ |
| HTTP Router | [Chi](https://github.com/go-chi/chi) |
| Database | PostgreSQL via [pgx](https://github.com/jackc/pgx) |
| Migrations | [goose](https://github.com/pressly/goose) |
| Logging | slog (structured) |
| Cache | In-memory CacheProvider |
| Auth | JWT + API key (bcrypt) |

---

## Related Repositories

| Repository | Description |
|---|---|
| [flagbridge/admin](https://github.com/flagbridge/admin) | Admin dashboard UI (Next.js, Tailwind) |
| [flagbridge/docs](https://github.com/flagbridge/docs) | Documentation site (Docusaurus) |
| [flagbridge/sdk-node](https://github.com/flagbridge/sdk-node) | Node.js SDK (`@flagbridge/sdk-node`) |
| [flagbridge/sdk-react](https://github.com/flagbridge/sdk-react) | React SDK (`@flagbridge/sdk-react`) |
| [flagbridge/sdk-go](https://github.com/flagbridge/sdk-go) | Go SDK |
| [flagbridge/sdk-python](https://github.com/flagbridge/sdk-python) | Python SDK |
| [flagbridge/openfeature-provider](https://github.com/flagbridge/openfeature-provider) | OpenFeature provider |
| [flagbridge/plugin-sdk](https://github.com/flagbridge/plugin-sdk) | Plugin development kit (`@flagbridge/plugin-sdk`) |
| [flagbridge/helm-charts](https://github.com/flagbridge/helm-charts) | Helm charts for Kubernetes deployment |

---

## Contributing

Contributions are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on branch naming, commit format, and how to run the test suite locally.

## License

Apache 2.0 — see [LICENSE](LICENSE) for details.
