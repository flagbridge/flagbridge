# Architecture Decision Records (ADRs)

Decisões arquiteturais do FlagBridge API (Go). Cada ADR documenta **uma** decisão com contexto, alternativas consideradas, e consequências.

## Formato

Seguimos [MADR](https://adr.github.io/madr/) simplificado:

```markdown
# ADR-000N: Título curto

- Status: Proposed | Accepted | Deprecated | Superseded by ADR-000X
- Date: YYYY-MM-DD
- Deciders: quem decidiu
- Scope: sprint/feature afetada

## Context
Por que essa decisão é necessária agora?

## Decision
O que foi decidido?

## Alternatives considered
Outras opções + motivo de rejeição.

## Consequences
Positivas, negativas/trade-offs.

## Implementation notes
Hints concretos pra quem for codar.

## Related
Links pra ClickUp, outras ADRs, docs.
```

## Regras

1. **Uma decisão por ADR.** Se precisa decidir duas coisas, dois ADRs.
2. **Numeração sequencial, nunca reuso.** Superseded ≠ deletado.
3. **Status `Accepted`** é pré-requisito pra codar a feature relacionada.
4. **Rascunhar ANTES de codar.** Não pular planejamento estruturante.
5. **Curto:** ~1 página. Se precisa de 3 páginas, a decisão não está madura.

## Índice

### Aceitas

| ADR | Título | Scope | Data |
|---|---|---|---|
| 0001 | Pro Gating via Dogfooding (`_flagbridge` internal project) | Infra / Pricing | (pendente formalização — existe em `CONTEXT.md`) |
| [0002](./0002-ai-prompt-context-struct.md) | AI PromptContext struct — versioned, dual render (XML/JSON), token budget | Sprint 5 — AI Layer | 2026-04-17 |
| [0003](./0003-ai-rate-limiter-postgres.md) | AI rate limiter — Postgres counter (CE 100/mo) | Sprint 5 — AI Layer | 2026-04-17 |
| [0004](./0004-ai-key-encryption-aes-gcm.md) | AI provider key encryption — AES-256-GCM + env var master key | Sprint 5 — AI Layer | 2026-04-17 |
| [0005](./0005-ai-streaming-sse-always.md) | AI completions streaming — SSE sempre, OpenAI-compatible event format | Sprint 5 — AI Layer | 2026-04-17 |
| [0006](./0006-ai-provider-interface.md) | AI Provider interface interna (não passthrough) | Sprint 5 — AI Layer | 2026-04-17 |

### Deprecated / Superseded

_(nenhum ainda)_

## Propondo uma nova ADR

1. Copie o template acima
2. Numeração = próximo número livre (zero-padded: `0007`)
3. Nome do arquivo: `NNNN-slug-kebab.md`
4. Abra PR com status `Proposed`
5. Merge quando consensus → atualize status pra `Accepted` e adicione ao índice
