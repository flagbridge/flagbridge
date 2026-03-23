FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /flagbridge-api ./cmd/server

FROM alpine:3.21

RUN apk --no-cache add ca-certificates
COPY --from=builder /flagbridge-api /usr/local/bin/flagbridge-api

EXPOSE 8080
ENTRYPOINT ["flagbridge-api"]
