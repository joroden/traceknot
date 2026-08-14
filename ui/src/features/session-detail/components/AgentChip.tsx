import type { MouseEvent } from "react";

interface AgentChipProps {
  name: string | null;
  onClick?: (event: MouseEvent<HTMLButtonElement>) => void;
}

export function AgentChip({ name, onClick }: AgentChipProps) {
  const label = name ?? "agent";
  const classes =
    "inline-flex items-center rounded border border-violet-500/30 bg-violet-500/10 px-1.5 py-0.5 font-mono text-xs text-violet-300 light:border-violet-500/30 light:bg-violet-500/10 light:text-violet-600";
  if (!onClick) {
    return <span className={classes}>{label}</span>;
  }
  return (
    <button
      type="button"
      onClick={onClick}
      title={`Focus this agent's subtree`}
      className={`${classes} cursor-pointer transition-colors hover:bg-violet-500/20`}
    >
      {label}
    </button>
  );
}
