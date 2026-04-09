.PHONY: dev dev-down api web test test-integration lint migrate

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

test-integration:
	docker compose -f apps/api/docker-compose.test.yml up -d --wait
	cd apps/api && DATABASE_URL="postgres://flagbridge_test:flagbridge_test@localhost:5433/flagbridge_test?sslmode=disable" \
		go test -tags=integration -count=1 -timeout=120s ./...
	docker compose -f apps/api/docker-compose.test.yml down

lint:
	cd apps/api && golangci-lint run
	cd apps/web && pnpm lint

migrate:
	cd apps/api && go run ./cmd/migrate up
