# FlagBridge — Product Document / Documento de Produto

> 🇺🇸 **"The bridge between feature flags and product strategy."**  
> 🇧🇷 **"A ponte entre feature flags e estratégia de produto."**

**Version / Versão:** 2.0  
**Author / Autor:** Gabriel Gripp  
**Date / Data:** March / Março 2026  
**Status:** Draft / MVP Planning

---

# Table of Contents / Índice

1. [Overview / Visão Geral](#1-overview--visão-geral)
2. [Competitive Analysis / Análise Competitiva](#2-competitive-analysis--análise-competitiva)
3. [MVP Definition — Community vs Pro / Definição de MVP](#3-mvp-definition--definição-de-mvp)
4. [Technical Architecture / Arquitetura Técnica](#4-technical-architecture--arquitetura-técnica)
5. [Admin Panel — Technical & Product Documentation / Painel Admin](#5-admin-panel--painel-admin)
6. [Plugin System & SDK / Sistema de Plugins & SDK](#6-plugin-system--sdk--sistema-de-plugins--sdk)
7. [Plugin Marketplace / Marketplace de Plugins](#7-plugin-marketplace--marketplace-de-plugins)
8. [Go-to-Market](#8-go-to-market)
9. [Roadmap](#9-roadmap)
10. [Risks & Mitigations / Riscos e Mitigações](#10-risks--mitigations--riscos-e-mitigações)

---

# 1. Overview / Visão Geral

## 🇺🇸 English

FlagBridge is an open-core Feature Flag Management platform that goes beyond simple toggle on/off: it connects feature flags to product planning, technical observability, and lifecycle management. The key differentiator is being **product-first with solid infrastructure** — every flag has business context, impact metrics, and clear lifecycle rules.

### The Problem

The current feature flag market is dominated by **infra-first** tools:

- **Engineering** creates flags without product context — nobody knows *why* a flag exists
- **Product** has no visibility into flag states or their real impact
- **Zombie flags** accumulate (100% ON for months, never removed), generating technical debt
- **No tool** natively connects a flag to a product hypothesis, experiment, or OKR
- Enterprise tools like LaunchDarkly charge median annual contracts of ~$72k (Vendr data, 2026), unaffordable for most companies
- **No tool** offers a plugin ecosystem — customization requires forking or waiting for vendor roadmap

### The Solution

FlagBridge offers:

1. **Feature Flag Management** — create, toggle, evaluate, rollout strategies
2. **Product Context Cards** — each flag linked to hypothesis, owner, success metrics, deadline
3. **Technical Dashboard** — adoption rate, error rate by variant, latency impact, stale flag detection
4. **Lifecycle Automation** — cleanup alerts, auto-archival, tech debt tracking
5. **OpenFeature Compatible** — official provider for the CNCF OpenFeature standard, zero vendor lock-in
6. **Plugin Ecosystem & Marketplace** — extensible architecture where developers can build, publish, and sell plugins

### Open-Core Model

| Layer | Distribution | Price |
|-------|-------------|-------|
| **Community Edition** | Open source (Apache 2.0) | Free |
| **Pro Edition** | Self-hosted plugin OR SaaS | $X/mo |
| **Enterprise** | Managed SaaS + support | Custom |

**The self-hosted plugin advantage:** the customer already runs FlagBridge CE, purchases the Pro license, runs `docker pull flagbridge/pro-plugin && docker compose up -d`, and restarts. No migration, no cloud switch, no downtime. The plugin injects Pro modules into the same deployment.

---

## 🇧🇷 Português (Brasil)

FlagBridge é uma plataforma open-core de Feature Flag Management que vai além do toggle on/off: conecta feature flags a planejamento de produto, observabilidade técnica e lifecycle management. O diferencial é ser **product-first com infra sólida** — cada flag tem contexto de negócio, métricas de impacto e regras de ciclo de vida claras.

### O Problema

O mercado atual de feature flags é dominado por ferramentas **infra-first**:

- **Engenharia** cria flags sem contexto de produto — ninguém sabe *por que* uma flag existe
- **Produto** não tem visibilidade sobre o estado das flags nem sobre seu impacto real
- **Flags zumbis** se acumulam (100% ON há meses, nunca removidas), gerando dívida técnica
- **Nenhuma ferramenta** conecta nativamente uma flag a uma hipótese de produto, experimento ou OKR
- Ferramentas enterprise como LaunchDarkly cobram contratos medianos de ~$72k/ano (dados Vendr, 2026), inacessíveis para a maioria das empresas
- **Nenhuma ferramenta** oferece um ecossistema de plugins — customização exige fork ou esperar o roadmap do vendor

### A Solução

FlagBridge oferece:

1. **Feature Flag Management** — create, toggle, evaluate, rollout strategies
2. **Product Context Cards** — cada flag atrelada a hipótese, owner, métricas de sucesso, prazo
3. **Technical Dashboard** — adoption rate, error rate por variante, latency impact, stale flag detection
4. **Lifecycle Automation** — alertas de cleanup, archival automático, tracking de dívida técnica
5. **OpenFeature Compatible** — provider oficial para o padrão CNCF OpenFeature, zero vendor lock-in
6. **Ecossistema de Plugins & Marketplace** — arquitetura extensível onde devs podem criar, publicar e vender plugins

### Modelo Open-Core

| Camada | Distribuição | Preço |
|--------|-------------|-------|
| **Community Edition** | Open source (Apache 2.0) | Grátis |
| **Pro Edition** | Plugin self-hosted OU SaaS | $X/mês |
| **Enterprise** | SaaS managed + suporte | Custom |

**A sacada do plugin self-hosted:** o cliente já roda o FlagBridge CE, compra a licença Pro, faz `docker pull flagbridge/pro-plugin && docker compose up -d` e reinicia. Sem migração, sem trocar de cloud, sem downtime. O plugin injeta os módulos Pro no mesmo deployment.

---

# 2. Competitive Analysis / Análise Competitiva

## 🇺🇸 Landscape

| Feature | FlagBridge | Unleash | LaunchDarkly | Flagsmith | PostHog FF |
|---------|-----------|---------|-------------|-----------|-----------|
| Open Source Core | ✅ Apache 2.0 | ✅ Apache 2.0 | ❌ | ✅ BSD 3 | ✅ MIT |
| Self-hosted | ✅ | ✅ | ❌ | ✅ | ✅ |
| SaaS hosted | ✅ | ✅ | ✅ | ✅ | ✅ |
| Zero-migration plugin upgrade | ✅ | ❌ | ❌ | ❌ | ❌ |
| Product Context (hypothesis, owner, OKR) | ✅ | ❌ | ❌ | ❌ | Partial |
| Technical Dashboard (observability) | ✅ Pro | Basic | ✅ (expensive) | Basic | ✅ |
| Lifecycle/Cleanup Automation | ✅ Pro | ❌ | Partial | ❌ | ❌ |
| Plugin Ecosystem & Marketplace | ✅ | ❌ | ❌ | ❌ | ❌ |
| OpenFeature Provider | ✅ | ✅ | ✅ | ✅ | ✅ |
| Entry paid price | ~$29/mo | $80/mo (5 seats) | $120/mo+ | $45/mo | $0 (bundled) |

### Positioning / Posicionamento

FlagBridge occupies the space between:

- **Unleash** — good open-source, weak in product/analytics, no plugin ecosystem
- **LaunchDarkly** — powerful but expensive and SaaS-only, no extensibility
- **PostHog** — good analytics but FF is secondary, no standalone FF product
- **None of them** offer a Plugin Marketplace — FlagBridge creates a new category

**Tagline:** *"Feature flags with product intelligence. Open source. Extensible. Your data, your rules."*

### Exploitable gaps in competitors / Gaps exploráveis nos concorrentes

- **Unleash:** Users complain about lack of multiple projects in open-source version and absence of analytics integration
- **LaunchDarkly:** Pricing model based on service connections + MAU generates unpredictable costs; median contracts of $72k/year push away startups
- **Flagsmith:** No real-time flag sync in free version
- **All of them:** None offer a "Product Context Card" natively; none offer extensibility via plugins

---

# 3. MVP Definition / Definição de MVP

## 3.1 Community Edition (Open Source)

### 🇺🇸 English

The CE must be **genuinely useful** on its own — not a crippled version that forces upgrades. This generates community, contributions, and trust.

#### Core Features (CE)

**Flag Management**
- CRUD for feature flags (boolean, string, number, JSON)
- Unlimited environments (dev, staging, production)
- Unlimited projects
- Basic targeting rules (user ID, percentage rollout, custom attributes)
- Instant kill switch
- Full REST API + Admin UI (Next.js)
- Basic audit log (who changed what, when)

**SDK & Integrations**
- Official SDKs: Go, Node.js/TypeScript, React (client-side), Python
- Official OpenFeature Provider
- Webhooks for integrations
- CLI (`flagbridge`) for automation

**Product Context (basic)**
- Description and owner field per flag
- Customizable tags/labels
- External link (for Jira, Linear, Notion, etc.)
- Flag status: `planning` → `active` → `rolled-out` → `archived`

**Dashboard (basic)**
- Flag list per project/environment
- Status overview (how many active, stale, archived)
- Flag age tracking (created date, last modified)

**Plugin System (basic)**
- Plugin runtime for loading community plugins
- Plugin API hooks for UI extensions and API middleware
- Plugin CLI for scaffolding new plugins (`flagbridge plugin create`)

---

### 🇧🇷 Português (Brasil)

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

**Sistema de Plugins (básico)**
- Runtime de plugins para carregar plugins da comunidade
- Plugin API com hooks para extensão de UI e API middleware
- Plugin CLI para scaffolding de novos plugins (`flagbridge plugin create`)

---

## 3.2 Pro Edition (Plugin)

### 🇺🇸 English

Pro is what transforms FlagBridge from "another FF tool" into a **product intelligence platform**.

#### Pro Features

**Advanced Product Context Cards**
- Structured fields: hypothesis, success metrics (KPIs), go/no-go criteria
- Decision workflow: experiment → analyze → decide (rollout/rollback/iterate)
- Bidirectional integration with Linear, Jira, Notion (auto status sync)
- Visual timeline of each flag's lifecycle
- Flag ↔ OKR/initiative linking

**Advanced Technical Dashboard**
- Adoption rate per flag (% of requests evaluating the flag)
- Error rate by variant (flag × errors correlation)
- Latency impact per flag (before/after rollout)
- Stale flag detection with "cleanup urgency" score
- SDK version tracking (which apps use which version)
- Real-time evaluation stream

**Lifecycle & Cleanup Automation**
- Configurable rules: "if flag is 100% ON for X days → notify owner"
- Alerts via Slack, email, webhook
- Cleanup suggestions with code link (GitHub/GitLab integration)
- Auto-archive expired flags
- "Technical debt score" per project

**Governance & Security**
- SSO (SAML, OIDC)
- Granular RBAC (per project, environment, flag)
- Change requests with approval workflow
- Detailed audit log with change diffs
- Scheduled flag changes (schedule rollout for specific date/time)

**Plugin Marketplace Access**
- Access to premium plugins in the marketplace
- Plugin analytics (usage, performance impact)
- Priority plugin support

**Basic Experimentation**
- A/B split by percentage
- Basic conversion metrics (via event webhook)
- Data export to analytics tools

#### Plugin Delivery Model

```bash
# Self-hosted: upgrade from CE to Pro
docker pull flagbridge/pro-plugin:latest
docker compose up -d --force-recreate

# The plugin detects the license and activates Pro modules
# Existing data is preserved — zero migration
```

---

### 🇧🇷 Português (Brasil)

O Pro é o que transforma FlagBridge de "mais um FF tool" em **plataforma de product intelligence**.

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

**Acesso ao Plugin Marketplace**
- Acesso a plugins premium no marketplace
- Plugin analytics (uso, impacto de performance)
- Suporte prioritário para plugins

**Experimentation (básico)**
- A/B split por percentage
- Métricas de conversão básicas (via webhook de eventos)
- Export de dados para ferramentas de analytics

#### Modelo de entrega do Plugin

```bash
# Self-hosted: upgrade de CE para Pro
docker pull flagbridge/pro-plugin:latest
docker compose up -d --force-recreate

# O plugin detecta a licença e ativa os módulos Pro
# Dados existentes são preservados — zero migration
```

---

## 3.3 Enterprise (SaaS Managed)

### 🇺🇸 English

- Everything in Pro
- Managed hosting by FlagBridge (multi-tenant or dedicated)
- SLA with uptime guarantee
- Priority support (Slack, email, call)
- Custom integrations
- SOC 2 compliance (roadmap)
- Data residency options
- White-label Admin UI option
- Dedicated plugin marketplace instance

### 🇧🇷 Português (Brasil)

- Tudo do Pro
- Hosting gerenciado pela FlagBridge (multi-tenant ou dedicated)
- SLA com uptime guarantee
- Suporte prioritário (Slack, email, call)
- Custom integrations
- SOC 2 compliance (roadmap)
- Data residency options
- Opção de white-label do Admin UI
- Instância dedicada do marketplace de plugins

---

# 4. Technical Architecture / Arquitetura Técnica

## 4.1 Architecture Overview / Visão Geral da Arquitetura

```
┌──────────────────────────────────────────────────────────────┐
│                      Client Applications                      │
│           (React, Node, Go, Python — via SDKs)                │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────┐      ┌────────────────────────────────┐   │
│  │  FlagBridge   │      │   FlagBridge Edge               │   │
│  │  SDK          │─────▶│   (optional, for scale)         │   │
│  └──────────────┘      └────────────┬───────────────────┘   │
│                                     │                        │
├─────────────────────────────────────┼────────────────────────┤
│                                     ▼                        │
│  ┌────────────────────────────────────────────────────────┐  │
│  │              FlagBridge API Server (Go)                 │  │
│  │                                                        │  │
│  │  ┌──────────┐ ┌───────────┐ ┌──────────┐ ┌─────────┐  │  │
│  │  │ Flag     │ │ Product   │ │ Dashboard│ │ Plugin  │  │  │
│  │  │ Eval     │ │ Context   │ │ & Metrics│ │ Runtime │  │  │
│  │  │ Engine   │ │ Module    │ │ Module   │ │ Engine  │  │  │
│  │  └──────────┘ └───────────┘ └──────────┘ └─────────┘  │  │
│  │    ▲ CE          ▲ Pro         ▲ Pro        ▲ CE+Pro   │  │
│  └────┼─────────────┼─────────────┼────────────┼─────────┘  │
│       │             │             │            │             │
│  ┌────┴─────────────┴─────────────┴────────────┴─────────┐  │
│  │                    PostgreSQL                          │  │
│  │  flags, rules, evaluations, product_cards, metrics,   │  │
│  │  plugins, marketplace, audit_log                      │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌───────────────────────────────────────────────────────┐  │
│  │              Admin UI (Next.js + TypeScript)            │  │
│  │                                                       │  │
│  │  ┌─────────────┐ ┌────────────┐ ┌──────────────────┐  │  │
│  │  │ Dashboard   │ │ Flag       │ │ Plugin Manager   │  │  │
│  │  │ & Analytics │ │ Manager &  │ │ & Marketplace    │  │  │
│  │  │             │ │ Product    │ │ (install, config,│  │  │
│  │  │             │ │ Cards      │ │  develop, sell)  │  │  │
│  │  └─────────────┘ └────────────┘ └──────────────────┘  │  │
│  │                                                       │  │
│  │  ┌─────────────┐ ┌────────────┐ ┌──────────────────┐  │  │
│  │  │ Settings &  │ │ Audit Log  │ │ Developer Portal │  │  │
│  │  │ Team Mgmt   │ │ & Activity │ │ (Plugin SDK docs,│  │  │
│  │  │             │ │            │ │  API explorer)   │  │  │
│  │  └─────────────┘ └────────────┘ └──────────────────┘  │  │
│  └───────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────┘
```

## 4.2 Tech Stack / Stack Técnica

| Component / Componente | Technology / Tecnologia | Justification / Justificativa |
|----------------------|----------------------|------------------------------|
| API Server | **Go** | Single binary, high performance, low memory footprint, ideal for self-hosted / Binary único, alta performance, baixo consumo de memória |
| Admin UI | **Next.js 15 + TypeScript** | SSR/SSG for landing pages, App Router for admin SPA, i18n built-in for bilingual support / SSR/SSG para landing pages, App Router para admin SPA, i18n nativo para suporte bilíngue |
| Admin UI Components | **Next.js + Tailwind CSS + Radix UI** | Accessible, composable components with design system flexibility / Componentes acessíveis e composáveis com flexibilidade de design system |
| Admin UI State | **TanStack Query + Zustand** | Server state caching + lightweight client state / Cache de server state + client state leve |
| Admin UI Charts | **Recharts** or **Tremor** | Dashboard visualizations optimized for Next.js / Visualizações de dashboard otimizadas para Next.js |
| Landing Page | **Next.js (SSG)** | Same codebase as admin, bilingual (en/pt) with `next-intl`, SEO-optimized / Mesmo codebase do admin, bilíngue com `next-intl`, otimizado para SEO |
| Database | **PostgreSQL** | Reliable, JSONB for targeting rules, extensible / Confiável, JSONB para targeting rules |
| Cache | **In-memory (Go)** + Redis (optional) | Flag evaluation in < 1ms; Redis for clusters / Avaliação em < 1ms; Redis para clusters |
| Edge Proxy | **Go** (FlagBridge Edge) | For high volume — caches flags close to client / Para alto volume — cacheia flags perto do client |
| SDK Transport | **SSE** (Server-Sent Events) | Real-time flag updates without WebSocket complexity / Updates real-time sem complexidade de WebSocket |
| Plugin Runtime | **Go (backend)** + **Next.js (UI)** | Sandboxed plugin execution with defined hook points / Execução sandboxed com hook points definidos |
| Containerization | **Docker** | Multi-stage Dockerfile / Dockerfile multi-stage |
| Orchestration | **Docker Compose** (default) / Helm (K8s) | 90% of self-hosted users use Compose / 90% dos self-hosted users usam Compose |

### 4.2.1 Next.js Architecture Detail / Detalhe da Arquitetura Next.js

```
flagbridge-ui/
├── src/
│   ├── app/                          # Next.js App Router
│   │   ├── [locale]/                 # i18n: /en/... and /pt/...
│   │   │   ├── (marketing)/          # Landing page (SSG)
│   │   │   │   ├── page.tsx          # Homepage
│   │   │   │   ├── pricing/
│   │   │   │   ├── docs/
│   │   │   │   └── blog/
│   │   │   ├── (admin)/              # Admin panel (authenticated)
│   │   │   │   ├── dashboard/
│   │   │   │   ├── projects/
│   │   │   │   │   └── [projectSlug]/
│   │   │   │   │       ├── flags/
│   │   │   │   │       │   └── [flagKey]/
│   │   │   │   │       │       ├── page.tsx         # Flag detail
│   │   │   │   │       │       ├── product-card/    # Product context
│   │   │   │   │       │       └── metrics/         # Technical dashboard
│   │   │   │   │       ├── lifecycle/
│   │   │   │   │       └── settings/
│   │   │   │   ├── plugins/
│   │   │   │   │   ├── installed/       # Installed plugins
│   │   │   │   │   ├── marketplace/     # Browse & install
│   │   │   │   │   └── develop/         # Plugin dev tools
│   │   │   │   ├── marketplace/
│   │   │   │   │   ├── browse/          # Browse plugins
│   │   │   │   │   ├── publish/         # Publish your plugin
│   │   │   │   │   └── earnings/        # Developer earnings
│   │   │   │   ├── settings/
│   │   │   │   │   ├── team/
│   │   │   │   │   ├── billing/
│   │   │   │   │   ├── api-keys/
│   │   │   │   │   └── integrations/
│   │   │   │   └── audit-log/
│   │   │   └── (developer)/            # Developer portal
│   │   │       ├── docs/               # Plugin SDK docs
│   │   │       ├── api-explorer/       # Interactive API docs
│   │   │       └── sandbox/            # Plugin testing sandbox
│   │   └── api/                        # Next.js API routes (BFF)
│   │       ├── auth/
│   │       └── proxy/                  # Proxy to Go API
│   ├── components/
│   │   ├── ui/                         # Design system (Radix + Tailwind)
│   │   ├── flags/                      # Flag-specific components
│   │   ├── dashboard/                  # Chart components
│   │   ├── plugins/                    # Plugin UI components
│   │   └── marketplace/                # Marketplace components
│   ├── lib/
│   │   ├── api/                        # API client (TanStack Query)
│   │   ├── i18n/                       # Internationalization config
│   │   ├── plugin-host/                # Plugin UI host runtime
│   │   └── auth/                       # Auth utilities
│   ├── messages/
│   │   ├── en.json                     # English translations
│   │   └── pt.json                  # Brazilian Portuguese translations
│   └── styles/
│       └── globals.css                 # Tailwind + custom tokens
├── public/
├── next.config.ts
├── tailwind.config.ts
└── package.json
```

**Key architectural decisions / Decisões arquiteturais importantes:**

1. **Route Groups**: `(marketing)` for SSG landing pages, `(admin)` for authenticated app, `(developer)` for developer portal — all in the same Next.js app
2. **i18n via `next-intl`**: All routes prefixed with locale (`/en/dashboard`, `/pt/dashboard`), SSG pages fully bilingual for SEO
3. **BFF Pattern**: Next.js API routes act as Backend-for-Frontend, proxying to Go API with auth token injection
4. **Plugin UI Host**: Plugins render into designated `<PluginSlot />` components using a sandboxed iframe or Module Federation approach

---

## 4.3 Data Model / Modelo de Dados

```sql
-- ========================
-- CORE TABLES (CE)
-- ========================

-- Projects / Projetos
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
    name        VARCHAR(100) NOT NULL,
    slug        VARCHAR(100) NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(project_id, slug)
);

-- Feature Flags
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

-- Flag States (per environment)
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

-- Audit Log
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
-- PLUGIN TABLES (CE)
-- ========================

-- Installed Plugins / Plugins Instalados
CREATE TABLE plugins (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug            VARCHAR(255) NOT NULL UNIQUE,
    name            VARCHAR(255) NOT NULL,
    version         VARCHAR(50) NOT NULL,
    author          VARCHAR(255),
    description     TEXT,
    source          VARCHAR(20) DEFAULT 'marketplace', -- marketplace, local, git
    manifest        JSONB NOT NULL, -- Full plugin manifest
    config          JSONB DEFAULT '{}', -- User-configured settings
    enabled         BOOLEAN DEFAULT TRUE,
    installed_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Plugin Hook Registrations / Registro de Hooks de Plugins
CREATE TABLE plugin_hooks (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plugin_id   UUID REFERENCES plugins(id) ON DELETE CASCADE,
    hook_point  VARCHAR(100) NOT NULL, -- e.g. 'flag.beforeEvaluate', 'ui.flagDetail.sidebar'
    handler     VARCHAR(255) NOT NULL, -- Handler reference in the plugin
    priority    INTEGER DEFAULT 100,
    enabled     BOOLEAN DEFAULT TRUE
);

-- ========================
-- PRO TABLES
-- ========================

-- Product Context Cards (Pro)
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

-- Flag Evaluation Metrics (Pro)
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

-- Lifecycle Rules (Pro)
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
-- MARKETPLACE TABLES (Pro/Enterprise)
-- ========================

-- Marketplace Listings / Listagens do Marketplace
CREATE TABLE marketplace_listings (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plugin_slug     VARCHAR(255) NOT NULL UNIQUE,
    developer_id    UUID REFERENCES users(id),
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    long_description TEXT,
    category        VARCHAR(50), -- integration, analytics, security, ui, automation
    tags            TEXT[],
    icon_url        TEXT,
    screenshots     TEXT[],
    repository_url  TEXT,
    documentation_url TEXT,
    pricing_type    VARCHAR(20) DEFAULT 'free', -- free, one_time, subscription
    price_cents     INTEGER DEFAULT 0,
    currency        VARCHAR(3) DEFAULT 'USD',
    status          VARCHAR(20) DEFAULT 'draft', -- draft, in_review, published, suspended
    downloads       INTEGER DEFAULT 0,
    avg_rating      DECIMAL(2,1) DEFAULT 0,
    review_count    INTEGER DEFAULT 0,
    published_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Marketplace Versions / Versões no Marketplace
CREATE TABLE marketplace_versions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    listing_id      UUID REFERENCES marketplace_listings(id),
    version         VARCHAR(50) NOT NULL,
    changelog       TEXT,
    min_flagbridge  VARCHAR(50), -- Minimum FlagBridge version
    package_url     TEXT NOT NULL,
    package_hash    VARCHAR(64) NOT NULL, -- SHA-256
    status          VARCHAR(20) DEFAULT 'pending', -- pending, approved, rejected
    reviewed_by     UUID REFERENCES users(id),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(listing_id, version)
);

-- Reviews / Avaliações
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

-- Plugin Purchases / Compras de Plugins
CREATE TABLE marketplace_purchases (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    listing_id      UUID REFERENCES marketplace_listings(id),
    buyer_id        UUID REFERENCES users(id),
    price_cents     INTEGER NOT NULL,
    currency        VARCHAR(3) NOT NULL,
    payment_provider VARCHAR(50), -- stripe
    payment_id      VARCHAR(255),
    status          VARCHAR(20) DEFAULT 'completed',
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Developer Payouts / Pagamentos a Desenvolvedores
CREATE TABLE developer_payouts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    developer_id    UUID REFERENCES users(id),
    amount_cents    INTEGER NOT NULL,
    currency        VARCHAR(3) NOT NULL,
    period_start    TIMESTAMPTZ NOT NULL,
    period_end      TIMESTAMPTZ NOT NULL,
    status          VARCHAR(20) DEFAULT 'pending', -- pending, processing, paid
    payout_ref      VARCHAR(255),
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
```

## 4.4 SDK Design

### 🇺🇸 English / 🇧🇷 Português

SDKs follow the OpenFeature standard, with a proprietary layer for FlagBridge-specific features.

```typescript
// @flagbridge/sdk-node — Usage example / Exemplo de uso

import { FlagBridge } from '@flagbridge/sdk-node';

const fb = new FlagBridge({
  serverUrl: 'https://flags.mycompany.com',
  apiKey: 'fb_sk_...',
  environment: 'production',
});

// Simple evaluation (boolean)
const isEnabled = await fb.isEnabled('checkout-v2', {
  userId: 'user_123',
  attributes: { plan: 'pro', country: 'BR' },
});

// Variant evaluation (string)
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
// @flagbridge/sdk-go — Usage example / Exemplo de uso

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

## 4.5 API Endpoints

```
# ========================
# FLAG MANAGEMENT (CE)
# ========================
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

# SSE (real-time updates)
GET    /api/v1/sse/:environment

# ========================
# PRODUCT CONTEXT (Pro)
# ========================
GET    /api/v1/projects/:project/flags/:key/product-card
PUT    /api/v1/projects/:project/flags/:key/product-card

# ========================
# METRICS & DASHBOARD (Pro)
# ========================
GET    /api/v1/projects/:project/flags/:key/metrics
GET    /api/v1/projects/:project/dashboard/overview
GET    /api/v1/projects/:project/lifecycle/stale

# ========================
# PLUGIN SYSTEM (CE+Pro)
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
POST   /api/v1/marketplace/listings            # Publish
PATCH  /api/v1/marketplace/listings/:slug       # Update
POST   /api/v1/marketplace/listings/:slug/review
GET    /api/v1/marketplace/developer/earnings
POST   /api/v1/marketplace/purchase/:slug

# ========================
# ADMIN
# ========================
GET    /api/v1/audit-log
GET    /api/v1/health
```

## 4.6 Deployment Infrastructure / Infraestrutura de Deploy

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
      - NEXT_PUBLIC_DEFAULT_LOCALE=en  # or pt
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

**Self-hosted with Pro Plugin:**
```yaml
# docker-compose.pro.yml (override)
services:
  flagbridge:
    image: flagbridge/flagbridge-pro:latest
    environment:
      - FB_LICENSE_KEY=fb_lic_...
```

**Upgrade path:** `docker compose -f docker-compose.yml -f docker-compose.pro.yml up -d`

---

# 5. Admin Panel / Painel Admin

## 🇺🇸 Admin Panel — Technical & Product Documentation

The FlagBridge Admin is a full-featured Next.js application that serves as the central hub for flag management, product intelligence, plugin management, and marketplace access.

### 5.1 Admin Sections

#### Dashboard
- **Overview cards**: Total flags, active flags, stale flags, flags by status
- **Activity feed**: Recent flag changes, deployments, plugin installs
- **Health metrics** (Pro): Evaluation volume, error rates, SDK connection count
- **Technical debt score** (Pro): Aggregate score based on stale flags, flags without owners, flags without product cards

#### Flag Manager
- **Flag list**: Filterable by project, environment, status, tags, owner
- **Flag detail page**: Toggle, targeting rules editor (visual builder), environment comparison
- **Product Context Card** (tab within flag detail): Hypothesis, success metrics, go/no-go criteria, decision history, OKR link
- **Metrics tab** (Pro): Adoption chart, error rate chart, latency impact, evaluation breakdown by variant
- **Lifecycle tab** (Pro): Timeline from creation to archival, cleanup reminders, code references

#### Plugin Manager
- **Installed plugins**: List with enable/disable toggle, config panel, health status
- **Marketplace browser**: Search, filter by category, install with one click
- **Plugin development**: Sandbox environment, logs, hot-reload for local dev

#### Settings
- **Team management**: Invite members, assign roles (Admin, Editor, Viewer)
- **API keys**: Create, rotate, revoke keys per environment
- **Integrations**: Slack, Linear, Jira, GitHub, GitLab webhooks
- **Billing** (Pro/SaaS): Plan management, invoices, usage metrics

#### Audit Log
- **Comprehensive history**: Every flag change, toggle, user action, plugin install
- **Filterable**: By user, action type, date range, entity
- **Diff viewer** (Pro): Side-by-side comparison of flag state changes

#### Developer Portal
- **Plugin SDK documentation**: Full API reference, guides, tutorials
- **Interactive API explorer**: Try endpoints with live data (Swagger/OpenAPI)
- **Plugin sandbox**: Test your plugin in an isolated environment
- **Marketplace publisher**: Submit, manage, and track your plugins

### 5.2 Bilingual Support / Suporte Bilíngue

The entire Admin UI uses `next-intl` for full internationalization:

```typescript
// src/messages/en.json (excerpt)
{
  "dashboard": {
    "title": "Dashboard",
    "totalFlags": "Total Flags",
    "activeFlags": "Active Flags",
    "staleFlags": "Stale Flags",
    "technicalDebt": "Technical Debt Score"
  },
  "flags": {
    "create": "Create Flag",
    "status": {
      "planning": "Planning",
      "active": "Active",
      "rolledOut": "Rolled Out",
      "archived": "Archived"
    }
  },
  "plugins": {
    "installed": "Installed Plugins",
    "marketplace": "Marketplace",
    "develop": "Develop Plugins"
  }
}
```

```typescript
// src/messages/pt.json (excerpt)
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
  }
}
```

---

## 🇧🇷 Painel Admin — Documentação Técnica e de Produto

O FlagBridge Admin é uma aplicação Next.js completa que serve como hub central para gerenciamento de flags, inteligência de produto, gerenciamento de plugins e acesso ao marketplace.

### 5.1 Seções do Admin

#### Dashboard
- **Cards de overview**: Total de flags, flags ativas, flags obsoletas, flags por status
- **Feed de atividade**: Mudanças recentes em flags, deploys, instalações de plugins
- **Métricas de saúde** (Pro): Volume de avaliações, taxas de erro, contagem de conexões SDK
- **Score de dívida técnica** (Pro): Score agregado baseado em flags obsoletas, flags sem owner, flags sem product cards

#### Gerenciador de Flags
- **Lista de flags**: Filtrável por projeto, environment, status, tags, owner
- **Página de detalhe**: Toggle, editor visual de targeting rules, comparação entre environments
- **Product Context Card** (aba dentro do detalhe): Hipótese, métricas de sucesso, critérios go/no-go, histórico de decisões, link para OKR
- **Aba de Métricas** (Pro): Gráfico de adoção, gráfico de taxa de erro, impacto de latência, breakdown por variante
- **Aba de Lifecycle** (Pro): Timeline da criação ao archival, lembretes de cleanup, referências no código

#### Gerenciador de Plugins
- **Plugins instalados**: Lista com toggle ativar/desativar, painel de configuração, status de saúde
- **Navegador do Marketplace**: Busca, filtro por categoria, instalação com um clique
- **Desenvolvimento de plugins**: Ambiente sandbox, logs, hot-reload para dev local

#### Configurações
- **Gerenciamento de time**: Convites, atribuição de roles (Admin, Editor, Viewer)
- **API keys**: Criação, rotação, revogação de chaves por environment
- **Integrações**: Slack, Linear, Jira, GitHub, GitLab webhooks
- **Billing** (Pro/SaaS): Gestão de plano, faturas, métricas de uso

#### Audit Log
- **Histórico completo**: Toda mudança de flag, toggle, ação de usuário, instalação de plugin
- **Filtrável**: Por usuário, tipo de ação, intervalo de data, entidade
- **Visualizador de diff** (Pro): Comparação lado a lado de mudanças no estado das flags

#### Portal do Desenvolvedor
- **Documentação do Plugin SDK**: Referência completa da API, guias, tutoriais
- **API explorer interativo**: Testar endpoints com dados reais (Swagger/OpenAPI)
- **Plugin sandbox**: Testar plugin em ambiente isolado
- **Publisher do Marketplace**: Submeter, gerenciar e acompanhar seus plugins

---

# 6. Plugin System & SDK / Sistema de Plugins & SDK

## 🇺🇸 English

### 6.1 Plugin Architecture

FlagBridge plugins are self-contained packages that extend the platform through well-defined hook points. Plugins can extend both the backend (Go API) and the frontend (Next.js Admin UI).

#### Plugin Manifest (`flagbridge-plugin.json`)

```json
{
  "name": "@myplugin/slack-alerts",
  "version": "1.0.0",
  "displayName": "Slack Alerts for Flag Changes",
  "description": "Send rich Slack notifications when flags are toggled or modified",
  "author": {
    "name": "Jane Developer",
    "email": "jane@example.com",
    "url": "https://github.com/janedeveloper"
  },
  "license": "MIT",
  "minFlagBridgeVersion": "1.0.0",
  "category": "integration",
  "tags": ["slack", "notifications", "alerts"],
  "pricing": {
    "type": "free"
  },
  "hooks": {
    "backend": [
      {
        "point": "flag.afterToggle",
        "handler": "handlers/on-flag-toggle"
      },
      {
        "point": "flag.afterUpdate",
        "handler": "handlers/on-flag-update"
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
        "label": "Slack Webhook URL",
        "required": true,
        "secret": true
      },
      "channel": {
        "type": "string",
        "label": "Default Channel",
        "default": "#feature-flags"
      },
      "notifyOnToggle": {
        "type": "boolean",
        "label": "Notify on flag toggle",
        "default": true
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

#### Backend Hook Points

| Hook Point | Trigger | Use Cases |
|-----------|---------|-----------|
| `flag.beforeEvaluate` | Before flag evaluation | Custom targeting logic, A/B routing |
| `flag.afterEvaluate` | After flag evaluation | Logging, analytics tracking |
| `flag.beforeToggle` | Before flag is toggled | Approval gates, validation |
| `flag.afterToggle` | After flag is toggled | Notifications, sync to external systems |
| `flag.beforeUpdate` | Before flag metadata changes | Validation, required fields |
| `flag.afterUpdate` | After flag metadata changes | Audit, sync, notifications |
| `flag.onStaleDetected` | When stale flag detected | Custom cleanup actions |
| `project.onCreate` | New project created | Auto-setup, templates |
| `api.middleware` | Every API request | Auth extensions, rate limiting |
| `evaluation.metrics` | Evaluation metrics collected | Custom analytics, export |

#### Frontend (UI) Hook Points

| Hook Point | Location | Use Cases |
|-----------|----------|-----------|
| `dashboard.widget` | Dashboard page | Custom metric cards, charts |
| `flag.detail.sidebar` | Flag detail sidebar | Related info, external data |
| `flag.detail.tab` | Flag detail tabs | Custom tabs (e.g. analytics) |
| `flag.list.column` | Flag list table | Custom columns |
| `settings.integrations.panel` | Settings > Integrations | Plugin configuration UI |
| `navigation.sidebar` | Main sidebar nav | New navigation items |
| `global.banner` | Top of page | Alerts, announcements |

#### Plugin SDK (for developers)

```typescript
// @flagbridge/plugin-sdk — Building a backend hook

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
              text: `🚩 *${flag.name}* was ${newState ? 'enabled' : 'disabled'}` +
                    ` in *${environment.name}* by ${toggledBy.name}`
            }
          }
        ]
      })
    });
  }
}
```

```tsx
// @flagbridge/plugin-sdk — Building a UI extension (Next.js component)

import { PluginUIComponent, useFlagContext, usePluginConfig } from '@flagbridge/plugin-sdk/react';

export const SlackActivityWidget: PluginUIComponent = () => {
  const flag = useFlagContext();
  const config = usePluginConfig();

  return (
    <div className="p-4 border rounded-lg">
      <h3 className="font-semibold text-sm text-gray-500">
        Slack Notifications
      </h3>
      <p className="text-sm mt-1">
        Posting to <code>{config.channel}</code>
      </p>
      <p className="text-xs text-gray-400 mt-2">
        Last notified: {flag.lastToggleAt}
      </p>
    </div>
  );
};
```

#### Plugin CLI

```bash
# Scaffold a new plugin
npx @flagbridge/create-plugin my-awesome-plugin

# Start local development with hot-reload
cd my-awesome-plugin
flagbridge plugin dev

# Run tests
flagbridge plugin test

# Build for distribution
flagbridge plugin build

# Publish to marketplace
flagbridge plugin publish
```

### 6.2 Plugin Security Model

- **Sandboxed execution**: Backend plugins run in isolated goroutines with limited syscall access
- **Permission system**: Plugins declare required permissions; admin must approve
- **UI sandboxing**: Frontend plugins render in sandboxed iframes or use strict CSP
- **Config encryption**: Secret config values (API keys, tokens) encrypted at rest
- **Code signing**: Marketplace plugins must be signed; hash verified on install
- **Review process**: All marketplace submissions reviewed before publishing

---

## 🇧🇷 Português (Brasil)

### 6.1 Arquitetura de Plugins

Plugins do FlagBridge são pacotes autocontidos que estendem a plataforma através de hook points bem definidos. Plugins podem estender tanto o backend (Go API) quanto o frontend (Next.js Admin UI).

#### Manifesto do Plugin (`flagbridge-plugin.json`)

*(Mesmo formato JSON mostrado na seção em inglês acima — o manifesto é sempre em inglês por convenção técnica, mas `displayName` e `description` suportam i18n)*

```json
{
  "displayName": {
    "en": "Slack Alerts for Flag Changes",
    "pt": "Alertas Slack para Mudanças em Flags"
  },
  "description": {
    "en": "Send rich Slack notifications when flags are toggled or modified",
    "pt": "Envie notificações ricas no Slack quando flags são alternadas ou modificadas"
  }
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

#### Hook Points de Frontend (UI)

| Hook Point | Local | Casos de Uso |
|-----------|-------|-------------|
| `dashboard.widget` | Página do Dashboard | Cards de métricas customizados, gráficos |
| `flag.detail.sidebar` | Sidebar do detalhe da flag | Info relacionada, dados externos |
| `flag.detail.tab` | Tabs do detalhe da flag | Tabs customizadas (ex: analytics) |
| `flag.list.column` | Tabela da lista de flags | Colunas customizadas |
| `settings.integrations.panel` | Configurações > Integrações | UI de configuração do plugin |
| `navigation.sidebar` | Sidebar de navegação | Novos itens de navegação |
| `global.banner` | Topo da página | Alertas, anúncios |

### 6.2 Modelo de Segurança de Plugins

- **Execução sandboxed**: Plugins de backend rodam em goroutines isoladas com acesso limitado a syscalls
- **Sistema de permissões**: Plugins declaram permissões necessárias; admin deve aprovar
- **Sandboxing de UI**: Plugins de frontend renderizam em iframes sandboxed ou com CSP restrito
- **Criptografia de config**: Valores secretos (API keys, tokens) criptografados em repouso
- **Assinatura de código**: Plugins do marketplace devem ser assinados; hash verificado na instalação
- **Processo de revisão**: Todas as submissões ao marketplace são revisadas antes da publicação

---

# 7. Plugin Marketplace / Marketplace de Plugins

## 🇺🇸 English

### 7.1 Vision

The FlagBridge Plugin Marketplace transforms FlagBridge from a product into a **platform ecosystem** — similar to Shopify App Store, WordPress Plugins, or Figma Community. Third-party developers can build, publish, and monetize plugins that extend FlagBridge for specific use cases.

### 7.2 Marketplace Features

#### For Plugin Users (Buyers)
- **Browse & Search**: Filter by category (integration, analytics, security, UI, automation), rating, price, compatibility
- **One-click install**: Install directly from the Admin UI, both for self-hosted and SaaS
- **Reviews & Ratings**: 5-star rating with written reviews
- **Compatibility check**: Automatic verification against current FlagBridge version
- **Plugin configuration**: Configure installed plugins through Admin UI
- **Auto-updates**: Optional automatic updates for installed plugins (self-hosted respects user control)

#### For Plugin Developers (Sellers)
- **Developer Portal**: Full documentation, API explorer, plugin SDK
- **Plugin Sandbox**: Test environment to develop and debug plugins
- **Publish workflow**: Submit → automated checks → manual review → published
- **Pricing options**: Free, one-time purchase, or monthly subscription
- **Revenue split**: 80% developer / 20% FlagBridge (industry standard)
- **Earnings dashboard**: Track installs, revenue, ratings, support tickets
- **Developer tiers**: Verified Developer badge for trusted publishers
- **Analytics**: Installation stats, usage metrics, error rates

#### For FlagBridge (Platform)
- **Review pipeline**: Automated security scan + manual code review
- **Revenue stream**: 20% commission on paid plugins
- **Ecosystem growth**: More plugins → more value → more users → more plugin developers
- **Quality control**: Minimum quality standards, security requirements, regular audits

### 7.3 Marketplace Categories

| Category | Examples |
|----------|---------|
| **Integrations** | Slack, Discord, Linear, Jira, PagerDuty, Datadog |
| **Analytics** | Mixpanel connector, Amplitude connector, custom dashboards |
| **Security** | Advanced RBAC, IP whitelisting, compliance reports |
| **UI Extensions** | Custom flag detail widgets, theme packs, dashboard layouts |
| **Automation** | Auto-remediation, scheduled rollouts, CI/CD bridges |
| **Data** | Export connectors (BigQuery, Snowflake), migration tools |

### 7.4 Revenue Model

```
Plugin Sale ($10.00)
├── Developer receives: $8.00 (80%)
├── FlagBridge receives: $2.00 (20%)
└── Payment processing: deducted from FlagBridge share

Subscription Plugin ($5.00/mo)
├── Developer receives: $4.00/mo (80%)
├── FlagBridge receives: $1.00/mo (20%)
└── Recurring billing managed by Stripe Connect
```

### 7.5 Technical Implementation

- **Payment**: Stripe Connect for marketplace payments and developer payouts
- **Package registry**: Private npm-like registry for plugin packages
- **CI/CD**: Automated testing pipeline for submitted plugins
- **CDN**: Plugin packages served via CDN for fast installation worldwide
- **Versioning**: Semantic versioning enforced, with compatibility matrix

---

## 🇧🇷 Português (Brasil)

### 7.1 Visão

O FlagBridge Plugin Marketplace transforma o FlagBridge de um produto em um **ecossistema de plataforma** — similar à Shopify App Store, WordPress Plugins ou Figma Community. Desenvolvedores third-party podem construir, publicar e monetizar plugins que estendem o FlagBridge para casos de uso específicos.

### 7.2 Features do Marketplace

#### Para Usuários de Plugins (Compradores)
- **Navegar & Buscar**: Filtrar por categoria (integração, analytics, segurança, UI, automação), rating, preço, compatibilidade
- **Instalação com um clique**: Instalar direto do Admin UI, tanto para self-hosted quanto SaaS
- **Reviews & Ratings**: Avaliação 5 estrelas com reviews escritas
- **Verificação de compatibilidade**: Verificação automática contra versão atual do FlagBridge
- **Configuração de plugins**: Configurar plugins instalados pelo Admin UI
- **Auto-updates**: Atualizações automáticas opcionais (self-hosted respeita controle do usuário)

#### Para Desenvolvedores de Plugins (Vendedores)
- **Portal do Desenvolvedor**: Documentação completa, API explorer, Plugin SDK
- **Plugin Sandbox**: Ambiente de teste para desenvolver e debugar plugins
- **Workflow de publicação**: Submit → checks automatizados → review manual → publicado
- **Opções de preço**: Gratuito, compra única, ou assinatura mensal
- **Split de receita**: 80% desenvolvedor / 20% FlagBridge (padrão da indústria)
- **Dashboard de ganhos**: Acompanhar instalações, receita, ratings, tickets de suporte
- **Tiers de desenvolvedor**: Badge de Desenvolvedor Verificado para publishers confiáveis
- **Analytics**: Estatísticas de instalação, métricas de uso, taxas de erro

#### Para FlagBridge (Plataforma)
- **Pipeline de review**: Scan de segurança automatizado + review manual de código
- **Fonte de receita**: 20% de comissão em plugins pagos
- **Crescimento do ecossistema**: Mais plugins → mais valor → mais usuários → mais devs de plugins
- **Controle de qualidade**: Padrões mínimos de qualidade, requisitos de segurança, auditorias regulares

### 7.3 Categorias do Marketplace

| Categoria | Exemplos |
|-----------|---------|
| **Integrações** | Slack, Discord, Linear, Jira, PagerDuty, Datadog |
| **Analytics** | Conector Mixpanel, conector Amplitude, dashboards customizados |
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

# 8. Go-to-Market

## 🇺🇸 English

### 8.1 Target Audience

**Primary:** Startups and scale-ups (50-500 devs) that:
- Already use feature flags (Unleash CE, in-house) but lack observability and product context
- Cannot afford LaunchDarkly (~$72k/year median)
- Value self-hosted and data sovereignty
- Want extensibility without vendor lock-in

**Secondary:** Product teams at larger companies that:
- Want visibility into flag states without depending on engineering
- Need to link flags to OKRs and roadmap
- Want a plugin ecosystem for custom integrations

### 8.2 Pricing

| Plan | Price | Target |
|------|-------|--------|
| **Community** | Free (forever) | Individual devs, small teams, evaluation |
| **Pro** (self-hosted) | $29/mo (≤10 seats) / $79/mo (≤50 seats) / $149/mo (unlimited) | Scale-ups, mid-size teams |
| **Pro** (SaaS) | $49/mo (≤10 seats) / $99/mo (≤50 seats) / $199/mo (unlimited) | Teams that don't want to manage infra |
| **Enterprise** | Custom (from $500/mo) | Large companies, compliance, SLA |

**Flat-rate per seats** — no MAU or service connection surprises. Direct differentiator against LaunchDarkly.

Self-hosted is cheaper than SaaS because the customer pays for infra. This incentivizes self-hosted (less hosting cost for us) and attracts the data sovereignty audience.

### 8.3 Acquisition Channels

**Developer-first growth:**
1. **GitHub** — well-documented open-source repo, attractive README, contributing guide
2. **Hacker News / Reddit / Dev.to** — launch posts
3. **Honest comparisons** — blog posts: "FlagBridge vs Unleash vs LaunchDarkly"
4. **Docker Hub** — optimized official image, one-liner setup
5. **OpenFeature ecosystem** — official provider listed on OpenFeature site
6. **YouTube / blog** — tutorials in English and Portuguese

**Product-led growth:**
- CE is genuinely useful → team adopts → grows → needs observability/governance → upgrades to Pro
- Discrete banner in Admin UI: "Unlock product dashboards and lifecycle automation"
- 14-day Pro trial (self-hosted and SaaS)

**Ecosystem growth (post-marketplace):**
- Plugin developers attract their own users to FlagBridge
- Marketplace becomes a discovery channel
- Revenue share incentivizes quality plugins

### 8.4 Success Metrics (Year 1)

| Metric | 6 months | 12 months |
|--------|----------|-----------|
| GitHub stars | 500 | 2,000 |
| Docker pulls | 5,000 | 25,000 |
| Active CE installs | 200 | 1,000 |
| Pro paying customers | 10 | 50 |
| MRR | $500 | $3,000 |
| Community members (Discord) | 100 | 500 |
| Published plugins | — | 20 |
| Marketplace GMV | — | $500/mo |

---

## 🇧🇷 Português (Brasil)

### 8.1 Público-alvo

**Primário:** Startups e scale-ups (50-500 devs) que:
- Já usam feature flags (Unleash CE, in-house) mas sentem falta de observabilidade e contexto de produto
- Não podem pagar LaunchDarkly (~$72k/ano median)
- Valorizam self-hosted e data sovereignty
- Querem extensibilidade sem vendor lock-in

**Secundário:** Times de produto em empresas maiores que:
- Querem visibilidade sobre o estado das flags sem depender de engenharia
- Precisam vincular flags a OKRs e roadmap
- Querem um ecossistema de plugins para integrações customizadas

### 8.2 Pricing

| Plano | Preço | Target |
|-------|-------|--------|
| **Community** | Grátis (forever) | Devs individuais, times pequenos, avaliação |
| **Pro** (self-hosted) | $29/mês (≤10 seats) / $79/mês (≤50 seats) / $149/mês (unlimited) | Scale-ups, times médios |
| **Pro** (SaaS) | $49/mês (≤10 seats) / $99/mês (≤50 seats) / $199/mês (unlimited) | Times que não querem gerenciar infra |
| **Enterprise** | Custom (a partir de $500/mês) | Grandes empresas, compliance, SLA |

**Modelo flat-rate por seats** — sem surpresas de MAU ou service connections. Diferencial direto contra LaunchDarkly.

### 8.3 Canais de Aquisição

**Developer-first growth:**
1. **GitHub** — repo open-source bem documentado, README atraente
2. **Hacker News / Reddit / Dev.to** — posts de lançamento
3. **Comparações honestas** — blog posts: "FlagBridge vs Unleash vs LaunchDarkly"
4. **Docker Hub** — imagem oficial otimizada, one-liner de setup
5. **OpenFeature ecosystem** — provider oficial listado no site do OpenFeature
6. **YouTube / blog** — tutoriais em inglês e português

**Product-led growth:**
- CE é genuinamente útil → time adota → cresce → precisa de observabilidade/governance → upgrade Pro
- Banner discreto no Admin UI
- Trial de 14 dias do Pro

**Crescimento do ecossistema (pós-marketplace):**
- Desenvolvedores de plugins atraem seus próprios usuários para o FlagBridge
- Marketplace se torna canal de descoberta
- Revenue share incentiva plugins de qualidade

### 8.4 Métricas de Sucesso (Ano 1)

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

---

# 9. Roadmap

## 🇺🇸 English / 🇧🇷 Português

### Q2 2026 — Foundation / Fundação
- [ ] Public repo on GitHub / Repo público no GitHub
- [ ] Core Go API server: flag management + eval engine
- [ ] Next.js Admin UI: flag CRUD, dashboard, settings
- [ ] Bilingual support (en/pt) from day 1 / Suporte bilíngue desde o dia 1
- [ ] Node.js SDK + Go SDK
- [ ] Docker image + Docker Compose
- [ ] Docs site (docs.flagbridge.io) — bilingual / bilíngue
- [ ] OpenFeature Provider (Node.js)
- [ ] Plugin runtime (basic) — load/enable/disable plugins

### Q3 2026 — Community Launch / Lançamento Community
- [ ] Product Context Cards (CE basic version / versão CE básica)
- [ ] CLI (`flagbridge`)
- [ ] React SDK (client-side) + Python SDK
- [ ] Helm chart for Kubernetes
- [ ] Webhooks
- [ ] Plugin CLI (`flagbridge plugin create`)
- [ ] Plugin SDK v1 (backend hooks + UI extensions)
- [ ] 3-5 official plugins (Slack, Linear, GitHub)
- [ ] Launch: Hacker News + Product Hunt + Reddit

### Q4 2026 — Pro Launch / Lançamento Pro
- [ ] Pro Plugin architecture / Arquitetura do Plugin Pro
- [ ] Technical Dashboard (metrics, observability)
- [ ] Advanced Product Cards (hypothesis, KPIs, decision workflow)
- [ ] Lifecycle automation (stale detection, cleanup alerts)
- [ ] SSO (OIDC)
- [ ] Billing system + license key management
- [ ] Integrations: Slack, Linear, Jira (bidirectional / bidirecional)
- [ ] Plugin Marketplace v1 (browse, install, free plugins only / apenas plugins gratuitos)

### Q1 2027 — Scale & Marketplace / Escala & Marketplace
- [ ] FlagBridge Edge (proxy for high volume / proxy para alto volume)
- [ ] Basic A/B experimentation
- [ ] Granular RBAC
- [ ] Change requests + approval workflow
- [ ] GitHub/GitLab integration (code cleanup suggestions)
- [ ] SaaS managed offering
- [ ] Plugin Marketplace v2: paid plugins, Stripe Connect, developer payouts
- [ ] Developer Portal with sandbox
- [ ] Landing page redesign (bilingual, showcasing marketplace / bilíngue, destacando marketplace)

### Q2 2027 — Ecosystem / Ecossistema
- [ ] Plugin Marketplace v3: subscriptions, verified developers, analytics
- [ ] White-label option (Enterprise)
- [ ] Advanced experimentation (statistical significance)
- [ ] SOC 2 compliance (start / início)
- [ ] Mobile-responsive Admin UI
- [ ] Community plugin hackathon

---

# 10. Risks & Mitigations / Riscos e Mitigações

## 🇺🇸 English

| Risk | Impact | Mitigation |
|------|--------|-----------|
| Saturated FF tools market | High | Clear differentiation via Product Context + Plugin Marketplace — no one else does both |
| Unleash adds similar features | Medium | Speed of execution + product-first vs infra-first focus |
| Difficulty monetizing open-source | High | CE is genuinely useful but Pro solves real pains (observability, governance); Marketplace creates additional revenue stream |
| Complexity of maintaining multiple SDKs | Medium | Start with 4 SDKs (Go, Node, React, Python); OpenFeature reduces surface area |
| Limited time (side project) | High | Lean MVP; prioritize core → launch → iterate with feedback |
| Plugin security vulnerabilities | Medium | Sandboxed execution, permission system, mandatory review process, code signing |
| Marketplace chicken-and-egg problem | High | Build 5-10 official plugins first; offer early developer incentives (reduced commission, featured placement) |
| Bilingual content maintenance | Low | `next-intl` with structured JSON; all content created in both languages from day 1 |

## 🇧🇷 Português (Brasil)

| Risco | Impacto | Mitigação |
|-------|---------|-----------|
| Mercado saturado de FF tools | Alto | Diferenciação clara via Product Context + Plugin Marketplace — ninguém faz os dois |
| Unleash adiciona features similares | Médio | Velocidade de execução + foco product-first vs infra-first |
| Dificuldade de monetizar open-source | Alto | CE genuinamente útil mas Pro resolve dores reais; Marketplace cria fonte de receita adicional |
| Complexidade de manter múltiplos SDKs | Médio | Começar com 4 SDKs; OpenFeature reduz superfície |
| Tempo limitado (side project) | Alto | MVP enxuto; priorizar core → launch → iterate com feedback |
| Vulnerabilidades de segurança em plugins | Médio | Execução sandboxed, sistema de permissões, review obrigatório, assinatura de código |
| Problema de chicken-and-egg do marketplace | Alto | Construir 5-10 plugins oficiais primeiro; oferecer incentivos a early developers |
| Manutenção de conteúdo bilíngue | Baixo | `next-intl` com JSON estruturado; todo conteúdo criado em ambas as línguas desde o dia 1 |

---

# Next Steps / Próximos Passos

## 🇺🇸 English

1. **Register domains** — flagbridge.io and flagbridge.dev
2. **Create GitHub org** — github.com/flagbridge
3. **Reserve npm scope** — @flagbridge
4. **Init monorepo** — `flagbridge/flagbridge` (api + ui + docs + plugin-sdk)
5. **Prototype** — Go API server with flag CRUD + eval endpoint
6. **Design** — Next.js Admin UI wireframes (dashboard + flag detail + product card + plugin manager)
7. **Landing page** — Bilingual Next.js SSG site with positioning, features, pricing
8. **Plugin SDK scaffold** — Define hook points, manifest format, CLI

## 🇧🇷 Português (Brasil)

1. **Registrar domínios** — flagbridge.io e flagbridge.dev
2. **Criar org no GitHub** — github.com/flagbridge
3. **Reservar scope no npm** — @flagbridge
4. **Iniciar monorepo** — `flagbridge/flagbridge` (api + ui + docs + plugin-sdk)
5. **Prototipar** — API server Go com flag CRUD + eval endpoint
6. **Design** — Wireframes do Admin UI Next.js (dashboard + detalhe da flag + product card + plugin manager)
7. **Landing page** — Site Next.js SSG bilíngue com posicionamento, features, pricing
8. **Scaffold do Plugin SDK** — Definir hook points, formato do manifesto, CLI

---

*FlagBridge — Bridging the gap between feature flags and product strategy.*  
*FlagBridge — A ponte entre feature flags e estratégia de produto.*
