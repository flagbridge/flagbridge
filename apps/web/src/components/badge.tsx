export function Badge({
  children,
  variant = "default",
}: {
  children: React.ReactNode;
  variant?: "default" | "blue" | "green";
}) {
  const colors = {
    default: "bg-[#1E293B] text-[#94A3B8] border-[#334155]",
    blue: "bg-[#1E3A5F] text-[#60A5FA] border-[#60A5FA]/20",
    green: "bg-[#052E16] text-[#22C55E] border-[#22C55E]/20",
  };

  return (
    <span
      className={`inline-block rounded-full border px-2 py-1 text-xs uppercase tracking-wide ${colors[variant]}`}
    >
      {children}
    </span>
  );
}
