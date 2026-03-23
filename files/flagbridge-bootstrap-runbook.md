# FlagBridge — Bootstrap Automation Runbook

> **Tool: Claude Code (Terminal CLI)**
> **Why Claude Code and not Cowork?**
> Cowork is great for file/task management on your desktop, but this workflow requires terminal access to `gh` CLI, `npm`, `git`, and shell commands — Claude Code is the right tool because it runs in your terminal, executes commands directly, and can orchestrate multi-step workflows with your approval at each stage.

---

## Prerequisites

Before running this, make sure you have installed:

```bash
# Claude Code
npm install -g @anthropic-ai/claude-code

# GitHub CLI
brew install gh  # or: sudo apt install gh
gh auth login    # Authenticate with your GitHub account

# Node.js (for npm scope reservation)
node --version   # Must be >= 18

# Docker (for later local dev)
docker --version
docker compose version
```

---

## How to Run

1. Create a working directory: `mkdir ~/flagbridge-bootstrap && cd ~/flagbridge-bootstrap`
2. Save the `CLAUDE.md` file below into this directory
3. Run: `claude`
4. In the Claude Code session, say: **"Execute the FlagBridge bootstrap plan step by step. Ask for my approval before any payment or destructive action."**
5. Claude Code will execute each step, pausing for your confirmation at payment checkpoints

---

## CLAUDE.md (paste this into ~/flagbridge-bootstrap/CLAUDE.md)

```markdown
# FlagBridge Bootstrap Project

## Context
FlagBridge is an open-core Feature Flag Management platform.
This project bootstraps all the infrastructure: GitHub org, repos, npm scope, domain registration, and initial monorepo scaffolding.

## Important Rules
- ALWAYS ask for explicit human approval before any action that involves payment, account creation on external services, or irreversible operations
- Show a clear summary of what will happen BEFORE executing each step
- Use `gh` CLI for all GitHub operations
- Use `npm` for npm scope operations
- For domain registration: only open the browser with the URL and instructions — do NOT attempt to automate payment
- Log every action to ./bootstrap-log.md with timestamp
- If any step fails, stop and ask for guidance before continuing

## Execution Plan

### Phase 1: GitHub Organization & Repos

Step 1.1 — Create GitHub Organization
- Create org "flagbridge" on GitHub using `gh` CLI
- If org name is taken, suggest alternatives and ask user

Step 1.2 — Create repositories in the org:
- flagbridge/flagbridge (main monorepo — public)
- flagbridge/docs (documentation site — public)
- flagbridge/sdk-node (Node.js/TypeScript SDK — public)
- flagbridge/sdk-go (Go SDK — public)
- flagbridge/sdk-python (Python SDK — public)
- flagbridge/openfeature-provider (OpenFeature Provider — public)
- flagbridge/plugin-sdk (Plugin SDK for developers — public)
- flagbridge/helm-charts (Kubernetes Helm charts — public)
- flagbridge/flagbridge-pro (Pro plugin — private)

Step 1.3 — Configure org settings:
- Set org description: "Feature flags with product intelligence. Open source."
- Set org URL: https://flagbridge.io
- Enable Discussions on main repo
- Add topics to main repo: feature-flags, feature-management, open-source, openfeature, golang, nextjs

### Phase 2: npm Scope Reservation

Step 2.1 — Reserve @flagbridge npm scope
- Run: npm login (if not authenticated)
- Create a placeholder package to reserve the scope
- Packages to eventually publish:
  - @flagbridge/sdk-node
  - @flagbridge/sdk-react
  - @flagbridge/openfeature-provider
  - @flagbridge/plugin-sdk
  - @flagbridge/create-plugin (CLI scaffolding tool)
  - @flagbridge/cli

Step 2.2 — Create initial placeholder package
- Create a minimal package.json for @flagbridge/sdk-node
- Publish as 0.0.1 with "placeholder" tag to reserve the name

### Phase 3: Domain Registration [MANUAL APPROVAL REQUIRED]

⚠️ PAYMENT CHECKPOINT — Do NOT automate payment. Only provide instructions.

Step 3.1 — Check domain availability
- Open browser to Cloudflare Registrar: https://dash.cloudflare.com
- Domains to check and register:
  - flagbridge.io (primary)
  - flagbridge.dev (secondary, HTTPS-only by default)
  - flagbridge.com (if available, defensive registration)

Step 3.2 — Provide user with step-by-step instructions:
1. Go to https://dash.cloudflare.com/domains
2. Search for "flagbridge.io"
3. If available, add to cart
4. Repeat for flagbridge.dev
5. Complete checkout (payment is manual)
6. After purchase, configure DNS:
   - A record → pointing to future server IP (placeholder: 127.0.0.1)
   - CNAME docs.flagbridge.io → flagbridge.github.io (for docs site)

Step 3.3 — Wait for user confirmation that domains are registered before proceeding

### Phase 4: Monorepo Scaffolding

Step 4.1 — Clone the main repo and set up monorepo structure:

```
flagbridge/
├── CLAUDE.md                    # Claude Code project instructions
├── README.md                    # Project overview (bilingual EN/pt)
├── LICENSE                      # Apache 2.0
├── CONTRIBUTING.md              # Contribution guide
├── CONTRIBUTORS.md              # Contributors list
├── docker-compose.yml           # Local dev environment
├── docker-compose.pro.yml       # Pro plugin overlay
├── Makefile                     # Common commands
├── .github/
│   ├── ISSUE_TEMPLATE/
│   ├── PULL_REQUEST_TEMPLATE.md
│   ├── workflows/
│   │   ├── ci.yml               # CI pipeline
│   │   └── release.yml          # Release pipeline
│   └── FUNDING.yml
├── apps/
│   ├── api/                     # Go API server
│   │   ├── cmd/
│   │   │   └── server/
│   │   │       └── main.go
│   │   ├── internal/
│   │   │   ├── flag/            # Flag domain
│   │   │   ├── eval/            # Evaluation engine
│   │   │   ├── project/         # Project domain
│   │   │   ├── auth/            # Authentication
│   │   │   ├── plugin/          # Plugin runtime
│   │   │   ├── integration/     # Integration layer
│   │   │   └── api/             # HTTP handlers
│   │   ├── pkg/
│   │   │   ├── database/        # PostgreSQL client
│   │   │   ├── config/          # Configuration
│   │   │   └── middleware/      # HTTP middleware
│   │   ├── migrations/          # SQL migrations
│   │   ├── go.mod
│   │   ├── go.sum
│   │   └── Dockerfile
│   └── web/                     # Next.js Admin UI + Landing Page
│       ├── src/
│       │   ├── app/
│       │   │   ├── [locale]/
│       │   │   │   ├── (marketing)/
│       │   │   │   ├── (admin)/
│       │   │   │   └── (developer)/
│       │   │   └── api/
│       │   ├── components/
│       │   │   └── ui/
│       │   ├── lib/
│       │   ├── messages/
│       │   │   ├── en.json
│       │   │   └── pt.json
│       │   └── styles/
│       ├── public/
│       ├── next.config.ts
│       ├── tailwind.config.ts
│       ├── tsconfig.json
│       ├── package.json
│       └── Dockerfile
├── packages/
│   ├── sdk-node/                # @flagbridge/sdk-node
│   │   ├── src/
│   │   ├── package.json
│   │   └── tsconfig.json
│   ├── sdk-react/               # @flagbridge/sdk-react
│   │   ├── src/
│   │   ├── package.json
│   │   └── tsconfig.json
│   ├── plugin-sdk/              # @flagbridge/plugin-sdk
│   │   ├── src/
│   │   ├── package.json
│   │   └── tsconfig.json
│   ├── openfeature-provider/    # @flagbridge/openfeature-provider
│   │   ├── src/
│   │   ├── package.json
│   │   └── tsconfig.json
│   ├── cli/                     # @flagbridge/cli
│   │   ├── src/
│   │   ├── package.json
│   │   └── tsconfig.json
│   └── create-plugin/           # @flagbridge/create-plugin
│       ├── templates/
│       ├── src/
│       ├── package.json
│       └── tsconfig.json
└── docs/                        # Documentation source
    ├── docs/
    │   ├── getting-started/
    │   ├── guides/
    │   ├── api-reference/
    │   ├── sdk/
    │   ├── plugins/
    │   └── integrations/
    ├── i18n/
    │   ├── en/
    │   └── pt/
    └── docusaurus.config.ts
```

Step 4.2 — Initialize Go module for API server:
- `cd apps/api && go mod init github.com/flagbridge/flagbridge/apps/api`
- Add initial dependencies: chi router, pgx, zerolog

Step 4.3 — Initialize Next.js app for Admin UI:
- `cd apps/web && npx create-next-app@latest . --typescript --tailwind --app --src-dir`
- Install: next-intl, @radix-ui/themes, @tanstack/react-query, zustand
- Set up i18n with en.json and pt.json base files

Step 4.4 — Initialize SDK packages with TypeScript:
- Each package gets a base package.json, tsconfig.json, and src/index.ts
- Use TypeScript project references for monorepo

Step 4.5 — Create docker-compose.yml for local dev:
- flagbridge-api (Go, port 8080)
- flagbridge-web (Next.js, port 3000)  
- postgres (port 5432)
- redis (port 6379, optional)

Step 4.6 — Create initial SQL migrations:
- 001_create_users.sql
- 002_create_projects.sql
- 003_create_environments.sql
- 004_create_flags.sql
- 005_create_flag_states.sql
- 006_create_audit_log.sql
- 007_create_plugins.sql

Step 4.7 — Create README.md (bilingual):
- English and Portuguese sections
- Quick start guide
- Architecture overview
- Contributing link
- License badge

Step 4.8 — Create initial GitHub Actions CI:
- Go: lint, test, build
- Next.js: lint, build
- Docker: build image

### Phase 5: Initial Commit & Push

Step 5.1 — Stage all files, create initial commit:
- Message: "feat: initial monorepo scaffolding for FlagBridge"
- Push to main branch

Step 5.2 — Create initial GitHub Issues for MVP tracking:
- "Set up Go API server with health endpoint"
- "Set up Next.js Admin UI shell with i18n"
- "Implement flag CRUD API"
- "Implement flag evaluation engine"
- "Create Node.js SDK"
- "Create Go SDK"
- "Design Admin UI dashboard"
- "Set up Docker Compose for local development"
- "Create documentation site"
- "Implement webhook integration (CE)"

Step 5.3 — Create GitHub Project board "FlagBridge MVP" and add issues

### Phase 6: Verification Checklist

After all steps, verify:
- [ ] GitHub org exists: https://github.com/flagbridge
- [ ] Main repo exists: https://github.com/flagbridge/flagbridge
- [ ] All repos created and accessible
- [ ] npm scope @flagbridge is reserved
- [ ] Monorepo structure is correct
- [ ] Go module initializes without errors
- [ ] Next.js app builds without errors  
- [ ] Docker Compose starts without errors
- [ ] CI pipeline runs on push
- [ ] README is bilingual and complete
- [ ] Domains registered (manual confirmation)

## Payment Checkpoints Summary
These are the ONLY steps that require payment:
1. Domain registration (flagbridge.io, flagbridge.dev) — ~$50-80/year total
2. npm org/scope — free for public packages
3. GitHub org — free for public repos

Everything else in this plan is free.
```

---

## Running the Bootstrap

```bash
# 1. Create working directory
mkdir ~/flagbridge-bootstrap && cd ~/flagbridge-bootstrap

# 2. Save the CLAUDE.md above into this directory

# 3. Start Claude Code
claude

# 4. In Claude Code, say:
# "Execute the FlagBridge bootstrap plan step by step.
#  Ask for my approval before any payment or destructive action.
#  Log each completed step to bootstrap-log.md."
```

Claude Code will:
1. Read the CLAUDE.md and understand the full plan
2. Execute Phase 1 (GitHub) — creating org and repos via `gh` CLI
3. Execute Phase 2 (npm) — reserving the scope
4. **PAUSE at Phase 3** — give you instructions to manually register domains with payment
5. After your confirmation, proceed to Phase 4 (monorepo scaffolding)
6. Make the initial commit and push
7. Create issues and project board
8. Run the verification checklist

**Estimated total time:** 15-20 minutes of your time (mostly waiting + domain payment)

---

## Post-Bootstrap: What's Next

After the bootstrap is complete, you can use Claude Code for the actual development:

```bash
cd ~/flagbridge-bootstrap/flagbridge

# Start building the Go API
claude "Implement the flag CRUD endpoints in apps/api following the API design in the product doc. Use chi router, pgx for PostgreSQL, and zerolog for logging."

# Start building the Next.js Admin UI
claude "Set up the Admin UI shell in apps/web with sidebar navigation, dark theme, and bilingual support using next-intl. Follow the Stitch designs as reference."

# Build the Node.js SDK
claude "Implement the @flagbridge/sdk-node package with flag evaluation, SSE real-time updates, and OpenFeature provider compatibility."
```

---

## Why Claude Code Over Cowork

| Aspect | Claude Code (CLI) | Cowork (Desktop) |
|--------|-------------------|-------------------|
| Terminal access (`gh`, `npm`, `git`) | ✅ Native | ❌ No terminal |
| File creation & editing | ✅ | ✅ |
| Run shell commands | ✅ | ❌ |
| Multi-step orchestration | ✅ With CLAUDE.md plan | ✅ But limited to file tasks |
| GitHub API integration | ✅ Via `gh` CLI + MCP | ❌ |
| npm publish | ✅ | ❌ |
| Docker commands | ✅ | ❌ |
| Approval checkpoints | ✅ Built-in permission system | ✅ |
| Visual diffs | ❌ (text-based) | ✅ |
| Best for | **Infrastructure bootstrap, coding, CLI workflows** | **File management, document editing, non-dev tasks** |

For this bootstrap workflow, Claude Code is the clear winner. Cowork would be great for tasks like "organize my FlagBridge product docs" or "update all my Notion pages" — but not for creating GitHub orgs and pushing code.

---

*After bootstrap, you'll have a fully set up GitHub organization, reserved npm scope, registered domains, and a monorepo ready for development. Total investment: ~$50-80 for domains + your time.*
