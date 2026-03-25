import Image from "next/image";
import Link from "next/link";
import { useTranslations } from "next-intl";
import { LocaleSwitcher } from "@/components/locale-switcher";

function GithubIcon({ size = 20 }: { size?: number }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="currentColor">
      <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" />
    </svg>
  );
}

function TwitterIcon({ size = 20 }: { size?: number }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="currentColor">
      <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z" />
    </svg>
  );
}

function LinkedInIcon({ size = 20 }: { size?: number }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="currentColor">
      <path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433c-1.144 0-2.063-.926-2.063-2.065 0-1.138.92-2.063 2.063-2.063 1.14 0 2.064.925 2.064 2.063 0 1.139-.925 2.065-2.064 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z" />
    </svg>
  );
}

export function Footer() {
  const t = useTranslations("footer");

  return (
    <footer className="border-t border-[#1E293B] bg-[#0F172A]">
      <div className="mx-auto max-w-[1200px] px-4 py-12 sm:px-6 lg:px-8">
        <div className="mb-8 grid grid-cols-2 gap-8 md:grid-cols-4">
          {/* Brand */}
          <div>
            <h4 className="mb-4 font-semibold text-[#F8FAFC]">{t("brand")}</h4>
            <ul className="space-y-2">
              <li><a href="#about" className="text-sm text-[#94A3B8] transition-colors hover:text-white">{t("about")}</a></li>
              <li><a href="#blog" className="text-sm text-[#94A3B8] transition-colors hover:text-white">{t("blog")}</a></li>
              <li><a href="#careers" className="text-sm text-[#94A3B8] transition-colors hover:text-white">{t("careers")}</a></li>
              <li><a href="#press" className="text-sm text-[#94A3B8] transition-colors hover:text-white">{t("pressKit")}</a></li>
            </ul>
          </div>

          {/* Product */}
          <div>
            <h4 className="mb-4 font-semibold text-[#F8FAFC]">{t("product")}</h4>
            <ul className="space-y-2">
              <li><a href="#features" className="text-sm text-[#94A3B8] transition-colors hover:text-white">{t("features")}</a></li>
              <li><a href="#pricing" className="text-sm text-[#94A3B8] transition-colors hover:text-white">{t("pricing")}</a></li>
              <li><a href="#security" className="text-sm text-[#94A3B8] transition-colors hover:text-white">{t("security")}</a></li>
              <li><a href="#roadmap" className="text-sm text-[#94A3B8] transition-colors hover:text-white">{t("roadmap")}</a></li>
            </ul>
          </div>

          {/* Developers */}
          <div>
            <h4 className="mb-4 font-semibold text-[#F8FAFC]">{t("developersTitle")}</h4>
            <ul className="space-y-2">
              <li><a href="https://docs.flagbridge.io" className="text-sm text-[#94A3B8] transition-colors hover:text-white">{t("documentation")}</a></li>
              <li><a href="#api" className="text-sm text-[#94A3B8] transition-colors hover:text-white">{t("apiReference")}</a></li>
              <li><Link href="/developers" className="text-sm text-[#94A3B8] transition-colors hover:text-white">{t("sdks")}</Link></li>
              <li><a href="#guides" className="text-sm text-[#94A3B8] transition-colors hover:text-white">{t("guides")}</a></li>
            </ul>
          </div>

          {/* Company */}
          <div>
            <h4 className="mb-4 font-semibold text-[#F8FAFC]">{t("company")}</h4>
            <ul className="space-y-2">
              <li><a href="#contact" className="text-sm text-[#94A3B8] transition-colors hover:text-white">{t("contact")}</a></li>
              <li><a href="#privacy" className="text-sm text-[#94A3B8] transition-colors hover:text-white">{t("privacy")}</a></li>
              <li><a href="#terms" className="text-sm text-[#94A3B8] transition-colors hover:text-white">{t("terms")}</a></li>
              <li><a href="https://github.com/flagbridge/flagbridge/blob/main/LICENSE" target="_blank" rel="noopener noreferrer" className="text-sm text-[#94A3B8] transition-colors hover:text-white">{t("license")}</a></li>
            </ul>
          </div>
        </div>

        {/* Bottom Section */}
        <div className="flex flex-col items-center gap-4 border-t border-[#1E293B] pt-8 md:flex-row md:justify-between">
          <div className="flex items-center gap-3">
            <Image src="/brand/logo-textual-light.svg" alt="FlagBridge" width={130} height={100} />
            <span className="text-sm text-[#94A3B8]">
              &copy; {new Date().getFullYear()} {t("openSourceMit")}
            </span>
          </div>

          <div className="flex items-center gap-4">
            <LocaleSwitcher />
            <a
              href="https://github.com/flagbridge/flagbridge"
              target="_blank"
              rel="noopener noreferrer"
              className="text-[#94A3B8] transition-colors hover:text-white"
            >
              <GithubIcon size={20} />
            </a>
            <a
              href="https://twitter.com/flagbridge"
              target="_blank"
              rel="noopener noreferrer"
              className="text-[#94A3B8] transition-colors hover:text-white"
            >
              <TwitterIcon size={20} />
            </a>
            <a
              href="https://linkedin.com/company/flagbridge"
              target="_blank"
              rel="noopener noreferrer"
              className="text-[#94A3B8] transition-colors hover:text-white"
            >
              <LinkedInIcon size={20} />
            </a>
          </div>
        </div>
      </div>
    </footer>
  );
}
