import type { Metadata } from "next";
import { Analytics } from "@vercel/analytics/next";
import { GoogleTagManager } from "@next/third-parties/google";
import { Inter } from "next/font/google";
import { NextIntlClientProvider } from "next-intl";
import { getMessages, setRequestLocale } from "next-intl/server";
import { notFound } from "next/navigation";
import { routing } from "@/i18n/routing";
import "../globals.css";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "FlagBridge — Feature Flags with Product Intelligence",
  description:
    "Open-source feature flag management platform with product intelligence. Self-hosted, extensible, OpenFeature compatible.",
  metadataBase: new URL("https://flagbridge.io"),
  openGraph: {
    title: "FlagBridge — Feature Flags with Product Intelligence",
    description:
      "Open-source feature flag management with product intelligence. Ship features safely and understand their impact.",
    url: "https://flagbridge.io",
    siteName: "FlagBridge",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: "FlagBridge",
    description: "Feature flags with product intelligence. Open source.",
  },
};

export function generateStaticParams() {
  return routing.locales.map((locale) => ({ locale }));
}

export default async function LocaleLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;

  if (!routing.locales.includes(locale as "en" | "pt")) {
    notFound();
  }

  setRequestLocale(locale);
  const messages = await getMessages();

  return (
    <html lang={locale} className="dark">
      <body className={`${inter.className} antialiased`}>
        <GoogleTagManager gtmId="GTM-W93TCMBF" />
        <NextIntlClientProvider messages={messages}>
          {children}
        </NextIntlClientProvider>
        <Analytics />
      </body>
    </html>
  );
}
