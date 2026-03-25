import { getRequestConfig } from "next-intl/server";
import { cookies } from "next/headers";

const SUPPORTED_LOCALES = ["en", "pt"] as const;
type Locale = (typeof SUPPORTED_LOCALES)[number];

export default getRequestConfig(async () => {
  const cookieStore = await cookies();
  const cookieLocale = cookieStore.get("NEXT_LOCALE")?.value;

  const locale: Locale =
    cookieLocale && SUPPORTED_LOCALES.includes(cookieLocale as Locale)
      ? (cookieLocale as Locale)
      : "en";

  return {
    locale,
    messages: (await import(`@/messages/${locale}.json`)).default,
  };
});
