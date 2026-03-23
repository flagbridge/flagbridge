# FlagBridge — Product Document

> **"The bridge between feature flags and product strategy."**

**Version:** 3.0  
**Author:** Gabriel Gripp  
**Date:** March 2026  
**Status:** Draft / MVP Planning

---

## Table of Contents

1. [Overview](#1-overview)
2. [Competitive Analysis](#2-competitive-analysis)
3. [MVP Definition — Community vs Pro](#3-mvp-definition--community-vs-pro)
4. [Technical Architecture](#4-technical-architecture)
5. [Admin Panel — Technical & Product Documentation](#5-admin-panel--technical--product-documentation)
6. [Plugin System & SDK](#6-plugin-system--sdk)
7. [Plugin Marketplace](#7-plugin-marketplace)
8. [Integrations Ecosystem](#8-integrations-ecosystem)
9. [Go-to-Market](#9-go-to-market)
10. [Roadmap](#10-roadmap)
11. [Risks & Mitigations](#11-risks--mitigations)

---

## 1. Overview

FlagBridge is an open-core Feature Flag Management platform that goes beyond simple toggle on/off: it connects feature flags to product planning, technical observability, lifecycle management, and a rich integrations ecosystem. The key differentiator is being **product-first with solid infrastructure** — every flag has business context, impact metrics, clear lifecycle rules, and deep connectivity to the tools teams already use.

### 1.1 The Problem

The current feature flag market is dominated by **infra-first** tools:

- **Engineering** creates flags without product context — nobody knows *why* a flag exists
- **Product** has no visibility into flag states or their real impact
- **Zombie flags** accumulate (100% ON for months, never removed), generating technical debt
- **No tool** natively connects a flag to a product hypothesis, experiment, or OKR
- Enterprise tools like LaunchDarkly charge median annual contracts of ~$72k (Vendr data, 2026), unaffordable for most companies
- **No tool** offers a plugin ecosystem — customization requires forking or waiting for vendor roadmap
- **No tool** bridges feature flags to messaging platforms, event analytics, or technical queues natively — teams build custom glue code every time

### 1.2 The Solution

FlagBridge offers:

1. **Feature Flag Management** — create, toggle, evaluate, rollout strategies
2. **Product Context Cards** — each flag linked to hypothesis, owner, success metrics, deadline
3. **Technical Dashboard** — adoption rate, error rate by variant, latency impact, stale flag detection
4. **Lifecycle Automation** — cleanup alerts, auto-archival, tech debt tracking
5. **OpenFeature Compatible** — official provider for the CNCF OpenFeature standard, zero vendor lock-in
6. **Plugin Ecosystem & Marketplace** — extensible architecture where developers can build, publish, and sell plugins
7. **Integrations Hub** — native connectors for messaging platforms (Resend, RD Station, SendGrid, Mailchimp), event analytics (Mixpanel, Amplitude, Segment, PostHog), and technical queues (SQS, Kafka, RabbitMQ, NATS)

### 1.3 Open-Core Model

| Layer | Distribution | Price |
|-------|-------------|-------|
| **Community Edition** | Open source (Apache 2.0) | Free |
| **Pro Edition** | Self-hosted plugin OR SaaS | $X/mo |
| **Enterprise** | Managed SaaS + support | Custom |

**The self-hosted plugin advantage:** the customer already runs FlagBridge CE, purchases the Pro license, runs `docker pull flagbridge/pro-plugin && docker compose up -d`, and restarts. No migration, no cloud switch, no downtime. The plugin injects Pro modules into the same deployment.

---

## 2. Competitive Analysis

### 2.1 Landscape

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
| Messaging Integrations (Resend, RD Station) | ✅ | ❌ | ❌ | ❌ | ❌ |
| Event Analytics (Mixpanel, Amplitude) | ✅ | ❌ | Partial | ❌ | Native |
| Technical Queues (SQS, Kafka) | ✅ Pro | ❌ | ❌ | ❌ | ❌ |
| OpenFeature Provider | ✅ | ✅ | ✅ | ✅ | ✅ |
| Entry paid price | ~$29/mo | $80/mo (5 seats) | $120/mo+ | $45/mo | $0 (bundled) |

### 2.2 Positioning

FlagBridge occupies the space between:

- **Unleash** — good open-source, weak in product/analytics, no plugin ecosystem, no integrations beyond basic webhooks
- **LaunchDarkly** — powerful but expensive and SaaS-only, no extensibility, limited messaging/queue support
- **PostHog** — good analytics but FF is secondary, no standalone FF product, no messaging integrations
- **None of them** offer a Plugin Marketplace or native bridges to messaging, analytics, and queue infrastructure

**Tagline:** *"Feature flags with product intelligence. Open source. Extensible. Your data, your rules."*

### 2.3 Exploitable Gaps in Competitors

- **Unleash:** Users complain about lack of multiple projects in open-source version and absence of analytics integration; no way to trigger notifications on flag events
- **LaunchDarkly:** Pricing model based on service connections + MAU generates unpredictable costs; median contracts of $72k/year push away startups; integrations are enterprise-gated
- **Flagsmith:** No real-time flag sync in free version; no event pipeline support
- **All of them:** None offer "Product Context Cards" natively; none offer extensibility via plugins; none bridge to messaging platforms or technical queues

---

## 3. MVP Definition — Community vs Pro

### 3.1 Community Edition (Open Source)

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
- Webhooks for integrations (generic HTTP POST on flag events)
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

**Integrations (basic)**
- Webhook events for all flag lifecycle events
- Slack notification (official CE plugin)
- Generic event emitter (JSON payload on flag changes)

### 3.2 Pro Edition (Plugin)

Pro is what transforms FlagBridge from "another FF tool" into a **product intelligence platform with enterprise-grade connectivity**.

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

**Messaging & Communication Integrations (Pro Plugin)**
- Native connectors for email/messaging platforms
- Trigger campaigns or transactional messages based on flag state changes
- See Section 8.1 for full details

**Event Analytics & Tracking Integrations (Pro Plugin)**
- Stream flag evaluation data to analytics platforms
- Correlate flag rollouts with product metrics
- See Section 8.2 for full details

**Technical Queue Integrations (Pro Plugin)**
- Publish flag events to message queues for downstream processing
- Consume queue messages to trigger flag state changes
- See Section 8.3 for full details

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

The plugin works as a **module that extends the CE**:
- Registers additional routes on the API
- Adds React components to the Next.js Admin UI
- Creates extra PostgreSQL tables (automatic migration)
- Validates license via license key (offline-first, with optional heartbeat)

### 3.3 Enterprise (SaaS Managed)

- Everything in Pro
- Managed hosting by FlagBridge (multi-tenant or dedicated)
- SLA with uptime guarantee
- Priority support (Slack, email, call)
- Custom integrations
- SOC 2 compliance (roadmap)
- Data residency options
- White-label Admin UI option
- Dedicated plugin marketplace instance
- Custom queue and messaging connector development

---

## 4. Technical Architecture

### 4.1 Architecture Overview

```
┌──────────────────────────────────────────────────────────────────────┐
│                        Client Applications                            │
│             (React, Node, Go, Python — via SDKs)                      │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌──────────────┐       ┌──────────────────────────────────┐        │
│  │  FlagBridge   │       │   FlagBridge Edge                 │        │
│  │  SDK          │──────▶│   (optional, for scale)           │        │
│  └──────────────┘       └──────────────┬─────────────────┘        │
│                                        │                            │
├────────────────────────────────────────┼────────────────────────────┤
│                                        ▼                            │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │                FlagBridge API Server (Go)                     │   │
│  │                                                              │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌───────────────┐  │   │
│  │  │ Flag     │ │ Product  │ │Dashboard │ │ Plugin        │  │   │
│  │  │ Eval     │ │ Context  │ │& Metrics │ │ Runtime       │  │   │
│  │  │ Engine   │ │ Module   │ │ Module   │ │ Engine        │  │   │
│  │  └──────────┘ └──────────┘ └──────────┘ └───────────────┘  │   │
│  │    ▲ CE          ▲ Pro        ▲ Pro        ▲ CE+Pro        │   │
│  │                                                              │   │
│  │  ┌──────────────────────────────────────────────────────┐   │   │
│  │  │              Integrations Layer (Pro)                  │   │   │
│  │  │                                                       │   │   │
│  │  │  ┌────────────┐  ┌──────────────┐  ┌──────────────┐  │   │   │
│  │  │  │ Messaging  │  │ Analytics &  │  │ Technical    │  │   │   │
│  │  │  │ Connectors │  │ Event        │  │ Queue        │  │   │   │
│  │  │  │            │  │ Connectors   │  │ Connectors   │  │   │   │
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
│  │                Admin UI (Next.js + TypeScript)                 │   │
│  │                                                               │   │
│  │  ┌────────────┐ ┌────────────┐ ┌───────────────────────────┐ │   │
│  │  │ Dashboard  │ │ Flag       │ │ Plugin Manager &          │ │   │
│  │  │ & Analytics│ │ Manager &  │ │ Marketplace               │ │   │
│  │  │            │ │ Product    │ │                           │ │   │
│  │  │            │ │ Cards      │ │                           │ │   │
│  │  └────────────┘ └────────────┘ └───────────────────────────┘ │   │
│  │                                                               │   │
│  │  ┌────────────┐ ┌────────────┐ ┌───────────────────────────┐ │   │
│  │  │ Settings & │ │ Audit Log  │ │ Integrations Hub          │ │   │
│  │  │ Team Mgmt  │ │ & Activity │ │ (Messaging, Analytics,    │ │   │
│  │  │            │ │            │ │  Queues config)           │ │   │
│  │  └────────────┘ └────────────┘ └───────────────────────────┘ │   │
│  └──────────────────────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────────────────────┘
```

### 4.2 Tech Stack

| Component | Technology | Justification |
|-----------|-----------|--------------|
| API Server | **Go** | Single binary, high performance, low memory footprint, ideal for self-hosted |
| Admin UI | **Next.js 15 + TypeScript** | SSR/SSG for landing pages, App Router for admin SPA, i18n built-in for bilingual support |
| Admin UI Components | **Next.js + Tailwind CSS + Radix UI** | Accessible, composable components with design system flexibility |
| Admin UI State | **TanStack Query + Zustand** | Server state caching + lightweight client state |
| Admin UI Charts | **Recharts** or **Tremor** | Dashboard visualizations optimized for Next.js |
| Landing Page | **Next.js (SSG)** | Same codebase as admin, bilingual (en/pt) with `next-intl`, SEO-optimized |
| Database | **PostgreSQL** | Reliable, JSONB for targeting rules, extensible |
| Cache | **In-memory (Go)** + Redis (optional) | Flag evaluation in < 1ms; Redis for clusters |
| Edge Proxy | **Go** (FlagBridge Edge) | For high volume — caches flags close to client |
| SDK Transport | **SSE** (Server-Sent Events) | Real-time flag updates without WebSocket complexity |
| Plugin Runtime | **Go (backend)** + **Next.js (UI)** | Sandboxed plugin execution with defined hook points |
| Integration Layer | **Go** with adapter pattern | Unified interface for messaging, analytics, and queue connectors |
| Containerization | **Docker** | Multi-stage Dockerfile |
| Orchestration | **Docker Compose** (default) / Helm (K8s) | 90% of self-hosted users use Compose |

### 4.2.1 Next.js Architecture Detail

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
│   │   │   │   ├── integrations/        # Integrations Hub
│   │   │   │   │   ├── messaging/       # Resend, RD Station, etc.
│   │   │   │   │   ├── analytics/       # Mixpanel, Amplitude, etc.
│   │   │   │   │   ├── queues/          # SQS, Kafka, etc.
│   │   │   │   │   └── webhooks/        # Generic webhooks
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
│   │   │   └── (developer)/            # Developer portal
│   │   │       ├── docs/
│   │   │       ├── api-explorer/
│   │   │       └── sandbox/
│   │   └── api/                        # Next.js API routes (BFF)
│   │       ├── auth/
│   │       └── proxy/
│   ├── components/
│   │   ├── ui/                         # Design system (Radix + Tailwind)
│   │   ├── flags/
│   │   ├── dashboard/
│   │   ├── plugins/
│   │   ├── marketplace/
│   │   └── integrations/               # Integration config components
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

**Key architectural decisions:**

1. **Route Groups**: `(marketing)` for SSG landing pages, `(admin)` for authenticated app, `(developer)` for developer portal — all in the same Next.js app
2. **i18n via `next-intl`**: All routes prefixed with locale (`/en/dashboard`, `/pt/dashboard`), SSG pages fully bilingual for SEO
3. **BFF Pattern**: Next.js API routes act as Backend-for-Frontend, proxying to Go API with auth token injection
4. **Plugin UI Host**: Plugins render into designated `<PluginSlot />` components using a sandboxed iframe or Module Federation approach
5. **Integrations Hub**: Dedicated admin section for managing all external service connections with unified config UI

### 4.3 Data Model (Core)

```sql
-- ========================
-- CORE TABLES (CE)
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
-- PLUGIN TABLES (CE)
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
-- PRO TABLES
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
-- INTEGRATIONS TABLES (Pro)
-- ========================

CREATE TABLE integrations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id      UUID REFERENCES projects(id),
    type            VARCHAR(50) NOT NULL, -- messaging, analytics, queue
    provider        VARCHAR(50) NOT NULL, -- resend, mixpanel, sqs, etc.
    name            VARCHAR(255) NOT NULL,
    config          JSONB NOT NULL, -- Provider-specific config (encrypted secrets)
    enabled         BOOLEAN DEFAULT TRUE,
    health_status   VARCHAR(20) DEFAULT 'unknown', -- healthy, degraded, down, unknown
    last_health_check TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE integration_rules (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    integration_id  UUID REFERENCES integrations(id) ON DELETE CASCADE,
    trigger_event   VARCHAR(100) NOT NULL, -- flag.toggled, flag.rolledOut, flag.staleDetected
    conditions      JSONB DEFAULT '{}', -- Optional conditions (specific flags, envs, etc.)
    action_config   JSONB NOT NULL, -- Provider-specific action (send email, publish event, etc.)
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
-- MARKETPLACE TABLES (Pro/Enterprise)
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

### 4.4 SDK Design

SDKs follow the OpenFeature standard, with a proprietary layer for FlagBridge-specific features.

```typescript
// @flagbridge/sdk-node — Usage example

import { FlagBridge } from '@flagbridge/sdk-node';

const fb = new FlagBridge({
  serverUrl: 'https://flags.mycompany.com',
  apiKey: 'fb_sk_...',
  environment: 'production',
});

// Simple boolean evaluation
const isEnabled = await fb.isEnabled('checkout-v2', {
  userId: 'user_123',
  attributes: { plan: 'pro', country: 'BR' },
});

// String variant evaluation
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
// @flagbridge/sdk-go — Usage example

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

### 4.5 API Endpoints

```
# ========================
# FLAG MANAGEMENT (CE)
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
# INTEGRATIONS (Pro)
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

GET    /api/v1/integrations/providers          # List available providers
GET    /api/v1/integrations/providers/:provider # Provider config schema

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

### 4.6 Deployment Infrastructure

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
      - NEXT_PUBLIC_DEFAULT_LOCALE=en
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

## 5. Admin Panel — Technical & Product Documentation

The FlagBridge Admin is a full-featured Next.js application that serves as the central hub for flag management, product intelligence, plugin management, integrations configuration, and marketplace access.

### 5.1 Admin Sections

#### Dashboard
- **Overview cards**: Total flags, active flags, stale flags, flags by status
- **Activity feed**: Recent flag changes, deployments, plugin installs, integration events
- **Health metrics** (Pro): Evaluation volume, error rates, SDK connection count, integration health
- **Technical debt score** (Pro): Aggregate score based on stale flags, flags without owners, flags without product cards

#### Flag Manager
- **Flag list**: Filterable by project, environment, status, tags, owner
- **Flag detail page**: Toggle, targeting rules editor (visual builder), environment comparison
- **Product Context Card** (tab within flag detail): Hypothesis, success metrics, go/no-go criteria, decision history, OKR link
- **Metrics tab** (Pro): Adoption chart, error rate chart, latency impact, evaluation breakdown by variant
- **Lifecycle tab** (Pro): Timeline from creation to archival, cleanup reminders, code references
- **Integrations tab** (Pro): Which integrations are connected to this flag, event history

#### Integrations Hub
- **Messaging**: Configure Resend, RD Station, SendGrid, Mailchimp, Brevo, Customer.io connections
- **Analytics & Events**: Configure Mixpanel, Amplitude, Segment, PostHog, GA4, Rudderstack connections
- **Technical Queues**: Configure SQS, Kafka, RabbitMQ, NATS, Redis Streams, GCP Pub/Sub connections
- **Webhooks**: Generic webhook management
- **Per-integration**: Connection config, trigger rules, event log, health status, test button

#### Plugin Manager
- **Installed plugins**: List with enable/disable toggle, config panel, health status
- **Marketplace browser**: Search, filter by category, install with one click
- **Plugin development**: Sandbox environment, logs, hot-reload for local dev

#### Settings
- **Team management**: Invite members, assign roles (Admin, Editor, Viewer)
- **API keys**: Create, rotate, revoke keys per environment
- **Billing** (Pro/SaaS): Plan management, invoices, usage metrics

#### Audit Log
- **Comprehensive history**: Every flag change, toggle, user action, plugin install, integration event
- **Filterable**: By user, action type, date range, entity
- **Diff viewer** (Pro): Side-by-side comparison of flag state changes

#### Developer Portal
- **Plugin SDK documentation**: Full API reference, guides, tutorials
- **Interactive API explorer**: Try endpoints with live data (Swagger/OpenAPI)
- **Plugin sandbox**: Test your plugin in an isolated environment
- **Marketplace publisher**: Submit, manage, and track your plugins
- **Integration SDK docs**: How to build custom integration connectors

### 5.2 Bilingual Support

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
  },
  "integrations": {
    "title": "Integrations",
    "messaging": "Messaging & Email",
    "analytics": "Analytics & Events",
    "queues": "Technical Queues",
    "webhooks": "Webhooks",
    "health": {
      "healthy": "Healthy",
      "degraded": "Degraded",
      "down": "Down"
    },
    "test": "Test Connection",
    "eventLog": "Event Log"
  }
}
```

---

## 6. Plugin System & SDK

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
| `integration.onEvent` | When integration event fires | Cross-integration workflows |

#### Frontend (UI) Hook Points

| Hook Point | Location | Use Cases |
|-----------|----------|-----------|
| `dashboard.widget` | Dashboard page | Custom metric cards, charts |
| `flag.detail.sidebar` | Flag detail sidebar | Related info, external data |
| `flag.detail.tab` | Flag detail tabs | Custom tabs (e.g. analytics) |
| `flag.list.column` | Flag list table | Custom columns |
| `settings.integrations.panel` | Settings > Integrations | Plugin configuration UI |
| `integrations.provider.config` | Integration config page | Custom provider config UI |
| `navigation.sidebar` | Main sidebar nav | New navigation items |
| `global.banner` | Top of page | Alerts, announcements |

#### Plugin SDK

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
      <h3 className="font-semibold text-sm text-gray-500">Slack Notifications</h3>
      <p className="text-sm mt-1">Posting to <code>{config.channel}</code></p>
      <p className="text-xs text-gray-400 mt-2">Last notified: {flag.lastToggleAt}</p>
    </div>
  );
};
```

#### Plugin CLI

```bash
npx @flagbridge/create-plugin my-awesome-plugin  # Scaffold
cd my-awesome-plugin
flagbridge plugin dev       # Local dev with hot-reload
flagbridge plugin test      # Run tests
flagbridge plugin build     # Build for distribution
flagbridge plugin publish   # Publish to marketplace
```

### 6.2 Plugin Security Model

- **Sandboxed execution**: Backend plugins run in isolated goroutines with limited syscall access
- **Permission system**: Plugins declare required permissions; admin must approve
- **UI sandboxing**: Frontend plugins render in sandboxed iframes or use strict CSP
- **Config encryption**: Secret config values (API keys, tokens) encrypted at rest
- **Code signing**: Marketplace plugins must be signed; hash verified on install
- **Review process**: All marketplace submissions reviewed before publishing

---

## 7. Plugin Marketplace

### 7.1 Vision

The FlagBridge Plugin Marketplace transforms FlagBridge from a product into a **platform ecosystem** — similar to Shopify App Store, WordPress Plugins, or Figma Community. Third-party developers can build, publish, and monetize plugins that extend FlagBridge for specific use cases.

### 7.2 Marketplace Features

#### For Plugin Users (Buyers)
- Browse & search by category, rating, price, compatibility
- One-click install from the Admin UI (self-hosted and SaaS)
- Reviews & ratings (5-star with written reviews)
- Compatibility check against current FlagBridge version
- Plugin configuration through Admin UI
- Optional auto-updates (self-hosted respects user control)

#### For Plugin Developers (Sellers)
- Developer Portal with full documentation, API explorer, Plugin SDK
- Plugin Sandbox for development and debugging
- Publish workflow: Submit → automated checks → manual review → published
- Pricing options: free, one-time purchase, or monthly subscription
- Revenue split: 80% developer / 20% FlagBridge
- Earnings dashboard: installs, revenue, ratings, support tickets
- Verified Developer badge for trusted publishers

#### For FlagBridge (Platform)
- Review pipeline: automated security scan + manual code review
- Revenue stream: 20% commission on paid plugins
- Ecosystem flywheel: more plugins → more value → more users → more developers
- Quality control: minimum quality standards, security requirements, regular audits

### 7.3 Marketplace Categories

| Category | Examples |
|----------|---------|
| **Integrations** | Slack, Discord, Linear, Jira, PagerDuty, Datadog |
| **Messaging** | Resend, RD Station, SendGrid, Mailchimp, Customer.io, Brevo |
| **Analytics** | Mixpanel, Amplitude, Segment, PostHog, GA4, Rudderstack |
| **Queues** | SQS, Kafka, RabbitMQ, NATS, Redis Streams, GCP Pub/Sub |
| **Security** | Advanced RBAC, IP whitelisting, compliance reports |
| **UI Extensions** | Custom widgets, theme packs, dashboard layouts |
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

## 8. Integrations Ecosystem

This section details the three categories of external service integrations that differentiate FlagBridge from competitors: **Messaging & Communication**, **Analytics & Event Tracking**, and **Technical Queues**.

### 8.1 Messaging & Communication Integrations

Connect flag lifecycle events to communication platforms for automated notifications, campaigns, and stakeholder updates.

#### Supported Providers

| Provider | Type | Use Cases |
|----------|------|-----------|
| **Resend** | Transactional email | Notify stakeholders when flags roll out; send digests of stale flags |
| **RD Station** | Marketing automation (Brazil-focused) | Trigger nurturing flows based on feature access; segment users by active flags |
| **SendGrid** | Transactional + marketing email | Automated rollout reports; user communication on feature launches |
| **Mailchimp** | Marketing campaigns | Audience segmentation by feature flags; launch announcements |
| **Brevo (Sendinblue)** | Multi-channel (email, SMS, WhatsApp) | Multi-channel notifications on flag events; SMS alerts for critical rollbacks |
| **Customer.io** | Behavioral messaging | Trigger messages based on flag evaluation patterns; onboarding flows by feature access |

#### How It Works

```
┌─────────────┐     ┌──────────────────┐     ┌─────────────────┐
│ Flag Event   │────▶│ Integration Rule  │────▶│ Messaging       │
│              │     │                  │     │ Provider        │
│ flag.toggled │     │ IF flag.key      │     │                 │
│ flag.rolledOut│    │   matches "x"    │     │ Resend.send()   │
│ flag.stale   │     │ AND env = prod   │     │ RDStation.push()│
│              │     │ THEN action...   │     │ SendGrid.send() │
└─────────────┘     └──────────────────┘     └─────────────────┘
```

#### Configuration Example (Admin UI)

```json
{
  "provider": "resend",
  "config": {
    "apiKey": "re_...",
    "fromEmail": "flags@mycompany.com",
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
        "cc": ["product-team@mycompany.com"],
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

#### RD Station Deep Integration

Given FlagBridge's Brazilian roots, RD Station gets a first-class integration:

- **Contact segmentation**: Automatically tag RD Station contacts based on which feature flags are active for them
- **Conversion events**: Push feature flag evaluations as conversion events to RD Station for funnel analysis
- **Lead scoring**: Adjust lead scores based on feature adoption patterns
- **Nurturing flows**: Trigger specific email flows when a user gains access to a new feature via flag rollout

### 8.2 Analytics & Event Tracking Integrations

Stream flag evaluation data to analytics platforms, enabling teams to correlate feature rollouts with product metrics, user behavior, and business outcomes.

#### Supported Providers

| Provider | Type | Use Cases |
|----------|------|-----------|
| **Mixpanel** | Product analytics | Track flag evaluations as events; build funnels per flag variant; measure feature adoption |
| **Amplitude** | Product analytics | Correlate feature rollouts with user behavior; experiment analysis dashboards |
| **Segment** | Customer data platform | Route flag events to any Segment destination; unified user profiles with flag data |
| **PostHog** | Open-source analytics | Feature flag correlation with session recordings; funnel analysis per variant |
| **Google Analytics 4** | Web analytics | Custom events for flag evaluations; conversion tracking per flag variant |
| **Rudderstack** | Open-source CDP | Self-hosted event pipeline; route flag data to warehouses and analytics tools |

#### How It Works

```
┌─────────────────┐     ┌───────────────────┐     ┌──────────────────┐
│ Flag Evaluation  │────▶│ Analytics Adapter  │────▶│ Analytics        │
│                  │     │                   │     │ Provider         │
│ flag: checkout-v2│     │ Transform to      │     │                  │
│ variant: "new"  │     │ provider format   │     │ mixpanel.track() │
│ userId: "u_123" │     │                   │     │ amplitude.log()  │
│ timestamp: ...  │     │ Apply sampling    │     │ segment.track()  │
│                  │     │ & batching        │     │ posthog.capture()│
└─────────────────┘     └───────────────────┘     └──────────────────┘
```

#### Event Schema

Every flag evaluation generates a standardized event that is adapted to each provider:

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
    "flag_owner": "gabriel@company.com",
    "product_hypothesis": "New checkout flow increases conversion by 15%"
  }
}
```

**Key features:**

- **Sampling**: Configurable sampling rate (e.g., send 10% of evaluations) to control costs on high-volume flags
- **Batching**: Events are batched (configurable: every N events or every N seconds) before sending to reduce API calls
- **Enrichment**: Events are automatically enriched with product context (hypothesis, owner, status) from Product Cards
- **Selective tracking**: Choose which flags and environments to track — no need to send everything

#### Configuration Example

```json
{
  "provider": "mixpanel",
  "config": {
    "projectToken": "mp_...",
    "apiSecret": "...",
    "region": "US"
  },
  "rules": [
    {
      "trigger": "flag.evaluated",
      "conditions": {
        "environment": "production",
        "flagTags": ["experiment"]
      },
      "action": {
        "eventName": "Feature Flag Evaluated",
        "sampling": 0.1,
        "batchSize": 100,
        "batchIntervalMs": 5000,
        "includeProductContext": true,
        "userIdField": "userId"
      }
    },
    {
      "trigger": "flag.rolledOut",
      "action": {
        "eventName": "Feature Rollout Complete",
        "sampling": 1.0
      }
    }
  ]
}
```

### 8.3 Technical Queue Integrations (Pro Plugin)

Publish flag events to message queues for downstream processing by microservices, data pipelines, and custom automation. This is delivered as a **Pro plugin** because it targets engineering teams with complex infrastructure needs.

#### Supported Providers

| Provider | Type | Use Cases |
|----------|------|-----------|
| **Amazon SQS** | Managed queue (AWS) | Decouple flag events from downstream consumers; trigger Lambda functions on flag changes |
| **Apache Kafka** | Distributed streaming | Stream flag evaluations for real-time processing; event sourcing for flag state changes |
| **RabbitMQ** | Message broker | Route flag events to specific consumers via exchanges; reliable delivery with acknowledgments |
| **NATS** | Lightweight messaging | Low-latency flag event distribution; pub/sub for microservices |
| **Redis Streams** | In-memory streaming | High-throughput flag evaluation streaming; consumer groups for parallel processing |
| **GCP Pub/Sub** | Managed messaging (GCP) | Event-driven architectures on Google Cloud; trigger Cloud Functions on flag events |

#### How It Works

```
┌─────────────────┐     ┌──────────────────┐     ┌──────────────────┐
│ Flag Event       │────▶│ Queue Adapter     │────▶│ Queue Provider   │
│                  │     │                  │     │                  │
│ flag.toggled     │     │ Serialize event  │     │ SQS.sendMessage()│
│ flag.evaluated   │     │ Route to queue   │     │ kafka.produce()  │
│ flag.stale       │     │ Retry on failure │     │ rabbit.publish() │
│ flag.rolledOut   │     │                  │     │ nats.publish()   │
└─────────────────┘     └──────────────────┘     └──────────────────┘

                         ┌──────────────────┐
                         │ Downstream       │
                         │ Consumers        │
                         │                  │
                         │ Lambda functions │
                         │ Data pipelines   │
                         │ Microservices    │
                         │ Alert systems    │
                         └──────────────────┘
```

#### Use Cases

**Event-driven flag reactions**: Microservices subscribe to flag change events and automatically adjust behavior (cache invalidation, config reload, service mesh updates)

**Data pipeline integration**: Stream flag evaluations to data warehouses via Kafka → Spark/Flink → BigQuery/Snowflake for offline analysis

**Audit & compliance**: Immutable event log of all flag changes in a durable queue, satisfying compliance requirements

**Cross-service coordination**: When a flag is toggled in FlagBridge, publish to SQS/Kafka so that all dependent services react without polling the FlagBridge API

**Custom alerting**: Consume flag events from a queue and feed them into custom alerting pipelines (PagerDuty, OpsGenie, custom Slack bots)

#### Configuration Example (SQS)

```json
{
  "provider": "sqs",
  "config": {
    "region": "us-east-1",
    "accessKeyId": "AKIA...",
    "secretAccessKey": "...",
    "queueUrl": "https://sqs.us-east-1.amazonaws.com/123456789/flagbridge-events"
  },
  "rules": [
    {
      "trigger": "flag.toggled",
      "conditions": {
        "environment": "production"
      },
      "action": {
        "messageAttributes": {
          "eventType": "flag.toggled",
          "source": "flagbridge"
        },
        "delaySeconds": 0
      }
    },
    {
      "trigger": "flag.evaluated",
      "conditions": {
        "flagTags": ["critical-path"]
      },
      "action": {
        "sampling": 0.01,
        "batchSize": 500,
        "batchIntervalMs": 10000
      }
    }
  ]
}
```

#### Configuration Example (Kafka)

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

### 8.4 Integration Architecture (Go)

All integrations share a unified adapter interface in Go:

```go
// pkg/integrations/adapter.go

package integrations

type EventType string

const (
    FlagToggled      EventType = "flag.toggled"
    FlagEvaluated    EventType = "flag.evaluated"
    FlagRolledOut    EventType = "flag.rolledOut"
    FlagStaleDetected EventType = "flag.staleDetected"
    FlagUpdated      EventType = "flag.updated"
    FlagArchived     EventType = "flag.archived"
)

// IntegrationEvent is the standardized event structure
type IntegrationEvent struct {
    Type        EventType              `json:"type"`
    Timestamp   time.Time              `json:"timestamp"`
    Flag        FlagSnapshot           `json:"flag"`
    Environment EnvironmentSnapshot    `json:"environment"`
    Actor       *ActorSnapshot         `json:"actor,omitempty"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Adapter is the interface all integration providers must implement
type Adapter interface {
    // Name returns the provider identifier (e.g., "resend", "mixpanel", "sqs")
    Name() string

    // Type returns the integration category
    Type() IntegrationType // messaging, analytics, queue

    // ConfigSchema returns the JSON schema for provider configuration
    ConfigSchema() json.RawMessage

    // Connect establishes a connection with the provider
    Connect(config json.RawMessage) error

    // Send delivers an event to the provider
    Send(ctx context.Context, event IntegrationEvent, actionConfig json.RawMessage) error

    // SendBatch delivers a batch of events (for providers that support batching)
    SendBatch(ctx context.Context, events []IntegrationEvent, actionConfig json.RawMessage) error

    // HealthCheck verifies the connection is alive
    HealthCheck(ctx context.Context) HealthStatus

    // Close cleanly shuts down the connection
    Close() error
}

// Registry manages available integration adapters
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

This adapter pattern means:

1. **Adding a new provider is simple** — implement the `Adapter` interface
2. **Community can contribute providers** — via the plugin system
3. **All providers share the same event format** — standardized `IntegrationEvent`
4. **Health checks are uniform** — every provider reports health the same way
5. **Batching and sampling are handled at the integration layer** — providers don't need to implement this

---

## 9. Go-to-Market

### 9.1 Target Audience

**Primary:** Startups and scale-ups (50-500 devs) that:
- Already use feature flags (Unleash CE, in-house) but lack observability and product context
- Cannot afford LaunchDarkly (~$72k/year median)
- Value self-hosted and data sovereignty
- Want extensibility without vendor lock-in
- Already use tools like Mixpanel, Segment, or Kafka and want their FF tool to integrate natively

**Secondary:** Product teams at larger companies that:
- Want visibility into flag states without depending on engineering
- Need to link flags to OKRs and roadmap
- Want a plugin ecosystem for custom integrations
- Need to trigger communications (Resend, RD Station) based on feature rollouts

### 9.2 Pricing

| Plan | Price | Target |
|------|-------|--------|
| **Community** | Free (forever) | Individual devs, small teams, evaluation |
| **Pro** (self-hosted) | $29/mo (≤10 seats) / $79/mo (≤50 seats) / $149/mo (unlimited) | Scale-ups, mid-size teams |
| **Pro** (SaaS) | $49/mo (≤10 seats) / $99/mo (≤50 seats) / $199/mo (unlimited) | Teams that don't want to manage infra |
| **Enterprise** | Custom (from $500/mo) | Large companies, compliance, SLA |

**Flat-rate per seats** — no MAU or service connection surprises. Direct differentiator against LaunchDarkly.

Self-hosted is cheaper than SaaS because the customer pays for infra. This incentivizes self-hosted (less hosting cost for us) and attracts the data sovereignty audience.

### 9.3 Acquisition Channels

**Developer-first growth:**
1. GitHub — well-documented open-source repo, attractive README, contributing guide
2. Hacker News / Reddit / Dev.to — launch posts
3. Honest comparisons — blog posts: "FlagBridge vs Unleash vs LaunchDarkly"
4. Docker Hub — optimized official image, one-liner setup
5. OpenFeature ecosystem — official provider listed on OpenFeature site
6. YouTube / blog — tutorials in English and Portuguese

**Product-led growth:**
- CE is genuinely useful → team adopts → grows → needs observability/governance → upgrades to Pro
- Discrete banner in Admin UI: "Unlock product dashboards and lifecycle automation"
- 14-day Pro trial (self-hosted and SaaS)

**Ecosystem growth (post-marketplace):**
- Plugin developers attract their own users to FlagBridge
- Marketplace becomes a discovery channel
- Revenue share incentivizes quality plugins

**Integration-driven growth:**
- Teams searching for "Mixpanel feature flags" or "Kafka feature flag events" find FlagBridge
- RD Station integration positions FlagBridge uniquely in the Brazilian market
- Integration partners may co-market FlagBridge on their directories

### 9.4 Success Metrics (Year 1)

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
| Active integrations (Pro) | — | 150 |

---

## 10. Roadmap

### Q2 2026 — Foundation
- [ ] Public repo on GitHub
- [ ] Core Go API server: flag management + eval engine
- [ ] Next.js Admin UI: flag CRUD, dashboard, settings (bilingual from day 1)
- [ ] Node.js SDK + Go SDK
- [ ] Docker image + Docker Compose
- [ ] Docs site (docs.flagbridge.io) — bilingual
- [ ] OpenFeature Provider (Node.js)
- [ ] Plugin runtime (basic)
- [ ] Webhook integration (CE)

### Q3 2026 — Community Launch
- [ ] Product Context Cards (CE basic version)
- [ ] CLI (`flagbridge`)
- [ ] React SDK (client-side) + Python SDK
- [ ] Helm chart for Kubernetes
- [ ] Plugin CLI + Plugin SDK v1
- [ ] 3-5 official plugins (Slack, Linear, GitHub)
- [ ] Slack notification integration (CE)
- [ ] Launch: Hacker News + Product Hunt + Reddit

### Q4 2026 — Pro Launch
- [ ] Pro Plugin architecture
- [ ] Technical Dashboard (metrics, observability)
- [ ] Advanced Product Cards
- [ ] Lifecycle automation
- [ ] SSO (OIDC)
- [ ] Billing + license key management
- [ ] **Messaging integrations**: Resend, SendGrid, RD Station
- [ ] **Analytics integrations**: Mixpanel, Segment
- [ ] Plugin Marketplace v1 (free plugins only)
- [ ] Integrations Hub in Admin UI

### Q1 2027 — Scale & Marketplace
- [ ] FlagBridge Edge (proxy for high volume)
- [ ] Basic A/B experimentation
- [ ] Granular RBAC + change requests
- [ ] **Queue integrations**: SQS, Kafka, RabbitMQ
- [ ] **Analytics integrations**: Amplitude, PostHog, GA4, Rudderstack
- [ ] **Messaging integrations**: Mailchimp, Brevo, Customer.io
- [ ] SaaS managed offering
- [ ] Plugin Marketplace v2: paid plugins, Stripe Connect
- [ ] Developer Portal with sandbox
- [ ] Bilingual landing page redesign

### Q2 2027 — Ecosystem
- [ ] Plugin Marketplace v3: subscriptions, verified developers
- [ ] **Queue integrations**: NATS, Redis Streams, GCP Pub/Sub
- [ ] Integration SDK (for community-built connectors)
- [ ] White-label option (Enterprise)
- [ ] Advanced experimentation
- [ ] SOC 2 compliance (start)
- [ ] Community plugin hackathon

---

## 11. Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|-----------|
| Saturated FF tools market | High | Clear differentiation via Product Context + Plugin Marketplace + Integrations Ecosystem — no one else does all three |
| Unleash adds similar features | Medium | Speed of execution + product-first focus; integration ecosystem as moat |
| Difficulty monetizing open-source | High | CE is genuinely useful but Pro solves real pains; Marketplace creates additional revenue; integrations drive Pro adoption |
| Complexity of maintaining multiple SDKs | Medium | Start with 4 SDKs; OpenFeature reduces surface area |
| Limited time (side project) | High | Lean MVP; prioritize core → launch → iterate with feedback |
| Plugin security vulnerabilities | Medium | Sandboxed execution, permission system, mandatory review, code signing |
| Marketplace chicken-and-egg | High | Build 5-10 official plugins first; early developer incentives |
| Integration maintenance burden | Medium | Adapter pattern isolates providers; community can contribute via plugin system; start with 6 providers and expand |
| Third-party API changes | Low | Adapter abstraction means changes are localized; version pinning on provider SDKs |

---

## Next Steps

1. **Register domains** — flagbridge.io and flagbridge.dev
2. **Create GitHub org** — github.com/flagbridge
3. **Reserve npm scope** — @flagbridge
4. **Init monorepo** — `flagbridge/flagbridge` (api + ui + docs + plugin-sdk)
5. **Prototype** — Go API server with flag CRUD + eval endpoint
6. **Design** — Next.js Admin UI wireframes (dashboard + flag detail + product card + integrations hub)
7. **Landing page** — Bilingual Next.js SSG site with positioning, features, pricing
8. **Plugin SDK scaffold** — Define hook points, manifest format, CLI
9. **Integration adapter scaffold** — Implement `Adapter` interface + first provider (Resend or Slack)

---

*FlagBridge — Bridging the gap between feature flags and product strategy.*
