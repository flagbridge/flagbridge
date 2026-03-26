import { useTranslations } from "next-intl";
import type { Metadata } from "next";
import { Header } from "@/components/header";
import { Footer } from "@/components/footer";
import { Badge } from "@/components/badge";
import { CodeBlock } from "@/components/code-block";
import { SdkTabs } from "@/components/sdk-tabs";

export const metadata: Metadata = {
  title: "Developers — FlagBridge",
  description:
    "Feature flag management built for developers. SDKs for every stack, OpenFeature support, testing API, plugin ecosystem, and more.",
  openGraph: {
    title: "Developers — FlagBridge",
    description:
      "Feature flag management built for developers. SDKs, testing API, plugins, and more.",
    url: "https://flagbridge.io/developers",
  },
};

/* ─── 1. Hero ─── */
function DeveloperHero() {
  const t = useTranslations("developers");
  return (
    <section className="bg-[#0F172A] px-4 pt-32 pb-24 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-[1200px]">
        <div className="mx-auto mb-16 max-w-4xl text-center">
          <h1 className="mb-6 text-5xl font-bold text-[#F8FAFC] sm:text-6xl">
            {t("heroTitle")}
          </h1>
          <p className="mb-10 text-xl text-[#E2E8F0]">
            {t("heroSubtitle")}
          </p>
          <div className="flex flex-col items-center justify-center gap-4 sm:flex-row">
            <button className="w-full rounded-lg bg-[#3B82F6] px-6 py-3 font-medium text-white transition-colors hover:bg-blue-600 sm:w-auto">
              {t("getStartedFree")}
            </button>
            <button className="w-full rounded-lg border border-[#475569] bg-transparent px-6 py-3 font-medium text-[#E2E8F0] transition-colors hover:border-[#64748B] sm:w-auto">
              {t("buildPlugin")}
            </button>
            <a href="#sales" className="font-medium text-[#3B82F6] transition-colors hover:text-blue-400">
              {t("talkToSales")} →
            </a>
          </div>
        </div>

        {/* Terminal Mockup */}
        <div className="mx-auto max-w-4xl">
          <div className="overflow-hidden rounded-xl border border-[#334155] bg-[#1E293B] shadow-2xl">
            <div className="flex items-center gap-2 border-b border-[#334155] bg-[#0F172A] px-4 py-3">
              <div className="flex gap-2">
                <div className="h-3 w-3 rounded-full bg-[#EF4444]" />
                <div className="h-3 w-3 rounded-full bg-[#F97316]" />
                <div className="h-3 w-3 rounded-full bg-[#22C55E]" />
              </div>
              <div className="flex-1 text-center text-sm text-[#64748B]">terminal</div>
            </div>
            <div className="p-6 font-mono text-sm">
              <div className="text-[#22D3EE]">$ npm install @flagbridge/sdk-node</div>
              <div className="mt-3 text-[#22D3EE]">$ flagbridge init --project my-app</div>
              <div className="mt-2 text-[#22C55E]">✓ Connected to FlagBridge</div>
              <div className="text-[#22C55E]">✓ 3 flags synced</div>
              <div className="mt-3 text-[#22D3EE]">$ flagbridge toggle checkout-v2 --env production --on</div>
              <div className="mt-2 text-[#22C55E]">✓ checkout-v2 is now ON in production</div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

/* ─── 2. Quick Start ─── */
function QuickStart() {
  const t = useTranslations("developers");
  const steps = [
    { num: "1", title: t("step1Title"), code: "docker compose up -d", desc: t("step1Desc") },
    { num: "2", title: t("step2Title"), code: `POST /api/v1/projects/\nmy-app/flags`, desc: t("step2Desc") },
    { num: "3", title: t("step3Title"), code: `const enabled = \n  fb.isEnabled('flag')`, desc: t("step3Desc") },
  ];

  return (
    <section className="bg-[#111827] px-4 py-24 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-[1200px]">
        <div className="mb-16 text-center">
          <h2 className="mb-4 text-4xl font-bold text-[#F8FAFC]">{t("quickStartTitle")}</h2>
        </div>

        <div className="grid grid-cols-1 gap-8 md:grid-cols-3">
          {steps.map((step) => (
            <div key={step.num} className="rounded-xl border border-[#334155] bg-[#1E293B] p-5">
              <div className="mb-4 flex h-12 w-12 items-center justify-center rounded-lg bg-[#3B82F6] text-xl font-bold text-white">
                {step.num}
              </div>
              <h3 className="mb-4 text-xl font-semibold text-[#F8FAFC]">{step.title}</h3>
              <div className="mb-3 rounded-lg border border-[#334155] bg-[#0F172A] p-4">
                <pre className="font-mono text-sm text-[#22D3EE]">{step.code}</pre>
              </div>
              <p className="text-sm text-[#94A3B8]">{step.desc}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

/* ─── 3. SDKs & OpenFeature ─── */
function SdkSection() {
  const t = useTranslations("developers");
  return (
    <section className="bg-[#0F172A] px-4 py-24 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-[1200px]">
        <div className="mb-16 text-center">
          <h2 className="mb-4 text-4xl font-bold text-[#F8FAFC]">{t("sdkTitle")}</h2>
        </div>

        <div className="grid items-start gap-12 md:grid-cols-2">
          <div>
            <p className="mb-8 text-lg text-[#E2E8F0]">
              {t("sdkDescription")}
            </p>
            <div className="flex flex-wrap gap-3">
              <Badge>Node.js</Badge>
              <Badge>Go</Badge>
              <Badge>Python</Badge>
              <Badge>React</Badge>
              <Badge variant="blue">OpenFeature</Badge>
            </div>
          </div>
          <SdkTabs />
        </div>
      </div>
    </section>
  );
}

/* ─── 4. Testing API ─── */
function TestingApi() {
  const t = useTranslations("developers");

  const checkItems = [
    { title: t("testFeature1"), desc: t("testFeature1Desc") },
    { title: t("testFeature2"), desc: t("testFeature2Desc") },
    { title: t("testFeature3"), desc: t("testFeature3Desc") },
  ];

  return (
    <section className="bg-[#111827] px-4 py-24 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-[1200px]">
        <div className="mb-16 text-center">
          <h2 className="mb-4 text-4xl font-bold text-[#F8FAFC]">{t("testingTitle")}</h2>
        </div>

        <div className="grid items-start gap-12 md:grid-cols-2">
          <div>
            <p className="mb-8 text-lg text-[#E2E8F0]">
              {t("testingDescription")}
            </p>
            <div className="space-y-4">
              {checkItems.map((item) => (
                <div key={item.title} className="flex items-start gap-3">
                  <svg width={20} height={20} fill="none" viewBox="0 0 24 24" stroke="#22C55E" strokeWidth={2} className="mt-0.5 flex-shrink-0">
                    <path strokeLinecap="round" strokeLinejoin="round" d="m4.5 12.75 6 6 9-13.5" />
                  </svg>
                  <div>
                    <div className="mb-1 font-medium text-[#F8FAFC]">{item.title}</div>
                    <div className="text-sm text-[#94A3B8]">{item.desc}</div>
                  </div>
                </div>
              ))}
            </div>
          </div>

          <div className="min-w-0">
            <CodeBlock
              language="typescript"
              code={`const session = await fb.testing.createSession({
  environment: 'staging',
  ttl: '30m',
  overrides: {
    'checkout-v2': true,
    'payment-method': 'pix-only',
  },
});`}
            />
            <div className="mt-4 flex gap-3">
              <Badge variant="green">CE</Badge>
              <Badge variant="blue">Pro</Badge>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

/* ─── 5. API & Webhooks ─── */
function ApiWebhooks() {
  const t = useTranslations("developers");

  const endpoints = [
    { label: "Evaluate", path: "GET /api/v1/evaluate/:key" },
    { label: "Testing", path: "POST /api/v1/testing/sessions" },
    { label: "Webhooks", path: "POST /api/v1/webhooks" },
    { label: "Flags", path: "GET /api/v1/projects/:id/flags" },
    { label: "Integrations", path: "GET /api/v1/integrations" },
  ];

  return (
    <section className="bg-[#0F172A] px-4 py-24 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-[1200px]">
        <div className="grid gap-12 md:grid-cols-2">
          {/* REST API */}
          <div className="rounded-xl border border-[#334155] bg-[#1E293B] p-5">
            <h3 className="mb-6 text-2xl font-bold text-[#F8FAFC]">{t("restApiTitle")}</h3>
            <div className="mb-6 space-y-4">
              {endpoints.map((ep) => (
                <div key={ep.label}>
                  <div className="mb-2 text-xs uppercase tracking-wide text-[#94A3B8]">{ep.label}</div>
                  <div className="font-mono text-sm text-[#E2E8F0]">{ep.path}</div>
                </div>
              ))}
            </div>
            <button className="flex w-full items-center justify-center gap-2 rounded-lg bg-[#3B82F6] px-4 py-2 text-white transition-colors hover:bg-blue-600">
              {t("openApiExplorer")}
              <svg width={16} height={16} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 6H5.25A2.25 2.25 0 0 0 3 8.25v10.5A2.25 2.25 0 0 0 5.25 21h10.5A2.25 2.25 0 0 0 18 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" />
              </svg>
            </button>
          </div>

          {/* Webhooks */}
          <div className="rounded-xl border border-[#334155] bg-[#1E293B] p-5">
            <h3 className="mb-6 text-2xl font-bold text-[#F8FAFC]">{t("webhooksTitle")}</h3>
            <div className="mb-4 rounded-lg border border-[#334155] bg-[#0F172A] p-4">
              <pre className="overflow-x-auto font-mono text-xs leading-relaxed text-[#22D3EE]">
{`{
  "event": "flag.toggled",
  "timestamp": "2026-03-25T10:30:00Z",
  "data": {
    "flagKey": "checkout-v2",
    "environment": "production",
    "enabled": true,
    "triggeredBy": "alice@example.com"
  }
}`}
              </pre>
            </div>
            <p className="text-sm text-[#94A3B8]">
              {t("webhooksDescription")}
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}

/* ─── 6. Managed Integrations ─── */
function ManagedIntegrations() {
  const t = useTranslations("developers");

  const integrations = [
    { name: "Mixpanel", desc: t("integrationMixpanel"), icon: "📊" },
    { name: "Customer.io", desc: t("integrationCustomerio"), icon: "📧" },
    { name: "Amplitude", desc: t("integrationAmplitude"), icon: "📈" },
    { name: "Segment", desc: t("integrationSegment"), icon: "🔄" },
    { name: "Datadog", desc: t("integrationDatadog"), icon: "🐕" },
    { name: "Slack", desc: t("integrationSlack"), icon: "💬" },
  ];

  return (
    <section className="bg-[#111827] px-4 py-24 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-[1200px]">
        <div className="mb-16 text-center">
          <h2 className="mb-4 text-4xl font-bold text-[#F8FAFC]">{t("integrationsTitle")}</h2>
          <p className="mx-auto max-w-3xl text-xl text-[#E2E8F0]">
            {t("integrationsSubtitle")}
          </p>
        </div>

        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {integrations.map((integration) => (
            <div
              key={integration.name}
              className="rounded-xl border border-[#334155] bg-[#1E293B] p-5 transition-colors hover:border-[#475569]"
            >
              <div className="mb-3 flex items-start justify-between">
                <span className="text-3xl">{integration.icon}</span>
                <Badge variant="blue">Pro</Badge>
              </div>
              <h3 className="mb-2 text-lg font-semibold text-[#F8FAFC]">{integration.name}</h3>
              <p className="text-sm text-[#94A3B8]">{integration.desc}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

/* ─── 7. Build Plugins ─── */
function BuildPlugins() {
  const t = useTranslations("developers");

  const valueProps = [
    { icon: "🔌", title: t("pluginSdk"), desc: t("pluginSdkDesc") },
    { icon: "👥", title: t("pluginAudience"), desc: t("pluginAudienceDesc") },
    { icon: "💵", title: t("pluginRevenue"), desc: t("pluginRevenueDesc") },
  ];

  return (
    <section className="bg-[#0F172A] px-4 py-24 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-[1200px]">
        <div className="mb-16 text-center">
          <h2 className="mb-4 text-4xl font-bold text-[#F8FAFC]">{t("buildPluginsTitle")}</h2>
        </div>

        <div className="mb-12 grid gap-8 md:grid-cols-3">
          {valueProps.map((prop) => (
            <div key={prop.title} className="rounded-xl border border-[#334155] bg-[#1E293B] p-5 text-center">
              <div className="mb-4 text-4xl">{prop.icon}</div>
              <h3 className="mb-2 text-lg font-semibold text-[#F8FAFC]">{prop.title}</h3>
              <p className="text-sm text-[#94A3B8]">{prop.desc}</p>
            </div>
          ))}
        </div>

        <div className="mb-8 grid gap-8 md:grid-cols-2">
          <CodeBlock
            language="json"
            code={`{
  "name": "my-plugin",
  "version": "1.0.0",
  "hooks": ["flag.created", "flag.toggled"],
  "permissions": ["flags:read", "webhooks:write"],
  "ui": {
    "settings": "./settings.tsx",
    "dashboard": "./widget.tsx"
  }
}`}
          />
          <div className="overflow-hidden rounded-xl border border-[#334155] bg-[#1E293B]">
            <div className="p-5">
              <pre className="font-mono text-sm text-[#22D3EE]">
{`$ npx @flagbridge/cli create-plugin

✓ Plugin scaffolded
✓ Dev server running on :3001
✓ Hot reload enabled

→ Edit src/index.ts to get started`}
              </pre>
            </div>
          </div>
        </div>

        <div className="flex flex-wrap justify-center gap-4">
          <button className="rounded-lg bg-[#3B82F6] px-6 py-3 font-medium text-white transition-colors hover:bg-blue-600">
            {t("startBuilding")}
          </button>
          <button className="rounded-lg border border-[#475569] bg-transparent px-6 py-3 font-medium text-[#E2E8F0] transition-colors hover:border-[#64748B]">
            {t("pluginSdkRef")}
          </button>
          <button className="rounded-lg border border-[#475569] bg-transparent px-6 py-3 font-medium text-[#E2E8F0] transition-colors hover:border-[#64748B]">
            {t("publishingGuidelines")}
          </button>
        </div>
      </div>
    </section>
  );
}

/* ─── 8. Why Engineering Teams Switch ─── */
function WhySwitch() {
  const t = useTranslations("developers");

  const reasons = [
    { icon: "🏗️", title: t("whySwitch1"), desc: t("whySwitch1Desc") },
    { icon: "💰", title: t("whySwitch2"), desc: t("whySwitch2Desc") },
    { icon: "🔓", title: t("whySwitch3"), desc: t("whySwitch3Desc") },
    { icon: "🧪", title: t("whySwitch4"), desc: t("whySwitch4Desc") },
  ];

  return (
    <section className="bg-[#111827] px-4 py-24 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-[1200px]">
        <div className="mb-16 text-center">
          <h2 className="mb-4 text-4xl font-bold text-[#F8FAFC]">{t("whySwitchTitle")}</h2>
        </div>

        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
          {reasons.map((reason) => (
            <div key={reason.title} className="rounded-xl border border-[#334155] bg-[#1E293B] p-5">
              <div className="mb-4 text-4xl">{reason.icon}</div>
              <h3 className="mb-2 text-lg font-semibold text-[#F8FAFC]">{reason.title}</h3>
              <p className="text-sm text-[#94A3B8]">{reason.desc}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

/* ─── 9. Community + Footer CTA ─── */
function CommunityCta() {
  const t = useTranslations("developers");

  const communityLinks = [
    { name: "GitHub", desc: t("communityGithub"), icon: "github" },
    { name: "Discord", desc: t("communityDiscord"), icon: "discord" },
    { name: t("communityDocsTitle"), desc: t("communityDocs"), icon: "docs" },
    { name: "Blog", desc: t("communityBlog"), icon: "blog" },
    { name: "YouTube", desc: t("communityYoutube"), icon: "youtube" },
    { name: t("communityRoadmapTitle"), desc: t("communityRoadmap"), icon: "roadmap" },
  ];

  return (
    <section className="bg-[#0F172A] px-4 py-24 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-[1200px]">
        <div className="mb-16 text-center">
          <h2 className="mb-4 text-4xl font-bold text-[#F8FAFC]">{t("communityTitle")}</h2>
        </div>

        <div className="mb-16 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {communityLinks.map((link) => (
            <a
              key={link.name}
              href="#"
              className="flex items-start gap-4 rounded-xl border border-[#334155] bg-[#1E293B] p-5 transition-colors hover:border-[#475569]"
            >
              <div className="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-lg bg-[#3B82F6]/10">
                <svg width={20} height={20} fill="none" viewBox="0 0 24 24" stroke="#3B82F6" strokeWidth={1.5}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 6H5.25A2.25 2.25 0 0 0 3 8.25v10.5A2.25 2.25 0 0 0 5.25 21h10.5A2.25 2.25 0 0 0 18 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" />
                </svg>
              </div>
              <div>
                <h3 className="mb-1 text-lg font-semibold text-[#F8FAFC]">{link.name}</h3>
                <p className="text-sm text-[#94A3B8]">{link.desc}</p>
              </div>
            </a>
          ))}
        </div>

        {/* Final CTA Banner */}
        <div className="rounded-xl border border-[#334155] bg-[#1E293B] p-12 text-center">
          <h2 className="mb-4 text-3xl font-bold text-[#F8FAFC]">{t("ctaTitle")}</h2>
          <p className="mx-auto mb-8 max-w-2xl text-lg text-[#E2E8F0]">
            {t("ctaSubtitle")}
          </p>
          <div className="flex flex-col justify-center gap-4 sm:flex-row">
            <button className="flex items-center justify-center gap-2 rounded-lg bg-[#3B82F6] px-8 py-4 font-medium text-white transition-colors hover:bg-blue-600">
              {t("getStartedFree")}
              <svg width={20} height={20} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 4.5 21 12m0 0-7.5 7.5M21 12H3" />
              </svg>
            </button>
            <button className="rounded-lg border border-[#475569] bg-transparent px-8 py-4 font-medium text-[#E2E8F0] transition-colors hover:border-[#64748B]">
              {t("scheduleDemo")}
            </button>
          </div>
        </div>
      </div>
    </section>
  );
}

/* ─── Page ─── */
export default function DevelopersPage() {
  return (
    <main className="min-h-screen bg-[#0F172A]">
      <Header />
      <DeveloperHero />
      <QuickStart />
      <SdkSection />
      <TestingApi />
      <ApiWebhooks />
      <ManagedIntegrations />
      <BuildPlugins />
      <WhySwitch />
      <CommunityCta />
      <Footer />
    </main>
  );
}
