import { useState } from "react";
import { ChevronDown, ChevronRight } from "lucide-react";
import { JsonString, LONG_STRING_CHARS } from "./JsonString";

export type JsonValue = string | number | boolean | null | JsonValue[] | { [key: string]: JsonValue };

export function JsonPair({ name, value }: { name: string; value: JsonValue }) {
  const [open, setOpen] = useState(false);
  const isObject = value !== null && typeof value === "object";

  if (isObject) {
    const entries = Object.entries(value as Record<string, JsonValue>);
    const summary = Array.isArray(value) ? `[${entries.length}]` : `{${entries.length}}`;
    return (
      <div>
        <button
          type="button"
          onClick={() => setOpen((prev) => !prev)}
          className="flex w-full cursor-pointer items-center gap-1.5 rounded px-1 py-0.5 text-left font-mono text-xs hover:bg-zinc-800/60 light:hover:bg-zinc-100"
        >
          {open ? (
            <ChevronDown className="h-3 w-3 shrink-0 text-zinc-500" />
          ) : (
            <ChevronRight className="h-3 w-3 shrink-0 text-zinc-500" />
          )}
          <span className="text-violet-300 light:text-violet-600">{name}</span>
          <span className="text-zinc-500">{summary}</span>
        </button>
        {open && (
          <div className="ml-4 border-l border-zinc-800 pl-2 light:border-zinc-200">
            {entries.map(([key, child]) => (
              <JsonPair key={key} name={key} value={child} />
            ))}
          </div>
        )}
      </div>
    );
  }

  const rendered =
    typeof value === "string" ? (
      value.length > LONG_STRING_CHARS ? (
        <JsonString value={value} />
      ) : (
        <span className="text-emerald-300 light:text-emerald-700">{value}</span>
      )
    ) : value === null ? (
      <span className="text-zinc-500">null</span>
    ) : typeof value === "number" ? (
      <span className="text-red-300 light:text-red-600">{value}</span>
    ) : (
      <span className="text-amber-300 light:text-amber-600">{String(value)}</span>
    );

  return (
    <div className="flex items-start gap-1.5 px-1 py-0.5 font-mono text-xs">
      <span className="shrink-0 pl-4 text-violet-300 light:text-violet-600">{name}</span>
      <div className="min-w-0 flex-1">{rendered}</div>
    </div>
  );
}
