# ADR-0003: AI rate limiter — Postgres counter (CE 100/mo)

- **Status:** Accepted
- **Date:** 2026-04-17
- **Deciders:** Bridge (CTO), Gabriel
- **Scope:** Sprint 5 — AI Layer (task `86e0xbf0v` — streaming + rate limiter CE)

## Context

CE tier limita AI usage a **100 requests/mês por project**. Pro tier é ilimitado (gate via `pro.ai_unlimited` flag, dogfooding ADR-0001).

Requisitos:
- **Durável**: reinícios de servidor (Fly.io faz deploy multiple vezes por dia) não podem resetar contador
- **Preciso**: contagem mensal — nunca passar o limite, nunca subcontar
- **Barato**: CE é free; não pode adicionar custo significativo de infra só pra esse rate limit
- **Simples**: self-hosted precisa rodar com dependências existentes (Postgres via Supabase já obrigatório)

## Decision

Tabela Postgres dedicada com contador por (project, mês):

```sql
CREATE TABLE ai_usage_counters (
    project_id uuid NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    year_month text NOT NULL, -- 'YYYY-MM', ex: '2026-05'
    count int NOT NULL DEFAULT 0,
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (project_id, year_month)
);

CREATE INDEX idx_ai_usage_counters_updated_at ON ai_usage_counters(updated_at);
```

Incremento via UPSERT atômico no path do request:

```sql
INSERT INTO ai_usage_counters (project_id, year_month, count)
VALUES ($1, $2, 1)
ON CONFLICT (project_id, year_month)
DO UPDATE SET count = ai_usage_counters.count + 1, updated_at = now()
RETURNING count;
```

Se `count > 100` **e** project não tem flag `pro.ai_unlimited` ativa → retornar HTTP 429 com header `Retry-After: <segundos até o primeiro dia do próximo mês UTC>`.

Limite lido de env var `FB_AI_CE_MONTHLY_LIMIT` (default 100) — ajustável sem redeploy.

## Alternatives considered

| Opção | Rejeição |
|---|---|
| In-memory map com mutex | Fly.io faz reinício/scale out → contador perdido. Grátis infinito pra quem faz deploy toda hora. |
| Redis (incr + expire) | Adiciona dependência operacional nova (backup, auth, connection pool) pra uma feature mensal. Overhead não justificado. |
| Token bucket com refill contínuo | Complexidade desnecessária — limite mensal é literal mês-calendar, não "100 nos últimos 30 dias". |
| Contador no `projects` table com coluna `ai_usage_count` | Precisa zerar manualmente todo mês. Perde histórico (quantos meses ficou no limite?). |
| Incrementar via API externa (Upstash ratelimit) | Custo + dependência externa num code path que já é CE-free. |

## Consequences

### Positive
- Zero dependência nova. Postgres já obrigatório pra operação.
- Histórico preservado por mês — analytics futuro trivial (`SELECT * FROM ai_usage_counters WHERE year_month = '2026-05' ORDER BY count DESC`)
- UPSERT atômico — sem race condition mesmo com alta concorrência
- Reset mensal implícito: entrada do mês novo começa em 0 sem cron job

### Negative / Trade-offs
- +1 query UPSERT por AI request. Aceitável porque AI requests naturalmente tem throughput baixo (LLM calls custam e são user-initiated).
- Contador pode ficar ligeiramente estourado em alta concorrência extrema (projeto fazendo 10 AI calls simultâneas no request 100). **Aceito**: passar de 100 pra 102 não muda o mundo, e UPSERT serializa mesmo em Postgres.
- Admin UI futura pra ver uso (Sprint 8+) precisa query adicional. OK.

## Implementation notes

- Localização sugerida: `internal/ai/ratelimit/` no repo Go
- Middleware `requireAIQuota`: roda **antes** de forward pro provider. Ordem: `auth → requireAIConfig → requireAIQuota → handler`.
- Teste crítico: 101º request retorna 429 com `Retry-After` válido (segundos até `next_month.first_day.utc()`)
- Migration: `migrations/008_ai_usage_counters.sql`
- Observabilidade: logar `project_id, year_month, count, limit` em cada 429 para dashboards futuros

## Related

- ClickUp task: [`86e0xbf0v` — streaming + rate limiter CE](https://app.clickup.com/t/86e0xbf0v)
- ADR relacionados: [0001 Pro Gating Dogfooding](../../CONTEXT.md#pro-gating--adr-001-dogfooding) (como checar `pro.ai_unlimited`), [0005 Streaming SSE](./0005-ai-streaming-sse-always.md)
