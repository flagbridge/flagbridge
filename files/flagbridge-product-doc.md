# FlagBridge — Product Document

> **"The bridge between feature flags and product strategy."**

**Version:** 1.0  
**Author:** Gabriel Gripp  
**Date:** March 2026  
**Status:** Draft / MVP Planning

---

## 1. Visão Geral

FlagBridge é uma plataforma open-core de Feature Flag Management que vai além do toggle on/off: conecta feature flags a planejamento de produto, observabilidade técnica e lifecycle management. O diferencial é ser **product-first com infra sólida** — cada flag tem contexto de negócio, métricas de impacto e regras de ciclo de vida claras.

### 1.1 O Problema

O mercado atual de feature flags é dominado por ferramentas **infra-first**:

- **Engenharia** cria flags sem contexto de produto — ninguém sabe *por que* uma flag existe
- **Produto** não tem visibilidade sobre o estado das flags nem sobre seu impacto real
- **Flags zumbis** se acumulam (100% ON há meses, nunca removidas), gerando dívida técnica
- **Nenhuma ferramenta** conecta nativamente uma flag a uma hipótese de produto, experimento ou OKR
- Ferramentas enterprise como LaunchDarkly cobram contratos medianos de ~$72k/ano (dados Vendr, 2026), inacessíveis para a maioria das empresas

### 1.2 A Solução

FlagBridge oferece:

1. **Feature Flag Management** — create, toggle, evaluate, rollout strategies
2. **Product Context Cards** — cada flag atrelada a hipótese, owner, métricas de sucesso, prazo
3. **Technical Dashboard** — adoption rate, error rate por variante, latency impact, stale flag detection
4. **Lifecycle Automation** — alertas de cleanup, archival automático, tracking de dívida técnica
5. **OpenFeature Compatible** — provider oficial para o padrão CNCF OpenFeature, zero vendor lock-in

### 1.3 Modelo Open-Core

| Camada | Distribuição | Preço |
|--------|-------------|-------|
| **Community Edition** | Open source (Apache 2.0) | Grátis |
| **Pro Edition** | Plugin self-hosted OU SaaS | $X/mês |
| **Enterprise** | SaaS managed + suporte | Custom |

**A sacada do plugin self-hosted:** o cliente já roda o FlagBridge CE, compra a licença Pro, faz `docker pull flagbridge/pro-plugin && docker compose up -d` e reinicia. Sem migração, sem trocar de cloud, sem downtime. O plugin injeta os módulos Pro no mesmo deployment.

---

## 2. Análise Competitiva

### 2.1 Landscape

| Feature | FlagBridge | Unleash | LaunchDarkly | Flagsmith | PostHog FF |
|---------|-----------|---------|-------------|-----------|-----------|
| Open Source Core | ✅ Apache 2.0 | ✅ Apache 2.0 | ❌ | ✅ BSD 3 | ✅ MIT |
| Self-hosted | ✅ | ✅ | ❌ | ✅ | ✅ |
| SaaS hosted | ✅ | ✅ | ✅ | ✅ | ✅ |
| Plugin upgrade (zero-migration) | ✅ | ❌ | ❌ | ❌ | ❌ |
| Product Context (hipótese, owner, OKR) | ✅ | ❌ | ❌ | ❌ | Parcial |
| Technical Dashboard (observabilidade) | ✅ Pro | Básico | ✅ (caro) | Básico | ✅ |
| Lifecycle/Cleanup Automation | ✅ Pro | ❌ | Parcial | ❌ | ❌ |
| OpenFeature Provider | ✅ | ✅ | ✅ | ✅ | ✅ |
| Preço entrada (pago) | ~$29/mês | $80/mês (5 seats) | $120/mês+ | $45/mês | $0 (bundled) |

### 2.2 Posicionamento

FlagBridge ocupa o espaço entre:

- **Unleash** (bom open-source, fraco em produto/analytics)
- **LaunchDarkly** (poderoso mas caro e SaaS-only)
- **PostHog** (bom em analytics mas FF é secundário)

**Tagline de posicionamento:** *"Feature flags with product intelligence. Open source. Self-hosted or cloud. Your data, your rules."*

### 2.3 Gaps exploráveis nos concorrentes

- **Unleash:** Usuários reclamam da falta de múltiplos projetos na versão open-source e da ausência de integração com analytics
- **LaunchDarkly:** Modelo de pricing baseado em service connections + MAU gera custos imprevisíveis; contratos medianos de $72k/ano afastam startups e scale-ups
- **Flagsmith:** Não oferece real-time flag sync na versão gratuita
- **Todos:** Nenhum oferece "Product Context Card" nativo — o vínculo entre flag e estratégia de produto é sempre externo (Jira, Notion, etc.)

---

## 3. Definição de MVP — Community vs Pro

### 3.1 Community Edition (Open Source)

O CE deve ser **genuinamente útil** sozinho — não uma versão castrada que força upgrade. Isso gera comunidade, contribuições e confiança.

#### Core Features (CE)

**Flag Management**
- CRUD de feature flags (boolean, string, number, JSON)
- Environments (dev, staging, production) — ilimitados
- Projetos ilimitados
- Targeting rules básicas (user ID, percentage rollout, custom attributes)
- Kill switch instantâneo
- API REST completa + Admin UI (React)
- Audit log básico (quem mudou o quê, quando)

**SDK & Integrations**
- SDKs oficiais: Go, Node.js/TypeScript, React (client-side), Python
- OpenFeature Provider oficial
- Webhooks para integrações
- CLI (`flagbridge`) para automação

**Product Context (básico)**
- Campo de descrição e owner para cada flag
- Tags/labels customizáveis
- Link externo (para Jira, Linear, Notion, etc.)
- Status da flag: `planning` → `active` → `rolled-out` → `archived`

**Dashboard (básico)**
- Lista de flags por projeto/environment
- Status overview (quantas ativas, stale, archived)
- Flag age tracking (quando foi criada, última alteração)

#### Tech Specs (CE)

- **Backend:** Go (API server) — single binary, fácil de deployar
- **Frontend:** React + TypeScript (Admin UI)
- **Database:** PostgreSQL
- **Cache/Eval:** In-memory com sync via SSE/WebSocket
- **Deploy:** Docker image oficial, Docker Compose, Helm chart
- **Auth:** Local users + API keys

### 3.2 Pro Edition (Plugin)

O Pro é o que transforma FlagBridge de "mais um FF tool" em **plataforma de product intelligence**.

#### Pro Features

**Product Context Cards (avançado)**
- Campos estruturados: hipótese, métricas de sucesso (KPIs), critérios de go/no-go
- Workflow de decisão: experiment → analyze → decide (rollout/rollback/iterate)
- Integração bidirecional com Linear, Jira, Notion (sync automático de status)
- Timeline visual do lifecycle de cada flag
- Vínculo flag ↔ OKR/initiative

**Technical Dashboard (avançado)**
- Adoption rate por flag (% de requests que avaliam a flag)
- Error rate por variante (correlação flag × errors)
- Latency impact por flag (antes/depois do rollout)
- Stale flag detection com score de "urgência de cleanup"
- SDK version tracking (quais apps usam qual versão)
- Real-time evaluation stream

**Lifecycle & Cleanup Automation**
- Regras configuráveis: "se flag está 100% ON há X dias → notificar owner"
- Alertas via Slack, email, webhook
- Cleanup suggestions com link para o código (integração GitHub/GitLab)
- Auto-archive de flags expiradas
- "Technical debt score" por projeto

**Governance & Security**
- SSO (SAML, OIDC)
- RBAC granular (por projeto, environment, flag)
- Change requests com approval workflow
- Audit log detalhado com diff de mudanças
- Scheduled flag changes (agendar rollout para data/hora)

**Experimentation (básico)**
- A/B split por percentage
- Métricas de conversão básicas (via webhook de eventos)
- Export de dados para ferramentas de analytics

#### Modelo de entrega do Plugin

```
# Self-hosted: upgrade de CE para Pro
docker pull flagbridge/pro-plugin:latest
docker compose up -d --force-recreate

# O plugin detecta a licença e ativa os módulos Pro
# Dados existentes são preservados — zero migration
```

O plugin funciona como um **módulo que estende o CE**:

- Registra rotas adicionais na API
- Adiciona componentes React no Admin UI
- Cria tabelas extras no PostgreSQL (migration automática)
- Valida licença via license key (offline-first, com heartbeat opcional)

### 3.3 Enterprise (SaaS Managed)

- Tudo do Pro
- Hosting gerenciado pela FlagBridge (multi-tenant ou dedicated)
- SLA com uptime guarantee
- Suporte prioritário (Slack, email, call)
- Custom integrations
- SOC 2 compliance (roadmap)
- Data residency options

---

## 4. Arquitetura Técnica

### 4.1 Overview

```
┌─────────────────────────────────────────────────────┐
│                    Client Apps                       │
│  (React, Node, Go, Python — via SDKs)               │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ┌─────────────┐    ┌──────────────────────────┐   │
│  │  FlagBridge  │    │   FlagBridge Edge        │   │
│  │  SDK         │───▶│   (optional, for scale)  │   │
│  └─────────────┘    └──────────┬───────────────┘   │
│                                │                    │
├────────────────────────────────┼────────────────────┤
│                                ▼                    │
│  ┌──────────────────────────────────────────────┐  │
│  │           FlagBridge API Server (Go)          │  │
│  │                                              │  │
│  │  ┌──────────┐  ┌───────────┐  ┌──────────┐  │  │
│  │  │ Flag     │  │ Product   │  │ Dashboard │  │  │
│  │  │ Eval     │  │ Context   │  │ & Metrics │  │  │
│  │  │ Engine   │  │ Module    │  │ Module    │  │  │
│  │  └──────────┘  └───────────┘  └──────────┘  │  │
│  │          ▲ CE        ▲ Pro        ▲ Pro      │  │
│  └──────────┼───────────┼────────────┼──────────┘  │
│             │           │            │              │
│  ┌──────────┴───────────┴────────────┴──────────┐  │
│  │              PostgreSQL                       │  │
│  │  flags, rules, evaluations, product_cards,   │  │
│  │  metrics, audit_log                          │  │
│  └──────────────────────────────────────────────┘  │
│                                                     │
│  ┌──────────────────────────────────────────────┐  │
│  │           Admin UI (React + TS)               │  │
│  │  Dashboard, Flag Manager, Product Cards,     │  │
│  │  Settings, Lifecycle Timeline                │  │
│  └──────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────┘
```

### 4.2 Stack Técnica

| Componente | Tecnologia | Justificativa |
|-----------|-----------|--------------|
| API Server | **Go** | Single binary, alta performance, baixo consumo de memória, ideal para self-hosted |
| Admin UI | **React + TypeScript** | Ecossistema maduro, fácil de contribuir, component library rica |
| Database | **PostgreSQL** | Confiável, JSONB para targeting rules, extensível |
| Cache | **In-memory (Go)** + Redis (opcional) | Avaliação de flags em < 1ms; Redis para clusters |
| Edge Proxy | **Go** (FlagBridge Edge) | Para alto volume — cacheia flags perto do client |
| SDK Transport | **SSE** (Server-Sent Events) | Real-time flag updates sem complexidade de WebSocket |
| Containerização | **Docker** | Single Dockerfile multi-stage para tudo |
| Orquestração | **Docker Compose** (padrão) / Helm chart (K8s) | 90% dos self-hosted users usam Compose |

### 4.3 Data Model (Core)

```sql
-- Projetos
CREATE TABLE projects (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    slug        VARCHAR(255) NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW()
);

-- Environments
CREATE TABLE environments (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id  UUID REFERENCES projects(id),
    name        VARCHAR(100) NOT NULL, -- dev, staging, production
    slug        VARCHAR(100) NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(project_id, slug)
);

-- Feature Flags
CREATE TABLE flags (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id      UUID REFERENCES projects(id),
    key             VARCHAR(255) NOT NULL, -- e.g. "checkout-v2"
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    flag_type       VARCHAR(20) NOT NULL DEFAULT 'boolean', -- boolean, string, number, json
    owner_id        UUID REFERENCES users(id),
    status          VARCHAR(20) DEFAULT 'planning', -- planning, active, rolled_out, archived
    tags            TEXT[], -- PostgreSQL array
    external_link   TEXT, -- Link para Jira, Linear, etc.
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(project_id, key)
);

-- Flag States (por environment)
CREATE TABLE flag_states (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flag_id         UUID REFERENCES flags(id),
    environment_id  UUID REFERENCES environments(id),
    enabled         BOOLEAN DEFAULT FALSE,
    default_value   JSONB, -- Valor default quando não há targeting match
    targeting_rules JSONB, -- Array de rules com conditions + value
    rollout_pct     INTEGER DEFAULT 0, -- 0-100
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_by      UUID REFERENCES users(id),
    UNIQUE(flag_id, environment_id)
);

-- Audit Log
CREATE TABLE audit_log (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type VARCHAR(50) NOT NULL, -- flag, project, environment
    entity_id   UUID NOT NULL,
    action      VARCHAR(50) NOT NULL, -- created, updated, toggled, archived
    actor_id    UUID REFERENCES users(id),
    diff        JSONB, -- O que mudou
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

-- === PRO TABLES === --

-- Product Context Cards (Pro)
CREATE TABLE product_cards (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flag_id         UUID REFERENCES flags(id) UNIQUE,
    hypothesis      TEXT,
    success_metrics JSONB, -- [{name: "conversion_rate", target: ">5%"}]
    go_nogo_criteria TEXT,
    decision        VARCHAR(20), -- NULL, rollout, rollback, iterate
    decision_date   TIMESTAMPTZ,
    decided_by      UUID REFERENCES users(id),
    okr_link        TEXT,
    initiative_link TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Flag Evaluations Metrics (Pro)
CREATE TABLE flag_evaluations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flag_id         UUID REFERENCES flags(id),
    environment_id  UUID REFERENCES environments(id),
    variant         VARCHAR(255),
    count           BIGINT DEFAULT 0,
    error_count     BIGINT DEFAULT 0,
    bucket_start    TIMESTAMPTZ NOT NULL, -- Bucketed per hour
    bucket_end      TIMESTAMPTZ NOT NULL,
    UNIQUE(flag_id, environment_id, variant, bucket_start)
);

-- Lifecycle Rules (Pro)
CREATE TABLE lifecycle_rules (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id      UUID REFERENCES projects(id),
    rule_type       VARCHAR(50) NOT NULL, -- stale_detection, auto_archive, cleanup_reminder
    conditions      JSONB NOT NULL, -- {"days_since_100pct": 30}
    actions         JSONB NOT NULL, -- {"notify": ["slack", "email"], "auto_archive": true}
    enabled         BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
```

### 4.4 SDK Design

Os SDKs seguem o padrão OpenFeature, com um layer proprietário para features específicas do FlagBridge.

```typescript
// @flagbridge/sdk-node — Exemplo de uso

import { FlagBridge } from '@flagbridge/sdk-node';

const fb = new FlagBridge({
  serverUrl: 'https://flags.mycompany.com',
  apiKey: 'fb_sk_...',
  environment: 'production',
});

// Avaliação simples (boolean)
const isEnabled = await fb.isEnabled('checkout-v2', {
  userId: 'user_123',
  attributes: { plan: 'pro', country: 'BR' },
});

// Avaliação com variante (string)
const variant = await fb.getString('homepage-hero', 'default', {
  userId: 'user_123',
});

// OpenFeature compatible
import { OpenFeature } from '@openfeature/server-sdk';
import { FlagBridgeProvider } from '@flagbridge/openfeature-provider';

OpenFeature.setProvider(new FlagBridgeProvider({
  serverUrl: 'https://flags.mycompany.com',
  apiKey: 'fb_sk_...',
}));

const client = OpenFeature.getClient();
const value = await client.getBooleanValue('checkout-v2', false);
```

```go
// @flagbridge/sdk-go — Exemplo de uso

package main

import (
    "context"
    "fmt"
    flagbridge "github.com/flagbridge/sdk-go"
)

func main() {
    client, _ := flagbridge.NewClient(flagbridge.Config{
        ServerURL:   "https://flags.mycompany.com",
        APIKey:      "fb_sk_...",
        Environment: "production",
    })
    defer client.Close()

    ctx := flagbridge.NewEvalContext("user_123", map[string]interface{}{
        "plan":    "pro",
        "country": "BR",
    })

    enabled, _ := client.IsEnabled(context.Background(), "checkout-v2", ctx)
    fmt.Println("checkout-v2 enabled:", enabled)
}
```

### 4.5 API Design (principais endpoints)

```
# Flag Management
GET    /api/v1/projects/:project/flags
POST   /api/v1/projects/:project/flags
GET    /api/v1/projects/:project/flags/:key
PATCH  /api/v1/projects/:project/flags/:key
DELETE /api/v1/projects/:project/flags/:key

# Flag Evaluation (SDK endpoint)
POST   /api/v1/evaluate
POST   /api/v1/evaluate/batch

# Flag State (per environment)
GET    /api/v1/projects/:project/flags/:key/states/:env
PUT    /api/v1/projects/:project/flags/:key/states/:env

# Product Cards (Pro)
GET    /api/v1/projects/:project/flags/:key/product-card
PUT    /api/v1/projects/:project/flags/:key/product-card

# Metrics (Pro)
GET    /api/v1/projects/:project/flags/:key/metrics
GET    /api/v1/projects/:project/dashboard/overview

# Lifecycle (Pro)
GET    /api/v1/projects/:project/lifecycle/stale
GET    /api/v1/projects/:project/lifecycle/rules
POST   /api/v1/projects/:project/lifecycle/rules

# Server-Sent Events (real-time updates)
GET    /api/v1/sse/:environment

# Admin
GET    /api/v1/audit-log
GET    /api/v1/health
```

### 4.6 Infraestrutura de Deployment

**Self-hosted mínimo (CE):**
```yaml
# docker-compose.yml
version: '3.8'
services:
  flagbridge:
    image: flagbridge/flagbridge:latest
    ports:
      - "8080:8080"  # API + Admin UI
    environment:
      - DATABASE_URL=postgres://fb:fb@db:5432/flagbridge
      - FB_API_KEY_SALT=your-secret-salt
    depends_on:
      - db

  db:
    image: postgres:16-alpine
    volumes:
      - pgdata:/var/lib/postgresql/data
    environment:
      - POSTGRES_USER=fb
      - POSTGRES_PASSWORD=fb
      - POSTGRES_DB=flagbridge

volumes:
  pgdata:
```

**Self-hosted com Pro Plugin:**
```yaml
# docker-compose.pro.yml (override)
version: '3.8'
services:
  flagbridge:
    image: flagbridge/flagbridge-pro:latest  # Imagem com plugin integrado
    environment:
      - FB_LICENSE_KEY=fb_lic_...  # Ativa módulos Pro
      # Tudo mais igual ao CE
```

**Upgrade path:** `docker compose -f docker-compose.yml -f docker-compose.pro.yml up -d`

---

## 5. Go-to-Market

### 5.1 Público-alvo

**Primário:** Startups e scale-ups (50-500 devs) que:
- Já usam feature flags (Unleash CE, in-house) mas sentem falta de observabilidade e contexto de produto
- Não podem pagar LaunchDarkly ($72k/ano median)
- Valorizam self-hosted e data sovereignty

**Secundário:** Times de produto em empresas maiores que:
- Querem visibilidade sobre o estado das flags sem depender de engenharia
- Precisam vincular flags a OKRs e roadmap

### 5.2 Pricing (proposta inicial)

| Plano | Preço | Target |
|-------|-------|--------|
| **Community** | Grátis (forever) | Devs individuais, times pequenos, avaliação |
| **Pro** (self-hosted plugin) | $29/mês (até 10 seats) / $79/mês (até 50 seats) / $149/mês (unlimited) | Scale-ups, times médios |
| **Pro** (SaaS hosted) | $49/mês (até 10 seats) / $99/mês (até 50 seats) / $199/mês (unlimited) | Times que não querem gerenciar infra |
| **Enterprise** | Custom (a partir de $500/mês) | Grandes empresas, compliance, SLA |

**Modelo flat-rate por seats** — sem surpresas de MAU ou service connections. Esse é um diferencial direto contra LaunchDarkly.

**Nota:** Self-hosted é mais barato que SaaS porque o cliente paga a infra. Isso incentiva self-hosted (menos custo de hosting pra nós) e atrai o público que valoriza data sovereignty.

### 5.3 Canais de Aquisição

**Developer-first growth:**
1. **GitHub** — repo open-source bem documentado, README atraente, contributing guide
2. **Hacker News / Reddit / Dev.to** — launch posts: "We built an open-source feature flag platform with product intelligence"
3. **Comparações honestas** — blog posts tipo "FlagBridge vs Unleash vs LaunchDarkly" com tabela real
4. **Docker Hub** — imagem oficial otimizada, one-liner de setup
5. **OpenFeature ecosystem** — provider oficial listado no site do OpenFeature
6. **YouTube / blog** — tutoriais: "Como implementar feature flags com contexto de produto em 10 minutos"

**Product-led growth:**
- CE é genuinamente útil → time adota → cresce → precisa de observabilidade/governance → upgrade Pro
- Banner discreto no Admin UI: "Unlock product dashboards and lifecycle automation — Upgrade to Pro"
- Trial de 14 dias do Pro (self-hosted e SaaS)

**Community building:**
- Discord/Slack community
- Roadmap público (Linear ou GitHub Projects)
- Changelog público
- Contributors recognized (CONTRIBUTORS.md, swag)

### 5.4 Métricas de Sucesso (primeiro ano)

| Métrica | Meta 6 meses | Meta 12 meses |
|---------|-------------|---------------|
| GitHub stars | 500 | 2,000 |
| Docker pulls | 5,000 | 25,000 |
| CE installs ativos | 200 | 1,000 |
| Pro paying customers | 10 | 50 |
| MRR | $500 | $3,000 |
| Community members (Discord) | 100 | 500 |

### 5.5 Roadmap (quarters)

**Q2 2026 — Foundation**
- [ ] Repo público no GitHub
- [ ] Core API server (Go) com flag management + eval engine
- [ ] Admin UI (React) com CRUD de flags
- [ ] SDK Node.js + SDK Go
- [ ] Docker image + Docker Compose
- [ ] Docs site (docs.flagbridge.io)
- [ ] OpenFeature Provider (Node.js)

**Q3 2026 — Community Launch**
- [ ] Product Context Cards (versão CE básica)
- [ ] CLI (`flagbridge`)
- [ ] SDK React (client-side) + SDK Python
- [ ] Helm chart para Kubernetes
- [ ] Webhooks
- [ ] Launch no Hacker News + Product Hunt + Reddit

**Q4 2026 — Pro Launch**
- [ ] Pro Plugin architecture
- [ ] Technical Dashboard (metrics, observability)
- [ ] Product Cards avançado (hipótese, KPIs, workflow de decisão)
- [ ] Lifecycle automation (stale detection, cleanup alerts)
- [ ] SSO (OIDC)
- [ ] Billing system + license key management
- [ ] Integrações: Slack, Linear, Jira

**Q1 2027 — Scale**
- [ ] FlagBridge Edge (proxy para alto volume)
- [ ] A/B experimentation básico
- [ ] RBAC granular
- [ ] Change requests + approval workflow
- [ ] GitHub/GitLab integration (code cleanup suggestions)
- [ ] SaaS managed offering

---

## 6. Diferenciadores-chave (resumo)

1. **Product Context nativo** — nenhum concorrente faz isso. A ponte entre eng e produto está no core, não num plugin.
2. **Plugin zero-migration** — upgrade de CE para Pro sem trocar de deployment, database ou cloud. Um `docker pull` e restart.
3. **Flat-rate pricing** — sem surpresas de MAU ou service connections. Preço previsível.
4. **OpenFeature-first** — zero vendor lock-in desde o dia 1.
5. **Single binary Go** — deploy trivial, baixo consumo de recursos, ideal para self-hosted.
6. **Lifecycle automation** — o único que trata flag cleanup como first-class feature, não afterthought.

---

## 7. Riscos e Mitigações

| Risco | Impacto | Mitigação |
|-------|---------|-----------|
| Mercado saturado de FF tools | Alto | Diferenciação clara via Product Context — ninguém faz isso |
| Unleash adiciona features similares | Médio | Velocidade de execução + foco product-first vs infra-first |
| Dificuldade de monetizar open-source | Alto | CE genuinamente útil mas Pro resolve dores reais (observabilidade, governance) |
| Complexidade de manter múltiplos SDKs | Médio | Começar com 4 SDKs (Go, Node, React, Python); OpenFeature reduz superfície |
| Tempo limitado (side project) | Alto | MVP enxuto; priorizar core → launch → iterate com feedback |

---

## 8. Next Steps Imediatos

1. **Registrar domínio** — flagbridge.io e flagbridge.dev
2. **Criar org no GitHub** — github.com/flagbridge
3. **Reservar no npm** — @flagbridge (scope)
4. **Iniciar repo** — `flagbridge/flagbridge` (monorepo: api + ui + docs)
5. **Prototipar** — API server Go com flag CRUD + eval endpoint
6. **Design** — Wireframe do Admin UI (dashboard + flag detail + product card)

---

*FlagBridge — Bridging the gap between feature flags and product strategy.*
