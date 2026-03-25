import { useTranslations } from "next-intl";
import type { Metadata } from "next";
import { Header } from "@/components/header";
import { Footer } from "@/components/footer";
import { Badge } from "@/components/badge";
import { CategoryFilter } from "@/components/category-filter";

export const metadata: Metadata = {
  title: "Marketplace — FlagBridge",
  description:
    "Extend FlagBridge with plugins for messaging, analytics, queues, and more. Install in one click or build your own and get paid.",
  openGraph: {
    title: "Marketplace — FlagBridge",
    description:
      "Plugins for messaging, analytics, queues, and more. Install in one click. Build your own and get paid.",
    url: "https://flagbridge.io/marketplace",
  },
};

/* ─── 1. Hero + Search ─── */
function MarketplaceHero() {
  const t = useTranslations("marketplace");
  return (
    <section className="bg-[#0F172A] px-4 pt-32 pb-24 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-[1200px]">
        <div className="mx-auto max-w-4xl text-center">
          <h1 className="mb-6 text-5xl font-bold text-[#F8FAFC] sm:text-6xl">
            {t("heroTitle")}
          </h1>
          <p className="mb-10 text-xl text-[#E2E8F0]">
            {t("heroSubtitle")}
          </p>
          <CategoryFilter
            searchPlaceholder={t("searchPlaceholder")}
            statsText={t("statsText")}
          />
        </div>
      </div>
    </section>
  );
}

/* ─── 2. Featured ─── */
function Featured() {
  const t = useTranslations("marketplace");

  const plugins = [
    {
      icon: "💬",
      name: "Slack Alerts",
      author: "FlagBridge",
      verified: true,
      price: "Free" as const,
      rating: 4.9,
      installs: "2.1k",
      category: "Messaging",
      desc: t("featuredSlackDesc"),
    },
    {
      icon: "📊",
      name: "Mixpanel Analytics",
      author: "FlagBridge",
      verified: true,
      price: "Free" as const,
      rating: 4.7,
      installs: "1.4k",
      category: "Analytics",
      desc: t("featuredMixpanelDesc"),
    },
    {
      icon: "🔄",
      name: "Kafka Event Stream",
      author: "FlagBridge",
      verified: true,
      price: "Pro" as const,
      rating: 4.8,
      installs: "890",
      category: "Queues",
      desc: t("featuredKafkaDesc"),
    },
  ];

  return (
    <section className="bg-[#111827] px-4 py-24 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-[1200px]">
        <h2 className="mb-10 text-3xl font-bold text-[#F8FAFC]">{t("featuredTitle")}</h2>

        <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
          {plugins.map((plugin) => (
            <div
              key={plugin.name}
              className="group cursor-pointer rounded-xl border border-[#334155] bg-[#1E293B] p-5 transition-colors hover:border-[#475569]"
            >
              <div className="mb-4 flex items-start justify-between">
                <div className="flex h-12 w-12 items-center justify-center rounded-lg border border-[#334155] bg-[#0F172A] text-2xl">
                  {plugin.icon}
                </div>
                <Badge variant={plugin.price === "Free" ? "green" : "blue"}>{plugin.price}</Badge>
              </div>

              <h3 className="mb-1 text-lg font-semibold text-[#F8FAFC]">{plugin.name}</h3>
              <div className="mb-3 flex items-center gap-2">
                <span className="text-sm text-[#94A3B8]">by {plugin.author}</span>
                {plugin.verified && (
                  <svg width={16} height={16} fill="none" viewBox="0 0 24 24" stroke="#3B82F6" strokeWidth={2}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="m4.5 12.75 6 6 9-13.5" />
                  </svg>
                )}
              </div>

              <p className="mb-4 min-h-[40px] text-sm text-[#94A3B8]">{plugin.desc}</p>

              <div className="mb-4 flex items-center justify-between">
                <div className="flex items-center gap-3 text-sm">
                  <div className="flex items-center gap-1 text-[#F97316]">
                    <svg width={14} height={14} viewBox="0 0 24 24" fill="#F97316">
                      <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
                    </svg>
                    <span className="text-[#E2E8F0]">{plugin.rating}</span>
                  </div>
                  <div className="flex items-center gap-1 text-[#94A3B8]">
                    <svg width={14} height={14} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                      <path strokeLinecap="round" strokeLinejoin="round" d="M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5M16.5 12 12 16.5m0 0L7.5 12m4.5 4.5V3" />
                    </svg>
                    <span>{plugin.installs}</span>
                  </div>
                </div>
                <Badge>{plugin.category}</Badge>
              </div>

              <button className="w-full rounded-lg bg-[#3B82F6] px-4 py-2 font-medium text-white transition-colors hover:bg-blue-600">
                {t("install")}
              </button>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

/* ─── 3. Browse by Category ─── */
function BrowseCategories() {
  const t = useTranslations("marketplace");

  const categories = [
    { name: t("catMessaging"), count: 8, icon: "💬" },
    { name: t("catAnalytics"), count: 7, icon: "📊" },
    { name: t("catQueues"), count: 6, icon: "🔄" },
    { name: t("catProjectMgmt"), count: 5, icon: "📋" },
    { name: t("catCommunication"), count: 4, icon: "💡" },
    { name: t("catSecurity"), count: 3, icon: "🔒" },
    { name: t("catCiCd"), count: 4, icon: "🚀" },
    { name: t("catUiExtensions"), count: 5, icon: "🎨" },
    { name: t("catDataExport"), count: 4, icon: "📁" },
    { name: t("catAutomation"), count: 3, icon: "⚡" },
  ];

  return (
    <section className="bg-[#0F172A] px-4 py-24 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-[1200px]">
        <h2 className="mb-10 text-3xl font-bold text-[#F8FAFC]">{t("browseByCategoryTitle")}</h2>

        <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-5">
          {categories.map((category) => (
            <button
              key={category.name}
              className="group rounded-xl border border-[#334155] bg-[#1E293B] p-5 text-left transition-colors hover:border-[#475569]"
            >
              <div className="mb-3 text-3xl">{category.icon}</div>
              <h3 className="mb-1 text-sm font-semibold text-[#F8FAFC] transition-colors group-hover:text-[#3B82F6]">
                {category.name}
              </h3>
              <p className="text-xs text-[#64748B]">{category.count} plugins</p>
            </button>
          ))}
        </div>
      </div>
    </section>
  );
}

/* ─── 4. For Engineering Teams ─── */
function ForEngineering() {
  const t = useTranslations("marketplace");

  const features = [
    { icon: "🔄", title: t("engFeature1"), desc: t("engFeature1Desc") },
    { icon: "🔍", title: t("engFeature2"), desc: t("engFeature2Desc") },
    { icon: "📈", title: t("engFeature3"), desc: t("engFeature3Desc") },
  ];

  return (
    <section className="bg-[#111827] px-4 py-24 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-[1200px]">
        <h2 className="mb-10 text-3xl font-bold text-[#F8FAFC]">{t("forEngineeringTitle")}</h2>

        <div className="grid gap-8 md:grid-cols-5">
          <div className="md:col-span-2">
            <p className="text-lg leading-relaxed text-[#E2E8F0]">
              {t("forEngineeringDesc")}
            </p>
          </div>

          <div className="space-y-4 md:col-span-3">
            {features.map((feature) => (
              <div
                key={feature.title}
                className="rounded-xl border border-[#334155] bg-[#1E293B] p-5 transition-colors hover:border-[#475569]"
              >
                <div className="flex items-start gap-4">
                  <div className="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-lg bg-[#3B82F6]/10 text-xl">
                    {feature.icon}
                  </div>
                  <div>
                    <h3 className="mb-2 text-lg font-semibold text-[#F8FAFC]">{feature.title}</h3>
                    <p className="text-sm text-[#94A3B8]">{feature.desc}</p>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="mt-8">
          <a href="#engineering" className="inline-flex items-center gap-2 font-medium text-[#3B82F6] transition-colors hover:text-blue-400">
            {t("seeAllEngineering")}
            <svg width={18} height={18} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 4.5 21 12m0 0-7.5 7.5M21 12H3" />
            </svg>
          </a>
        </div>
      </div>
    </section>
  );
}

/* ─── 5. For Product & Business Teams ─── */
function ForProduct() {
  const t = useTranslations("marketplace");

  const features = [
    { icon: "🎯", title: t("prodFeature1"), desc: t("prodFeature1Desc") },
    { icon: "📧", title: t("prodFeature2"), desc: t("prodFeature2Desc") },
    { icon: "📊", title: t("prodFeature3"), desc: t("prodFeature3Desc") },
  ];

  return (
    <section className="bg-[#0F172A] px-4 py-24 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-[1200px]">
        <h2 className="mb-10 text-3xl font-bold text-[#F8FAFC]">{t("forProductTitle")}</h2>

        <div className="grid gap-8 md:grid-cols-5">
          <div className="md:col-span-2">
            <p className="text-lg leading-relaxed text-[#E2E8F0]">
              {t("forProductDesc")}
            </p>
          </div>

          <div className="space-y-4 md:col-span-3">
            {features.map((feature) => (
              <div
                key={feature.title}
                className="rounded-xl border border-[#334155] bg-[#1E293B] p-5 transition-colors hover:border-[#475569]"
              >
                <div className="flex items-start gap-4">
                  <div className="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-lg bg-[#3B82F6]/10 text-xl">
                    {feature.icon}
                  </div>
                  <div>
                    <h3 className="mb-2 text-lg font-semibold text-[#F8FAFC]">{feature.title}</h3>
                    <p className="text-sm text-[#94A3B8]">{feature.desc}</p>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="mt-8">
          <a href="#product" className="inline-flex items-center gap-2 font-medium text-[#3B82F6] transition-colors hover:text-blue-400">
            {t("seeAllProduct")}
            <svg width={18} height={18} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 4.5 21 12m0 0-7.5 7.5M21 12H3" />
            </svg>
          </a>
        </div>
      </div>
    </section>
  );
}

/* ─── 6. For Plugin Developers ─── */
function ForPluginDevelopers() {
  const t = useTranslations("marketplace");
  return (
    <section className="bg-[#111827] px-4 py-24 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-[1200px]">
        <div className="mb-12 text-center">
          <h2 className="mb-4 text-4xl font-bold text-[#F8FAFC]">{t("forDevTitle")}</h2>
          <p className="mx-auto max-w-3xl text-xl text-[#E2E8F0]">
            {t("forDevSubtitle")}
          </p>
        </div>

        {/* Stats */}
        <div className="mb-12 grid grid-cols-1 gap-8 md:grid-cols-3">
          <div className="text-center">
            <div className="mb-2 text-5xl font-bold text-[#3B82F6]">80%</div>
            <div className="text-[#94A3B8]">{t("revenueShare")}</div>
          </div>
          <div className="text-center">
            <div className="mb-2 text-5xl font-bold text-[#3B82F6]">$0</div>
            <div className="text-[#94A3B8]">{t("costToPublish")}</div>
          </div>
          <div className="text-center">
            <div className="mb-2 text-5xl font-bold text-[#3B82F6]">1</div>
            <div className="text-[#94A3B8]">{t("commandToPublish")}</div>
          </div>
        </div>

        {/* Steps */}
        <div className="mb-12 flex flex-col items-center justify-center gap-4 md:flex-row">
          {[t("stepScaffold"), t("stepBuild"), t("stepPublish"), t("stepEarn")].map((step, idx) => (
            <div key={step} className="flex items-center gap-4">
              <div className="rounded-xl border border-[#334155] bg-[#1E293B] px-6 py-4">
                <div className="font-semibold text-[#F8FAFC]">{idx + 1}. {step}</div>
              </div>
              {idx < 3 && (
                <svg width={24} height={24} fill="none" viewBox="0 0 24 24" stroke="#475569" strokeWidth={2} className="hidden rotate-0 md:block">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 4.5 21 12m0 0-7.5 7.5M21 12H3" />
                </svg>
              )}
            </div>
          ))}
        </div>

        {/* CTAs */}
        <div className="flex flex-col justify-center gap-4 sm:flex-row">
          <button className="rounded-lg bg-[#3B82F6] px-6 py-3 font-medium text-white transition-colors hover:bg-blue-600">
            {t("startBuildingPlugins")}
          </button>
          <button className="rounded-lg border border-[#475569] bg-transparent px-6 py-3 font-medium text-[#E2E8F0] transition-colors hover:border-[#64748B]">
            {t("pluginSdkDocs")}
          </button>
          <button className="rounded-lg border border-[#475569] bg-transparent px-6 py-3 font-medium text-[#E2E8F0] transition-colors hover:border-[#64748B]">
            {t("publishingGuidelines")}
          </button>
        </div>
      </div>
    </section>
  );
}

/* ─── 7. Social Proof ─── */
function SocialProof() {
  const t = useTranslations("marketplace");

  const testimonials = [
    {
      quote: t("testimonial1"),
      author: "Sarah Chen",
      company: "Engineering Lead, Acme Corp",
    },
    {
      quote: t("testimonial2"),
      author: "Marcus Rodriguez",
      company: "Staff Engineer, TechFlow",
    },
    {
      quote: t("testimonial3"),
      author: "Emily Thompson",
      company: "Product Manager, DataHub",
    },
  ];

  return (
    <section className="bg-[#0F172A] px-4 py-24 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-[1200px]">
        {/* Stats Bar */}
        <div className="mb-16 text-center">
          <p className="text-sm text-[#64748B]">{t("socialStatsBar")}</p>
        </div>

        {/* Testimonials */}
        <div className="mb-16 grid grid-cols-1 gap-6 md:grid-cols-3">
          {testimonials.map((testimonial) => (
            <div
              key={testimonial.author}
              className="rounded-xl border border-[#334155] bg-[#1E293B] p-5"
            >
              <p className="mb-4 italic text-[#E2E8F0]">&quot;{testimonial.quote}&quot;</p>
              <div>
                <div className="text-sm font-medium text-[#F8FAFC]">{testimonial.author}</div>
                <div className="text-xs text-[#64748B]">{testimonial.company}</div>
              </div>
            </div>
          ))}
        </div>

        {/* Logo Bar */}
        <div>
          <p className="mb-6 text-center text-sm text-[#64748B]">{t("trustedBy")}</p>
          <div className="flex flex-wrap justify-center gap-8">
            {[1, 2, 3, 4, 5, 6].map((i) => (
              <div
                key={i}
                className="h-12 w-32 rounded-lg border border-[#334155] bg-[#1E293B]"
              />
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}

/* ─── 8. Footer CTA ─── */
function FooterCta() {
  const t = useTranslations("marketplace");
  return (
    <section className="bg-[#111827] px-4 py-24 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-[1200px]">
        <div className="grid gap-12 md:grid-cols-2">
          {/* Left Side */}
          <div className="text-center md:border-r md:border-[#334155] md:pr-12 md:text-left">
            <h3 className="mb-6 text-3xl font-bold text-[#F8FAFC]">{t("ctaLeftTitle")}</h3>
            <div className="flex flex-col gap-4 sm:flex-row">
              <button className="flex items-center justify-center gap-2 rounded-lg bg-[#3B82F6] px-6 py-3 font-medium text-white transition-colors hover:bg-blue-600">
                {t("browseMarketplace")}
                <svg width={18} height={18} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 4.5 21 12m0 0-7.5 7.5M21 12H3" />
                </svg>
              </button>
              <button className="rounded-lg border border-[#475569] bg-transparent px-6 py-3 font-medium text-[#E2E8F0] transition-colors hover:border-[#64748B]">
                {t("getStartedFree")}
              </button>
            </div>
          </div>

          {/* Right Side */}
          <div className="text-center md:pl-12 md:text-left">
            <h3 className="mb-6 text-3xl font-bold text-[#F8FAFC]">{t("ctaRightTitle")}</h3>
            <div className="flex flex-col gap-4 sm:flex-row">
              <button className="flex items-center justify-center gap-2 rounded-lg bg-[#3B82F6] px-6 py-3 font-medium text-white transition-colors hover:bg-blue-600">
                {t("createFirstPlugin")}
                <svg width={18} height={18} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="m3.75 13.5 10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75Z" />
                </svg>
              </button>
              <button className="flex items-center justify-center gap-2 rounded-lg border border-[#475569] bg-transparent px-6 py-3 font-medium text-[#E2E8F0] transition-colors hover:border-[#64748B]">
                <svg width={18} height={18} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M8.625 12a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm0 0H8.25m4.125 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm0 0H12m4.125 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm0 0h-.375M21 12c0 4.556-4.03 8.25-9 8.25a9.764 9.764 0 0 1-2.555-.337A5.972 5.972 0 0 1 5.41 20.97a5.969 5.969 0 0 1-.474-.065 4.48 4.48 0 0 0 .978-2.025c.09-.457-.133-.901-.467-1.226C3.93 16.178 3 14.189 3 12c0-4.556 4.03-8.25 9-8.25s9 3.694 9 8.25Z" />
                </svg>
                {t("joinDiscord")}
              </button>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

/* ─── Page ─── */
export default function MarketplacePage() {
  return (
    <main className="min-h-screen bg-[#0F172A]">
      <Header />
      <MarketplaceHero />
      <Featured />
      <BrowseCategories />
      <ForEngineering />
      <ForProduct />
      <ForPluginDevelopers />
      <SocialProof />
      <FooterCta />
      <Footer />
    </main>
  );
}
