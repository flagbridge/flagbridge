# ADR-0006: AI Provider interface interna (não passthrough)

- **Status:** Accepted
- **Date:** 2026-04-17
- **Deciders:** Bridge (CTO), Flag (CPO), Gabriel
- **Scope:** Sprint 5 — AI Layer (task `86e0xbf0k` — provider router Anthropic + OpenAI)

## Context

Sprint 5 implementa 2 providers LLM (Anthropic, OpenAI). Sprint 7 adicionará Ollama (self-hosted LLM — ADR para CPO: diferencial único vs LaunchDarkly/Unleash/GrowthBook). Pro roadmap: Gemini, Mistral, Bedrock.

Cada provider tem:
- Formato de request diferente (headers, body schema)
- Formato de streaming diferente (ADR-0005 resolve só o formato **nosso**, não deles)
- Autenticação diferente (Bearer token, api-key header, custom)
- Error shapes incompatíveis

Decisão-chave: **nossa API expõe shape interno abstrato, ou faz passthrough de um provider nativo?**

## Decision

Interface interna abstrata:

```go
package ai

type Provider interface {
    // Complete executa um prompt e retorna chunks via canal.
    // Fecha o canal ao terminar (sucesso ou erro).
    // Respeita ctx.Done() — cancelamento propaga pro HTTP upstream.
    Complete(ctx context.Context, req CompleteRequest) (<-chan Event, error)
}

type CompleteRequest struct {
    Model     string           // ex: "claude-sonnet-4-6", "gpt-4o", "llama3.2"
    System    string           // system prompt — pode ser ""
    Messages  []Message        // histórico da conversa
    MaxTokens int              // 0 = default do provider
}

type Message struct {
    Role    string // "user" | "assistant"
    Content string
}

type Event struct {
    Delta string         // token text — pode ser ""
    Usage *UsageMetrics  // non-nil apenas no último evento bem-sucedido
    Error *ProviderError // non-nil se erro
}

type UsageMetrics struct {
    InputTokens  int
    OutputTokens int
}

type ProviderError struct {
    Code    string // ex: "rate_limit", "invalid_key", "timeout", "provider_5xx"
    Message string // human-readable, NÃO inclui conteúdo do prompt
}
```

Zero passthrough. Cada provider faz:
1. Converter `CompleteRequest` → payload nativo
2. Fazer HTTP call (streaming)
3. Parse eventos nativos → emitir `Event` via canal

Sprint 5: implementar `AnthropicProvider` e `OpenAIProvider`. Interface preparada pra Ollama (Sprint 7) sem refactor.

**Registry simples** pra router escolher provider baseado em `ai_configs.provider`:

```go
type Registry struct {
    providers map[string]Provider
}

func (r *Registry) Get(name string) (Provider, error) { ... }
```

## Alternatives considered

| Opção | Rejeição |
|---|---|
| **Passthrough** (expor shape Anthropic diretamente, com provider=openai só remapeando headers) | Lock-in imediato. Breaking change no Anthropic = breaking change na nossa API. Clientes tiram dependency em detalhes de implementação. |
| LangChain Go / AI SDK abstraction | Dependência pesada. Cobrem 50+ providers mas carregam muito código. Nosso escopo é 3 providers na feature, 5-6 max Pro. Custo > benefício. |
| Protocol Buffers + múltiplos endpoints (`/ai/anthropic`, `/ai/openai`) | Força cliente a conhecer diferenças. Contradiz ADR-0005 (endpoint único). |
| OpenAI shape como "nosso" shape (adapter reverso nos outros providers) | Próximo da opção escolhida mas torna Anthropic-como-provider mais verboso (forced chat.completions shape onde Anthropic tem Messages API). |

## Consequences

### Positive
- Adicionar provider novo = implementar interface + registrar. Zero mudança no handler/middleware/rate limiter.
- Migração entre providers transparente pro cliente — user troca `provider=anthropic` → `provider=openai` no config, zero mudança de API
- Testável em isolamento — `FakeProvider` implementa interface, middleware/handler testáveis sem HTTP real
- Error normalization: cliente vê `{"error": {"code": "rate_limit"}}` independente do provider (OpenAI 429 vs Anthropic `type: rate_limit_error`)
- Ollama plug-in em Sprint 7 é **~4-6 horas** (API é OpenAI-compatible — literalmente `OpenAIProvider` com `baseURL` diferente)

### Negative / Trade-offs
- Features novas de provider específico (ex: Anthropic extended thinking, OpenAI function calling) precisam ser abstraídas ou escondidas. **Aceito** pra V1 — escopo Cmd+K não exige essas features.
- ~200-300 LOC de adapter code por provider (dois na Sprint 5, mais depois)
- Breaking changes na nossa `Event` struct quebram TODOS os providers. Mitigado: `Event.Delta + Usage + Error` é shape mínimo e estável — adições futuras via campos opcionais, não mudanças.

## Implementation notes

- Localização: `internal/ai/` (interface), `internal/ai/provider/anthropic/`, `internal/ai/provider/openai/`
- Sprint 5 providers:
  - `anthropic`: Messages API streaming (`POST /v1/messages` com `stream: true`)
  - `openai`: Chat Completions streaming (`POST /v1/chat/completions` com `stream: true`)
- Parse SSE dos providers usa `bufio.Scanner` com `ScanLines` + separação manual de frames (diferentes shapes)
- Contexto `context.Context` obrigatório — cancelamento de HTTP downstream (Sprint 5 requirement via ADR-0005)
- Teste comum `internal/ai/provider/testcommon/`: uma suite de table tests que **qualquer** implementação de Provider deve passar (mock server que simula streaming/errors/disconnects)
- `FakeProvider` em `internal/ai/testing/`: implementação in-memory pra tests upstream (handler, middleware)
- Rate limiter (ADR-0003) chama provider via `Complete` → cuida do counter antes de abrir stream
- **Timeout**: `http.Client` com `Timeout: 0` pro streaming (nunca corta pelo cliente HTTP) mas respeita `ctx.Done()` via `context.WithTimeout` quando configurado por projeto

## Related

- ClickUp task Sprint 5: [`86e0xbf0k` provider router Anthropic + OpenAI](https://app.clickup.com/t/86e0xbf0k)
- ClickUp task Sprint 7: [`86e0yq95h` Ollama provider](https://app.clickup.com/t/86e0yq95h)
- ADR relacionados: [0002 PromptContext](./0002-ai-prompt-context-struct.md) (input), [0005 streaming format](./0005-ai-streaming-sse-always.md) (output shape)
- Pattern reference: `internal/auth/AuthProvider`, `internal/cache/CacheProvider` — mesma filosofia de interface + múltiplos implementations
