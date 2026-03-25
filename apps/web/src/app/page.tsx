import { useTranslations } from "next-intl";
import type { Metadata } from "next";
import { Header } from "@/components/header";
import { Footer } from "@/components/footer";

export const metadata: Metadata = {
  title: "FlagBridge — Feature Flags with Product Intelligence",
  description:
    "Open-source feature flag platform with product intelligence, plugin ecosystem, and OpenFeature support. Self-hosted or cloud.",
  openGraph: {
    title: "FlagBridge — Feature Flags with Product Intelligence",
    description:
      "Open-source feature flag management with product intelligence. Ship features safely and understand their impact.",
    url: "https://flagbridge.io",
  },
};

/* ─── Hero ─── */
function Hero() {
  const t = useTranslations("hero");
  return (
    <section className="bg-[#0F172A] px-4 pt-32 pb-24 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-[1200px]">
        <div className="mx-auto max-w-4xl text-center">
          {/* Badge */}
          <div className="mb-8 inline-flex items-center gap-2 rounded-full border border-[#334155] bg-[#1E293B]/50 px-4 py-2">
            <svg width={16} height={16} viewBox="0 0 24 24" fill="#3B82F6">
              <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
            </svg>
            <span className="text-sm text-[#94A3B8]">{t("badge")}</span>
          </div>

          <h1 className="mb-6 text-5xl font-bold leading-tight text-white sm:text-6xl lg:text-7xl">
            {t("titlePart1")}
            <span className="block text-[#3B82F6]">{t("titleHighlight")}</span>
          </h1>

          <p className="mx-auto mb-10 max-w-2xl text-xl text-[#94A3B8]">
            {t("description")}
          </p>

          {/* CTA Buttons */}
          <div className="mb-12 flex flex-col items-center justify-center gap-4 sm:flex-row">
            <a
              href="https://github.com/flagbridge/flagbridge"
              target="_blank"
              rel="noopener noreferrer"
              className="flex w-full items-center justify-center gap-2 rounded-lg bg-[#3B82F6] px-8 py-4 font-medium text-white transition-colors hover:bg-blue-600 sm:w-auto"
            >
              {t("cta")}
              <svg width={20} height={20} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 4.5 21 12m0 0-7.5 7.5M21 12H3" />
              </svg>
            </a>
            <a
              href="https://github.com/flagbridge/flagbridge"
              target="_blank"
              rel="noopener noreferrer"
              className="flex w-full items-center justify-center gap-2 rounded-lg border border-[#334155] bg-[#1E293B] px-8 py-4 font-medium text-white transition-colors hover:bg-[#253548] sm:w-auto"
            >
              <svg width={20} height={20} viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" />
              </svg>
              {t("learnMore")}
            </a>
          </div>

          {/* Stats */}
          <div className="mb-16 flex flex-wrap justify-center gap-8 text-sm text-[#94A3B8]">
            <div className="flex items-center gap-2">
              <div className="h-2 w-2 rounded-full bg-[#3B82F6]" />
              <span>100% Open Source</span>
            </div>
            <div className="flex items-center gap-2">
              <div className="h-2 w-2 rounded-full bg-[#3B82F6]" />
              <span>Self-hosted or Cloud</span>
            </div>
            <div className="flex items-center gap-2">
              <div className="h-2 w-2 rounded-full bg-[#3B82F6]" />
              <span>Enterprise Ready</span>
            </div>
          </div>
        </div>

        {/* Dashboard Mockup */}
        <div className="mx-auto max-w-5xl">
          <div className="relative">
            <div className="absolute inset-0 rounded-full bg-[#3B82F6]/20 blur-3xl" />
            <div className="relative overflow-hidden rounded-xl border border-[#334155] bg-[#1E293B] shadow-2xl">
              <div className="flex items-center gap-2 border-b border-[#334155] bg-[#0F172A] px-4 py-3">
                <div className="flex gap-2">
                  <div className="h-3 w-3 rounded-full bg-[#EF4444]" />
                  <div className="h-3 w-3 rounded-full bg-[#F97316]" />
                  <div className="h-3 w-3 rounded-full bg-[#22C55E]" />
                </div>
                <div className="flex-1 text-center text-sm text-[#64748B]">feature-flags-dashboard</div>
              </div>
              <div className="p-8">
                <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                  <div className="rounded-lg border border-[#334155] bg-[#0F172A]/50 p-6">
                    <div className="mb-4 flex items-center justify-between">
                      <h3 className="font-medium text-white">New Checkout Flow</h3>
                      <div className="rounded bg-[#22C55E]/20 px-2 py-1 text-xs text-[#22C55E]">Active</div>
                    </div>
                    <div className="space-y-2 text-sm text-[#94A3B8]">
                      <div className="flex justify-between">
                        <span>Rollout</span>
                        <span className="text-white">45%</span>
                      </div>
                      <div className="h-2 w-full rounded-full bg-[#334155]">
                        <div className="h-2 rounded-full bg-[#3B82F6]" style={{ width: "45%" }} />
                      </div>
                    </div>
                  </div>
                  <div className="rounded-lg border border-[#334155] bg-[#0F172A]/50 p-6">
                    <div className="mb-4 flex items-center justify-between">
                      <h3 className="font-medium text-white">Dark Mode Toggle</h3>
                      <div className="rounded bg-[#3B82F6]/20 px-2 py-1 text-xs text-[#60A5FA]">Testing</div>
                    </div>
                    <div className="space-y-2 text-sm text-[#94A3B8]">
                      <div className="flex justify-between">
                        <span>Rollout</span>
                        <span className="text-white">10%</span>
                      </div>
                      <div className="h-2 w-full rounded-full bg-[#334155]">
                        <div className="h-2 rounded-full bg-[#3B82F6]" style={{ width: "10%" }} />
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

/* ─── Features ─── */
function Features() {
  const t = useTranslations("features");

  const features = [
    { icon: "🚩", titleKey: "featureToggles", descKey: "featureTogglesDesc" },
    { icon: "👥", titleKey: "userTargeting", descKey: "userTargetingDesc" },
    { icon: "📊", titleKey: "abTesting", descKey: "abTestingDesc" },
    { icon: "🛡️", titleKey: "safeDeployments", descKey: "safeDeploymentsDesc" },
    { icon: "💻", titleKey: "sdkSupport", descKey: "sdkSupportDesc" },
    { icon: "⚡", titleKey: "realTimeUpdates", descKey: "realTimeUpdatesDesc" },
    { icon: "🌍", titleKey: "multiEnvironment", descKey: "multiEnvironmentDesc" },
    { icon: "📋", titleKey: "auditLogs", descKey: "auditLogsDesc" },
  ];

  return (
    <section id="features" className="bg-[#111827] px-4 py-24 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-[1200px]">
        <div className="mx-auto mb-16 max-w-3xl text-center">
          <h2 className="mb-4 text-4xl font-bold text-white sm:text-5xl">{t("title")}</h2>
          <p className="text-xl text-[#94A3B8]">{t("subtitle")}</p>
        </div>

        <div className="grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-4">
          {features.map((f) => (
            <div
              key={f.titleKey}
              className="rounded-xl border border-[#334155] bg-[#1E293B]/50 p-6 transition-colors hover:border-[#3B82F6]/50"
            >
              <div className="mb-4 flex h-12 w-12 items-center justify-center rounded-lg bg-[#3B82F6]/10 text-2xl">
                {f.icon}
              </div>
              <h3 className="mb-2 text-lg font-semibold text-white">{t(f.titleKey)}</h3>
              <p className="text-sm text-[#94A3B8]">{t(f.descKey)}</p>
            </div>
          ))}
        </div>

        {/* Code Example Section */}
        <div className="mt-20 grid items-center gap-8 md:grid-cols-2">
          <div>
            <h3 className="mb-4 text-3xl font-bold text-white">
              {t("codeExampleTitle")}
            </h3>
            <p className="mb-6 text-[#94A3B8]">
              {t("codeExampleSubtitle")}
            </p>
            <ul className="space-y-3">
              {["Type-safe SDKs", "Automatic updates", "Offline support", "Zero dependencies"].map((item) => (
                <li key={item} className="flex items-center gap-3 text-[#E2E8F0]">
                  <div className="flex h-5 w-5 items-center justify-center rounded-full bg-[#3B82F6]/20">
                    <div className="h-2 w-2 rounded-full bg-[#3B82F6]" />
                  </div>
                  {item}
                </li>
              ))}
            </ul>
          </div>
          <div className="overflow-hidden rounded-xl border border-[#334155] bg-[#0F172A]">
            <div className="border-b border-[#334155] bg-[#1E293B] px-4 py-2">
              <div className="text-xs text-[#94A3B8]">example.tsx</div>
            </div>
            <div className="p-6 font-mono text-sm">
              <pre className="text-[#E2E8F0]">
{`import { useFlag } from '@flagbridge/react';

function Checkout() {
  const newFlow = useFlag('new-checkout');

  return (
    <>
      {newFlow ? (
        <NewCheckoutFlow />
      ) : (
        <LegacyCheckout />
      )}
    </>
  );
}`}
              </pre>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

/* ─── Pricing ─── */
function Pricing() {
  const t = useTranslations("openCore");

  const plans = [
    {
      name: t("community"),
      subtitle: t("communitySubtitle"),
      price: t("communityPrice"),
      license: t("communityLicense"),
      features: Array.from({ length: 6 }, (_, i) => t(`communityFeatures.${i}`)),
      cta: "Deploy Now",
      highlighted: false,
    },
    {
      name: t("pro"),
      subtitle: t("proSubtitle"),
      price: "$29",
      period: "/month",
      features: Array.from({ length: 6 }, (_, i) => t(`proFeatures.${i}`)),
      cta: "Start Free Trial",
      highlighted: true,
    },
    {
      name: t("enterprise"),
      subtitle: t("enterpriseSubtitle"),
      price: t("enterprisePrice"),
      features: Array.from({ length: 6 }, (_, i) => t(`enterpriseFeatures.${i}`)),
      cta: "Contact Sales",
      highlighted: false,
    },
  ];

  return (
    <section id="pricing" className="bg-[#0F172A] px-4 py-24 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-[1200px]">
        <div className="mx-auto mb-16 max-w-3xl text-center">
          <h2 className="mb-4 text-4xl font-bold text-white sm:text-5xl">{t("title")}</h2>
          <p className="text-xl text-[#94A3B8]">{t("subtitle")}</p>
        </div>

        <div className="mx-auto grid max-w-6xl grid-cols-1 gap-8 md:grid-cols-3">
          {plans.map((plan) => (
            <div
              key={plan.name}
              className={`relative rounded-xl p-8 ${
                plan.highlighted
                  ? "bg-[#3B82F6] text-white"
                  : "border border-[#334155] bg-[#1E293B]/50 text-white"
              }`}
            >
              {plan.highlighted && (
                <div className="absolute -top-4 left-1/2 -translate-x-1/2 rounded-full border border-[#3B82F6] bg-[#0F172A] px-4 py-1 text-sm font-medium text-[#3B82F6]">
                  Most Popular
                </div>
              )}
              <div className="mb-6">
                <h3 className="mb-2 text-xl font-semibold">{plan.name}</h3>
                <div className="mb-2">
                  <span className="text-4xl font-bold">{plan.price}</span>
                  {plan.period && (
                    <span className={plan.highlighted ? "text-blue-100" : "text-[#94A3B8]"}>
                      {plan.period}
                    </span>
                  )}
                </div>
                <p className={plan.highlighted ? "text-blue-100" : "text-[#94A3B8]"}>
                  {plan.subtitle}
                </p>
              </div>

              <button
                className={`mb-8 w-full rounded-lg px-6 py-3 font-medium transition-colors ${
                  plan.highlighted
                    ? "bg-white text-[#3B82F6] hover:bg-blue-50"
                    : "bg-[#3B82F6] text-white hover:bg-blue-600"
                }`}
              >
                {plan.cta}
              </button>

              <ul className="space-y-4">
                {plan.features.map((feature) => (
                  <li key={feature} className="flex items-start gap-3">
                    <svg
                      width={20}
                      height={20}
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                      strokeWidth={2}
                      className={`mt-0.5 flex-shrink-0 ${plan.highlighted ? "text-white" : "text-[#3B82F6]"}`}
                    >
                      <path strokeLinecap="round" strokeLinejoin="round" d="m4.5 12.75 6 6 9-13.5" />
                    </svg>
                    <span className={plan.highlighted ? "text-blue-50" : "text-[#E2E8F0]"}>
                      {feature}
                    </span>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>

        <div className="mt-16 text-center">
          <p className="text-[#94A3B8]">
            All plans include unlimited API requests.{" "}
            <a href="#contact" className="text-[#3B82F6] hover:underline">Contact us</a>
          </p>
        </div>
      </div>
    </section>
  );
}

/* ─── Comparison ─── */
const comparisonKeys = [
  "openSourceCore",
  "selfHosted",
  "productContextCards",
  "pluginEcosystem",
  "lifecycleAutomation",
  "openfeatureProvider",
  "messagingIntegrations",
  "eventAnalyticsBridge",
] as const;

const comparisonData: Record<string, { fb: boolean; unleash: boolean; ld: boolean; flagsmith: boolean }> = {
  openSourceCore:        { fb: true,  unleash: true,  ld: false, flagsmith: true  },
  selfHosted:            { fb: true,  unleash: true,  ld: false, flagsmith: true  },
  productContextCards:   { fb: true,  unleash: false, ld: false, flagsmith: false },
  pluginEcosystem:       { fb: true,  unleash: false, ld: false, flagsmith: false },
  lifecycleAutomation:   { fb: true,  unleash: false, ld: false, flagsmith: false },
  openfeatureProvider:   { fb: true,  unleash: true,  ld: true,  flagsmith: true  },
  messagingIntegrations: { fb: true,  unleash: false, ld: false, flagsmith: false },
  eventAnalyticsBridge:  { fb: true,  unleash: false, ld: false, flagsmith: false },
};

function BoolCell({ value, highlight = false }: { value: boolean; highlight?: boolean }) {
  if (value) {
    return <span className={highlight ? "text-[#22C55E]" : "text-[#22C55E]"}>✓</span>;
  }
  return <span className="text-[#64748B]">✗</span>;
}

function Comparison() {
  const t = useTranslations("comparison");
  return (
    <section id="comparison" className="bg-[#111827] px-4 py-24 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-[1200px]">
        <div className="mx-auto mb-16 max-w-3xl text-center">
          <h2 className="mb-4 text-4xl font-bold text-white">{t("title")}</h2>
          <p className="text-xl text-[#94A3B8]">{t("subtitle")}</p>
        </div>

        <div className="overflow-hidden rounded-xl border border-[#334155] bg-[#1E293B]">
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="bg-[#111827] text-left text-xs uppercase tracking-wide text-[#94A3B8]">
                  <th className="px-6 py-4 font-medium">{t("feature")}</th>
                  <th className="px-6 py-4 font-medium text-[#3B82F6]">FlagBridge</th>
                  <th className="px-6 py-4 font-medium">Unleash</th>
                  <th className="px-6 py-4 font-medium">LaunchDarkly</th>
                  <th className="px-6 py-4 font-medium">Flagsmith</th>
                </tr>
              </thead>
              <tbody>
                {comparisonKeys.map((key) => {
                  const row = comparisonData[key];
                  return (
                    <tr key={key} className="border-b border-[#1E293B] transition-colors hover:bg-[#253548]">
                      <td className="px-6 py-4 text-[#E2E8F0]">{t(key)}</td>
                      <td className="px-6 py-4"><BoolCell value={row.fb} highlight /></td>
                      <td className="px-6 py-4"><BoolCell value={row.unleash} /></td>
                      <td className="px-6 py-4"><BoolCell value={row.ld} /></td>
                      <td className="px-6 py-4"><BoolCell value={row.flagsmith} /></td>
                    </tr>
                  );
                })}
                <tr className="border-b border-[#1E293B] transition-colors hover:bg-[#253548]">
                  <td className="px-6 py-4 text-[#E2E8F0]">{t("startingPrice")}</td>
                  <td className="px-6 py-4 font-bold text-[#22C55E]">{t("free")}</td>
                  <td className="px-6 py-4 text-[#94A3B8]">$80/mo</td>
                  <td className="px-6 py-4 text-[#94A3B8]">$120/mo+</td>
                  <td className="px-6 py-4 text-[#94A3B8]">$45/mo</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </section>
  );
}

/* ─── TechStack ─── */
function TechStack() {
  const t = useTranslations("techStack");

  const techs = [
    { name: "Go", desc: "API Server" },
    { name: "Next.js 15", desc: "Admin UI" },
    { name: "PostgreSQL", desc: "Database" },
    { name: "Redis", desc: "Cache & Pub/Sub" },
    { name: "TypeScript", desc: "SDKs" },
    { name: "Docker", desc: "Deployment" },
    { name: "Kubernetes", desc: "Helm Charts" },
    { name: "OpenFeature", desc: "Standard" },
  ];

  return (
    <section className="bg-[#0F172A] px-4 py-24 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-[1200px] text-center">
        <h2 className="mb-4 text-4xl font-bold text-white">{t("title")}</h2>
        <p className="mx-auto mb-12 max-w-2xl text-[#94A3B8]">{t("subtitle")}</p>

        <div className="grid grid-cols-2 gap-6 md:grid-cols-4">
          {techs.map((item) => (
            <div
              key={item.name}
              className="rounded-xl border border-[#334155] bg-[#1E293B] p-6 text-center"
            >
              <p className="text-base font-semibold text-[#F8FAFC]">{item.name}</p>
              <p className="mt-1 text-sm text-[#94A3B8]">{item.desc}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

/* ─── CTA ─── */
function CtaSection() {
  const t = useTranslations("cta");
  return (
    <section className="bg-[#111827] px-4 py-24 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-[1200px]">
        <div className="rounded-xl border border-[#334155] bg-[#1E293B] p-12 text-center">
          <h2 className="mb-4 text-3xl font-bold text-[#F8FAFC]">{t("title")}</h2>
          <p className="mx-auto mb-8 max-w-2xl text-lg text-[#E2E8F0]">{t("subtitle")}</p>
          <div className="flex flex-col justify-center gap-4 sm:flex-row">
            <a
              href="https://github.com/flagbridge/flagbridge"
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center justify-center gap-2 rounded-lg bg-[#3B82F6] px-8 py-4 font-medium text-white transition-colors hover:bg-blue-600"
            >
              {t("star")}
              <svg width={20} height={20} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 4.5 21 12m0 0-7.5 7.5M21 12H3" />
              </svg>
            </a>
            <a
              href="https://github.com/flagbridge/flagbridge/blob/main/CONTRIBUTING.md"
              target="_blank"
              rel="noopener noreferrer"
              className="rounded-lg border border-[#475569] bg-transparent px-8 py-4 font-medium text-[#E2E8F0] transition-colors hover:border-[#64748B]"
            >
              {t("contribute")}
            </a>
          </div>
        </div>
      </div>
    </section>
  );
}

/* ─── Page ─── */
export default function Home() {
  return (
    <main className="min-h-screen bg-[#0F172A]">
      <Header />
      <Hero />
      <Features />
      <Pricing />
      <Comparison />
      <TechStack />
      <CtaSection />
      <Footer />
    </main>
  );
}
