.PHONY: dev dev-down run test test-integration lint migrate

dev:
	docker compose up -d

dev-down:
	docker compose down

run:
	go run ./cmd/server

test:
	go test ./...

test-integration:
	docker compose -f docker-compose.test.yml up -d --wait
	DATABASE_URL="postgres://flagbridge_test:flagbridge_test@localhost:5433/flagbridge_test?sslmode=disable" \
		go test -tags=integration -count=1 -timeout=120s ./...
	docker compose -f docker-compose.test.yml down

lint:
	golangci-lint run

migrate:
	go run ./cmd/migrate up
