import { Bot, Circle, GitBranch, MessageSquare, Wrench, Zap } from "lucide-react";
import type { LucideIcon } from "lucide-react";

const KIND_STYLES: Record<string, { icon: LucideIcon; label: string; className: string }> = {
  agent: {
    icon: Bot,
    label: "Agent",
    className: "bg-violet-500/15 text-violet-400 light:bg-violet-500/10 light:text-violet-600",
  },
  subagent: {
    icon: GitBranch,
    label: "Subagent",
    className: "bg-emerald-500/15 text-emerald-400 light:bg-emerald-500/10 light:text-emerald-600",
  },
  tool_call: {
    icon: Wrench,
    label: "Tool",
    className: "bg-red-500/15 text-red-400 light:bg-red-500/10 light:text-red-600",
  },
  chat: {
    icon: MessageSquare,
    label: "Chat",
    className: "bg-zinc-500/15 text-zinc-400 light:bg-zinc-500/10 light:text-zinc-600",
  },
  event: {
    icon: Zap,
    label: "Event",
    className: "bg-amber-500/15 text-amber-400 light:bg-amber-500/10 light:text-amber-600",
  },
};

export function kindStyle(kind: string): {
  icon: LucideIcon;
  label: string;
  className: string;
} {
  return KIND_STYLES[kind] ?? {
    icon: Circle,
    label: kind || "Node",
    className: "bg-zinc-500/15 text-zinc-400 light:bg-zinc-500/10 light:text-zinc-600",
  };
}

export function KindBadge({ kind }: { kind: string }) {
  const style = kindStyle(kind);
  const Icon = style.icon;
  return (
    <span
      className={`inline-flex items-center gap-1 rounded px-1.5 py-0.5 text-xs font-medium ${style.className}`}
      title={style.label}
    >
      <Icon className="h-3 w-3" />
      <span className="sr-only">{style.label}</span>
    </span>
  );
}
