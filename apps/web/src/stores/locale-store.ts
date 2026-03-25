import { create } from "zustand";

export type Locale = "en" | "pt";

interface LocaleState {
  locale: Locale;
  setLocale: (locale: Locale) => void;
}

function getInitialLocale(): Locale {
  if (typeof document !== "undefined") {
    const match = document.cookie.match(/NEXT_LOCALE=(\w+)/);
    if (match && (match[1] === "en" || match[1] === "pt")) {
      return match[1];
    }
  }
  return "en";
}

export const useLocaleStore = create<LocaleState>((set) => ({
  locale: getInitialLocale(),
  setLocale: (locale: Locale) => {
    document.cookie = `NEXT_LOCALE=${locale};path=/;max-age=31536000;SameSite=Lax`;
    set({ locale });
  },
}));
