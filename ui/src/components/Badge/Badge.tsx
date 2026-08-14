import type { ReactNode } from "react";

export type BadgeTone = "neutral" | "violet" | "amber" | "emerald" | "red";

export interface BadgeProps {
  children: ReactNode;
  tone?: BadgeTone;
}

const TONE_CLASSES: Record<BadgeTone, string> = {
  neutral:
    "border-zinc-700 bg-zinc-800/60 text-zinc-300 light:border-zinc-300 light:bg-zinc-200 light:text-zinc-700",
  violet:
    "border-violet-700 bg-violet-900/40 text-violet-300 light:border-violet-300 light:bg-violet-100 light:text-violet-700",
  amber:
    "border-amber-700 bg-amber-900/30 text-amber-300 light:border-amber-300 light:bg-amber-100 light:text-amber-700",
  emerald:
    "border-emerald-700 bg-emerald-900/30 text-emerald-300 light:border-emerald-300 light:bg-emerald-100 light:text-emerald-700",
  red: "border-red-700 bg-red-900/30 text-red-300 light:border-red-300 light:bg-red-100 light:text-red-700",
};

export function Badge({ children, tone = "neutral" }: BadgeProps) {
  return (
    <span
      className={`inline-flex shrink-0 items-center rounded border px-1.5 py-0.5 font-mono text-xs uppercase tracking-wide ${TONE_CLASSES[tone]}`}
    >
      {children}
    </span>
  );
}
