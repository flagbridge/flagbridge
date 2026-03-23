import Image from "next/image";
import Link from "next/link";

function NavBar() {
  return (
    <nav className="fixed top-0 z-50 w-full bg-[#0b1326]/80 backdrop-blur-xl">
      <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
        <Link href="/" className="flex items-center gap-3">
          <Image
            src="/brand/logo-square.svg"
            alt="FlagBridge"
            width={36}
            height={36}
          />
          <span className="text-lg font-semibold tracking-tight text-[#dae2fd]">
            FlagBridge
          </span>
        </Link>
        <div className="hidden items-center gap-8 md:flex">
          <a
            href="#features"
            className="text-sm text-[#c2c6d6] transition hover:text-[#adc6ff]"
          >
            Features
          </a>
          <a
            href="#comparison"
            className="text-sm text-[#c2c6d6] transition hover:text-[#adc6ff]"
          >
            Compare
          </a>
          <a
            href="#open-core"
            className="text-sm text-[#c2c6d6] transition hover:text-[#adc6ff]"
          >
            Pricing
          </a>
          <a
            href="https://github.com/flagbridge/flagbridge"
            target="_blank"
            rel="noopener noreferrer"
            className="text-sm text-[#c2c6d6] transition hover:text-[#adc6ff]"
          >
            GitHub
          </a>
        </div>
        <a
          href="https://github.com/flagbridge/flagbridge"
          target="_blank"
          rel="noopener noreferrer"
          className="rounded-md bg-gradient-to-br from-[#adc6ff] to-[#4d8eff] px-4 py-2 text-sm font-medium text-[#002e6a] transition hover:opacity-90"
        >
          Get Started
        </a>
      </div>
    </nav>
  );
}

function Hero() {
  return (
    <section className="relative flex min-h-screen flex-col items-center justify-center overflow-hidden px-6 pt-24">
      {/* Ambient glow */}
      <div className="pointer-events-none absolute left-1/2 top-1/3 -translate-x-1/2 -translate-y-1/2">
        <div className="h-[600px] w-[600px] rounded-full bg-[#3B82F6]/8 blur-[120px]" />
      </div>

      <div className="relative z-10 mx-auto max-w-4xl text-center">
        <div className="mb-8 inline-flex items-center gap-2 rounded-full bg-[#171f33] px-4 py-1.5 text-sm text-[#adc6ff]">
          <span className="inline-block h-2 w-2 rounded-full bg-[#4d8eff]" />
          Open Source &middot; Apache 2.0
        </div>

        <h1 className="text-5xl font-bold leading-tight tracking-tight text-[#dae2fd] md:text-7xl">
          Feature flags with{" "}
          <span className="bg-gradient-to-r from-[#adc6ff] to-[#4d8eff] bg-clip-text text-transparent">
            product intelligence
          </span>
        </h1>

        <p className="mx-auto mt-6 max-w-2xl text-lg leading-relaxed text-[#c2c6d6]">
          FlagBridge goes beyond simple toggles. Connect feature flags to
          product strategy, track impact metrics, automate lifecycle cleanup, and
          extend everything with plugins.
        </p>

        <div className="mt-10 flex flex-col items-center gap-4 sm:flex-row sm:justify-center">
          <a
            href="https://github.com/flagbridge/flagbridge"
            target="_blank"
            rel="noopener noreferrer"
            className="rounded-md bg-gradient-to-br from-[#adc6ff] to-[#4d8eff] px-8 py-3 text-base font-medium text-[#002e6a] transition hover:opacity-90"
          >
            View on GitHub
          </a>
          <a
            href="#features"
            className="rounded-md bg-[#222a3d] px-8 py-3 text-base font-medium text-[#adc6ff] transition hover:bg-[#2d3449]"
          >
            Learn More
          </a>
        </div>

        {/* Quick start terminal */}
        <div className="mx-auto mt-16 max-w-lg overflow-hidden rounded-lg bg-[#060e20]">
          <div className="flex items-center gap-2 bg-[#131b2e] px-4 py-2.5">
            <div className="h-3 w-3 rounded-full bg-[#424754]" />
            <div className="h-3 w-3 rounded-full bg-[#424754]" />
            <div className="h-3 w-3 rounded-full bg-[#424754]" />
            <span className="ml-2 text-xs text-[#8c909f]">terminal</span>
          </div>
          <div className="px-5 py-4 text-left font-mono text-sm text-[#c2c6d6]">
            <p>
              <span className="text-[#4d8eff]">$</span> git clone
              github.com/flagbridge/flagbridge
            </p>
            <p>
              <span className="text-[#4d8eff]">$</span> cd flagbridge
            </p>
            <p>
              <span className="text-[#4d8eff]">$</span> docker compose up -d
            </p>
            <p className="mt-2 text-[#8c909f]">
              # API at :8080 &middot; Admin UI at :3000
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}

const features = [
  {
    title: "Flag Management",
    description:
      "Create boolean, string, number, and JSON flags. Target by user ID, percentage rollout, or custom attributes. Instant kill switch.",
    icon: (
      <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
        <path strokeLinecap="round" strokeLinejoin="round" d="M3 3v1.5M3 21v-6m0 0 2.77-.693a9 9 0 0 1 6.208.682l.108.054a9 9 0 0 0 6.086.71l3.114-.732a48.524 48.524 0 0 1-.005-10.499l-3.11.732a9 9 0 0 1-6.085-.711l-.108-.054a9 9 0 0 0-6.208-.682L3 4.5M3 15V4.5" />
      </svg>
    ),
  },
  {
    title: "Product Context Cards",
    description:
      "Link every flag to a hypothesis, owner, success metrics, and deadline. No more zombie flags without business context.",
    icon: (
      <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
        <path strokeLinecap="round" strokeLinejoin="round" d="M3.75 9.776c.112-.017.227-.026.344-.026h15.812c.117 0 .232.009.344.026m-16.5 0a2.25 2.25 0 0 0-1.883 2.542l.857 6a2.25 2.25 0 0 0 2.227 1.932H19.05a2.25 2.25 0 0 0 2.227-1.932l.857-6a2.25 2.25 0 0 0-1.883-2.542m-16.5 0V6A2.25 2.25 0 0 1 6 3.75h3.879a1.5 1.5 0 0 1 1.06.44l2.122 2.12a1.5 1.5 0 0 0 1.06.44H18A2.25 2.25 0 0 1 20.25 9v.776" />
      </svg>
    ),
  },
  {
    title: "Technical Dashboard",
    description:
      "Adoption rate, error rate by variant, latency impact, stale flag detection, and SDK version tracking. All in real-time.",
    icon: (
      <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
        <path strokeLinecap="round" strokeLinejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 0 1 3 19.875v-6.75ZM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V8.625ZM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V4.125Z" />
      </svg>
    ),
  },
  {
    title: "Lifecycle Automation",
    description:
      'Configurable rules: "if flag is 100% ON for X days, notify owner." Auto-archive expired flags. Technical debt score per project.',
    icon: (
      <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
        <path strokeLinecap="round" strokeLinejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182" />
      </svg>
    ),
  },
  {
    title: "Plugin Ecosystem",
    description:
      "Build, publish, and install plugins. Extend the UI, add API middleware, create custom integrations. Plugin marketplace included.",
    icon: (
      <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
        <path strokeLinecap="round" strokeLinejoin="round" d="M14.25 6.087c0-.355.186-.676.401-.959.221-.29.349-.634.349-1.003 0-1.036-1.007-1.875-2.25-1.875s-2.25.84-2.25 1.875c0 .369.128.713.349 1.003.215.283.401.604.401.959v0a.64.64 0 0 1-.657.643 48.39 48.39 0 0 1-4.163-.3c.186 1.613.293 3.25.315 4.907a.656.656 0 0 1-.658.663v0c-.355 0-.676-.186-.959-.401a1.647 1.647 0 0 0-1.003-.349c-1.036 0-1.875 1.007-1.875 2.25s.84 2.25 1.875 2.25c.369 0 .713-.128 1.003-.349.283-.215.604-.401.959-.401v0c.31 0 .555.26.532.57a48.039 48.039 0 0 1-.642 5.056c1.518.19 3.058.309 4.616.354a.64.64 0 0 0 .657-.643v0c0-.355-.186-.676-.401-.959a1.647 1.647 0 0 1-.349-1.003c0-1.035 1.008-1.875 2.25-1.875 1.243 0 2.25.84 2.25 1.875 0 .369-.128.713-.349 1.003-.215.283-.4.604-.4.959v0c0 .333.277.599.61.58a48.1 48.1 0 0 0 5.427-.63 48.05 48.05 0 0 0 .582-4.717.532.532 0 0 0-.533-.57v0c-.355 0-.676.186-.959.401-.29.221-.634.349-1.003.349-1.035 0-1.875-1.007-1.875-2.25s.84-2.25 1.875-2.25c.37 0 .713.128 1.003.349.283.215.604.401.96.401v0a.656.656 0 0 0 .658-.663 48.422 48.422 0 0 0-.37-5.36c-1.886.342-3.81.574-5.766.689a.578.578 0 0 1-.61-.58v0Z" />
      </svg>
    ),
  },
  {
    title: "OpenFeature Compatible",
    description:
      "Official provider for the CNCF OpenFeature standard. Zero vendor lock-in. Switch providers without changing application code.",
    icon: (
      <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
        <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 16.875h3.375m0 0h3.375m-3.375 0V13.5m0 3.375v3.375M6 10.5h2.25a2.25 2.25 0 0 0 2.25-2.25V6a2.25 2.25 0 0 0-2.25-2.25H6A2.25 2.25 0 0 0 3.75 6v2.25A2.25 2.25 0 0 0 6 10.5Zm0 9.75h2.25A2.25 2.25 0 0 0 10.5 18v-2.25a2.25 2.25 0 0 0-2.25-2.25H6a2.25 2.25 0 0 0-2.25 2.25V18A2.25 2.25 0 0 0 6 20.25Zm9.75-9.75H18a2.25 2.25 0 0 0 2.25-2.25V6A2.25 2.25 0 0 0 18 3.75h-2.25A2.25 2.25 0 0 0 13.5 6v2.25a2.25 2.25 0 0 0 2.25 2.25Z" />
      </svg>
    ),
  },
];

function Features() {
  return (
    <section id="features" className="px-6 py-32">
      <div className="mx-auto max-w-6xl">
        <h2 className="text-center text-3xl font-bold tracking-tight text-[#dae2fd] md:text-4xl">
          Everything you need to ship with confidence
        </h2>
        <p className="mx-auto mt-4 max-w-2xl text-center text-[#c2c6d6]">
          From simple boolean toggles to advanced product intelligence — all
          open source.
        </p>

        <div className="mt-16 grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {features.map((f) => (
            <div
              key={f.title}
              className="rounded-lg bg-[#171f33] p-6 transition hover:bg-[#1c2438]"
            >
              <div className="mb-4 inline-flex h-10 w-10 items-center justify-center rounded-md bg-[#222a3d] text-[#adc6ff]">
                {f.icon}
              </div>
              <h3 className="text-lg font-semibold text-[#dae2fd]">
                {f.title}
              </h3>
              <p className="mt-2 text-sm leading-relaxed text-[#c2c6d6]">
                {f.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

function Comparison() {
  const rows = [
    { feature: "Open Source Core", fb: true, unleash: true, ld: false, flagsmith: true },
    { feature: "Self-hosted", fb: true, unleash: true, ld: false, flagsmith: true },
    { feature: "Product Context Cards", fb: true, unleash: false, ld: false, flagsmith: false },
    { feature: "Plugin Ecosystem", fb: true, unleash: false, ld: false, flagsmith: false },
    { feature: "Lifecycle Automation", fb: true, unleash: false, ld: false, flagsmith: false },
    { feature: "OpenFeature Provider", fb: true, unleash: true, ld: true, flagsmith: true },
    { feature: "Messaging Integrations", fb: true, unleash: false, ld: false, flagsmith: false },
    { feature: "Event Analytics Bridge", fb: true, unleash: false, ld: false, flagsmith: false },
    { feature: "Starting Price", fb: "Free", unleash: "$80/mo", ld: "$120/mo+", flagsmith: "$45/mo" },
  ];

  return (
    <section id="comparison" className="px-6 py-32">
      <div className="mx-auto max-w-4xl">
        <h2 className="text-center text-3xl font-bold tracking-tight text-[#dae2fd] md:text-4xl">
          How FlagBridge compares
        </h2>
        <p className="mx-auto mt-4 max-w-2xl text-center text-[#c2c6d6]">
          The only platform combining feature flags, product intelligence, and a
          plugin ecosystem.
        </p>

        <div className="mt-12 overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="text-left text-[#8c909f]">
                <th className="pb-4 pr-4 font-medium">Feature</th>
                <th className="pb-4 pr-4 font-medium text-[#adc6ff]">FlagBridge</th>
                <th className="pb-4 pr-4 font-medium">Unleash</th>
                <th className="pb-4 pr-4 font-medium">LaunchDarkly</th>
                <th className="pb-4 font-medium">Flagsmith</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((row) => (
                <tr
                  key={row.feature}
                  className="border-t border-[#424754]/15"
                >
                  <td className="py-3 pr-4 text-[#c2c6d6]">{row.feature}</td>
                  <td className="py-3 pr-4">
                    {typeof row.fb === "boolean" ? (
                      row.fb ? (
                        <span className="text-[#4d8eff]">Yes</span>
                      ) : (
                        <span className="text-[#8c909f]">No</span>
                      )
                    ) : (
                      <span className="font-medium text-[#adc6ff]">{row.fb}</span>
                    )}
                  </td>
                  <td className="py-3 pr-4">
                    {typeof row.unleash === "boolean" ? (
                      row.unleash ? (
                        <span className="text-[#c2c6d6]">Yes</span>
                      ) : (
                        <span className="text-[#8c909f]">No</span>
                      )
                    ) : (
                      <span className="text-[#c2c6d6]">{row.unleash}</span>
                    )}
                  </td>
                  <td className="py-3 pr-4">
                    {typeof row.ld === "boolean" ? (
                      row.ld ? (
                        <span className="text-[#c2c6d6]">Yes</span>
                      ) : (
                        <span className="text-[#8c909f]">No</span>
                      )
                    ) : (
                      <span className="text-[#c2c6d6]">{row.ld}</span>
                    )}
                  </td>
                  <td className="py-3">
                    {typeof row.flagsmith === "boolean" ? (
                      row.flagsmith ? (
                        <span className="text-[#c2c6d6]">Yes</span>
                      ) : (
                        <span className="text-[#8c909f]">No</span>
                      )
                    ) : (
                      <span className="text-[#c2c6d6]">{row.flagsmith}</span>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </section>
  );
}

function OpenCore() {
  return (
    <section id="open-core" className="px-6 py-32">
      <div className="mx-auto max-w-5xl">
        <h2 className="text-center text-3xl font-bold tracking-tight text-[#dae2fd] md:text-4xl">
          Open core. Your data, your rules.
        </h2>
        <p className="mx-auto mt-4 max-w-2xl text-center text-[#c2c6d6]">
          The Community Edition is genuinely useful on its own — not a crippled
          version that forces upgrades.
        </p>

        <div className="mt-12 grid gap-6 md:grid-cols-3">
          {/* CE */}
          <div className="rounded-lg bg-[#171f33] p-8">
            <h3 className="text-lg font-semibold text-[#dae2fd]">Community</h3>
            <p className="mt-1 text-sm text-[#8c909f]">Open Source</p>
            <p className="mt-4 text-3xl font-bold text-[#dae2fd]">Free</p>
            <p className="text-sm text-[#8c909f]">Apache 2.0</p>
            <ul className="mt-6 space-y-2 text-sm text-[#c2c6d6]">
              <li>Unlimited flags & environments</li>
              <li>Unlimited projects</li>
              <li>Targeting rules & kill switch</li>
              <li>All SDKs + OpenFeature</li>
              <li>Basic product context</li>
              <li>Plugin runtime</li>
              <li>Webhooks + Slack plugin</li>
              <li>Full REST API + Admin UI</li>
            </ul>
          </div>

          {/* Pro */}
          <div className="rounded-lg bg-[#171f33] p-8 ring-1 ring-[#4d8eff]/30">
            <h3 className="text-lg font-semibold text-[#adc6ff]">Pro</h3>
            <p className="mt-1 text-sm text-[#8c909f]">Self-hosted plugin</p>
            <p className="mt-4 text-3xl font-bold text-[#dae2fd]">Soon</p>
            <p className="text-sm text-[#8c909f]">Everything in Community, plus</p>
            <ul className="mt-6 space-y-2 text-sm text-[#c2c6d6]">
              <li>Advanced product context cards</li>
              <li>Technical dashboard & analytics</li>
              <li>Lifecycle & cleanup automation</li>
              <li>A/B testing (statistical significance)</li>
              <li>Scheduled flag rollouts</li>
              <li>Messaging integrations</li>
              <li>Event analytics bridge</li>
              <li>Technical queue connectors</li>
            </ul>
          </div>

          {/* Enterprise */}
          <div className="rounded-lg bg-[#171f33] p-8">
            <h3 className="text-lg font-semibold text-[#dae2fd]">Enterprise</h3>
            <p className="mt-1 text-sm text-[#8c909f]">Managed SaaS</p>
            <p className="mt-4 text-3xl font-bold text-[#dae2fd]">Custom</p>
            <p className="text-sm text-[#8c909f]">Everything in Pro, plus</p>
            <ul className="mt-6 space-y-2 text-sm text-[#c2c6d6]">
              <li>SSO / SAML</li>
              <li>Custom approval workflows</li>
              <li>Audit log export</li>
              <li>Dedicated infrastructure</li>
              <li>SLA & priority support</li>
              <li>Custom integrations</li>
            </ul>
          </div>
        </div>
      </div>
    </section>
  );
}

function TechStack() {
  return (
    <section className="px-6 py-32">
      <div className="mx-auto max-w-4xl text-center">
        <h2 className="text-3xl font-bold tracking-tight text-[#dae2fd] md:text-4xl">
          Built with a modern stack
        </h2>
        <p className="mx-auto mt-4 max-w-2xl text-[#c2c6d6]">
          Designed for performance, reliability, and developer experience.
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
          ].map((t) => (
            <div key={t.name} className="rounded-lg bg-[#171f33] p-5">
              <p className="text-base font-semibold text-[#dae2fd]">{t.name}</p>
              <p className="mt-1 text-xs text-[#8c909f]">{t.desc}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

function CTA() {
  return (
    <section className="px-6 py-32">
      <div className="mx-auto max-w-3xl rounded-2xl bg-gradient-to-br from-[#171f33] to-[#131b2e] p-12 text-center">
        <h2 className="text-3xl font-bold tracking-tight text-[#dae2fd] md:text-4xl">
          Ready to ship with confidence?
        </h2>
        <p className="mx-auto mt-4 max-w-lg text-[#c2c6d6]">
          Start with the open-source Community Edition. Self-host in minutes
          with Docker.
        </p>
        <div className="mt-8 flex flex-col items-center gap-4 sm:flex-row sm:justify-center">
          <a
            href="https://github.com/flagbridge/flagbridge"
            target="_blank"
            rel="noopener noreferrer"
            className="rounded-md bg-gradient-to-br from-[#adc6ff] to-[#4d8eff] px-8 py-3 text-base font-medium text-[#002e6a] transition hover:opacity-90"
          >
            Star on GitHub
          </a>
          <a
            href="https://github.com/flagbridge/flagbridge/blob/main/CONTRIBUTING.md"
            target="_blank"
            rel="noopener noreferrer"
            className="rounded-md bg-[#222a3d] px-8 py-3 text-base font-medium text-[#adc6ff] transition hover:bg-[#2d3449]"
          >
            Contribute
          </a>
        </div>
      </div>
    </section>
  );
}

function Footer() {
  return (
    <footer className="border-t border-[#424754]/15 px-6 py-12">
      <div className="mx-auto flex max-w-6xl flex-col items-center gap-6 md:flex-row md:justify-between">
        <div className="flex items-center gap-3">
          <Image
            src="/brand/logo-square.svg"
            alt="FlagBridge"
            width={24}
            height={24}
          />
          <span className="text-sm text-[#8c909f]">
            FlagBridge &copy; {new Date().getFullYear()}. Apache 2.0.
          </span>
        </div>
        <div className="flex gap-6 text-sm text-[#8c909f]">
          <a
            href="https://github.com/flagbridge"
            target="_blank"
            rel="noopener noreferrer"
            className="transition hover:text-[#c2c6d6]"
          >
            GitHub
          </a>
          <a
            href="https://github.com/flagbridge/flagbridge/discussions"
            target="_blank"
            rel="noopener noreferrer"
            className="transition hover:text-[#c2c6d6]"
          >
            Community
          </a>
          <a
            href="https://github.com/flagbridge/flagbridge/blob/main/LICENSE"
            target="_blank"
            rel="noopener noreferrer"
            className="transition hover:text-[#c2c6d6]"
          >
            License
          </a>
        </div>
      </div>
    </footer>
  );
}

export default function Home() {
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
