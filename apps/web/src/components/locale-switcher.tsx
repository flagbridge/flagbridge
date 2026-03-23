"use client";

import { useLocale } from "next-intl";
import { useRouter, usePathname } from "next/navigation";

export function LocaleSwitcher() {
  const locale = useLocale();
  const router = useRouter();
  const pathname = usePathname();

  function switchTo(newLocale: string) {
    const segments = pathname.split("/");
    segments[1] = newLocale;
    router.push(segments.join("/"));
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
