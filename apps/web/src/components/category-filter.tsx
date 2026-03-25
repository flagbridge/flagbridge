"use client";

import { useState } from "react";

const categories = [
  "All",
  "Messaging",
  "Analytics",
  "Queues",
  "Security",
  "UI Extensions",
  "Automation",
  "Data",
  "CI/CD",
];

function SearchIcon() {
  return (
    <svg
      width={20}
      height={20}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={1.5}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z"
      />
    </svg>
  );
}

export function CategoryFilter({
  searchPlaceholder,
  statsText,
}: {
  searchPlaceholder: string;
  statsText: string;
}) {
  const [selectedCategory, setSelectedCategory] = useState("All");

  return (
    <>
      {/* Search Bar */}
      <div className="relative mx-auto mb-8 max-w-2xl">
        <div className="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-[#64748B]">
          <SearchIcon />
        </div>
        <input
          type="text"
          placeholder={searchPlaceholder}
          className="w-full rounded-xl border border-[#334155] bg-[#1E293B] py-4 pl-12 pr-4 text-[#E2E8F0] placeholder-[#64748B] transition-colors focus:border-[#3B82F6] focus:outline-none"
        />
      </div>

      {/* Category Pills */}
      <div className="mb-6 flex flex-wrap justify-center gap-3 pb-4">
        {categories.map((category) => (
          <button
            key={category}
            onClick={() => setSelectedCategory(category)}
            className={`whitespace-nowrap rounded-full px-4 py-2 text-sm font-medium transition-colors ${
              selectedCategory === category
                ? "bg-[#3B82F6] text-white"
                : "border border-[#334155] bg-[#1E293B] text-[#94A3B8] hover:border-[#475569]"
            }`}
          >
            {category}
          </button>
        ))}
      </div>

      {/* Stats */}
      <p className="text-sm text-[#64748B]">{statsText}</p>
    </>
  );
}
