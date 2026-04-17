# ADR-0002: AI PromptContext struct — versioned, dual render, token budget

- **Status:** Accepted
- **Date:** 2026-04-17
- **Deciders:** Bridge (CTO), Gabriel
- **Scope:** Sprint 5 — AI Layer (task `86e0xbf0d` — context builder)

## Context

O AI proxy precisa serializar estado relevante do projeto (flags, rules, product cards, audit entries recentes, role do usuário) em formato consumível por LLMs de providers diferentes (Anthropic, OpenAI, futuramente Ollama).

Sem estrutura formal:
- Cada renderer criaria seu próprio formato ad-hoc → desync entre prompts
- Impossível versionar context payload (breaking changes silenciosas)
- Cliente grande com centenas de flags estoura context window → custos de API explodem
- Anthropic rende melhor com XML (best practice oficial deles). OpenAI prefere JSON. Sem abstração, precisa reimplementar o mesmo contexto duas vezes.

## Decision

Definir struct Go tipada `PromptContext` com campo `Version` explícito:

```go
type PromptContext struct {
    Version      string                `json:"version"`        // "1.0" — bump em breaking change
    Project      ProjectContext        `json:"project"`
    Flags        []FlagContext         `json:"flags"`          // truncado a 50, nota "truncated_at": N
    Rules        []RuleContext         `json:"rules,omitempty"`
    ProductCards []ProductCardContext  `json:"product_cards,omitempty"`
    RecentAudit  []AuditEntry          `json:"recent_audit,omitempty"` // max 20, sorted desc por created_at
    Role         string                `json:"role"`           // "engineer" | "product" | "viewer" | "admin"
}
```

**Renderização dual**: mesmo struct, dois renderers:
- `RenderXML(ctx PromptContext) string` — pra Anthropic
- `RenderJSON(ctx PromptContext) string` — pra OpenAI

Renderer escolhido pelo `Provider` (ver ADR-0006).

**Token budget estático** (Sprint 5 V1):
- Max 50 flags (sorted desc por `updated_at`)
- Max 20 audit entries (sorted desc por `created_at`)
- Sem product cards quando `role=engineer` (flag CE não tem cards de qualquer jeito)

Heurística adaptativa (tamanho real em tokens) → backlog Sprint 6+.

## Alternatives considered

| Opção | Rejeição |
|---|---|
| `map[string]any` | Falta type safety, schema evolution caótica, nenhum compilador ajuda em breaking change |
| JSON-only (sem XML) | Anthropic rende 15-20% melhor com XML em testes internos (consistente com doc oficial). Perdemos qualidade gratuita. |
| Passthrough: entregar raw DB rows pro LLM | Explode context window em 3 flags, vaza detalhes internos (UUIDs, timestamps irrelevantes), sem versionamento |
| Sem token budget (manda tudo) | Cliente com 500 flags gera prompt de 100k tokens → $1+/request só em input. Insustentável pra CE free. |

## Consequences

### Positive
- Type safety end-to-end (struct → renderer → provider)
- `Version` field permite adicionar campos sem breaking change em clientes que conheçam v1
- Token budget previsível e auditável (sabemos qual é o custo-teto por request)
- Renderer duplo preserva qualidade com cada provider

### Negative / Trade-offs
- Adicionar campo novo ao contexto exige tocar 2 renderers (XML + JSON) + bump de `Version`
- Truncation estática (50 flags) pode cortar flag importante. Mitigado por sort `updated_at DESC` (flags recentemente tocadas provavelmente relevantes).
- Renderer XML é código extra (~150 LOC). Custo pontual, não recorrente.

## Implementation notes

- Localização sugerida: `internal/ai/context/` no repo Go da API
- Fixtures em `internal/ai/context/testdata/` — snapshot tests comparam output XML/JSON render com fixture versionado
- Teste de budget: seed com 51 flags, assert `len(ctx.Flags) == 50` e presença de `truncated_at` nota
- Sanitização: remover campos `internal` (ex: `created_by_user_id`) antes de serializar — LLM não precisa e pode vazar PII

## Related

- ClickUp task: [`86e0xbf0d` — context builder](https://app.clickup.com/t/86e0xbf0d)
- ADR relacionados: [0006 Provider interface](./0006-ai-provider-interface.md) (quem consome o render)
- Doc externa: [Anthropic — Use XML tags to structure prompts](https://docs.anthropic.com/claude/docs/use-xml-tags)
