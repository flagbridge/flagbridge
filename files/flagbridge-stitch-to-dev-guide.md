# FlagBridge — Guia: Como extrair designs do Stitch para o desenvolvimento

> Este guia explica todas as formas de tirar seus layouts do Google Stitch
> e trazê-los para o codebase do FlagBridge para uso no desenvolvimento do Admin UI.

---

## Resumo Rápido: 4 caminhos de export do Stitch

| Método | O que você recebe | Melhor para |
|--------|-------------------|-------------|
| **1. Export HTML/CSS/Tailwind** | Código limpo, copiar e colar | Referência de estilo, tokens de design |
| **2. Export React/JSX** | Componentes React prontos | Base direta para componentes Next.js |
| **3. Export para Figma** | Layers editáveis com Auto Layout | Refinamento de design, handoff detalhado |
| **4. DESIGN.md + MCP Server** | Design system como markdown + conexão direta com Claude Code | Workflow automatizado design → código |

**Recomendação para o FlagBridge:** Use o caminho **4 (DESIGN.md + MCP)** como principal e o **2 (React/JSX)** como complemento.

---

## Método 1: Export HTML/CSS/Tailwind (Copy & Paste)

### Como fazer

1. No Stitch, com a tela desejada aberta, clique no ícone de código (`</>`) no canto da tela
2. O Stitch mostra o HTML + CSS (ou Tailwind) gerado para aquela tela
3. Copie o código

### Como usar no FlagBridge

O código HTML/CSS do Stitch serve como **referência visual** — não pra copiar direto no Next.js, mas pra extrair:

- **Tokens de design**: cores, espaçamentos, fontes, border-radius
- **Estrutura de layout**: como os componentes estão organizados
- **Classes Tailwind**: se você exportou com Tailwind, muitas classes podem ser reutilizadas direto

```bash
# Salve o HTML de cada tela num diretório de referência
mkdir -p ~/flagbridge/design-reference/stitch-exports
# Salve como: dashboard.html, flag-list.html, flag-detail.html, etc.
```

Depois, no Claude Code:
```
"Use os arquivos HTML em design-reference/stitch-exports/ como referência visual
para implementar os componentes Next.js do Admin UI. Extraia os tokens de design
(cores, fontes, espaçamentos) e converta a estrutura para componentes React com Tailwind."
```

---

## Método 2: Export React/JSX (Mais direto)

### Como fazer

1. No Stitch, com a tela aberta, procure a opção de export como **React/JSX**
2. O Stitch gera componentes React com JSX + Tailwind classes
3. Copie ou baixe o código

### Como usar no FlagBridge

Este é o caminho mais prático pra transferir direto:

```bash
# Salve cada componente exportado
mkdir -p ~/flagbridge/design-reference/react-exports

# Exemplos de arquivos:
# dashboard-overview.jsx
# flag-list-page.jsx
# flag-detail-page.jsx
# integrations-hub.jsx
# plugin-marketplace.jsx
# etc.
```

**Importante:** O React exportado pelo Stitch é estático (sem state, sem API calls, sem i18n). Ele serve como **esqueleto visual** que você vai evoluir:

1. Converter `.jsx` para `.tsx` (TypeScript)
2. Adicionar hooks de state (`useState`, `useQuery`)
3. Conectar com a API do FlagBridge
4. Adicionar i18n com `next-intl`
5. Dividir em componentes menores seguindo a estrutura do Next.js App Router

No Claude Code:
```
"Tenho componentes React exportados do Stitch em design-reference/react-exports/.
Converta cada um para TypeScript, divida em componentes menores seguindo a estrutura
do Next.js App Router em apps/web/src/app/[locale]/(admin)/, adicione i18n
com next-intl, e conecte com a API usando TanStack Query."
```

---

## Método 3: Export para Figma

### Como fazer

1. No Stitch, use a opção **"Export to Figma"**
2. O Stitch gera um arquivo Figma com:
   - Layers editáveis (não é imagem flat)
   - Auto Layout preservado
   - Componentes organizados
3. Abra no Figma e refine

### Como usar no FlagBridge

Figma é útil se você quiser:
- Refinar espaçamentos e alinhamentos com precisão pixel-perfect
- Criar um design system formal com componentes reutilizáveis
- Gerar specs de handoff para outros devs
- Fazer user testing com protótipos clicáveis

Depois de refinar no Figma, você tem duas opções pra trazer pro código:
- **Manual**: Use o Figma como referência visual enquanto implementa no Next.js
- **Figma → código via plugin**: Use plugins como Locofy ou Anima pra gerar React

Para o FlagBridge, Figma é o caminho se você quiser um polish extra antes de codar. Mas adiciona uma etapa a mais no workflow.

---

## Método 4: DESIGN.md + Stitch MCP Server (Recomendado)

### O que é

O update de Março 2026 do Stitch introduziu:

- **DESIGN.md**: Um arquivo markdown que descreve seu design system (cores, fontes, espaçamentos, componentes) de forma que qualquer ferramenta de coding AI entende
- **Stitch MCP Server**: Um servidor Model Context Protocol que conecta o Stitch diretamente a ferramentas como Claude Code, Cursor, e Gemini CLI

### Como fazer

**Passo 1 — Extrair DESIGN.md do Stitch**

1. No Stitch canvas, após finalizar todos os designs do FlagBridge
2. Procure a opção de exportar o **DESIGN.md** (nas configurações de Design System ou no menu de export)
3. O Stitch gera um markdown com:
   - Paleta de cores (hex, nomes semânticos)
   - Tipografia (fontes, tamanhos, pesos)
   - Espaçamentos (spacing scale)
   - Border radius, shadows
   - Componentes base (botões, cards, inputs, badges, etc.)
   - Layout patterns (grid, sidebar, etc.)

**Passo 2 — Colocar o DESIGN.md no monorepo**

```bash
# Copie o DESIGN.md para a raiz do monorepo
cp ~/Downloads/DESIGN.md ~/flagbridge/flagbridge/DESIGN.md

# Ou coloque dentro do app web
cp ~/Downloads/DESIGN.md ~/flagbridge/flagbridge/apps/web/DESIGN.md
```

**Passo 3 — Usar no Claude Code**

O Claude Code lê o DESIGN.md automaticamente (assim como lê o CLAUDE.md). Basta rodar:

```
"Implemente o Dashboard do FlagBridge Admin seguindo as especificações
do DESIGN.md. Use Next.js App Router, Tailwind CSS com os tokens de design
definidos, e Radix UI para componentes base."
```

O Claude Code vai:
1. Ler o DESIGN.md e entender cores, fontes, spacing
2. Aplicar consistentemente em todos os componentes que criar
3. Manter o design system coerente entre telas

**Passo 4 (Opcional) — Conectar via Stitch MCP Server**

Se quiser um workflow ainda mais integrado:

```bash
# No Claude Code, adicione o MCP server do Stitch
# (verifique a URL exata na documentação do Stitch)
claude --mcp-server "stitch://your-project-id"
```

Com o MCP conectado, o Claude Code pode:
- Ler designs direto do Stitch sem export manual
- Acessar componentes individuais por nome
- Buscar o design system atualizado em tempo real

---

## Workflow Completo Recomendado

```
┌─────────────────────────────────────────────────────────┐
│                    GOOGLE STITCH                         │
│                                                         │
│  1. Rodar os 28 prompts do guia de Stitch               │
│  2. Refinar cada tela com prompts incrementais          │
│  3. Revisar visualmente todas as telas                   │
│                                                         │
│  EXPORTS:                                               │
│  ├── DESIGN.md (design system completo)                 │
│  ├── React/JSX de cada tela (esqueleto visual)          │
│  └── Screenshots de referência (backup)                 │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────┐
│              MONOREPO FLAGBRIDGE                         │
│                                                         │
│  flagbridge/                                            │
│  ├── DESIGN.md          ← do Stitch                    │
│  ├── CLAUDE.md          ← instruções pro Claude Code   │
│  ├── design-reference/                                  │
│  │   ├── react-exports/ ← JSX do Stitch                │
│  │   └── screenshots/   ← PNGs de referência           │
│  └── apps/web/          ← Next.js Admin UI             │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────┐
│                   CLAUDE CODE                            │
│                                                         │
│  Lê automaticamente:                                    │
│  • CLAUDE.md (regras do projeto)                        │
│  • DESIGN.md (design system)                            │
│  • design-reference/ (componentes de referência)        │
│                                                         │
│  Prompt:                                                │
│  "Implemente o Admin UI do FlagBridge usando os         │
│   React exports do Stitch como referência visual,       │
│   o DESIGN.md como design system, e converta tudo       │
│   para Next.js + TypeScript + next-intl + TanStack      │
│   Query. Conecte com a Go API."                         │
│                                                         │
│  Output:                                                │
│  → Componentes Next.js production-ready                 │
│  → Bilíngue (EN/pt)                                 │
│  → Conectados com API                                   │
│  → Design system consistente                            │
└─────────────────────────────────────────────────────────┘
```

---

## Checklist: O que exportar de cada tela do Stitch

Para cada uma das telas geradas pelos 28 prompts:

- [ ] **Screenshot PNG** (para referência visual rápida)
- [ ] **React/JSX export** (esqueleto de componente)
- [ ] **HTML/Tailwind export** (tokens e classes de referência)

E uma vez para o projeto todo:

- [ ] **DESIGN.md** (design system completo extraído do Stitch)

### Nomes dos arquivos

```
design-reference/
├── screenshots/
│   ├── 01-dashboard.png
│   ├── 02-flag-list.png
│   ├── 03-flag-detail-targeting.png
│   ├── 04-flag-detail-product-card.png
│   ├── 05-flag-detail-metrics.png
│   ├── 06-flag-detail-lifecycle.png
│   ├── 07-integrations-hub.png
│   ├── 08-integrations-detail.png
│   ├── 09-plugins-installed.png
│   ├── 10-marketplace-browse.png
│   ├── 11-marketplace-detail.png
│   ├── 12-developer-portal.png
│   ├── 13-earnings.png
│   ├── 14-settings.png
│   ├── 15-audit-log.png
│   ├── 16-login.png
│   └── 17-onboarding.png
├── react-exports/
│   ├── Dashboard.jsx
│   ├── FlagList.jsx
│   ├── FlagDetailTargeting.jsx
│   ├── FlagDetailProductCard.jsx
│   ├── FlagDetailMetrics.jsx
│   ├── FlagDetailLifecycle.jsx
│   ├── IntegrationsHub.jsx
│   ├── IntegrationDetail.jsx
│   ├── PluginsInstalled.jsx
│   ├── MarketplaceBrowse.jsx
│   ├── MarketplaceDetail.jsx
│   ├── DeveloperPortal.jsx
│   ├── Earnings.jsx
│   ├── Settings.jsx
│   ├── AuditLog.jsx
│   ├── Login.jsx
│   └── Onboarding.jsx
└── DESIGN.md
```

---

## Dicas Importantes

1. **Não tente usar o HTML/JSX do Stitch como código final** — ele é estático, sem state, sem API. Use como referência visual e skeleton.

2. **O DESIGN.md é o asset mais valioso** — ele garante que o Claude Code vai produzir componentes visualmente consistentes com o que você desenhou no Stitch.

3. **Screenshots são seu backup** — se algo se perder no export, o screenshot garante que você sabe como era o design.

4. **Exporte ANTES de fazer mudanças grandes** — o Stitch pode perder contexto em edições complexas. Exporte cada tela assim que ela ficar boa.

5. **O Stitch MCP Server é o futuro** — quando estiver estável, vai ser possível ler designs direto do Stitch no Claude Code sem export manual. Fique de olho.

---

*Com esse pipeline, seus designs do Stitch viram componentes Next.js production-ready sem perder a fidelidade visual.*
