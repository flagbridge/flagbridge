.PHONY: dev dev-down api web test lint migrate

dev:
	docker compose up -d

dev-down:
	docker compose down

api:
	cd apps/api && go run ./cmd/server

web:
	cd apps/web && pnpm dev

test:
	cd apps/api && go test ./...
	cd apps/web && pnpm test

lint:
	cd apps/api && golangci-lint run
	cd apps/web && pnpm lint

migrate:
	cd apps/api && go run ./cmd/migrate up
