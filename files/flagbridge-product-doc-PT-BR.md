# FlagBridge — Documento de Produto

> **"A ponte entre feature flags e estratégia de produto."**

**Versão:** 3.0  
**Autor:** Gabriel Gripp  
**Data:** Março 2026  
**Status:** Draft / Planejamento de MVP

---

## Índice

1. [Visão Geral](#1-visão-geral)
2. [Análise Competitiva](#2-análise-competitiva)
3. [Definição de MVP — Community vs Pro](#3-definição-de-mvp--community-vs-pro)
4. [Arquitetura Técnica](#4-arquitetura-técnica)
5. [Painel Admin — Documentação Técnica e de Produto](#5-painel-admin--documentação-técnica-e-de-produto)
6. [Sistema de Plugins & SDK](#6-sistema-de-plugins--sdk)
7. [Marketplace de Plugins](#7-marketplace-de-plugins)
8. [Ecossistema de Integrações](#8-ecossistema-de-integrações)
9. [Go-to-Market](#9-go-to-market)
10. [Roadmap](#10-roadmap)
11. [Riscos e Mitigações](#11-riscos-e-mitigações)

---

## 1. Visão Geral

FlagBridge é uma plataforma open-core de Feature Flag Management que vai além do toggle on/off: conecta feature flags a planejamento de produto, observabilidade técnica, lifecycle management e um rico ecossistema de integrações. O diferencial é ser **product-first com infra sólida** — cada flag tem contexto de negócio, métricas de impacto, regras de ciclo de vida claras e conectividade profunda com as ferramentas que os times já usam.

### 1.1 O Problema

O mercado atual de feature flags é dominado por ferramentas **infra-first**:

- **Engenharia** cria flags sem contexto de produto — ninguém sabe *por que* uma flag existe
- **Produto** não tem visibilidade sobre o estado das flags nem sobre seu impacto real
- **Flags zumbis** se acumulam (100% ON há meses, nunca removidas), gerando dívida técnica
- **Nenhuma ferramenta** conecta nativamente uma flag a uma hipótese de produto, experimento ou OKR
- Ferramentas enterprise como LaunchDarkly cobram contratos medianos de ~$72k/ano (dados Vendr, 2026), inacessíveis para a maioria das empresas
- **Nenhuma ferramenta** oferece um ecossistema de plugins — customização exige fork ou esperar o roadmap do vendor
- **Nenhuma ferramenta** conecta feature flags a plataformas de mensageria, analytics de eventos ou filas técnicas nativamente — times constroem código de integração customizado toda vez

### 1.2 A Solução

FlagBridge oferece:

1. **Feature Flag Management** — create, toggle, evaluate, rollout strategies
2. **Product Context Cards** — cada flag atrelada a hipótese, owner, métricas de sucesso, prazo
3. **Technical Dashboard** — adoption rate, error rate por variante, latency impact, stale flag detection
4. **Lifecycle Automation** — alertas de cleanup, archival automático, tracking de dívida técnica
5. **OpenFeature Compatible** — provider oficial para o padrão CNCF OpenFeature, zero vendor lock-in
6. **Ecossistema de Plugins & Marketplace** — arquitetura extensível onde devs podem criar, publicar e vender plugins
7. **Hub de Integrações** — conectores nativos para plataformas de mensageria (Resend, RD Station, SendGrid, Mailchimp), analytics de eventos (Mixpanel, Amplitude, Segment, PostHog) e filas técnicas (SQS, Kafka, RabbitMQ, NATS)

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
| Upgrade via plugin (zero-migration) | ✅ | ❌ | ❌ | ❌ | ❌ |
| Product Context (hipótese, owner, OKR) | ✅ | ❌ | ❌ | ❌ | Parcial |
| Technical Dashboard (observabilidade) | ✅ Pro | Básico | ✅ (caro) | Básico | ✅ |
| Lifecycle/Cleanup Automation | ✅ Pro | ❌ | Parcial | ❌ | ❌ |
| Ecossistema de Plugins & Marketplace | ✅ | ❌ | ❌ | ❌ | ❌ |
| Integrações de Mensageria (Resend, RD Station) | ✅ | ❌ | ❌ | ❌ | ❌ |
| Analytics de Eventos (Mixpanel, Amplitude) | ✅ | ❌ | Parcial | ❌ | Nativo |
| Filas Técnicas (SQS, Kafka) | ✅ Pro | ❌ | ❌ | ❌ | ❌ |
| OpenFeature Provider | ✅ | ✅ | ✅ | ✅ | ✅ |
| Preço entrada (pago) | ~$29/mês | $80/mês (5 seats) | $120/mês+ | $45/mês | $0 (bundled) |

### 2.2 Posicionamento

FlagBridge ocupa o espaço entre:

- **Unleash** — bom open-source, fraco em produto/analytics, sem ecossistema de plugins, sem integrações além de webhooks básicos
- **LaunchDarkly** — poderoso mas caro e SaaS-only, sem extensibilidade, suporte limitado a mensageria/filas
- **PostHog** — bom em analytics mas FF é secundário, sem produto standalone de FF, sem integrações de mensageria
- **Nenhum deles** oferece Plugin Marketplace ou bridges nativos para mensageria, analytics e infraestrutura de filas

**Tagline:** *"Feature flags com inteligência de produto. Open source. Extensível. Seus dados, suas regras."*

### 2.3 Gaps Exploráveis nos Concorrentes

- **Unleash:** Usuários reclamam da falta de múltiplos projetos na versão open-source e da ausência de integração com analytics; sem forma de disparar notificações em eventos de flags
- **LaunchDarkly:** Modelo de pricing baseado em service connections + MAU gera custos imprevisíveis; contratos medianos de $72k/ano afastam startups; integrações são gated por plano enterprise
- **Flagsmith:** Não oferece real-time flag sync na versão gratuita; sem suporte a event pipeline
- **Todos:** Nenhum oferece "Product Context Cards" nativamente; nenhum oferece extensibilidade via plugins; nenhum conecta a plataformas de mensageria ou filas técnicas

---

## 3. Definição de MVP — Community vs Pro

### 3.1 Community Edition (Open Source)

O CE deve ser **genuinamente útil** sozinho — não uma versão castrada que força upgrade. Isso gera comunidade, contribuições e confiança.

#### Features Core (CE)

**Flag Management**
- CRUD de feature flags (boolean, string, number, JSON)
- Environments ilimitados (dev, staging, production)
- Projetos ilimitados
- Targeting rules básicas (user ID, percentage rollout, custom attributes)
- Kill switch instantâneo
- API REST completa + Admin UI (Next.js)
- Audit log básico (quem mudou o quê, quando)

**SDK & Integrações**
- SDKs oficiais: Go, Node.js/TypeScript, React (client-side), Python
- OpenFeature Provider oficial
- Webhooks para integrações (HTTP POST genérico em eventos de flag)
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

**Sistema de Plugins (básico)**
- Runtime de plugins para carregar plugins da comunidade
- Plugin API com hooks para extensão de UI e API middleware
- Plugin CLI para scaffolding de novos plugins (`flagbridge plugin create`)

**Integrações (básico)**
- Webhook events para todos os eventos de lifecycle de flags
- Notificação Slack (plugin oficial CE)
- Emissor de eventos genérico (payload JSON em mudanças de flags)

### 3.2 Pro Edition (Plugin)

O Pro é o que transforma FlagBridge de "mais um FF tool" em **plataforma de product intelligence com conectividade enterprise-grade**.

#### Features Pro

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

**Integrações de Mensageria & Comunicação (Plugin Pro)**
- Conectores nativos para plataformas de email/mensageria
- Disparar campanhas ou mensagens transacionais baseadas em mudanças de estado de flags
- Veja Seção 8.1 para detalhes completos

**Integrações de Analytics & Rastreamento de Eventos (Plugin Pro)**
- Streamar dados de avaliação de flags para plataformas de analytics
- Correlacionar rollouts de flags com métricas de produto
- Veja Seção 8.2 para detalhes completos

**Integrações de Filas Técnicas (Plugin Pro)**
- Publicar eventos de flags em filas de mensagens para processamento downstream
- Consumir mensagens de filas para disparar mudanças de estado de flags
- Veja Seção 8.3 para detalhes completos

**Acesso ao Plugin Marketplace**
- Acesso a plugins premium no marketplace
- Plugin analytics (uso, impacto de performance)
- Suporte prioritário para plugins

**Experimentation (básico)**
- A/B split por percentage
- Métricas de conversão básicas (via webhook de eventos)
- Export de dados para ferramentas de analytics

#### Modelo de Entrega do Plugin

```bash
# Self-hosted: upgrade de CE para Pro
docker pull flagbridge/pro-plugin:latest
docker compose up -d --force-recreate

# O plugin detecta a licença e ativa os módulos Pro
# Dados existentes são preservados — zero migration
```

O plugin funciona como um **módulo que estende o CE**:
- Registra rotas adicionais na API
- Adiciona componentes React no Admin UI Next.js
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
- Opção de white-label do Admin UI
- Instância dedicada do marketplace de plugins
- Desenvolvimento de conectores customizados de filas e mensageria

---

## 4. Arquitetura Técnica

### 4.1 Visão Geral da Arquitetura

```
┌──────────────────────────────────────────────────────────────────────┐
│                      Aplicações Cliente                               │
│             (React, Node, Go, Python — via SDKs)                      │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌──────────────┐       ┌──────────────────────────────────┐        │
│  │  FlagBridge   │       │   FlagBridge Edge                 │        │
│  │  SDK          │──────▶│   (opcional, para escala)         │        │
│  └──────────────┘       └──────────────┬─────────────────┘        │
│                                        │                            │
├────────────────────────────────────────┼────────────────────────────┤
│                                        ▼                            │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │              FlagBridge API Server (Go)                       │   │
│  │                                                              │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌───────────────┐  │   │
│  │  │ Flag     │ │ Product  │ │Dashboard │ │ Plugin        │  │   │
│  │  │ Eval     │ │ Context  │ │& Métricas│ │ Runtime       │  │   │
│  │  │ Engine   │ │ Module   │ │ Module   │ │ Engine        │  │   │
│  │  └──────────┘ └──────────┘ └──────────┘ └───────────────┘  │   │
│  │    ▲ CE          ▲ Pro        ▲ Pro        ▲ CE+Pro        │   │
│  │                                                              │   │
│  │  ┌──────────────────────────────────────────────────────┐   │   │
│  │  │            Camada de Integrações (Pro)                 │   │   │
│  │  │                                                       │   │   │
│  │  │  ┌────────────┐  ┌──────────────┐  ┌──────────────┐  │   │   │
│  │  │  │ Conectores │  │ Conectores   │  │ Conectores   │  │   │   │
│  │  │  │ Mensageria │  │ Analytics &  │  │ Filas        │  │   │   │
│  │  │  │            │  │ Eventos      │  │ Técnicas     │  │   │   │
│  │  │  │ Resend     │  │ Mixpanel     │  │ SQS          │  │   │   │
│  │  │  │ RD Station │  │ Amplitude    │  │ Kafka        │  │   │   │
│  │  │  │ SendGrid   │  │ Segment      │  │ RabbitMQ     │  │   │   │
│  │  │  │ Mailchimp  │  │ PostHog      │  │ NATS         │  │   │   │
│  │  │  │ Brevo      │  │ GA4          │  │ Redis Streams│  │   │   │
│  │  │  │ Customer.io│  │ Rudderstack  │  │ GCP Pub/Sub  │  │   │   │
│  │  │  └────────────┘  └──────────────┘  └──────────────┘  │   │   │
│  │  └──────────────────────────────────────────────────────┘   │   │
│  └──────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│  ┌───────────────────────────┴──────────────────────────────────┐   │
│  │                       PostgreSQL                              │   │
│  │  flags, rules, evaluations, product_cards, metrics,          │   │
│  │  plugins, marketplace, integrations, audit_log               │   │
│  └──────────────────────────────────────────────────────────────┘   │
│                                                                      │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │              Admin UI (Next.js + TypeScript)                   │   │
│  │                                                               │   │
│  │  ┌────────────┐ ┌────────────┐ ┌───────────────────────────┐ │   │
│  │  │ Dashboard  │ │ Gerenciador│ │ Gerenciador de Plugins    │ │   │
│  │  │ & Analytics│ │ de Flags & │ │ & Marketplace             │ │   │
│  │  │            │ │ Product    │ │                           │ │   │
│  │  │            │ │ Cards      │ │                           │ │   │
│  │  └────────────┘ └────────────┘ └───────────────────────────┘ │   │
│  │                                                               │   │
│  │  ┌────────────┐ ┌────────────┐ ┌───────────────────────────┐ │   │
│  │  │ Config. &  │ │ Audit Log  │ │ Hub de Integrações        │ │   │
│  │  │ Gestão de  │ │ & Atividade│ │ (Mensageria, Analytics,   │ │   │
│  │  │ Time       │ │            │ │  Filas - configuração)    │ │   │
│  │  └────────────┘ └────────────┘ └───────────────────────────┘ │   │
│  └──────────────────────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────────────────────┘
```

### 4.2 Stack Técnica

| Componente | Tecnologia | Justificativa |
|-----------|-----------|--------------|
| API Server | **Go** | Binary único, alta performance, baixo consumo de memória, ideal para self-hosted |
| Admin UI | **Next.js 15 + TypeScript** | SSR/SSG para landing pages, App Router para admin SPA, i18n nativo para suporte bilíngue |
| Componentes Admin UI | **Next.js + Tailwind CSS + Radix UI** | Componentes acessíveis e composáveis com flexibilidade de design system |
| Estado Admin UI | **TanStack Query + Zustand** | Cache de server state + client state leve |
| Gráficos Admin UI | **Recharts** ou **Tremor** | Visualizações de dashboard otimizadas para Next.js |
| Landing Page | **Next.js (SSG)** | Mesmo codebase do admin, bilíngue (en/pt) com `next-intl`, otimizado para SEO |
| Banco de Dados | **PostgreSQL** | Confiável, JSONB para targeting rules, extensível |
| Cache | **In-memory (Go)** + Redis (opcional) | Avaliação de flags em < 1ms; Redis para clusters |
| Edge Proxy | **Go** (FlagBridge Edge) | Para alto volume — cacheia flags perto do client |
| Transporte SDK | **SSE** (Server-Sent Events) | Updates real-time sem complexidade de WebSocket |
| Plugin Runtime | **Go (backend)** + **Next.js (UI)** | Execução sandboxed com hook points definidos |
| Camada de Integrações | **Go** com padrão adapter | Interface unificada para conectores de mensageria, analytics e filas |
| Containerização | **Docker** | Dockerfile multi-stage |
| Orquestração | **Docker Compose** (padrão) / Helm (K8s) | 90% dos usuários self-hosted usam Compose |

### 4.2.1 Detalhe da Arquitetura Next.js

```
flagbridge-ui/
├── src/
│   ├── app/                          # Next.js App Router
│   │   ├── [locale]/                 # i18n: /en/... e /pt/...
│   │   │   ├── (marketing)/          # Landing page (SSG)
│   │   │   │   ├── page.tsx          # Homepage
│   │   │   │   ├── pricing/
│   │   │   │   ├── docs/
│   │   │   │   └── blog/
│   │   │   ├── (admin)/              # Painel admin (autenticado)
│   │   │   │   ├── dashboard/
│   │   │   │   ├── projects/
│   │   │   │   │   └── [projectSlug]/
│   │   │   │   │       ├── flags/
│   │   │   │   │       │   └── [flagKey]/
│   │   │   │   │       │       ├── page.tsx         # Detalhe da flag
│   │   │   │   │       │       ├── product-card/    # Contexto de produto
│   │   │   │   │       │       └── metrics/         # Dashboard técnico
│   │   │   │   │       ├── lifecycle/
│   │   │   │   │       └── settings/
│   │   │   │   ├── plugins/
│   │   │   │   │   ├── installed/       # Plugins instalados
│   │   │   │   │   ├── marketplace/     # Navegar & instalar
│   │   │   │   │   └── develop/         # Ferramentas de dev
│   │   │   │   ├── integrations/        # Hub de Integrações
│   │   │   │   │   ├── messaging/       # Resend, RD Station, etc.
│   │   │   │   │   ├── analytics/       # Mixpanel, Amplitude, etc.
│   │   │   │   │   ├── queues/          # SQS, Kafka, etc.
│   │   │   │   │   └── webhooks/        # Webhooks genéricos
│   │   │   │   ├── marketplace/
│   │   │   │   │   ├── browse/
│   │   │   │   │   ├── publish/
│   │   │   │   │   └── earnings/
│   │   │   │   ├── settings/
│   │   │   │   │   ├── team/
│   │   │   │   │   ├── billing/
│   │   │   │   │   ├── api-keys/
│   │   │   │   │   └── integrations/
│   │   │   │   └── audit-log/
│   │   │   └── (developer)/            # Portal do desenvolvedor
│   │   │       ├── docs/
│   │   │       ├── api-explorer/
│   │   │       └── sandbox/
│   │   └── api/                        # Rotas API Next.js (BFF)
│   │       ├── auth/
│   │       └── proxy/
│   ├── components/
│   │   ├── ui/                         # Design system (Radix + Tailwind)
│   │   ├── flags/
│   │   ├── dashboard/
│   │   ├── plugins/
│   │   ├── marketplace/
│   │   └── integrations/               # Componentes de config de integrações
│   ├── lib/
│   │   ├── api/
│   │   ├── i18n/
│   │   ├── plugin-host/
│   │   └── auth/
│   ├── messages/
│   │   ├── en.json
│   │   └── pt.json
│   └── styles/
│       └── globals.css
├── next.config.ts
├── tailwind.config.ts
└── package.json
```

**Decisões arquiteturais importantes:**

1. **Route Groups**: `(marketing)` para landing pages SSG, `(admin)` para app autenticado, `(developer)` para portal do desenvolvedor — tudo no mesmo app Next.js
2. **i18n via `next-intl`**: Todas as rotas prefixadas com locale (`/en/dashboard`, `/pt/dashboard`), páginas SSG totalmente bilíngues para SEO
3. **Padrão BFF**: Rotas API do Next.js atuam como Backend-for-Frontend, fazendo proxy para Go API com injeção de token de auth
4. **Plugin UI Host**: Plugins renderizam em componentes `<PluginSlot />` designados usando iframe sandboxed ou Module Federation
5. **Hub de Integrações**: Seção dedicada no admin para gerenciar todas as conexões com serviços externos com UI de config unificada

### 4.3 Modelo de Dados

```sql
-- ========================
-- TABELAS CORE (CE)
-- ========================

CREATE TABLE projects (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    slug        VARCHAR(255) NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE environments (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id  UUID REFERENCES projects(id),
    name        VARCHAR(100) NOT NULL,
    slug        VARCHAR(100) NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(project_id, slug)
);

CREATE TABLE flags (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id      UUID REFERENCES projects(id),
    key             VARCHAR(255) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    flag_type       VARCHAR(20) NOT NULL DEFAULT 'boolean',
    owner_id        UUID REFERENCES users(id),
    status          VARCHAR(20) DEFAULT 'planning',
    tags            TEXT[],
    external_link   TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(project_id, key)
);

CREATE TABLE flag_states (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flag_id         UUID REFERENCES flags(id),
    environment_id  UUID REFERENCES environments(id),
    enabled         BOOLEAN DEFAULT FALSE,
    default_value   JSONB,
    targeting_rules JSONB,
    rollout_pct     INTEGER DEFAULT 0,
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_by      UUID REFERENCES users(id),
    UNIQUE(flag_id, environment_id)
);

CREATE TABLE audit_log (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type VARCHAR(50) NOT NULL,
    entity_id   UUID NOT NULL,
    action      VARCHAR(50) NOT NULL,
    actor_id    UUID REFERENCES users(id),
    diff        JSONB,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

-- ========================
-- TABELAS DE PLUGINS (CE)
-- ========================

CREATE TABLE plugins (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug            VARCHAR(255) NOT NULL UNIQUE,
    name            VARCHAR(255) NOT NULL,
    version         VARCHAR(50) NOT NULL,
    author          VARCHAR(255),
    description     TEXT,
    source          VARCHAR(20) DEFAULT 'marketplace',
    manifest        JSONB NOT NULL,
    config          JSONB DEFAULT '{}',
    enabled         BOOLEAN DEFAULT TRUE,
    installed_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE plugin_hooks (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plugin_id   UUID REFERENCES plugins(id) ON DELETE CASCADE,
    hook_point  VARCHAR(100) NOT NULL,
    handler     VARCHAR(255) NOT NULL,
    priority    INTEGER DEFAULT 100,
    enabled     BOOLEAN DEFAULT TRUE
);

-- ========================
-- TABELAS PRO
-- ========================

CREATE TABLE product_cards (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flag_id         UUID REFERENCES flags(id) UNIQUE,
    hypothesis      TEXT,
    success_metrics JSONB,
    go_nogo_criteria TEXT,
    decision        VARCHAR(20),
    decision_date   TIMESTAMPTZ,
    decided_by      UUID REFERENCES users(id),
    okr_link        TEXT,
    initiative_link TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE flag_evaluations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flag_id         UUID REFERENCES flags(id),
    environment_id  UUID REFERENCES environments(id),
    variant         VARCHAR(255),
    count           BIGINT DEFAULT 0,
    error_count     BIGINT DEFAULT 0,
    bucket_start    TIMESTAMPTZ NOT NULL,
    bucket_end      TIMESTAMPTZ NOT NULL,
    UNIQUE(flag_id, environment_id, variant, bucket_start)
);

CREATE TABLE lifecycle_rules (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id      UUID REFERENCES projects(id),
    rule_type       VARCHAR(50) NOT NULL,
    conditions      JSONB NOT NULL,
    actions         JSONB NOT NULL,
    enabled         BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- ========================
-- TABELAS DE INTEGRAÇÕES (Pro)
-- ========================

CREATE TABLE integrations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id      UUID REFERENCES projects(id),
    type            VARCHAR(50) NOT NULL, -- messaging, analytics, queue
    provider        VARCHAR(50) NOT NULL, -- resend, mixpanel, sqs, etc.
    name            VARCHAR(255) NOT NULL,
    config          JSONB NOT NULL, -- Config específica do provider (secrets criptografados)
    enabled         BOOLEAN DEFAULT TRUE,
    health_status   VARCHAR(20) DEFAULT 'unknown', -- healthy, degraded, down, unknown
    last_health_check TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE integration_rules (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    integration_id  UUID REFERENCES integrations(id) ON DELETE CASCADE,
    trigger_event   VARCHAR(100) NOT NULL, -- flag.toggled, flag.rolledOut, etc.
    conditions      JSONB DEFAULT '{}',
    action_config   JSONB NOT NULL,
    enabled         BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE integration_events (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    integration_id  UUID REFERENCES integrations(id),
    rule_id         UUID REFERENCES integration_rules(id),
    event_type      VARCHAR(100) NOT NULL,
    payload         JSONB,
    status          VARCHAR(20) NOT NULL, -- sent, failed, retrying
    error_message   TEXT,
    attempts        INTEGER DEFAULT 1,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- ========================
-- TABELAS DE MARKETPLACE (Pro/Enterprise)
-- ========================

CREATE TABLE marketplace_listings (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plugin_slug     VARCHAR(255) NOT NULL UNIQUE,
    developer_id    UUID REFERENCES users(id),
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    long_description TEXT,
    category        VARCHAR(50),
    tags            TEXT[],
    icon_url        TEXT,
    screenshots     TEXT[],
    repository_url  TEXT,
    documentation_url TEXT,
    pricing_type    VARCHAR(20) DEFAULT 'free',
    price_cents     INTEGER DEFAULT 0,
    currency        VARCHAR(3) DEFAULT 'USD',
    status          VARCHAR(20) DEFAULT 'draft',
    downloads       INTEGER DEFAULT 0,
    avg_rating      DECIMAL(2,1) DEFAULT 0,
    review_count    INTEGER DEFAULT 0,
    published_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE marketplace_versions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    listing_id      UUID REFERENCES marketplace_listings(id),
    version         VARCHAR(50) NOT NULL,
    changelog       TEXT,
    min_flagbridge  VARCHAR(50),
    package_url     TEXT NOT NULL,
    package_hash    VARCHAR(64) NOT NULL,
    status          VARCHAR(20) DEFAULT 'pending',
    reviewed_by     UUID REFERENCES users(id),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(listing_id, version)
);

CREATE TABLE marketplace_reviews (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    listing_id      UUID REFERENCES marketplace_listings(id),
    user_id         UUID REFERENCES users(id),
    rating          INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 5),
    title           VARCHAR(255),
    body            TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(listing_id, user_id)
);

CREATE TABLE marketplace_purchases (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    listing_id      UUID REFERENCES marketplace_listings(id),
    buyer_id        UUID REFERENCES users(id),
    price_cents     INTEGER NOT NULL,
    currency        VARCHAR(3) NOT NULL,
    payment_provider VARCHAR(50),
    payment_id      VARCHAR(255),
    status          VARCHAR(20) DEFAULT 'completed',
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE developer_payouts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    developer_id    UUID REFERENCES users(id),
    amount_cents    INTEGER NOT NULL,
    currency        VARCHAR(3) NOT NULL,
    period_start    TIMESTAMPTZ NOT NULL,
    period_end      TIMESTAMPTZ NOT NULL,
    status          VARCHAR(20) DEFAULT 'pending',
    payout_ref      VARCHAR(255),
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
```

### 4.4 Design dos SDKs

Os SDKs seguem o padrão OpenFeature, com uma camada proprietária para features específicas do FlagBridge.

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

// Compatível com OpenFeature
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
    fmt.Println("checkout-v2 habilitado:", enabled)
}
```

### 4.5 Endpoints da API

```
# ========================
# GERENCIAMENTO DE FLAGS (CE)
# ========================
GET    /api/v1/projects/:project/flags
POST   /api/v1/projects/:project/flags
GET    /api/v1/projects/:project/flags/:key
PATCH  /api/v1/projects/:project/flags/:key
DELETE /api/v1/projects/:project/flags/:key

POST   /api/v1/evaluate
POST   /api/v1/evaluate/batch

GET    /api/v1/projects/:project/flags/:key/states/:env
PUT    /api/v1/projects/:project/flags/:key/states/:env

GET    /api/v1/sse/:environment

# ========================
# CONTEXTO DE PRODUTO (Pro)
# ========================
GET    /api/v1/projects/:project/flags/:key/product-card
PUT    /api/v1/projects/:project/flags/:key/product-card

# ========================
# MÉTRICAS & DASHBOARD (Pro)
# ========================
GET    /api/v1/projects/:project/flags/:key/metrics
GET    /api/v1/projects/:project/dashboard/overview
GET    /api/v1/projects/:project/lifecycle/stale

# ========================
# INTEGRAÇÕES (Pro)
# ========================
GET    /api/v1/projects/:project/integrations
POST   /api/v1/projects/:project/integrations
GET    /api/v1/projects/:project/integrations/:id
PATCH  /api/v1/projects/:project/integrations/:id
DELETE /api/v1/projects/:project/integrations/:id
POST   /api/v1/projects/:project/integrations/:id/test
GET    /api/v1/projects/:project/integrations/:id/health
GET    /api/v1/projects/:project/integrations/:id/events

GET    /api/v1/projects/:project/integrations/:id/rules
POST   /api/v1/projects/:project/integrations/:id/rules
PATCH  /api/v1/projects/:project/integrations/:id/rules/:ruleId
DELETE /api/v1/projects/:project/integrations/:id/rules/:ruleId

GET    /api/v1/integrations/providers
GET    /api/v1/integrations/providers/:provider

# ========================
# SISTEMA DE PLUGINS (CE+Pro)
# ========================
GET    /api/v1/plugins
POST   /api/v1/plugins/install
DELETE /api/v1/plugins/:slug
PATCH  /api/v1/plugins/:slug/config
GET    /api/v1/plugins/:slug/status

# ========================
# MARKETPLACE (Pro)
# ========================
GET    /api/v1/marketplace/listings
GET    /api/v1/marketplace/listings/:slug
POST   /api/v1/marketplace/listings
PATCH  /api/v1/marketplace/listings/:slug
POST   /api/v1/marketplace/listings/:slug/review
GET    /api/v1/marketplace/developer/earnings
POST   /api/v1/marketplace/purchase/:slug

# ========================
# ADMIN
# ========================
GET    /api/v1/audit-log
GET    /api/v1/health
```

### 4.6 Infraestrutura de Deploy

**Self-hosted (CE):**
```yaml
# docker-compose.yml
services:
  flagbridge:
    image: flagbridge/flagbridge:latest
    ports:
      - "8080:8080"   # Go API
      - "3000:3000"   # Next.js Admin UI
    environment:
      - DATABASE_URL=postgres://fb:fb@db:5432/flagbridge
      - FB_API_KEY_SALT=your-secret-salt
      - NEXT_PUBLIC_API_URL=http://localhost:8080
      - NEXT_PUBLIC_DEFAULT_LOCALE=pt
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
services:
  flagbridge:
    image: flagbridge/flagbridge-pro:latest
    environment:
      - FB_LICENSE_KEY=fb_lic_...
```

**Caminho de upgrade:** `docker compose -f docker-compose.yml -f docker-compose.pro.yml up -d`

---

## 5. Painel Admin — Documentação Técnica e de Produto

O FlagBridge Admin é uma aplicação Next.js completa que serve como hub central para gerenciamento de flags, inteligência de produto, gerenciamento de plugins, configuração de integrações e acesso ao marketplace.

### 5.1 Seções do Admin

#### Dashboard
- **Cards de overview**: Total de flags, flags ativas, flags obsoletas, flags por status
- **Feed de atividade**: Mudanças recentes em flags, deploys, instalações de plugins, eventos de integrações
- **Métricas de saúde** (Pro): Volume de avaliações, taxas de erro, contagem de conexões SDK, saúde das integrações
- **Score de dívida técnica** (Pro): Score agregado baseado em flags obsoletas, flags sem owner, flags sem product cards

#### Gerenciador de Flags
- **Lista de flags**: Filtrável por projeto, environment, status, tags, owner
- **Página de detalhe**: Toggle, editor visual de targeting rules, comparação entre environments
- **Product Context Card** (aba): Hipótese, métricas de sucesso, critérios go/no-go, histórico de decisões, link para OKR
- **Aba de Métricas** (Pro): Gráfico de adoção, gráfico de taxa de erro, impacto de latência, breakdown por variante
- **Aba de Lifecycle** (Pro): Timeline da criação ao archival, lembretes de cleanup, referências no código
- **Aba de Integrações** (Pro): Quais integrações estão conectadas a esta flag, histórico de eventos

#### Hub de Integrações
- **Mensageria**: Configurar conexões Resend, RD Station, SendGrid, Mailchimp, Brevo, Customer.io
- **Analytics & Eventos**: Configurar conexões Mixpanel, Amplitude, Segment, PostHog, GA4, Rudderstack
- **Filas Técnicas**: Configurar conexões SQS, Kafka, RabbitMQ, NATS, Redis Streams, GCP Pub/Sub
- **Webhooks**: Gerenciamento de webhooks genéricos
- **Por integração**: Config de conexão, regras de trigger, log de eventos, status de saúde, botão de teste

#### Gerenciador de Plugins
- **Plugins instalados**: Lista com toggle ativar/desativar, painel de configuração, status de saúde
- **Navegador do Marketplace**: Busca, filtro por categoria, instalação com um clique
- **Desenvolvimento de plugins**: Ambiente sandbox, logs, hot-reload para dev local

#### Configurações
- **Gerenciamento de time**: Convites, atribuição de roles (Admin, Editor, Viewer)
- **API keys**: Criação, rotação, revogação de chaves por environment
- **Billing** (Pro/SaaS): Gestão de plano, faturas, métricas de uso

#### Audit Log
- **Histórico completo**: Toda mudança de flag, toggle, ação de usuário, instalação de plugin, evento de integração
- **Filtrável**: Por usuário, tipo de ação, intervalo de data, entidade
- **Visualizador de diff** (Pro): Comparação lado a lado de mudanças no estado das flags

#### Portal do Desenvolvedor
- **Documentação do Plugin SDK**: Referência completa da API, guias, tutoriais
- **API explorer interativo**: Testar endpoints com dados reais (Swagger/OpenAPI)
- **Plugin sandbox**: Testar plugin em ambiente isolado
- **Publisher do Marketplace**: Submeter, gerenciar e acompanhar seus plugins
- **Docs do SDK de Integrações**: Como construir conectores de integração customizados

### 5.2 Suporte Bilíngue

Todo o Admin UI usa `next-intl` para internacionalização completa:

```typescript
// src/messages/pt.json (trecho)
{
  "dashboard": {
    "title": "Painel",
    "totalFlags": "Total de Flags",
    "activeFlags": "Flags Ativas",
    "staleFlags": "Flags Obsoletas",
    "technicalDebt": "Score de Dívida Técnica"
  },
  "flags": {
    "create": "Criar Flag",
    "status": {
      "planning": "Planejamento",
      "active": "Ativa",
      "rolledOut": "Rollout Completo",
      "archived": "Arquivada"
    }
  },
  "plugins": {
    "installed": "Plugins Instalados",
    "marketplace": "Marketplace",
    "develop": "Desenvolver Plugins"
  },
  "integrations": {
    "title": "Integrações",
    "messaging": "Mensageria & Email",
    "analytics": "Analytics & Eventos",
    "queues": "Filas Técnicas",
    "webhooks": "Webhooks",
    "health": {
      "healthy": "Saudável",
      "degraded": "Degradado",
      "down": "Fora do ar"
    },
    "test": "Testar Conexão",
    "eventLog": "Log de Eventos"
  }
}
```

---

## 6. Sistema de Plugins & SDK

### 6.1 Arquitetura de Plugins

Plugins do FlagBridge são pacotes autocontidos que estendem a plataforma através de hook points bem definidos. Plugins podem estender tanto o backend (Go API) quanto o frontend (Next.js Admin UI).

#### Manifesto do Plugin (`flagbridge-plugin.json`)

```json
{
  "name": "@myplugin/slack-alerts",
  "version": "1.0.0",
  "displayName": "Alertas Slack para Mudanças em Flags",
  "description": "Envie notificações ricas no Slack quando flags são alternadas ou modificadas",
  "author": {
    "name": "Jane Developer",
    "email": "jane@example.com",
    "url": "https://github.com/janedeveloper"
  },
  "license": "MIT",
  "minFlagBridgeVersion": "1.0.0",
  "category": "integration",
  "tags": ["slack", "notificações", "alertas"],
  "pricing": {
    "type": "free"
  },
  "hooks": {
    "backend": [
      {
        "point": "flag.afterToggle",
        "handler": "handlers/on-flag-toggle"
      }
    ],
    "ui": [
      {
        "point": "settings.integrations.panel",
        "component": "components/SlackConfigPanel"
      },
      {
        "point": "flag.detail.sidebar",
        "component": "components/SlackActivityWidget"
      }
    ]
  },
  "config": {
    "schema": {
      "slackWebhookUrl": {
        "type": "string",
        "label": "URL do Webhook Slack",
        "required": true,
        "secret": true
      },
      "channel": {
        "type": "string",
        "label": "Canal Padrão",
        "default": "#feature-flags"
      }
    }
  },
  "permissions": [
    "flag:read",
    "flag:events",
    "settings:write"
  ]
}
```

#### Hook Points de Backend

| Hook Point | Gatilho | Casos de Uso |
|-----------|---------|-------------|
| `flag.beforeEvaluate` | Antes da avaliação da flag | Lógica de targeting customizada, roteamento A/B |
| `flag.afterEvaluate` | Após avaliação da flag | Logging, tracking de analytics |
| `flag.beforeToggle` | Antes do toggle da flag | Gates de aprovação, validação |
| `flag.afterToggle` | Após toggle da flag | Notificações, sync com sistemas externos |
| `flag.beforeUpdate` | Antes de mudança nos metadados | Validação, campos obrigatórios |
| `flag.afterUpdate` | Após mudança nos metadados | Audit, sync, notificações |
| `flag.onStaleDetected` | Quando flag obsoleta detectada | Ações de cleanup customizadas |
| `project.onCreate` | Novo projeto criado | Auto-setup, templates |
| `api.middleware` | Toda requisição API | Extensões de auth, rate limiting |
| `evaluation.metrics` | Métricas de avaliação coletadas | Analytics customizado, export |
| `integration.onEvent` | Quando evento de integração dispara | Workflows cross-integration |

#### Hook Points de Frontend (UI)

| Hook Point | Local | Casos de Uso |
|-----------|-------|-------------|
| `dashboard.widget` | Página do Dashboard | Cards de métricas customizados, gráficos |
| `flag.detail.sidebar` | Sidebar do detalhe da flag | Info relacionada, dados externos |
| `flag.detail.tab` | Tabs do detalhe da flag | Tabs customizadas (ex: analytics) |
| `flag.list.column` | Tabela da lista de flags | Colunas customizadas |
| `settings.integrations.panel` | Configurações > Integrações | UI de configuração do plugin |
| `integrations.provider.config` | Página de config de integração | UI customizada de config do provider |
| `navigation.sidebar` | Sidebar de navegação | Novos itens de navegação |
| `global.banner` | Topo da página | Alertas, anúncios |

#### Plugin SDK

```typescript
// @flagbridge/plugin-sdk — Construindo um hook de backend

import { FlagBridgePlugin, FlagToggleEvent } from '@flagbridge/plugin-sdk';

export default class SlackAlertsPlugin extends FlagBridgePlugin {
  name = 'slack-alerts';

  async onFlagToggle(event: FlagToggleEvent): Promise<void> {
    const config = this.getConfig();
    const { flag, environment, toggledBy, newState } = event;

    await fetch(config.slackWebhookUrl, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        channel: config.channel,
        blocks: [
          {
            type: 'section',
            text: {
              type: 'mrkdwn',
              text: `🚩 *${flag.name}* foi ${newState ? 'habilitada' : 'desabilitada'}` +
                    ` em *${environment.name}* por ${toggledBy.name}`
            }
          }
        ]
      })
    });
  }
}
```

```tsx
// @flagbridge/plugin-sdk — Construindo uma extensão de UI (componente Next.js)

import { PluginUIComponent, useFlagContext, usePluginConfig } from '@flagbridge/plugin-sdk/react';

export const SlackActivityWidget: PluginUIComponent = () => {
  const flag = useFlagContext();
  const config = usePluginConfig();

  return (
    <div className="p-4 border rounded-lg">
      <h3 className="font-semibold text-sm text-gray-500">Notificações Slack</h3>
      <p className="text-sm mt-1">Postando em <code>{config.channel}</code></p>
      <p className="text-xs text-gray-400 mt-2">Última notificação: {flag.lastToggleAt}</p>
    </div>
  );
};
```

#### Plugin CLI

```bash
npx @flagbridge/create-plugin meu-plugin-incrivel  # Scaffold
cd meu-plugin-incrivel
flagbridge plugin dev       # Dev local com hot-reload
flagbridge plugin test      # Rodar testes
flagbridge plugin build     # Build para distribuição
flagbridge plugin publish   # Publicar no marketplace
```

### 6.2 Modelo de Segurança de Plugins

- **Execução sandboxed**: Plugins de backend rodam em goroutines isoladas com acesso limitado a syscalls
- **Sistema de permissões**: Plugins declaram permissões necessárias; admin deve aprovar
- **Sandboxing de UI**: Plugins de frontend renderizam em iframes sandboxed ou com CSP restrito
- **Criptografia de config**: Valores secretos (API keys, tokens) criptografados em repouso
- **Assinatura de código**: Plugins do marketplace devem ser assinados; hash verificado na instalação
- **Processo de revisão**: Todas as submissões ao marketplace são revisadas antes da publicação

---

## 7. Marketplace de Plugins

### 7.1 Visão

O FlagBridge Plugin Marketplace transforma o FlagBridge de um produto em um **ecossistema de plataforma** — similar à Shopify App Store, WordPress Plugins ou Figma Community. Desenvolvedores third-party podem construir, publicar e monetizar plugins que estendem o FlagBridge.

### 7.2 Features do Marketplace

#### Para Usuários de Plugins (Compradores)
- Navegar & buscar por categoria, rating, preço, compatibilidade
- Instalação com um clique do Admin UI (self-hosted e SaaS)
- Reviews & ratings (5 estrelas com reviews escritas)
- Verificação de compatibilidade contra versão atual do FlagBridge
- Configuração de plugins pelo Admin UI
- Auto-updates opcionais (self-hosted respeita controle do usuário)

#### Para Desenvolvedores de Plugins (Vendedores)
- Portal do Desenvolvedor com documentação completa, API explorer, Plugin SDK
- Plugin Sandbox para desenvolvimento e debugging
- Workflow de publicação: Submit → checks automatizados → review manual → publicado
- Opções de preço: gratuito, compra única, ou assinatura mensal
- Split de receita: 80% desenvolvedor / 20% FlagBridge
- Dashboard de ganhos: instalações, receita, ratings, tickets de suporte
- Badge de Desenvolvedor Verificado para publishers confiáveis

#### Para FlagBridge (Plataforma)
- Pipeline de review: scan de segurança automatizado + review manual de código
- Fonte de receita: 20% de comissão em plugins pagos
- Flywheel do ecossistema: mais plugins → mais valor → mais usuários → mais desenvolvedores
- Controle de qualidade: padrões mínimos de qualidade, requisitos de segurança, auditorias regulares

### 7.3 Categorias do Marketplace

| Categoria | Exemplos |
|-----------|---------|
| **Integrações** | Slack, Discord, Linear, Jira, PagerDuty, Datadog |
| **Mensageria** | Resend, RD Station, SendGrid, Mailchimp, Customer.io, Brevo |
| **Analytics** | Mixpanel, Amplitude, Segment, PostHog, GA4, Rudderstack |
| **Filas** | SQS, Kafka, RabbitMQ, NATS, Redis Streams, GCP Pub/Sub |
| **Segurança** | RBAC avançado, IP whitelisting, relatórios de compliance |
| **Extensões de UI** | Widgets customizados, packs de tema, layouts de dashboard |
| **Automação** | Auto-remediação, rollouts agendados, bridges de CI/CD |
| **Dados** | Conectores de export (BigQuery, Snowflake), ferramentas de migração |

### 7.4 Modelo de Receita

```
Venda de Plugin ($10.00)
├── Desenvolvedor recebe: $8.00 (80%)
├── FlagBridge recebe: $2.00 (20%)
└── Processamento de pagamento: deduzido da parte FlagBridge

Plugin por Assinatura ($5.00/mês)
├── Desenvolvedor recebe: $4.00/mês (80%)
├── FlagBridge recebe: $1.00/mês (20%)
└── Cobrança recorrente gerenciada por Stripe Connect
```

### 7.5 Implementação Técnica

- **Pagamento**: Stripe Connect para pagamentos do marketplace e repasses a desenvolvedores
- **Registry de pacotes**: Registry privado tipo npm para pacotes de plugins
- **CI/CD**: Pipeline de testes automatizados para plugins submetidos
- **CDN**: Pacotes de plugins servidos via CDN para instalação rápida globalmente
- **Versionamento**: Semantic versioning obrigatório, com matriz de compatibilidade

---

## 8. Ecossistema de Integrações

Esta seção detalha as três categorias de integrações com serviços externos que diferenciam o FlagBridge dos concorrentes: **Mensageria & Comunicação**, **Analytics & Rastreamento de Eventos**, e **Filas Técnicas**.

### 8.1 Integrações de Mensageria & Comunicação

Conecte eventos de lifecycle de flags a plataformas de comunicação para notificações automatizadas, campanhas e atualizações para stakeholders.

#### Providers Suportados

| Provider | Tipo | Casos de Uso |
|----------|------|-------------|
| **Resend** | Email transacional | Notificar stakeholders quando flags fazem rollout; enviar digests de flags obsoletas |
| **RD Station** | Automação de marketing (foco Brasil) | Disparar fluxos de nurturing baseados em acesso a features; segmentar usuários por flags ativas |
| **SendGrid** | Email transacional + marketing | Relatórios de rollout automatizados; comunicação com usuários sobre lançamentos de features |
| **Mailchimp** | Campanhas de marketing | Segmentação de audiência por feature flags; anúncios de lançamento |
| **Brevo (Sendinblue)** | Multi-canal (email, SMS, WhatsApp) | Notificações multi-canal em eventos de flags; alertas SMS para rollbacks críticos |
| **Customer.io** | Mensageria comportamental | Disparar mensagens baseadas em padrões de avaliação de flags; fluxos de onboarding por acesso a features |

#### Como Funciona

```
┌─────────────┐     ┌──────────────────┐     ┌─────────────────┐
│ Evento Flag  │────▶│ Regra de         │────▶│ Provider de     │
│              │     │ Integração       │     │ Mensageria      │
│ flag.toggled │     │ SE flag.key      │     │                 │
│ flag.rolledOut│    │   bate com "x"   │     │ Resend.send()   │
│ flag.stale   │     │ E env = prod     │     │ RDStation.push()│
│              │     │ ENTÃO ação...    │     │ SendGrid.send() │
└─────────────┘     └──────────────────┘     └─────────────────┘
```

#### Exemplo de Configuração (Admin UI)

```json
{
  "provider": "resend",
  "config": {
    "apiKey": "re_...",
    "fromEmail": "flags@minhaempresa.com",
    "fromName": "FlagBridge"
  },
  "rules": [
    {
      "trigger": "flag.rolledOut",
      "conditions": {
        "environment": "production"
      },
      "action": {
        "template": "flag-rollout-notification",
        "to": "{{flag.owner.email}}",
        "cc": ["time-produto@minhaempresa.com"],
        "data": {
          "flagName": "{{flag.name}}",
          "rolledOutAt": "{{event.timestamp}}",
          "hypothesis": "{{flag.productCard.hypothesis}}"
        }
      }
    },
    {
      "trigger": "flag.staleDetected",
      "conditions": {
        "staleDays": { "gte": 30 }
      },
      "action": {
        "template": "stale-flag-reminder",
        "to": "{{flag.owner.email}}",
        "data": {
          "flagName": "{{flag.name}}",
          "staleSinceDays": "{{lifecycle.staleDays}}",
          "cleanupUrl": "{{flag.adminUrl}}"
        }
      }
    }
  ]
}
```

#### Integração Profunda com RD Station

Dado as raízes brasileiras do FlagBridge, RD Station recebe uma integração de primeira classe:

- **Segmentação de contatos**: Marcar automaticamente contatos do RD Station baseado em quais feature flags estão ativas para eles
- **Eventos de conversão**: Enviar avaliações de feature flags como eventos de conversão para o RD Station para análise de funil
- **Lead scoring**: Ajustar scores de leads baseado em padrões de adoção de features
- **Fluxos de nurturing**: Disparar fluxos de email específicos quando um usuário ganha acesso a uma nova feature via rollout de flag

### 8.2 Integrações de Analytics & Rastreamento de Eventos

Streamar dados de avaliação de flags para plataformas de analytics, permitindo que times correlacionem rollouts de features com métricas de produto, comportamento de usuários e resultados de negócio.

#### Providers Suportados

| Provider | Tipo | Casos de Uso |
|----------|------|-------------|
| **Mixpanel** | Analytics de produto | Rastrear avaliações de flags como eventos; construir funis por variante; medir adoção de features |
| **Amplitude** | Analytics de produto | Correlacionar rollouts de features com comportamento de usuários; dashboards de análise de experimentos |
| **Segment** | Customer data platform | Rotear eventos de flags para qualquer destino Segment; perfis unificados de usuários com dados de flags |
| **PostHog** | Analytics open-source | Correlação de feature flags com gravações de sessão; análise de funil por variante |
| **Google Analytics 4** | Web analytics | Eventos customizados para avaliações de flags; tracking de conversão por variante |
| **Rudderstack** | CDP open-source | Pipeline de eventos self-hosted; rotear dados de flags para warehouses e ferramentas de analytics |

#### Como Funciona

```
┌─────────────────┐     ┌───────────────────┐     ┌──────────────────┐
│ Avaliação Flag   │────▶│ Adapter Analytics  │────▶│ Provider         │
│                  │     │                   │     │ Analytics        │
│ flag: checkout-v2│     │ Transforma para   │     │                  │
│ variante: "new" │     │ formato provider  │     │ mixpanel.track() │
│ userId: "u_123" │     │                   │     │ amplitude.log()  │
│ timestamp: ...  │     │ Aplica sampling   │     │ segment.track()  │
│                  │     │ & batching        │     │ posthog.capture()│
└─────────────────┘     └───────────────────┘     └──────────────────┘
```

#### Schema de Eventos

Cada avaliação de flag gera um evento padronizado que é adaptado para cada provider:

```json
{
  "event": "flagbridge.flag_evaluated",
  "timestamp": "2026-03-22T14:30:00Z",
  "userId": "user_123",
  "properties": {
    "flag_key": "checkout-v2",
    "flag_name": "Checkout V2 Redesign",
    "variant": "new",
    "environment": "production",
    "project": "web-app",
    "evaluation_reason": "targeting_match",
    "flag_status": "active",
    "flag_owner": "gabriel@empresa.com",
    "product_hypothesis": "Novo fluxo de checkout aumenta conversão em 15%"
  }
}
```

**Features chave:**

- **Sampling**: Taxa de amostragem configurável (ex: enviar 10% das avaliações) para controlar custos em flags de alto volume
- **Batching**: Eventos são agrupados (configurável: a cada N eventos ou a cada N segundos) antes de enviar para reduzir chamadas de API
- **Enriquecimento**: Eventos são automaticamente enriquecidos com contexto de produto (hipótese, owner, status) dos Product Cards
- **Tracking seletivo**: Escolha quais flags e environments rastrear — não precisa enviar tudo

### 8.3 Integrações de Filas Técnicas (Plugin Pro)

Publicar eventos de flags em filas de mensagens para processamento downstream por microsserviços, pipelines de dados e automação customizada. Isso é entregue como **plugin Pro** porque mira times de engenharia com necessidades de infraestrutura complexas.

#### Providers Suportados

| Provider | Tipo | Casos de Uso |
|----------|------|-------------|
| **Amazon SQS** | Fila gerenciada (AWS) | Desacoplar eventos de flags de consumidores downstream; disparar funções Lambda em mudanças de flags |
| **Apache Kafka** | Streaming distribuído | Streamar avaliações de flags para processamento real-time; event sourcing para mudanças de estado |
| **RabbitMQ** | Message broker | Rotear eventos de flags para consumidores específicos via exchanges; entrega confiável com acknowledgments |
| **NATS** | Mensageria leve | Distribuição de eventos de flags com baixa latência; pub/sub para microsserviços |
| **Redis Streams** | Streaming in-memory | Streaming de avaliações de alta throughput; consumer groups para processamento paralelo |
| **GCP Pub/Sub** | Mensageria gerenciada (GCP) | Arquiteturas event-driven no Google Cloud; disparar Cloud Functions em eventos de flags |

#### Como Funciona

```
┌─────────────────┐     ┌──────────────────┐     ┌──────────────────┐
│ Evento Flag      │────▶│ Adapter Fila      │────▶│ Provider Fila    │
│                  │     │                  │     │                  │
│ flag.toggled     │     │ Serializar evento│     │ SQS.sendMessage()│
│ flag.evaluated   │     │ Rotear p/ fila   │     │ kafka.produce()  │
│ flag.stale       │     │ Retry em falha   │     │ rabbit.publish() │
│ flag.rolledOut   │     │                  │     │ nats.publish()   │
└─────────────────┘     └──────────────────┘     └──────────────────┘

                         ┌──────────────────┐
                         │ Consumidores      │
                         │ Downstream       │
                         │                  │
                         │ Funções Lambda   │
                         │ Pipelines dados  │
                         │ Microsserviços   │
                         │ Sistemas alerta  │
                         └──────────────────┘
```

#### Casos de Uso

**Reações event-driven a flags**: Microsserviços se inscrevem em eventos de mudança de flags e ajustam comportamento automaticamente (invalidação de cache, reload de config, updates de service mesh)

**Integração com pipeline de dados**: Streamar avaliações de flags para data warehouses via Kafka → Spark/Flink → BigQuery/Snowflake para análise offline

**Auditoria & compliance**: Log imutável de todas as mudanças de flags em fila durável, satisfazendo requisitos de compliance

**Coordenação cross-service**: Quando uma flag é alternada no FlagBridge, publicar no SQS/Kafka para que todos os serviços dependentes reajam sem fazer polling na API

**Alerting customizado**: Consumir eventos de flags de uma fila e alimentar pipelines de alerta customizados (PagerDuty, OpsGenie, bots Slack customizados)

#### Exemplo de Configuração (Kafka)

```json
{
  "provider": "kafka",
  "config": {
    "brokers": ["kafka-1:9092", "kafka-2:9092"],
    "clientId": "flagbridge-producer",
    "sasl": {
      "mechanism": "SCRAM-SHA-256",
      "username": "flagbridge",
      "password": "..."
    },
    "ssl": true
  },
  "rules": [
    {
      "trigger": "flag.toggled",
      "action": {
        "topic": "flagbridge.flag-events",
        "partitionKey": "{{flag.projectId}}"
      }
    },
    {
      "trigger": "flag.evaluated",
      "action": {
        "topic": "flagbridge.evaluations",
        "partitionKey": "{{flag.key}}",
        "sampling": 0.05,
        "batchSize": 1000,
        "compressionType": "gzip"
      }
    }
  ]
}
```

### 8.4 Arquitetura de Integrações (Go)

Todas as integrações compartilham uma interface adapter unificada em Go:

```go
// pkg/integrations/adapter.go

package integrations

type EventType string

const (
    FlagToggled       EventType = "flag.toggled"
    FlagEvaluated     EventType = "flag.evaluated"
    FlagRolledOut     EventType = "flag.rolledOut"
    FlagStaleDetected EventType = "flag.staleDetected"
    FlagUpdated       EventType = "flag.updated"
    FlagArchived      EventType = "flag.archived"
)

// IntegrationEvent é a estrutura de evento padronizada
type IntegrationEvent struct {
    Type        EventType              `json:"type"`
    Timestamp   time.Time              `json:"timestamp"`
    Flag        FlagSnapshot           `json:"flag"`
    Environment EnvironmentSnapshot    `json:"environment"`
    Actor       *ActorSnapshot         `json:"actor,omitempty"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Adapter é a interface que todos os providers devem implementar
type Adapter interface {
    Name() string
    Type() IntegrationType
    ConfigSchema() json.RawMessage
    Connect(config json.RawMessage) error
    Send(ctx context.Context, event IntegrationEvent, actionConfig json.RawMessage) error
    SendBatch(ctx context.Context, events []IntegrationEvent, actionConfig json.RawMessage) error
    HealthCheck(ctx context.Context) HealthStatus
    Close() error
}

// Registry gerencia adapters de integração disponíveis
type Registry struct {
    adapters map[string]Adapter
}

func (r *Registry) Register(adapter Adapter) {
    r.adapters[adapter.Name()] = adapter
}

func (r *Registry) Get(name string) (Adapter, bool) {
    a, ok := r.adapters[name]
    return a, ok
}
```

Este padrão adapter significa:

1. **Adicionar um novo provider é simples** — implemente a interface `Adapter`
2. **A comunidade pode contribuir providers** — via o sistema de plugins
3. **Todos os providers compartilham o mesmo formato de evento** — `IntegrationEvent` padronizado
4. **Health checks são uniformes** — todo provider reporta saúde da mesma forma
5. **Batching e sampling são tratados na camada de integração** — providers não precisam implementar isso

---

## 9. Go-to-Market

### 9.1 Público-alvo

**Primário:** Startups e scale-ups (50-500 devs) que:
- Já usam feature flags (Unleash CE, in-house) mas sentem falta de observabilidade e contexto de produto
- Não podem pagar LaunchDarkly (~$72k/ano)
- Valorizam self-hosted e data sovereignty
- Querem extensibilidade sem vendor lock-in
- Já usam ferramentas como Mixpanel, Segment ou Kafka e querem que sua ferramenta de FF integre nativamente

**Secundário:** Times de produto em empresas maiores que:
- Querem visibilidade sobre o estado das flags sem depender de engenharia
- Precisam vincular flags a OKRs e roadmap
- Querem um ecossistema de plugins para integrações customizadas
- Precisam disparar comunicações (Resend, RD Station) baseadas em rollouts de features

### 9.2 Pricing

| Plano | Preço | Target |
|-------|-------|--------|
| **Community** | Grátis (forever) | Devs individuais, times pequenos, avaliação |
| **Pro** (self-hosted) | $29/mês (≤10 seats) / $79/mês (≤50 seats) / $149/mês (unlimited) | Scale-ups, times médios |
| **Pro** (SaaS) | $49/mês (≤10 seats) / $99/mês (≤50 seats) / $199/mês (unlimited) | Times que não querem gerenciar infra |
| **Enterprise** | Custom (a partir de $500/mês) | Grandes empresas, compliance, SLA |

**Modelo flat-rate por seats** — sem surpresas de MAU ou service connections. Diferencial direto contra LaunchDarkly.

Self-hosted é mais barato que SaaS porque o cliente paga a infra. Isso incentiva self-hosted (menos custo de hosting pra nós) e atrai o público que valoriza data sovereignty.

### 9.3 Canais de Aquisição

**Developer-first growth:**
1. GitHub — repo open-source bem documentado, README atraente
2. Hacker News / Reddit / Dev.to — posts de lançamento
3. Comparações honestas — blog posts: "FlagBridge vs Unleash vs LaunchDarkly"
4. Docker Hub — imagem oficial otimizada, one-liner de setup
5. OpenFeature ecosystem — provider oficial listado no site do OpenFeature
6. YouTube / blog — tutoriais em inglês e português

**Product-led growth:**
- CE é genuinamente útil → time adota → cresce → precisa de observabilidade/governance → upgrade Pro
- Banner discreto no Admin UI
- Trial de 14 dias do Pro

**Crescimento do ecossistema (pós-marketplace):**
- Desenvolvedores de plugins atraem seus próprios usuários para o FlagBridge
- Marketplace se torna canal de descoberta
- Revenue share incentiva plugins de qualidade

**Crescimento via integrações:**
- Times buscando "Mixpanel feature flags" ou "Kafka feature flag events" encontram o FlagBridge
- Integração com RD Station posiciona o FlagBridge de forma única no mercado brasileiro
- Parceiros de integração podem co-divulgar o FlagBridge em seus diretórios

### 9.4 Métricas de Sucesso (Ano 1)

| Métrica | 6 meses | 12 meses |
|---------|---------|----------|
| GitHub stars | 500 | 2,000 |
| Docker pulls | 5,000 | 25,000 |
| CE installs ativos | 200 | 1,000 |
| Pro paying customers | 10 | 50 |
| MRR | $500 | $3,000 |
| Community (Discord) | 100 | 500 |
| Plugins publicados | — | 20 |
| Marketplace GMV | — | $500/mês |
| Integrações ativas (Pro) | — | 150 |

---

## 10. Roadmap

### Q2 2026 — Fundação
- [ ] Repo público no GitHub
- [ ] Core Go API server: flag management + eval engine
- [ ] Next.js Admin UI: flag CRUD, dashboard, settings (bilíngue desde o dia 1)
- [ ] SDK Node.js + SDK Go
- [ ] Docker image + Docker Compose
- [ ] Site de docs (docs.flagbridge.io) — bilíngue
- [ ] OpenFeature Provider (Node.js)
- [ ] Plugin runtime (básico)
- [ ] Integração webhook (CE)

### Q3 2026 — Lançamento Community
- [ ] Product Context Cards (versão CE básica)
- [ ] CLI (`flagbridge`)
- [ ] SDK React (client-side) + SDK Python
- [ ] Helm chart para Kubernetes
- [ ] Plugin CLI + Plugin SDK v1
- [ ] 3-5 plugins oficiais (Slack, Linear, GitHub)
- [ ] Integração notificação Slack (CE)
- [ ] Lançamento: Hacker News + Product Hunt + Reddit

### Q4 2026 — Lançamento Pro
- [ ] Arquitetura do Plugin Pro
- [ ] Technical Dashboard (métricas, observabilidade)
- [ ] Product Cards avançado
- [ ] Lifecycle automation
- [ ] SSO (OIDC)
- [ ] Billing + gerenciamento de license key
- [ ] **Integrações de mensageria**: Resend, SendGrid, RD Station
- [ ] **Integrações de analytics**: Mixpanel, Segment
- [ ] Plugin Marketplace v1 (apenas plugins gratuitos)
- [ ] Hub de Integrações no Admin UI

### Q1 2027 — Escala & Marketplace
- [ ] FlagBridge Edge (proxy para alto volume)
- [ ] Experimentação A/B básica
- [ ] RBAC granular + change requests
- [ ] **Integrações de filas**: SQS, Kafka, RabbitMQ
- [ ] **Integrações de analytics**: Amplitude, PostHog, GA4, Rudderstack
- [ ] **Integrações de mensageria**: Mailchimp, Brevo, Customer.io
- [ ] Oferta SaaS managed
- [ ] Plugin Marketplace v2: plugins pagos, Stripe Connect
- [ ] Portal do Desenvolvedor com sandbox
- [ ] Redesign da landing page bilíngue

### Q2 2027 — Ecossistema
- [ ] Plugin Marketplace v3: assinaturas, desenvolvedores verificados
- [ ] **Integrações de filas**: NATS, Redis Streams, GCP Pub/Sub
- [ ] SDK de Integrações (para conectores construídos pela comunidade)
- [ ] Opção white-label (Enterprise)
- [ ] Experimentação avançada
- [ ] SOC 2 compliance (início)
- [ ] Hackathon de plugins da comunidade

---

## 11. Riscos e Mitigações

| Risco | Impacto | Mitigação |
|-------|---------|-----------|
| Mercado saturado de FF tools | Alto | Diferenciação clara via Product Context + Plugin Marketplace + Ecossistema de Integrações — ninguém faz os três |
| Unleash adiciona features similares | Médio | Velocidade de execução + foco product-first; ecossistema de integrações como fosso competitivo |
| Dificuldade de monetizar open-source | Alto | CE genuinamente útil mas Pro resolve dores reais; Marketplace cria receita adicional; integrações impulsionam adoção do Pro |
| Complexidade de manter múltiplos SDKs | Médio | Começar com 4 SDKs; OpenFeature reduz superfície |
| Tempo limitado (side project) | Alto | MVP enxuto; priorizar core → launch → iterate com feedback |
| Vulnerabilidades de segurança em plugins | Médio | Execução sandboxed, sistema de permissões, review obrigatório, assinatura de código |
| Chicken-and-egg do marketplace | Alto | Construir 5-10 plugins oficiais primeiro; incentivos para early developers |
| Carga de manutenção de integrações | Médio | Padrão adapter isola providers; comunidade pode contribuir via sistema de plugins; começar com 6 providers e expandir |
| Mudanças em APIs de terceiros | Baixo | Abstração via adapter localiza mudanças; version pinning em SDKs de providers |

---

## Próximos Passos

1. **Registrar domínios** — flagbridge.io e flagbridge.dev
2. **Criar org no GitHub** — github.com/flagbridge
3. **Reservar scope no npm** — @flagbridge
4. **Iniciar monorepo** — `flagbridge/flagbridge` (api + ui + docs + plugin-sdk)
5. **Prototipar** — API server Go com flag CRUD + eval endpoint
6. **Design** — Wireframes do Admin UI Next.js (dashboard + detalhe da flag + product card + hub de integrações)
7. **Landing page** — Site Next.js SSG bilíngue com posicionamento, features, pricing
8. **Scaffold do Plugin SDK** — Definir hook points, formato do manifesto, CLI
9. **Scaffold do adapter de integrações** — Implementar interface `Adapter` + primeiro provider (Resend ou Slack)

---

*FlagBridge — A ponte entre feature flags e estratégia de produto.*
