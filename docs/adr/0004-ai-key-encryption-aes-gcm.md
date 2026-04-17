# ADR-0004: AI provider key encryption — AES-256-GCM + env var master key

- **Status:** Accepted
- **Date:** 2026-04-17
- **Deciders:** Bridge (CTO), Flag (CPO), Gabriel
- **Scope:** Sprint 5 — AI Layer (task `86e0xbf10` — AI config encrypted storage)

## Context

Cada project configura seu próprio AI provider (Anthropic, OpenAI) com API key real da organização. Essas keys são **bearer tokens altamente sensíveis** — vazamento = consumo ilimitado na conta do cliente.

Storage em plaintext no Postgres é inaceitável:
- Self-host: dump de banco → key vaza
- SaaS: breach do DB → keys de TODOS os clientes expostas
- Compliance: clientes enterprise (perfil ICP do FlagBridge) exigem encryption-at-rest pra provider credentials

Self-host complica: não podemos assumir AWS KMS / GCP KMS disponível. Solução precisa funcionar **100% self-contained** via Docker Compose.

CPO confirmou explicitamente: **cortar esta task não é opção** — "security é posicionamento, não feature".

## Decision

**AES-256-GCM** (authenticated encryption) com master key via env var `FB_AI_MASTER_KEY`:

```
FB_AI_MASTER_KEY=<32 bytes base64-encoded>
```

Gerada uma vez por instância. Sem `FB_AI_MASTER_KEY` configurada → feature AI inteiramente desabilitada (log warning no boot, endpoints `/api/v1/ai/*` retornam 503).

Schema:

```sql
CREATE TABLE ai_configs (
    project_id     uuid PRIMARY KEY REFERENCES projects(id) ON DELETE CASCADE,
    provider       text NOT NULL, -- 'anthropic' | 'openai' | 'ollama' (Sprint 7)
    model          text NOT NULL,
    encrypted_key  bytea NOT NULL,
    nonce          bytea NOT NULL, -- 12 bytes, random por write
    endpoint_url   text,           -- opcional, pra Ollama self-hosted futuro
    updated_at     timestamptz NOT NULL DEFAULT now()
);
```

**Requisitos de implementação:**
- Nonce: 12 bytes random (`crypto/rand`) **por cada write** — nunca reusar
- AEAD: `crypto/aes` + `cipher.NewGCM` (stdlib Go, sem deps externas)
- Decrypt: verifica authentication tag — fail-close se tag inválida (possível tampering)
- Constant-time compare: via `crypto/subtle.ConstantTimeCompare` em qualquer check de equality
- API: endpoints never retornam a key decrypted em response bodies (nem admin/logs). Só "key_configured: true/false" + `updated_at`.

**Rotação de master key**: fora de escopo Sprint 5. Doc técnico + script `scripts/rotate-ai-master-key.sh` ficam pra backlog Sprint 6+.

## Alternatives considered

| Opção | Rejeição |
|---|---|
| Plaintext no DB | Trivialmente inaceitável (compliance, segurança, posicionamento) |
| AWS KMS / GCP KMS | Auto-exclui self-hosted sem cloud AWS/GCP. Nosso ICP é self-host + on-prem + Fly.io misto — não todos tem cloud. |
| HashiCorp Vault | Adiciona dep operacional gigante (deploy + auth + rotação) pra CE. Overkill. |
| bcrypt/argon2 | Não servem — são **hashes** one-way. Precisamos **encryption** (roundtrip decrypt no uso). |
| AES-CBC | Não tem authentication. Tampering detectado? Não. GCM é a escolha moderna. |
| ChaCha20-Poly1305 | Equivalente de segurança a AES-GCM. AES-GCM escolhido por hardware acceleration AES-NI disponível em todo CPU moderno. |

## Consequences

### Positive
- Zero deps externas além de stdlib Go
- Self-host funciona out-of-the-box (só `FB_AI_MASTER_KEY` em docker-compose env)
- Authentication embutida (GCM detecta tampering) — defesa em profundidade
- Nonce único por write → mesmo key reescrita gera ciphertext diferente (IND-CPA security)

### Negative / Trade-offs
- Master key em env var é **single point of failure**. Se vazar, todas as keys do DB podem ser decrypted. Mitigações:
  - Doc explícita: master key vai em secret manager (Fly secrets, Docker secrets, 1Password CLI)
  - Rotação documentada como "re-encrypt em background job" quando chegar Sprint 6
- Backup do DB **sem** backup da master key = dados perdidos (by design)
- Sem envelope encryption (data key wrapping) — complexidade pra ganho marginal em escala atual

### Future (fora de escopo Sprint 5)
- `KeyProvider` interface (mesmo padrão de `AuthProvider`, `CacheProvider`) — permite KMS via adapter Pro em Sprint 8+
- Envelope encryption se/quando performance justificar

## Implementation notes

- Localização sugerida: `internal/ai/crypto/` (`encrypt.go`, `decrypt.go`, `key.go`)
- Migration: `migrations/009_ai_configs.sql`
- Teste obrigatório: encrypt → decrypt roundtrip com 3 plaintexts diferentes, assert todos distintos no `encrypted_key`+`nonce` pair
- Teste tampering: decrypt com ciphertext adulterado em 1 byte → erro (não silent success)
- Boot check: se `FB_AI_MASTER_KEY` ausente, log `slog.Warn("AI feature disabled: FB_AI_MASTER_KEY not configured")` e skip route registration
- Doc em `docs/self-hosting.md`: como gerar master key (`openssl rand -base64 32`)
- Audit log: registrar `ai_config.created`, `ai_config.updated`, `ai_config.deleted` (não o plaintext da key, só o fato)

## Related

- ClickUp task: [`86e0xbf10` — AI config encrypted storage](https://app.clickup.com/t/86e0xbf10)
- Depende desta ADR: [`86e0xbf27` — AI config UI (Sprint 6)](https://app.clickup.com/t/86e0xbf27)
- Pattern reference: `internal/auth/` (AuthProvider interface existente) — seguir padrão análogo se promover pra `KeyProvider` interface no futuro
