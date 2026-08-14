import { useState } from "react";
import { ChevronDown, ChevronRight } from "lucide-react";

interface ChatBlockProps {
  label: string;
  text: string;
  collapsible: boolean;
  muted: boolean;
}

export function ChatBlock({ label, text, collapsible, muted }: ChatBlockProps) {
  const [open, setOpen] = useState(!collapsible);
  return (
    <div
      className={`flex min-h-0 flex-col rounded border border-zinc-800 bg-zinc-800/40 light:border-zinc-200 light:bg-zinc-100/70 ${
        open ? "grow" : ""
      }`}
    >
      <div className="flex shrink-0 items-center gap-1.5 px-2.5 py-1.5">
        {collapsible && (
          <button
            type="button"
            onClick={() => setOpen((prev) => !prev)}
            className="cursor-pointer"
            aria-label={open ? "Collapse reasoning" : "Expand reasoning"}
          >
            {open ? (
              <ChevronDown className="h-3 w-3 text-zinc-500" />
            ) : (
              <ChevronRight className="h-3 w-3 text-zinc-500" />
            )}
          </button>
        )}
        <span className="text-xs font-medium uppercase tracking-wide text-zinc-500 light:text-zinc-500">
          {label}
        </span>
        <span className="ml-auto font-mono text-xs text-zinc-600 light:text-zinc-400">
          {text.length.toLocaleString()} chars
        </span>
      </div>
      {open && (
        <div
          className={`min-h-0 flex-1 overflow-y-auto whitespace-pre-wrap break-words border-t border-zinc-800 px-2.5 pb-2.5 pt-2 text-[13px] leading-relaxed light:border-zinc-200 ${
            muted ? "text-zinc-400 light:text-zinc-600" : "text-zinc-200 light:text-zinc-800"
          }`}
        >
          {text}
        </div>
      )}
    </div>
  );
}
