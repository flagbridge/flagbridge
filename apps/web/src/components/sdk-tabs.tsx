"use client";

import { useState } from "react";

const sdkLanguages = [
  {
    name: "Node.js",
    code: `import { FlagBridge } from '@flagbridge/sdk-node';

const fb = new FlagBridge({
  apiKey: process.env.FLAGBRIDGE_KEY
});

const enabled = await fb.isEnabled('checkout-v2', {
  userId: user.id
});`,
  },
  {
    name: "Go",
    code: `import "github.com/flagbridge/sdk-go"

fb := flagbridge.New(flagbridge.Config{
  APIKey: os.Getenv("FLAGBRIDGE_KEY"),
})

enabled := fb.IsEnabled("checkout-v2", flagbridge.Context{
  UserID: user.ID,
})`,
  },
  {
    name: "Python",
    code: `from flagbridge import FlagBridge

fb = FlagBridge(
  api_key=os.environ["FLAGBRIDGE_KEY"]
)

enabled = fb.is_enabled("checkout-v2",
  user_id=user.id
)`,
  },
  {
    name: "React",
    code: `import { useFlag } from '@flagbridge/react';

function CheckoutPage() {
  const enabled = useFlag('checkout-v2');

  return enabled ?
    <CheckoutV2 /> :
    <LegacyCheckout />;
}`,
  },
  {
    name: "OpenFeature",
    code: `import { OpenFeature } from '@openfeature/server-sdk';
import { FlagBridgeProvider } from '@flagbridge/openfeature';

OpenFeature.setProvider(new FlagBridgeProvider({
  apiKey: process.env.FLAGBRIDGE_KEY
}));

const client = OpenFeature.getClient();
const enabled = await client.getBooleanValue('checkout-v2', false);`,
  },
];

export function SdkTabs() {
  const [activeTab, setActiveTab] = useState(0);

  return (
    <div className="overflow-hidden rounded-xl border border-[#334155] bg-[#1E293B]">
      <div className="flex overflow-x-auto border-b border-[#334155]">
        {sdkLanguages.map((lang, idx) => (
          <button
            key={idx}
            onClick={() => setActiveTab(idx)}
            className={`whitespace-nowrap px-4 py-3 text-sm font-medium transition-colors ${
              activeTab === idx
                ? "border-b-2 border-[#3B82F6] bg-[#0F172A] text-[#F8FAFC]"
                : "text-[#94A3B8] hover:text-[#E2E8F0]"
            }`}
          >
            {lang.name}
          </button>
        ))}
      </div>
      <div className="p-5">
        <pre className="overflow-x-auto font-mono text-sm leading-relaxed text-[#22D3EE]">
          <code>{sdkLanguages[activeTab].code}</code>
        </pre>
      </div>
    </div>
  );
}
