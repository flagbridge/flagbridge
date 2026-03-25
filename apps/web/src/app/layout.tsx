import type { Metadata } from "next";
import { Analytics } from "@vercel/analytics/next";
import { GoogleTagManager } from "@next/third-parties/google";
import { Inter } from "next/font/google";
import { NextIntlClientProvider } from "next-intl";
import { getMessages, getLocale } from "next-intl/server";
import "./globals.css";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "FlagBridge — Feature Flags with Product Intelligence",
  description:
    "Open-source feature flag management platform with product intelligence. Self-hosted, extensible, OpenFeature compatible.",
  metadataBase: new URL("https://flagbridge.io"),
  icons: {
    icon: [{ url: "https://flagbridge.io/favicon.png", type: "image/png" }],
    apple: [
      { url: "https://flagbridge.io/apple-icon.png", type: "image/png" },
    ],
  },
  manifest: "https://flagbridge.io/manifest.json",
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

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const locale = await getLocale();
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
