import Image from "next/image";
import Link from "next/link";
import { useTranslations } from "next-intl";
import { setRequestLocale } from "next-intl/server";
import { LocaleSwitcher } from "@/components/locale-switcher";

function NavBar() {
  const t = useTranslations("nav");
  return (
    <nav className="fixed top-0 z-50 w-full bg-[#0b1326]/80 backdrop-blur-xl">
      <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
        <Link href="/" className="flex items-center gap-3">
          <Image
            src="/brand/logo-textual-light.svg"
            alt="FlagBridge"
            width={160}
            height={44}
            className="h-8 w-auto"
          />
        </Link>
        <div className="hidden items-center gap-8 md:flex">
          <a href="#features" className="text-sm text-[#c2c6d6] transition hover:text-[#adc6ff]">
            {t("features")}
          </a>
          <a href="#comparison" className="text-sm text-[#c2c6d6] transition hover:text-[#adc6ff]">
            {t("compare")}
          </a>
          <a href="#open-core" className="text-sm text-[#c2c6d6] transition hover:text-[#adc6ff]">
            {t("pricing")}
          </a>
          <a
            href="https://github.com/flagbridge/flagbridge"
            target="_blank"
            rel="noopener noreferrer"
            className="text-sm text-[#c2c6d6] transition hover:text-[#adc6ff]"
          >
            {t("github")}
          </a>
        </div>
        <div className="flex items-center gap-3">
          <LocaleSwitcher />
          <a
            href="https://github.com/flagbridge/flagbridge"
            target="_blank"
            rel="noopener noreferrer"
            className="rounded-md bg-gradient-to-br from-[#adc6ff] to-[#4d8eff] px-4 py-2 text-sm font-medium text-[#002e6a] transition hover:opacity-90"
          >
            {t("getStarted")}
          </a>
        </div>
      </div>
    </nav>
  );
}

function Hero() {
  const t = useTranslations("hero");
  return (
    <section className="relative flex min-h-screen flex-col items-center justify-center overflow-hidden px-6 pt-24">
      <div className="pointer-events-none absolute left-1/2 top-1/3 -translate-x-1/2 -translate-y-1/2">
        <div className="h-[600px] w-[600px] rounded-full bg-[#3B82F6]/8 blur-[120px]" />
      </div>

      <div className="relative z-10 mx-auto max-w-4xl text-center">
        <div className="mb-8 inline-flex items-center gap-2 rounded-full bg-[#171f33] px-4 py-1.5 text-sm text-[#adc6ff]">
          <span className="inline-block h-2 w-2 rounded-full bg-[#4d8eff]" />
          {t("badge")}
        </div>

        <h1 className="text-5xl font-bold leading-tight tracking-tight text-[#dae2fd] md:text-7xl">
          {t("titlePart1")}
          <span className="bg-gradient-to-r from-[#adc6ff] to-[#4d8eff] bg-clip-text text-transparent">
            {t("titleHighlight")}
          </span>
        </h1>

        <p className="mx-auto mt-6 max-w-2xl text-lg leading-relaxed text-[#c2c6d6]">
          {t("description")}
        </p>

        <div className="mt-10 flex flex-col items-center gap-4 sm:flex-row sm:justify-center">
          <a
            href="https://github.com/flagbridge/flagbridge"
            target="_blank"
            rel="noopener noreferrer"
            className="rounded-md bg-gradient-to-br from-[#adc6ff] to-[#4d8eff] px-8 py-3 text-base font-medium text-[#002e6a] transition hover:opacity-90"
          >
            {t("cta")}
          </a>
          <a
            href="#features"
            className="rounded-md bg-[#222a3d] px-8 py-3 text-base font-medium text-[#adc6ff] transition hover:bg-[#2d3449]"
          >
            {t("learnMore")}
          </a>
        </div>

        <div className="mx-auto mt-16 max-w-lg overflow-hidden rounded-lg bg-[#060e20]">
          <div className="flex items-center gap-2 bg-[#131b2e] px-4 py-2.5">
            <div className="h-3 w-3 rounded-full bg-[#424754]" />
            <div className="h-3 w-3 rounded-full bg-[#424754]" />
            <div className="h-3 w-3 rounded-full bg-[#424754]" />
            <span className="ml-2 text-xs text-[#8c909f]">terminal</span>
          </div>
          <div className="px-5 py-4 text-left font-mono text-sm text-[#c2c6d6]">
            <p><span className="text-[#4d8eff]">$</span> git clone github.com/flagbridge/flagbridge</p>
            <p><span className="text-[#4d8eff]">$</span> cd flagbridge</p>
            <p><span className="text-[#4d8eff]">$</span> docker compose up -d</p>
            <p className="mt-2 text-[#8c909f]">{t("terminalComment")}</p>
          </div>
        </div>
      </div>
    </section>
  );
}

const featureIcons = [
  <svg key="flag" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}><path strokeLinecap="round" strokeLinejoin="round" d="M3 3v1.5M3 21v-6m0 0 2.77-.693a9 9 0 0 1 6.208.682l.108.054a9 9 0 0 0 6.086.71l3.114-.732a48.524 48.524 0 0 1-.005-10.499l-3.11.732a9 9 0 0 1-6.085-.711l-.108-.054a9 9 0 0 0-6.208-.682L3 4.5M3 15V4.5" /></svg>,
  <svg key="card" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}><path strokeLinecap="round" strokeLinejoin="round" d="M3.75 9.776c.112-.017.227-.026.344-.026h15.812c.117 0 .232.009.344.026m-16.5 0a2.25 2.25 0 0 0-1.883 2.542l.857 6a2.25 2.25 0 0 0 2.227 1.932H19.05a2.25 2.25 0 0 0 2.227-1.932l.857-6a2.25 2.25 0 0 0-1.883-2.542m-16.5 0V6A2.25 2.25 0 0 1 6 3.75h3.879a1.5 1.5 0 0 1 1.06.44l2.122 2.12a1.5 1.5 0 0 0 1.06.44H18A2.25 2.25 0 0 1 20.25 9v.776" /></svg>,
  <svg key="chart" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}><path strokeLinecap="round" strokeLinejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 0 1 3 19.875v-6.75ZM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V8.625ZM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V4.125Z" /></svg>,
  <svg key="cycle" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}><path strokeLinecap="round" strokeLinejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182" /></svg>,
  <svg key="plug" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}><path strokeLinecap="round" strokeLinejoin="round" d="M14.25 6.087c0-.355.186-.676.401-.959.221-.29.349-.634.349-1.003 0-1.036-1.007-1.875-2.25-1.875s-2.25.84-2.25 1.875c0 .369.128.713.349 1.003.215.283.401.604.401.959v0a.64.64 0 0 1-.657.643 48.39 48.39 0 0 1-4.163-.3c.186 1.613.293 3.25.315 4.907a.656.656 0 0 1-.658.663v0c-.355 0-.676-.186-.959-.401a1.647 1.647 0 0 0-1.003-.349c-1.036 0-1.875 1.007-1.875 2.25s.84 2.25 1.875 2.25c.369 0 .713-.128 1.003-.349.283-.215.604-.401.959-.401v0c.31 0 .555.26.532.57a48.039 48.039 0 0 1-.642 5.056c1.518.19 3.058.309 4.616.354a.64.64 0 0 0 .657-.643v0c0-.355-.186-.676-.401-.959a1.647 1.647 0 0 1-.349-1.003c0-1.035 1.008-1.875 2.25-1.875 1.243 0 2.25.84 2.25 1.875 0 .369-.128.713-.349 1.003-.215.283-.4.604-.4.959v0c0 .333.277.599.61.58a48.1 48.1 0 0 0 5.427-.63 48.05 48.05 0 0 0 .582-4.717.532.532 0 0 0-.533-.57v0c-.355 0-.676.186-.959.401-.29.221-.634.349-1.003.349-1.035 0-1.875-1.007-1.875-2.25s.84-2.25 1.875-2.25c.37 0 .713.128 1.003.349.283.215.604.401.96.401v0a.656.656 0 0 0 .658-.663 48.422 48.422 0 0 0-.37-5.36c-1.886.342-3.81.574-5.766.689a.578.578 0 0 1-.61-.58v0Z" /></svg>,
  <svg key="of" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}><path strokeLinecap="round" strokeLinejoin="round" d="M13.5 16.875h3.375m0 0h3.375m-3.375 0V13.5m0 3.375v3.375M6 10.5h2.25a2.25 2.25 0 0 0 2.25-2.25V6a2.25 2.25 0 0 0-2.25-2.25H6A2.25 2.25 0 0 0 3.75 6v2.25A2.25 2.25 0 0 0 6 10.5Zm0 9.75h2.25A2.25 2.25 0 0 0 10.5 18v-2.25a2.25 2.25 0 0 0-2.25-2.25H6a2.25 2.25 0 0 0-2.25 2.25V18A2.25 2.25 0 0 0 6 20.25Zm9.75-9.75H18a2.25 2.25 0 0 0 2.25-2.25V6A2.25 2.25 0 0 0 18 3.75h-2.25A2.25 2.25 0 0 0 13.5 6v2.25a2.25 2.25 0 0 0 2.25 2.25Z" /></svg>,
];

const featureKeys = [
  { title: "flagManagement", desc: "flagManagementDesc" },
  { title: "productContext", desc: "productContextDesc" },
  { title: "techDashboard", desc: "techDashboardDesc" },
  { title: "lifecycle", desc: "lifecycleDesc" },
  { title: "plugins", desc: "pluginsDesc" },
  { title: "openfeature", desc: "openfeatureDesc" },
] as const;

function Features() {
  const t = useTranslations("features");
  return (
    <section id="features" className="px-6 py-32">
      <div className="mx-auto max-w-6xl">
        <h2 className="text-center text-3xl font-bold tracking-tight text-[#dae2fd] md:text-4xl">
          {t("title")}
        </h2>
        <p className="mx-auto mt-4 max-w-2xl text-center text-[#c2c6d6]">
          {t("subtitle")}
        </p>

        <div className="mt-16 grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {featureKeys.map((f, i) => (
            <div key={f.title} className="rounded-lg bg-[#171f33] p-6 transition hover:bg-[#1c2438]">
              <div className="mb-4 inline-flex h-10 w-10 items-center justify-center rounded-md bg-[#222a3d] text-[#adc6ff]">
                {featureIcons[i]}
              </div>
              <h3 className="text-lg font-semibold text-[#dae2fd]">{t(f.title)}</h3>
              <p className="mt-2 text-sm leading-relaxed text-[#c2c6d6]">{t(f.desc)}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

function BoolCell({ value }: { value: boolean }) {
  return <span className="text-base">{value ? "\u2705" : "\u274C"}</span>;
}

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
  openSourceCore:       { fb: true,  unleash: true,  ld: false, flagsmith: true  },
  selfHosted:           { fb: true,  unleash: true,  ld: false, flagsmith: true  },
  productContextCards:  { fb: true,  unleash: false, ld: false, flagsmith: false },
  pluginEcosystem:      { fb: true,  unleash: false, ld: false, flagsmith: false },
  lifecycleAutomation:  { fb: true,  unleash: false, ld: false, flagsmith: false },
  openfeatureProvider:  { fb: true,  unleash: true,  ld: true,  flagsmith: true  },
  messagingIntegrations:{ fb: true,  unleash: false, ld: false, flagsmith: false },
  eventAnalyticsBridge: { fb: true,  unleash: false, ld: false, flagsmith: false },
};

function Comparison() {
  const t = useTranslations("comparison");
  return (
    <section id="comparison" className="px-6 py-32">
      <div className="mx-auto max-w-4xl">
        <h2 className="text-center text-3xl font-bold tracking-tight text-[#dae2fd] md:text-4xl">
          {t("title")}
        </h2>
        <p className="mx-auto mt-4 max-w-2xl text-center text-[#c2c6d6]">
          {t("subtitle")}
        </p>

        <div className="mt-12 overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="text-left text-[#8c909f]">
                <th className="pb-4 pr-4 font-medium">{t("feature")}</th>
                <th className="pb-4 pr-4 font-medium text-[#adc6ff]">FlagBridge</th>
                <th className="pb-4 pr-4 font-medium">Unleash</th>
                <th className="pb-4 pr-4 font-medium">LaunchDarkly</th>
                <th className="pb-4 font-medium">Flagsmith</th>
              </tr>
            </thead>
            <tbody>
              {comparisonKeys.map((key) => {
                const row = comparisonData[key];
                return (
                  <tr key={key} className="border-t border-[#424754]/15">
                    <td className="py-3 pr-4 text-[#c2c6d6]">{t(key)}</td>
                    <td className="py-3 pr-4"><BoolCell value={row.fb} /></td>
                    <td className="py-3 pr-4"><BoolCell value={row.unleash} /></td>
                    <td className="py-3 pr-4"><BoolCell value={row.ld} /></td>
                    <td className="py-3"><BoolCell value={row.flagsmith} /></td>
                  </tr>
                );
              })}
              <tr className="border-t border-[#424754]/15">
                <td className="py-3 pr-4 text-[#c2c6d6]">{t("startingPrice")}</td>
                <td className="py-3 pr-4 font-medium text-[#adc6ff]">{t("free")}</td>
                <td className="py-3 pr-4 text-[#c2c6d6]">$80/mo</td>
                <td className="py-3 pr-4 text-[#c2c6d6]">$120/mo+</td>
                <td className="py-3 text-[#c2c6d6]">$45/mo</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </section>
  );
}

function OpenCore() {
  const t = useTranslations("openCore");
  const communityFeatures = Array.from({ length: 8 }, (_, i) => t(`communityFeatures.${i}`));
  const proFeatures = Array.from({ length: 8 }, (_, i) => t(`proFeatures.${i}`));
  const enterpriseFeatures = Array.from({ length: 6 }, (_, i) => t(`enterpriseFeatures.${i}`));

  return (
    <section id="open-core" className="px-6 py-32">
      <div className="mx-auto max-w-5xl">
        <h2 className="text-center text-3xl font-bold tracking-tight text-[#dae2fd] md:text-4xl">
          {t("title")}
        </h2>
        <p className="mx-auto mt-4 max-w-2xl text-center text-[#c2c6d6]">
          {t("subtitle")}
        </p>

        <div className="mt-12 grid gap-6 md:grid-cols-3">
          <div className="rounded-lg bg-[#171f33] p-8">
            <h3 className="text-lg font-semibold text-[#dae2fd]">{t("community")}</h3>
            <p className="mt-1 text-sm text-[#8c909f]">{t("communitySubtitle")}</p>
            <p className="mt-4 text-3xl font-bold text-[#dae2fd]">{t("communityPrice")}</p>
            <p className="text-sm text-[#8c909f]">{t("communityLicense")}</p>
            <ul className="mt-6 space-y-2 text-sm text-[#c2c6d6]">
              {communityFeatures.map((f) => <li key={f}>{f}</li>)}
            </ul>
          </div>

          <div className="rounded-lg bg-[#171f33] p-8 ring-1 ring-[#4d8eff]/30">
            <h3 className="text-lg font-semibold text-[#adc6ff]">{t("pro")}</h3>
            <p className="mt-1 text-sm text-[#8c909f]">{t("proSubtitle")}</p>
            <p className="mt-4 text-3xl font-bold text-[#dae2fd]">{t("proPrice")}</p>
            <p className="text-sm text-[#8c909f]">{t("proIncludes")}</p>
            <ul className="mt-6 space-y-2 text-sm text-[#c2c6d6]">
              {proFeatures.map((f) => <li key={f}>{f}</li>)}
            </ul>
          </div>

          <div className="rounded-lg bg-[#171f33] p-8">
            <h3 className="text-lg font-semibold text-[#dae2fd]">{t("enterprise")}</h3>
            <p className="mt-1 text-sm text-[#8c909f]">{t("enterpriseSubtitle")}</p>
            <p className="mt-4 text-3xl font-bold text-[#dae2fd]">{t("enterprisePrice")}</p>
            <p className="text-sm text-[#8c909f]">{t("enterpriseIncludes")}</p>
            <ul className="mt-6 space-y-2 text-sm text-[#c2c6d6]">
              {enterpriseFeatures.map((f) => <li key={f}>{f}</li>)}
            </ul>
          </div>
        </div>
      </div>
    </section>
  );
}

function TechStack() {
  const t = useTranslations("techStack");
  return (
    <section className="px-6 py-32">
      <div className="mx-auto max-w-4xl text-center">
        <h2 className="text-3xl font-bold tracking-tight text-[#dae2fd] md:text-4xl">
          {t("title")}
        </h2>
        <p className="mx-auto mt-4 max-w-2xl text-[#c2c6d6]">
          {t("subtitle")}
        </p>

        <div className="mt-12 grid grid-cols-2 gap-6 md:grid-cols-4">
          {[
            { name: "Go", desc: "API Server" },
            { name: "Next.js 15", desc: "Admin UI" },
            { name: "PostgreSQL", desc: "Database" },
            { name: "Redis", desc: "Cache & Pub/Sub" },
            { name: "TypeScript", desc: "SDKs" },
            { name: "Docker", desc: "Deployment" },
            { name: "Kubernetes", desc: "Helm Charts" },
            { name: "OpenFeature", desc: "Standard" },
          ].map((item) => (
            <div key={item.name} className="rounded-lg bg-[#171f33] p-5">
              <p className="text-base font-semibold text-[#dae2fd]">{item.name}</p>
              <p className="mt-1 text-xs text-[#8c909f]">{item.desc}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

function CTA() {
  const t = useTranslations("cta");
  return (
    <section className="px-6 py-32">
      <div className="mx-auto max-w-3xl rounded-2xl bg-gradient-to-br from-[#171f33] to-[#131b2e] p-12 text-center">
        <h2 className="text-3xl font-bold tracking-tight text-[#dae2fd] md:text-4xl">
          {t("title")}
        </h2>
        <p className="mx-auto mt-4 max-w-lg text-[#c2c6d6]">
          {t("subtitle")}
        </p>
        <div className="mt-8 flex flex-col items-center gap-4 sm:flex-row sm:justify-center">
          <a
            href="https://github.com/flagbridge/flagbridge"
            target="_blank"
            rel="noopener noreferrer"
            className="rounded-md bg-gradient-to-br from-[#adc6ff] to-[#4d8eff] px-8 py-3 text-base font-medium text-[#002e6a] transition hover:opacity-90"
          >
            {t("star")}
          </a>
          <a
            href="https://github.com/flagbridge/flagbridge/blob/main/CONTRIBUTING.md"
            target="_blank"
            rel="noopener noreferrer"
            className="rounded-md bg-[#222a3d] px-8 py-3 text-base font-medium text-[#adc6ff] transition hover:bg-[#2d3449]"
          >
            {t("contribute")}
          </a>
        </div>
      </div>
    </section>
  );
}

function Footer() {
  const t = useTranslations("footer");
  return (
    <footer className="border-t border-[#424754]/15 px-6 py-12">
      <div className="mx-auto flex max-w-6xl flex-col items-center gap-6 md:flex-row md:justify-between">
        <div className="flex items-center gap-3">
          <Image src="/brand/logo-square.svg" alt="FlagBridge" width={24} height={24} />
          <span className="text-sm text-[#8c909f]">
            FlagBridge &copy; {new Date().getFullYear()}. Apache 2.0.
          </span>
        </div>
        <div className="flex gap-6 text-sm text-[#8c909f]">
          <a href="https://github.com/flagbridge" target="_blank" rel="noopener noreferrer" className="transition hover:text-[#c2c6d6]">
            GitHub
          </a>
          <a href="https://github.com/flagbridge/flagbridge/discussions" target="_blank" rel="noopener noreferrer" className="transition hover:text-[#c2c6d6]">
            {t("community")}
          </a>
          <a href="https://github.com/flagbridge/flagbridge/blob/main/LICENSE" target="_blank" rel="noopener noreferrer" className="transition hover:text-[#c2c6d6]">
            {t("license")}
          </a>
        </div>
      </div>
    </footer>
  );
}

export default async function Home({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  setRequestLocale(locale);

  return (
    <main className="min-h-screen bg-[#0b1326]">
      <NavBar />
      <Hero />
      <Features />
      <Comparison />
      <OpenCore />
      <TechStack />
      <CTA />
      <Footer />
    </main>
  );
}
