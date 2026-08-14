import { useState } from "react";
import { ChevronDown, ChevronRight } from "lucide-react";

const DEFAULT_MAX_CHARS = 10_000;

interface CodeBlockProps {
  text: string;
  label?: string;
  maxChars?: number;
  fill?: boolean;
}

export function CodeBlock({ text, label, maxChars = DEFAULT_MAX_CHARS, fill }: CodeBlockProps) {
  const [expanded, setExpanded] = useState(false);
  const truncated = text.length > maxChars;
  const shown = expanded || !truncated ? text : text.slice(0, maxChars);

  return (
    <div className={fill ? "flex min-h-0 flex-1 flex-col" : undefined}>
      {label ? (
        <div className="mb-1 flex items-center gap-2 text-xs uppercase tracking-wide text-zinc-500 light:text-zinc-500">
          <span>{label}</span>
          <span className="font-mono normal-case tracking-normal">
            {text.length.toLocaleString()} chars
          </span>
        </div>
      ) : null}
      <pre
        className={`${fill ? "min-h-0 flex-1" : "max-h-[60vh]"} overflow-auto whitespace-pre-wrap break-words rounded border border-zinc-800 bg-zinc-950 p-2 font-mono text-xs leading-relaxed text-zinc-300 light:border-zinc-200 light:bg-zinc-50 light:text-zinc-700`}
      >
        {shown}
      </pre>
      {truncated && (
        <button
          type="button"
          onClick={() => setExpanded((prev) => !prev)}
          className="mt-1 inline-flex cursor-pointer items-center gap-1 text-xs text-violet-400 transition-colors hover:text-violet-300 light:text-violet-600 light:hover:text-violet-500"
        >
          {expanded ? <ChevronDown className="h-3 w-3" /> : <ChevronRight className="h-3 w-3" />}
          {expanded ? "Collapse" : `Show all (${text.length.toLocaleString()} chars)`}
        </button>
      )}
    </div>
  );
}
