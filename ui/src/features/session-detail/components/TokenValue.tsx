import { formatTokens } from "../../../utils/format";

export function TokenValue({
  aggregated,
  self,
  estimated,
  title,
}: {
  aggregated: number;
  self: number;
  estimated: number;
  title: string;
}) {
  const value = aggregated > 0 ? aggregated : self;
  const estimateOnly = value === 0 && estimated > 0;
  return (
    <div className="flex items-center justify-end pr-3 font-mono text-xs tabular-nums">
      <span
        className={estimateOnly ? "text-zinc-500" : "text-zinc-300 light:text-zinc-600"}
        title={`${title} — self: ${formatTokens(self)} · subtree: ${formatTokens(aggregated)}${estimateOnly ? " (estimated)" : ""}`}
      >
        {estimateOnly ? "~" : ""}
        {formatTokens(estimateOnly ? estimated : value)}
      </span>
    </div>
  );
}
