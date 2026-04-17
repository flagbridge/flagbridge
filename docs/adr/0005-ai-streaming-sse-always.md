# ADR-0005: AI completions streaming — SSE sempre, OpenAI-compatible event format

- **Status:** Accepted
- **Date:** 2026-04-17
- **Deciders:** Bridge (CTO), Gabriel
- **Scope:** Sprint 5 — AI Layer (tasks `86e0xbf0v` streaming + `86e0xbf1m` Cmd+K UI)

## Context

LLM completions são **inerentemente streaming** (geração token-by-token). UX responsiva (Cmd+K prompt bar) depende de mostrar tokens conforme chegam — usuário vê resposta "aparecendo" em vez de esperar 5-15s paradão.

Duas vias possíveis na API:
1. Endpoint sempre streaming (SSE)
2. Query param `?stream=true|false` — dois modos

Também: cada provider tem seu próprio event format. Anthropic: `message_start`, `content_block_delta`, etc. OpenAI: `delta.content` nested. Abstração é necessária — mas qual shape adotar como "nossa"?

## Decision

**Endpoint único sempre streaming:**

```
POST /api/v1/ai/completions
Response: text/event-stream
```

Nunca retorna JSON body. Sem query param `?stream=`.

**Formato de eventos — compatível com OpenAI Chat Completions streaming**:

```
data: {"delta":"Hello","usage":null}

data: {"delta":" world","usage":null}

data: {"delta":"","usage":{"input_tokens":42,"output_tokens":12}}

data: [DONE]

```

Campos:
- `delta` (string): próximo fragmento de texto. Pode ser `""` no evento final com metadata.
- `usage` (object|null): métricas de token. Presente apenas no último evento com conteúdo.
- `[DONE]` sentinel: termina stream.

**Providers que não fazem streaming nativo** (cenário hipotético): o router faz buffer da resposta completa e emite **um chunk único** com todo o texto, depois `[DONE]`. Cliente não precisa saber.

**Erros durante streaming**: evento especial `data: {"error":{"code":"...","message":"..."}}` antes de `[DONE]`. Cliente detecta `.error` na payload.

## Alternatives considered

| Opção | Rejeição |
|---|---|
| Query param `?stream=true|false` | Dobra surface area cliente. Cmd+K UI nunca vai querer non-streaming. Complexidade gratuita. |
| JSON response sempre | UX horrível — usuário espera 5-15s sem feedback |
| Formato Anthropic nativo (`message_start`/`content_block_delta`) | Nosso cliente precisa entender apenas um formato. OpenAI shape é universalmente conhecido (ecossistema, tutorials, debugger tools) |
| WebSocket | SSE é stateless (HTTP/1.1 long-lived GET) — mais simples, passa por proxies/CDN sem config. WebSocket exige handshake upgrade, sticky sessions, e pouco benefício pra unidirectional streaming. |
| Server-Sent Events custom com multiple `event:` types | Complexidade acima do necessário. Single `data:` event type com payload discriminável é suficiente. |
| gRPC streaming | Exige cliente especializado. Cmd+K UI no browser usaria grpc-web (outro layer). Ganho zero. |

## Consequences

### Positive
- Cliente consome **uma única interface** (EventSource ou fetch ReadableStream) — sem branching
- OpenAI-compatible = tooling ecossistema (Vercel AI SDK, LangChain JS, etc.) funciona out-of-the-box
- SSE passa por proxies, Cloudflare, Fly.io edge sem config especial (HTTP/1.1 padrão)
- Stateless — reconnect = novo request, não precisa state recovery
- Browser-native (`EventSource` API)

### Negative / Trade-offs
- Cliente que quer **resposta síncrona inteira** (ex: script CLI futuro) precisa acumular chunks. OK — script simples (`join('')` dos deltas).
- Connection longa (15-60s típico) — precisa configurar Fly.io `[services.ports]` sem timeouts agressivos
- Curl para debug requer `-N` flag (no buffering). Doc obrigatória.
- Sem suporte nativo a compressão dentro do stream (gzip ok em HTTP, mas por-frame não). Aceito — payloads são texto curto.

## Implementation notes

- Localização: `internal/ai/stream/` para serializer + `internal/ai/router.go` para fan-out
- Handler deve:
  - `w.Header().Set("Content-Type", "text/event-stream")`
  - `w.Header().Set("Cache-Control", "no-cache")`
  - `w.Header().Set("Connection", "keep-alive")` (implícito mas explícito ajuda proxies)
  - `http.Flusher` cast no `ResponseWriter` — `flusher.Flush()` após cada write
- **Client disconnect**: monitorar `r.Context().Done()`. Se cancelado, cancelar context do provider também (evita leak de request upstream).
- Teste crítico com `-race`: client disconnect no meio do stream → sem goroutine leak, sem panic
- Teste MSW (admin): mock de SSE retornando 3 deltas + usage + `[DONE]`. Assert UI mostra texto progressivo + métricas no final.
- Fly.io config: garantir `hard_limit` e `soft_limit` de timeout suficientes (min 60s, recomendado 120s)
- **Cache**: nunca cachear streaming (cliente ou edge) — CDN headers `Cache-Control: no-store`

## Related

- ClickUp tasks: [`86e0xbf0v` streaming + rate limit](https://app.clickup.com/t/86e0xbf0v), [`86e0xbf1m` Cmd+K UI](https://app.clickup.com/t/86e0xbf1m)
- ADR relacionados: [0003 rate limiter](./0003-ai-rate-limiter-postgres.md) (429 precisa vir antes do stream abrir), [0006 Provider interface](./0006-ai-provider-interface.md) (implementations convertem native → nosso shape)
- Doc externa: [OpenAI Chat streaming format](https://platform.openai.com/docs/api-reference/chat-streaming), [MDN EventSource](https://developer.mozilla.org/en-US/docs/Web/API/EventSource)
