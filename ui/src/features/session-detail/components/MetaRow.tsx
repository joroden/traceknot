import type { ReactNode } from "react";

export function MetaRow({ label, children }: { label: string; children: ReactNode }) {
  return (
    <div className="flex items-center gap-2">
      <span className="w-20 shrink-0 text-xs uppercase tracking-wide text-zinc-500 light:text-zinc-500">
        {label}
      </span>
      <span className="min-w-0 truncate font-mono text-xs text-zinc-200 light:text-zinc-800">
        {children}
      </span>
    </div>
  );
}
