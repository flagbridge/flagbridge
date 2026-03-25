"use client";

import { useRouter } from "next/navigation";
import { useLocaleStore, type Locale } from "@/stores/locale-store";

export function LocaleSwitcher() {
  const router = useRouter();
  const { locale, setLocale } = useLocaleStore();

  function switchTo(newLocale: Locale) {
    setLocale(newLocale);
    router.refresh();
  }

  return (
    <div className="flex items-center gap-1">
      <button
        onClick={() => switchTo("en")}
        className={`rounded px-1.5 py-0.5 text-base transition ${
          locale === "en" ? "bg-[#222a3d]" : "opacity-50 hover:opacity-100"
        }`}
        aria-label="English"
      >
        🇺🇸
      </button>
      <button
        onClick={() => switchTo("pt")}
        className={`rounded px-1.5 py-0.5 text-base transition ${
          locale === "pt" ? "bg-[#222a3d]" : "opacity-50 hover:opacity-100"
        }`}
        aria-label="Portugues"
      >
        🇧🇷
      </button>
    </div>
  );
}
