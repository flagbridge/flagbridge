export function CodeBlock({
  code,
  language = "bash",
}: {
  code: string;
  language?: string;
}) {
  return (
    <div className="overflow-hidden rounded-xl border border-[#334155] bg-[#1E293B]">
      <div className="flex items-center justify-between border-b border-[#334155] px-5 py-3">
        <span className="text-xs uppercase tracking-wide text-[#94A3B8]">
          {language}
        </span>
      </div>
      <div className="p-5">
        <pre className="overflow-x-auto font-mono text-sm leading-relaxed text-[#22D3EE]">
          <code>{code}</code>
        </pre>
      </div>
    </div>
  );
}
