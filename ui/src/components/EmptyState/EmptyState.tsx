import type { ReactNode } from "react";

export interface EmptyStateProps {
  title: string;
  action?: ReactNode;
}

export function EmptyState({ title, action }: EmptyStateProps) {
  return (
    <div className="flex min-h-[160px] flex-1 flex-col items-center justify-center gap-2 rounded-lg border border-dashed border-zinc-700 p-6 text-center light:border-zinc-300">
      <p className="text-sm text-zinc-400 light:text-zinc-500">{title}</p>
      {action ? <div className="flex gap-2">{action}</div> : null}
    </div>
  );
}
