import { AlertTriangle } from "lucide-react";
import { formatDeltaPct } from "../../utils/format";

export type StatCardTone = "neutral" | "attention";

export interface StatCardProps {
  label: string;
  value: string;
  deltaPct?: number | null;
  goodWhenDown?: boolean;
  note?: string;
  tone?: StatCardTone;
  onSelect?: () => void;
}

export function StatCard({
  label,
  value,
  deltaPct = null,
  goodWhenDown = false,
  note,
  tone = "neutral",
  onSelect,
}: StatCardProps) {
  const delta = formatDeltaPct(deltaPct);
  const positive = deltaPct !== null && deltaPct >= 0;
  const deltaGood = deltaPct === null ? null : goodWhenDown ? !positive : positive;
  const deltaClass =
    deltaGood === null
      ? "text-zinc-500"
      : deltaGood
        ? "text-emerald-500"
        : "text-red-500";
  const attention = tone === "attention";

  const containerClassName = `rounded-lg border bg-zinc-900 p-4 text-left light:bg-white ${
    attention ? "border-amber-800/60 light:border-amber-300" : "border-zinc-800 light:border-zinc-200"
  } ${onSelect ? "w-full cursor-pointer transition-colors hover:border-violet-600 light:hover:border-violet-400" : ""}`;

  const content = (
    <>
      <p className="text-xs text-zinc-400 light:text-zinc-500">{label}</p>
      <div className="mt-1.5 flex items-baseline gap-2">
        <span className="font-mono text-xl font-semibold text-zinc-100 light:text-zinc-900">
          {value}
        </span>
        {deltaPct !== null ? (
          <span className={`font-mono text-xs ${deltaClass}`}>{delta}</span>
        ) : null}
      </div>
      {note ? (
        <p
          className={`mt-1 flex items-center gap-1 text-2xs ${
            attention ? "text-amber-500 light:text-amber-600" : "text-zinc-500"
          }`}
        >
          {attention ? <AlertTriangle className="h-3 w-3 shrink-0" /> : null}
          {note}
        </p>
      ) : null}
    </>
  );

  if (onSelect) {
    return (
      <button type="button" onClick={onSelect} className={containerClassName}>
        {content}
      </button>
    );
  }
  return <div className={containerClassName}>{content}</div>;
}
