import { useState } from "react";
import { AlignLeft, ChevronUp, Copy, Quote } from "lucide-react";

export interface PromptBannerProps {
  prompt: string;
}

export function PromptBanner({ prompt }: PromptBannerProps) {
  const [expanded, setExpanded] = useState(false);
  const [copied, setCopied] = useState(false);
  const firstLine = prompt.trim().split("\n")[0] ?? "";
  const lineCount = prompt.trim().split("\n").length;

  const copyPrompt = async () => {
    await navigator.clipboard.writeText(prompt);
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  };

  return (
    <div className="rounded-lg border border-zinc-800 bg-zinc-900 light:border-zinc-200 light:bg-white">
      <div className="flex items-start gap-2.5 px-3 py-2.5">
        <span className="mt-0.5 inline-flex text-amber-500 light:text-amber-600">
          <Quote size={12} />
        </span>
        <div className="min-w-0 flex-1">
          <div className="font-mono text-xs font-bold uppercase tracking-wider text-amber-500 light:text-amber-600">
            Target prompt
          </div>
          <p className="mt-0.5 select-all truncate font-mono text-xs text-zinc-400 light:text-zinc-500">
            {firstLine}
          </p>
        </div>
        <button
          type="button"
          className="inline-flex flex-shrink-0 cursor-pointer items-center gap-1.5 rounded-md border border-zinc-800 bg-zinc-800 px-2 py-1 text-xs text-zinc-400 transition-colors hover:border-zinc-700 hover:text-zinc-100 light:border-zinc-200 light:bg-zinc-100 light:text-zinc-500 light:hover:text-zinc-900"
          onClick={() => setExpanded((value) => !value)}
        >
          {expanded ? <ChevronUp size={12} /> : <AlignLeft size={12} />}
          {expanded ? "Hide prompt" : "View full prompt"}
        </button>
      </div>
      {expanded ? (
        <div className="max-h-[220px] overflow-y-auto border-t border-zinc-800 px-3 py-2.5 light:border-zinc-200">
          <div className="mb-2 flex items-center justify-between font-mono text-xs text-zinc-400 light:text-zinc-500">
            <span>
              Full prompt ({lineCount} {lineCount === 1 ? "line" : "lines"})
            </span>
            <button type="button" className="inline-flex cursor-pointer items-center gap-1 rounded-md border border-zinc-800 bg-zinc-800 px-[7px] py-[3px] text-xs text-zinc-400 transition-colors hover:text-zinc-100 light:border-zinc-200 light:bg-zinc-100 light:text-zinc-500 light:hover:text-zinc-900" onClick={copyPrompt}>
              <Copy size={11} />
              {copied ? "Copied" : "Copy"}
            </button>
          </div>
          <pre className="m-0 select-all whitespace-pre-wrap break-words font-mono text-xs leading-relaxed text-zinc-100 light:text-zinc-900">
            {prompt}
          </pre>
        </div>
      ) : null}
    </div>
  );
}
