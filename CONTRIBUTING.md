# Contributing to FlagBridge

Thank you for your interest in contributing to FlagBridge!

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/flagbridge.git`
3. Create a branch: `git checkout -b feat/your-feature`
4. Make your changes
5. Run tests: `make test`
6. Commit using [Conventional Commits](https://www.conventionalcommits.org/): `git commit -m "feat: add new feature"`
7. Push and create a Pull Request

## Development Setup

```bash
# Start local environment
docker compose up -d

# Run API
cd apps/api && go run ./cmd/server

# Run Web
cd apps/web && pnpm dev
```

## Code Style

- **Go**: `gofmt` and `golangci-lint`
- **TypeScript**: ESLint + Prettier via project config
- **Commits**: Conventional Commits (feat, fix, chore, docs, refactor, test)

## Reporting Issues

Use [GitHub Issues](https://github.com/flagbridge/flagbridge/issues) with the provided templates.

## License

By contributing, you agree that your contributions will be licensed under the Apache 2.0 License.
