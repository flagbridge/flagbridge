<p  align="center">
  <img width="512" height="140" alt="textual-dark" src="https://github.com/user-attachments/assets/67852e6f-5815-482a-ad7a-a38b74eb5cd3" />
</p>

![textual-light-2](https://github.com/user-attachments/assets/96583357-5695-4a23-9700-6e1750966112)

# FlagBridge

**Feature flags with product intelligence. Open source.**

[![Website](https://img.shields.io/badge/website-flagbridge.io-blue)](https://flagbridge.io)
[![License](https://img.shields.io/badge/license-Apache%202.0-green)](LICENSE)

FlagBridge is an open-core feature flag management platform that combines powerful flag evaluation with product intelligence — helping teams not just toggle features, but understand their impact.

---

## Features

- **Flag Management** — Create, organize, and manage feature flags with a modern admin UI
- **Targeting Rules** — Target users by attributes, segments, and percentages
- **Real-time Evaluation** — Sub-millisecond flag evaluation with SSE streaming updates
- **Product Intelligence** — Understand how flags affect user behavior and metrics
- **Plugin System** — Extend FlagBridge with community and custom plugins
- **OpenFeature Compatible** — Works with the OpenFeature standard
- **Multi-language SDKs** — Node.js, Go, Python, and React SDKs
- **Self-hosted** — Run on your own infrastructure with Docker or Kubernetes

## Quick Start

```bash
# Clone the repository
git clone https://github.com/flagbridge/flagbridge.git
cd flagbridge

# Start the local environment
docker compose up -d

# API available at http://localhost:8080
# Admin UI available at http://localhost:3000
```

## Architecture

```
flagbridge/
├── apps/
│   ├── api/          # Go API server
│   └── web/          # Next.js Admin UI
├── packages/
│   ├── sdk-node/     # Node.js/TypeScript SDK
│   ├── sdk-react/    # React SDK
│   ├── plugin-sdk/   # Plugin development kit
│   └── ...           # More packages
└── docs/             # Documentation
```

## Tech Stack

| Component | Technology |
|-----------|------------|
| API | Go, Chi, PostgreSQL, Redis |
| Admin UI | Next.js 15, Tailwind CSS, Radix UI |
| SDKs | TypeScript, Go, Python |
| Infrastructure | Docker, Kubernetes (Helm) |

## Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

Apache 2.0 — see [LICENSE](LICENSE) for details.

---

# FlagBridge (pt)

**Feature flags com inteligencia de produto. Open source.**

FlagBridge e uma plataforma open-core de gerenciamento de feature flags que combina avaliacao de flags com inteligencia de produto — ajudando times nao apenas a alternar features, mas a entender seu impacto.

## Inicio Rapido

```bash
git clone https://github.com/flagbridge/flagbridge.git
cd flagbridge
docker compose up -d
```

## Contribuindo

Contribuicoes sao bem-vindas! Veja [CONTRIBUTING.md](CONTRIBUTING.md) para mais detalhes.

## Licenca

Apache 2.0 — veja [LICENSE](LICENSE).
