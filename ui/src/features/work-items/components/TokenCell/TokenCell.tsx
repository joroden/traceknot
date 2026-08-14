import { ArrowDown, ArrowUp } from "lucide-react";
import { formatTokens } from "../../../../utils/format";
import type { TokenBreakdown } from "../../api";

export interface TokenCellProps {
  tokens: TokenBreakdown;
  kind: "input" | "output";
}

function breakdownParts(tokens: TokenBreakdown, kind: "input" | "output"): string[] {
  if (kind === "input") {
    return [
      tokens.raw ? `${formatTokens(tokens.raw)} raw` : null,
      tokens.cached ? `${formatTokens(tokens.cached)} cached` : null,
      tokens.write ? `${formatTokens(tokens.write)} write` : null,
    ].filter((part): part is string => part !== null);
  }
  return tokens.reasoning ? [`${formatTokens(tokens.reasoning)} reasoning`] : [];
}

export function TokenCell({ tokens, kind }: TokenCellProps) {
  const Icon = kind === "input" ? ArrowUp : ArrowDown;
  const iconClass = kind === "input" ? "text-emerald-500" : "text-red-500";
  const breakdown = breakdownParts(tokens, kind).join(" · ");

  return (
    <div className="flex min-w-0 flex-col">
      <span className="flex items-center gap-1 font-mono text-xs text-zinc-200 light:text-zinc-800">
        <Icon className={`h-3 w-3 shrink-0 ${iconClass}`} />
        {formatTokens(tokens.total)}
      </span>
      {breakdown ? (
        <span title={breakdown} className="truncate text-xs text-zinc-500">
          {breakdown}
        </span>
      ) : null}
    </div>
  );
}
