"use client";

import { useState } from "react";
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

function MenuIcon() {
  return (
    <svg width={24} height={24} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
      <path strokeLinecap="round" strokeLinejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5" />
    </svg>
  );
}

function XIcon() {
  return (
    <svg width={24} height={24} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
      <path strokeLinecap="round" strokeLinejoin="round" d="M6 18 18 6M6 6l12 12" />
    </svg>
  );
}

export function Header() {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const t = useTranslations("nav");

  return (
    <header className="fixed top-0 left-0 right-0 z-50 h-16 border-b border-[#1E293B] bg-[#0F172A]/80 backdrop-blur">
      <nav className="mx-auto flex h-full max-w-[1200px] items-center justify-between px-4 sm:px-6 lg:px-8">
        {/* Logo */}
        <Link href="/" className="flex items-center gap-2">
          <Image
            src="/brand/logo-textual-light.svg"
            alt="FlagBridge"
            width={130}
            height={100}
          />
        </Link>

        {/* Desktop Navigation */}
        <div className="hidden items-center gap-8 md:flex">
          <Link href="/developers" className="text-sm text-[#94A3B8] transition-colors hover:text-white">
            {t("developers")}
          </Link>
          <Link href="/marketplace" className="text-sm text-[#94A3B8] transition-colors hover:text-white">
            {t("marketplace")}
          </Link>
          <a href="#pricing" className="text-sm text-[#94A3B8] transition-colors hover:text-white">
            {t("pricing")}
          </a>
          <a href="https://docs.flagbridge.io" className="text-sm text-[#94A3B8] transition-colors hover:text-white">
            {t("docs")}
          </a>
        </div>

        {/* Right side */}
        <div className="hidden items-center gap-4 md:flex">
          <LocaleSwitcher />
          <a
            href="https://github.com/flagbridge/flagbridge"
            target="_blank"
            rel="noopener noreferrer"
            className="text-[#94A3B8] transition-colors hover:text-white"
          >
            <GithubIcon size={20} />
          </a>
        </div>

        {/* Mobile menu button */}
        <button
          className="text-[#94A3B8] hover:text-white md:hidden"
          onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
        >
          {mobileMenuOpen ? <XIcon /> : <MenuIcon />}
        </button>
      </nav>

      {/* Mobile Navigation */}
      {mobileMenuOpen && (
        <div className="border-t border-[#1E293B] bg-[#0F172A] px-4 py-4 md:hidden">
          <div className="flex flex-col gap-4">
            <Link
              href="/developers"
              className="text-[#94A3B8] transition-colors hover:text-white"
              onClick={() => setMobileMenuOpen(false)}
            >
              {t("developers")}
            </Link>
            <Link
              href="/marketplace"
              className="text-[#94A3B8] transition-colors hover:text-white"
              onClick={() => setMobileMenuOpen(false)}
            >
              {t("marketplace")}
            </Link>
            <a href="#pricing" className="text-[#94A3B8] transition-colors hover:text-white">
              {t("pricing")}
            </a>
            <a href="https://docs.flagbridge.io" className="text-[#94A3B8] transition-colors hover:text-white">
              {t("docs")}
            </a>
            <div className="flex items-center gap-4 border-t border-[#1E293B] pt-4">
              <LocaleSwitcher />
              <a
                href="https://github.com/flagbridge/flagbridge"
                target="_blank"
                rel="noopener noreferrer"
                className="text-[#94A3B8] transition-colors hover:text-white"
              >
                <GithubIcon size={20} />
              </a>
            </div>
          </div>
        </div>
      )}
    </header>
  );
}
