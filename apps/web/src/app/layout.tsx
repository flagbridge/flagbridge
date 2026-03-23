import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";

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

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className="dark">
      <body className={`${inter.className} antialiased`}>{children}</body>
    </html>
  );
}
