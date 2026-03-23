# FlagBridge — Stitch Prompt Sequence for Admin UI Generation

> **How to use this document:**
> Execute each prompt **one at a time** in Google Stitch, in order.
> Wait for the result, verify it looks correct, then proceed to the next prompt.
> Do NOT combine multiple prompts — Stitch works best with incremental, focused changes.
> Save/screenshot after each step that works well.

---

## Phase 1: Initial App Setup & Vibe

### Prompt 1 — App Foundation

```
A professional, modern SaaS admin dashboard for a developer tool called "FlagBridge" — a feature flag management platform. Dark theme with deep navy (#0F172A) background, electric blue (#3B82F6) as primary accent, and subtle gray (#1E293B) cards. Clean, minimal, developer-focused aesthetic. Sans-serif font. The app should feel like a premium open-source tool — think Linear meets Vercel Dashboard.
```

---

## Phase 2: Sidebar Navigation

### Prompt 2 — Sidebar Structure

```
Add a left sidebar navigation with the FlagBridge logo at the top (a simple bridge icon with the text "FlagBridge"). The sidebar should have these navigation items with small icons: Dashboard, Projects, Flags, Integrations, Plugins, Marketplace, Audit Log, Settings. The sidebar should be narrow, dark, with the active item highlighted in the primary blue color. At the bottom of the sidebar, show a small user avatar with the name "Gabriel G." and a role badge "Admin".
```

---

## Phase 3: Dashboard Screen

### Prompt 3 — Dashboard Overview Cards

```
On the main Dashboard page, add a top row of four stat cards in a horizontal grid. The cards are: "Total Flags" showing 127 with a small flag icon, "Active Flags" showing 43 with a green dot, "Stale Flags" showing 12 with an orange warning icon, and "Tech Debt Score" showing 6.2/10 with a yellow gauge icon. Each card has a dark background with subtle border, white numbers, and gray labels. Below the number, show a small trend indicator like "+3 this week" in green or "-2 this week" in red.
```

### Prompt 4 — Dashboard Activity Feed & Chart

```
Below the stat cards on the Dashboard, add two columns. Left column (wider, about 60%): a line chart titled "Flag Evaluations (7 days)" showing evaluation volume over time with a smooth blue line on a dark background, with subtle grid lines. Right column (40%): an "Activity Feed" list showing recent actions like "checkout-v2 toggled ON in production — Gabriel G. — 2 min ago" and "new-pricing flag created — Ana M. — 15 min ago". Each activity item has a small colored icon indicating the action type. Show 5-6 items.
```

### Prompt 5 — Dashboard Integration Health

```
Below the chart and activity feed on the Dashboard, add a section titled "Integration Health". Show a horizontal row of 4 small status cards: "Mixpanel" with a green "Healthy" badge, "Resend" with a green "Healthy" badge, "Kafka" with a yellow "Degraded" badge, and "RD Station" with a green "Healthy" badge. Each card shows the provider name, a small provider icon, and the health status badge. Compact layout, same dark card style.
```

---

## Phase 4: Flags List Screen

### Prompt 6 — Flag List Page

```
Create a new screen: the Flags List page. At the top, show the page title "Flags" on the left, and a blue "Create Flag" button on the right. Below, add a filter bar with: a search input field, a dropdown for "Environment" (showing "Production"), a dropdown for "Status" (showing "All"), and a dropdown for "Owner" (showing "All"). Below the filters, show a table with columns: Flag Key, Name, Status, Environment, Owner, Last Modified, and a toggle switch. Show 6 rows of sample flag data like "checkout-v2", "new-pricing-page", "dark-mode", "onboarding-flow-v3", "ai-recommendations", "holiday-banner". Status badges should be colored: green for Active, blue for Planning, gray for Archived, orange for Rolled Out.
```

---

## Phase 5: Flag Detail Screen

### Prompt 7 — Flag Detail Header & Toggle

```
Create a new screen: the Flag Detail page for a flag called "checkout-v2". At the top, show a breadcrumb: "Project: Web App > Flags > checkout-v2". Below, show the flag name "Checkout V2 Redesign" as a large heading, with status badge "Active" in green, and an owner avatar labeled "Gabriel G." to the right. Next to the owner, show a large toggle switch that is currently ON (blue). Below the header, add a horizontal tab bar with tabs: "Targeting", "Product Card", "Metrics", "Lifecycle", "Integrations".
```

### Prompt 8 — Flag Detail Targeting Tab

```
On the Flag Detail page, show the "Targeting" tab as active. Display an "Environment" selector showing three pills: "Development", "Staging", "Production" (Production is selected/highlighted). Below, show the targeting rules section. First rule: "Default" — toggle ON, value "true". Second rule card: "Rule 1 — Beta Users" with conditions: "attribute 'plan' equals 'pro'" AND "attribute 'country' is one of 'BR, US, PT'" — serve variant "true", rollout 100%. Third rule card: "Rule 2 — Gradual Rollout" with condition "percentage rollout: 25%" — show a slider at 25%. Each rule card is a dark bordered card with clear hierarchy.
```

### Prompt 9 — Flag Detail Product Card Tab

```
On the Flag Detail page, switch to the "Product Card" tab. Show a structured form/display with these fields filled in: "Hypothesis" with text "The new checkout flow will increase conversion rate by 15% for Pro users in Brazil", "Success Metrics" showing two metric pills: "Conversion Rate > 5%" and "Cart Abandonment < 30%", "Go/No-Go Criteria" with text "Minimum 2 weeks of data with 10k evaluations", "Decision" showing a dropdown currently set to "Pending" in yellow, "OKR Link" showing a clickable link "Q2 2026 — Increase Revenue 20%", and "Owner" showing "Gabriel G." with avatar. At the bottom, a timeline showing: "Created Mar 1" → "Activated Mar 5" → "25% Rollout Mar 10" → "Now". Clean card layout with labels in gray and values in white.
```

### Prompt 10 — Flag Detail Metrics Tab

```
On the Flag Detail page, switch to the "Metrics" tab. Show three charts in a grid: Top left — "Adoption Rate" area chart showing percentage of requests evaluating this flag over 14 days, currently at 78%, blue fill. Top right — "Error Rate by Variant" bar chart comparing "control" (0.3% errors, gray bar) vs "new" (0.5% errors, blue bar). Bottom spanning full width — "Evaluations Over Time" line chart with two lines: "variant: control" in gray and "variant: new" in blue, showing volume over 14 days. Above the charts, show summary cards: "Total Evaluations: 245,302", "Unique Users: 18,430", "Error Delta: +0.2%". Dark theme charts with subtle grid lines.
```

### Prompt 11 — Flag Detail Lifecycle Tab

```
On the Flag Detail page, switch to the "Lifecycle" tab. Show a horizontal timeline at the top with milestones: "Created" (Mar 1, green dot), "Planning" (Mar 1-4, gray), "Activated" (Mar 5, blue dot), "25% Rollout" (Mar 10, blue dot), "Current" (Mar 22, pulsing blue dot). Below the timeline, show a "Cleanup Status" card with: "Days since creation: 21 days", "Days at current rollout: 12 days", "Cleanup urgency: Low" with a green badge. Below that, a "Code References" section showing: "Found in 3 files" with file paths like "src/checkout/CheckoutPage.tsx:42", "src/api/checkout.ts:15", "src/hooks/useCheckout.ts:8" — each as a clickable link with a GitHub icon. At the bottom, show active lifecycle rules: "Auto-notify if 100% for 30+ days" (enabled, green) and "Auto-archive after 90 days rolled out" (enabled, green).
```

---

## Phase 6: Integrations Hub Screen

### Prompt 12 — Integrations Hub Overview

```
Create a new screen: the Integrations Hub page. Title "Integrations" at top left. Below, show three category tabs: "Messaging & Email", "Analytics & Events", "Technical Queues". The "Messaging & Email" tab is active. Show a grid of integration cards (3 columns). Each card shows: a provider logo/icon, the provider name, a short description, connection status (green "Connected" or gray "Not Connected"), and a "Configure" or "Connect" button. Show these providers: Resend (Connected), RD Station (Connected), SendGrid (Not Connected), Mailchimp (Not Connected), Brevo (Not Connected), Customer.io (Not Connected). Cards have dark background with subtle border. Connected cards have a subtle green left border accent.
```

### Prompt 13 — Integrations Analytics Tab

```
On the Integrations Hub, switch to the "Analytics & Events" tab. Show the same card grid layout with these providers: Mixpanel (Connected, green accent), Amplitude (Not Connected), Segment (Connected, green accent), PostHog (Not Connected), Google Analytics 4 (Not Connected), Rudderstack (Not Connected). Same card style as messaging tab.
```

### Prompt 14 — Integrations Queues Tab

```
On the Integrations Hub, switch to the "Technical Queues" tab. Show the card grid with: Amazon SQS (Connected, green accent), Apache Kafka (Connected, green accent), RabbitMQ (Not Connected), NATS (Not Connected), Redis Streams (Not Connected), GCP Pub/Sub (Not Connected). Same card style. Add a small "Pro" badge on the tab title to indicate this is a Pro feature.
```

### Prompt 15 — Integration Detail Page (Resend)

```
Create a new screen: the Integration Detail page for Resend. At the top, show a breadcrumb: "Integrations > Messaging > Resend". Below, the Resend logo, name, connection status "Connected" (green), and a "Test Connection" button. Then show two sections. First section "Configuration": a form with fields for API Key (masked with dots), From Email ("flags@company.com"), and From Name ("FlagBridge"). Second section "Trigger Rules": a list of configured rules. Rule 1: "When flag is rolled out in production → Send email to flag owner using template 'rollout-notification'". Rule 2: "When stale flag detected (30+ days) → Send reminder email to flag owner". Each rule has an enable/disable toggle and an "Edit" button. At the bottom, an "Event Log" section showing the last 5 events with timestamp, status (green "Sent" or red "Failed"), and event description.
```

---

## Phase 7: Plugin Manager Screen

### Prompt 16 — Installed Plugins Page

```
Create a new screen: the Plugin Manager page. Title "Plugins" with two tabs: "Installed" (active) and "Marketplace". Show a list of installed plugins as horizontal cards. Each card shows: plugin icon, plugin name, version, author, a short description, an enable/disable toggle, a "Configure" button, and a health status indicator (green dot). Show these plugins: "Slack Alerts v1.2.0 by FlagBridge — Send Slack notifications on flag changes" (enabled), "Linear Sync v1.0.0 by FlagBridge — Bidirectional sync with Linear issues" (enabled), "GitHub Code Refs v0.9.0 by Community — Find flag references in GitHub repos" (enabled), "Datadog Metrics v1.1.0 by Community — Export flag metrics to Datadog" (disabled, gray). At the top right, show a "Browse Marketplace" button.
```

---

## Phase 8: Marketplace Screen

### Prompt 17 — Marketplace Browse Page

```
Create a new screen: the Plugin Marketplace page. Title "Marketplace" at top left. Below, show a search bar and filter dropdowns for "Category" (All, Integration, Analytics, Security, UI, Automation, Data) and "Price" (All, Free, Paid). Below filters, show a grid of plugin cards (3 columns). Each card shows: plugin icon, plugin name, author, short description, category badge, star rating (out of 5), install count, and price tag ("Free" in green or "$9/mo" in blue). Show 6 plugins: "PagerDuty Alerts — Free — 4.8 stars — 1.2k installs", "BigQuery Export — $9/mo — 4.5 stars — 340 installs", "Advanced RBAC — $15/mo — 4.7 stars — 890 installs", "Custom Themes — Free — 4.2 stars — 2.1k installs", "Jira Sync — Free — 4.6 stars — 1.8k installs", "Snowflake Connector — $12/mo — 4.4 stars — 220 installs". Cards have dark background, subtle border, and a blue "Install" button.
```

### Prompt 18 — Marketplace Plugin Detail

```
Create a new screen: the Marketplace Plugin Detail page for "PagerDuty Alerts". At the top, show the plugin icon (large), name "PagerDuty Alerts", author "FlagBridge Official" with a blue verified badge, version "2.1.0", and a large blue "Install" button. Below, show rating "4.8" with stars, "1,247 installs", "Last updated: Mar 15, 2026". Then a horizontal tab bar: "Overview" (active), "Reviews", "Changelog", "Permissions". The Overview tab shows a description paragraph, a list of features with checkmarks, 3 screenshot thumbnails of the plugin in action, and a "Compatibility" section showing "FlagBridge >= 1.2.0". On the right sidebar, show plugin metadata: Category "Integration", License "MIT", Price "Free", and links to Documentation and Repository.
```

---

## Phase 9: Developer Portal

### Prompt 19 — Developer Portal

```
Create a new screen: the Developer Portal page. Title "Developer Portal" with subtitle "Build, publish, and monetize plugins for FlagBridge". Show three large cards in a row: Card 1 "Documentation" with a book icon, subtitle "Plugin SDK reference, guides, and tutorials", and a "View Docs" button. Card 2 "API Explorer" with a terminal icon, subtitle "Interactive API documentation with live examples", and a "Open Explorer" button. Card 3 "Plugin Sandbox" with a beaker/flask icon, subtitle "Test your plugin in an isolated environment", and a "Launch Sandbox" button. Below the cards, show a section "Your Published Plugins" with a table: Plugin Name, Version, Status (Published/In Review/Draft), Downloads, Rating, Revenue. Show 2 sample rows. At the bottom right, a "Publish New Plugin" button.
```

---

## Phase 10: Earnings Dashboard (Marketplace Developer)

### Prompt 20 — Developer Earnings Page

```
Create a new screen: the Developer Earnings page. Title "Earnings" with subtitle "Track your plugin revenue". Show four stat cards at the top: "Total Revenue" showing "$1,247.00", "This Month" showing "$189.00", "Total Installs" showing "3,412", and "Average Rating" showing "4.6 stars". Below, show a line chart "Monthly Revenue" with data for the last 6 months. Below the chart, show a table "Recent Transactions" with columns: Date, Plugin, Customer, Amount, Status (Completed/Pending). Show 5 rows of sample data. On the right sidebar, show "Payout Settings" with current payout method (Stripe), next payout date, and a "Manage Payouts" button. The revenue split "80% you / 20% FlagBridge" should be shown as a small info badge.
```

---

## Phase 11: Settings Screen

### Prompt 21 — Settings Page

```
Create a new screen: the Settings page. Show a left sidebar within the page with setting categories: "General", "Team", "API Keys", "Billing", "Integrations". The "General" tab is active. Show a form with: "Project Name" text input (value: "Web App"), "Default Environment" dropdown (value: "Production"), "Default Language" dropdown (value: "English"), and "Timezone" dropdown (value: "America/Sao_Paulo"). Below, a section "Danger Zone" with a red-bordered card containing "Delete Project" with a red button and a warning text. Clean form layout with labels above inputs, dark inputs with subtle borders.
```

### Prompt 22 — Team Settings

```
On the Settings page, switch to the "Team" tab. Show a "Team Members" section with an "Invite Member" button at the top right. Below, a table with columns: Avatar, Name, Email, Role (dropdown showing Admin/Editor/Viewer), Last Active, and a three-dot menu for actions. Show 4 team members: "Gabriel Gripp — gabriel@company.com — Admin — Active now", "Ana Martins — ana@company.com — Editor — 2 hours ago", "Carlos Silva — carlos@company.com — Viewer — 1 day ago", "Maria Santos — maria@company.com — Editor — 3 hours ago". Below the table, show "Pending Invitations" with 1 pending invite showing email and a "Resend" link.
```

### Prompt 23 — API Keys Settings

```
On the Settings page, switch to the "API Keys" tab. Show a "Generate New Key" button at the top right. Below, a table of existing API keys with columns: Key Name, Key Preview (showing first 8 chars then masked like "fb_sk_a3f2..."), Environment, Created, Last Used, and actions (Copy, Rotate, Revoke). Show 3 keys: "Production SDK — fb_sk_prod... — Production — Mar 1, 2026 — 2 min ago", "Staging SDK — fb_sk_stag... — Staging — Feb 15, 2026 — 1 hour ago", "CI/CD Pipeline — fb_sk_cicd... — All — Jan 10, 2026 — 5 min ago". Each row has subtle action icons on hover.
```

---

## Phase 12: Audit Log Screen

### Prompt 24 — Audit Log Page

```
Create a new screen: the Audit Log page. Title "Audit Log" at the top. Show filters: a date range picker, a "User" dropdown (All), an "Action" dropdown (All, Created, Updated, Toggled, Archived, Deleted), and an "Entity" dropdown (All, Flag, Project, Integration, Plugin). Below, show a timeline-style list of audit events. Each event shows: timestamp, user avatar and name, action description with colored verb (blue for "toggled", green for "created", orange for "updated", red for "deleted"), and the entity name as a clickable link. Show 8 events with varied actions. For "toggled" and "updated" events, show a small "View diff" link. Dark cards, clean typography, subtle timestamps in gray.
```

---

## Phase 13: Bilingual Support Indicator

### Prompt 25 — Language Switcher

```
In the top-right corner of the header bar (visible on all screens), add a language switcher showing a small globe icon with "EN" text. When conceptually clicked, it would switch to "pt". Place it next to a notification bell icon and the user avatar. The language switcher should be subtle and minimal, not distracting from the main content. Also add a small "Pro" badge next to the FlagBridge logo in the sidebar to indicate the current plan.
```

---

## Phase 14: Mobile Responsive View

### Prompt 26 — Mobile Dashboard View

```
Create a mobile-responsive version of the Dashboard screen. The sidebar collapses into a hamburger menu at the top. The four stat cards stack vertically (2 per row on mobile). The chart takes full width. The activity feed appears below the chart as a scrollable list. Keep the same dark theme and colors. Compact spacing suitable for mobile screens.
```

---

## Phase 15: Login / Onboarding

### Prompt 27 — Login Screen

```
Create a Login screen for FlagBridge. Center of the page: the FlagBridge logo (bridge icon + "FlagBridge" text) above a login form. The form has: Email input, Password input, a "Sign In" blue button, and a "Sign in with SSO" secondary button below. Under the form, a "Forgot password?" link and a "Don't have an account? Start free" link. The background is a very dark gradient with a subtle geometric pattern. Clean, minimal, trustworthy feel — like a developer tool login page.
```

### Prompt 28 — Empty State / Onboarding

```
Create a screen showing the Dashboard in an empty state for a new user. Instead of charts and data, show a friendly onboarding card in the center with the title "Welcome to FlagBridge!", subtitle "Let's set up your first feature flag in under 2 minutes." Below, show 3 numbered steps as a horizontal progress bar: "1. Create a Project" (active, blue), "2. Add your first Flag" (gray), "3. Install an SDK" (gray). Below the steps, a large blue "Create Your First Project" button. The sidebar shows the same navigation but with subtle notification dots suggesting where to start. Encouraging, clean, not overwhelming.
```

---

# Quick Reference: Stitch Best Practices Applied

1. **One major change per prompt** — every prompt above focuses on one screen or one component
2. **Specific UI/UX keywords used** — "navigation bar", "card layout", "toggle switch", "breadcrumb", "tab bar", "stat cards"
3. **Colors defined early** — dark navy #0F172A, blue #3B82F6, gray #1E293B established in Prompt 1
4. **Vibe set from the start** — "professional, modern SaaS", "Linear meets Vercel Dashboard"
5. **Incremental complexity** — starts with shell/vibe, then sidebar, then screen by screen
6. **Real sample data** — flags, users, metrics all have realistic names and numbers
7. **Visual hierarchy described** — font sizes, spacing, column widths, color emphasis all specified

---

# Post-Generation Refinement Prompts

Use these if specific elements need adjustment after the initial generation:

```
Refinement: Make the sidebar 20% narrower and increase the main content area width.
```

```
Refinement: Change the chart color from blue to a gradient from blue to purple.
```

```
Refinement: Add a "Last 7 days / Last 30 days / Last 90 days" date range selector above the Dashboard charts.
```

```
Refinement: On the Flags List page, add a "Bulk Actions" dropdown next to the search bar with options "Archive Selected" and "Change Owner".
```

```
Refinement: On all screens, make the page titles slightly larger and add a subtle underline accent in the primary blue color.
```

```
Refinement: Switch all text to Portuguese. Button labels, navigation items, page titles, descriptions — everything should be in Brazilian Portuguese.
```

---

*Total: 28 primary prompts + 6 refinement prompts = complete FlagBridge Admin UI in Stitch*
