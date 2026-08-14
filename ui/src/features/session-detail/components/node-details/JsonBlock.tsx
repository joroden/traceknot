import { CodeBlock } from "./CodeBlock";
import { JsonPair, type JsonValue } from "./JsonPair";

function parseJSON(raw: string | null): JsonValue | null {
  if (!raw) {
    return null;
  }
  try {
    return JSON.parse(raw) as JsonValue;
  } catch {
    return null;
  }
}

interface JsonBlockProps {
  raw: string | null;
  label?: string;
  fill?: boolean;
}

export function JsonBlock({ raw, label, fill }: JsonBlockProps) {
  const parsed = parseJSON(raw);
  if (parsed === null) {
    return raw ? <CodeBlock text={raw} label={label} fill={fill} /> : null;
  }
  if (typeof parsed !== "object") {
    return <CodeBlock text={JSON.stringify(parsed, null, 2)} label={label} fill={fill} />;
  }
  const entries = Object.entries(parsed as Record<string, JsonValue>);
  return (
    <div className={fill ? "flex min-h-0 flex-1 flex-col" : undefined}>
      {label ? (
        <div className="mb-1 text-xs uppercase tracking-wide text-zinc-500 light:text-zinc-500">{label}</div>
      ) : null}
      <div
        className={`${fill ? "min-h-0 flex-1" : ""} overflow-y-auto rounded border border-zinc-800 bg-zinc-950 p-2 light:border-zinc-200 light:bg-zinc-50`}
      >
        {entries.length === 0 ? (
          <span className="font-mono text-xs text-zinc-500">{Array.isArray(parsed) ? "[]" : "{}"}</span>
        ) : (
          entries.map(([key, value]) => <JsonPair key={key} name={key} value={value} />)
        )}
      </div>
    </div>
  );
}
