import { useState } from "react";
import { ChevronDown, ChevronRight } from "lucide-react";

const DEFAULT_MAX_LINES = 400;

export function looksLikeDiff(text: string): boolean {
  return (
    text.startsWith("diff --git") ||
    /^\+\+\+ |^--- |^@@ /.test(text) ||
    text.includes("\n@@ ")
  );
}

function lineClass(line: string): string {
  if (line.startsWith("+++") || line.startsWith("---")) {
    return "text-zinc-500";
  }
  if (line.startsWith("@")) {
    return "text-amber-300 light:text-amber-600";
  }
  if (line.startsWith("+")) {
    return "bg-emerald-500/10 text-emerald-300 light:bg-emerald-500/10 light:text-emerald-700";
  }
  if (line.startsWith("-")) {
    return "bg-red-500/10 text-red-300 light:bg-red-500/10 light:text-red-700";
  }
  return "";
}

interface DiffBlockProps {
  text: string;
  label?: string;
  maxLines?: number;
  fill?: boolean;
}

export function DiffBlock({ text, label, maxLines = DEFAULT_MAX_LINES, fill }: DiffBlockProps) {
  const [expanded, setExpanded] = useState(false);
  const lines = text.split("\n");
  const truncated = lines.length > maxLines;
  const shown = expanded || !truncated ? lines : lines.slice(0, maxLines);

  return (
    <div className={fill ? "flex min-h-0 flex-1 flex-col" : undefined}>
      {label ? (
        <div className="mb-1 flex items-center gap-2 text-xs uppercase tracking-wide text-zinc-500 light:text-zinc-500">
          <span>{label}</span>
          <span className="font-mono normal-case tracking-normal">
            {lines.length.toLocaleString()} lines
          </span>
        </div>
      ) : null}
      <pre
        className={`${fill ? "min-h-0 flex-1" : "max-h-[60vh]"} overflow-auto whitespace-pre-wrap break-words rounded border border-zinc-800 bg-zinc-950 p-2 font-mono text-xs leading-relaxed text-zinc-300 light:border-zinc-200 light:bg-zinc-50 light:text-zinc-700`}
      >
        {shown.map((line, index) => (
          <div key={index} className={lineClass(line)}>
            {line || " "}
          </div>
        ))}
      </pre>
      {truncated && (
        <button
          type="button"
          onClick={() => setExpanded((prev) => !prev)}
          className="mt-1 inline-flex cursor-pointer items-center gap-1 text-xs text-violet-400 transition-colors hover:text-violet-300 light:text-violet-600 light:hover:text-violet-500"
        >
          {expanded ? <ChevronDown className="h-3 w-3" /> : <ChevronRight className="h-3 w-3" />}
          {expanded ? "Collapse" : `Show all (${lines.length.toLocaleString()} lines)`}
        </button>
      )}
    </div>
  );
}
