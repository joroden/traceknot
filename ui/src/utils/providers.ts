import type { BadgeTone } from "../components/Badge";

const PROVIDER_COLORS: Record<string, { badge: BadgeTone; bar: string }> = {
  claude: { badge: "violet", bar: "bg-violet-600" },
  codex: { badge: "emerald", bar: "bg-emerald-600" },
  copilot: { badge: "neutral", bar: "bg-zinc-500" },
};

export function providerBadgeTone(provider: string): BadgeTone {
  return PROVIDER_COLORS[provider]?.badge ?? "neutral";
}

export function providerBarClass(provider: string): string {
  return PROVIDER_COLORS[provider]?.bar ?? "bg-zinc-500";
}
