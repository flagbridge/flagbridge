import type { Metadata } from "next";

export const metadata: Metadata = {
  icons: {
    icon: [{ url: "https://flagbridge.io/favicon.png", type: "image/png" }],
    apple: [
      { url: "https://flagbridge.io/apple-icon.png", type: "image/png" },
    ],
  },
  manifest: "https://flagbridge.io/manifest.json",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return children;
}
