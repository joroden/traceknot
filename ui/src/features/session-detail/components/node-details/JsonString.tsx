import { useState } from "react";
import { CodeBlock } from "./CodeBlock";

export const LONG_STRING_CHARS = 160;

export function JsonString({ value }: { value: string }) {
  const [open, setOpen] = useState(false);
  return (
    <div>
      <button
        type="button"
        onClick={() => setOpen((prev) => !prev)}
        className="cursor-pointer text-emerald-300 hover:underline light:text-emerald-700"
      >
        {open ? (
          <span className="text-zinc-400">– {value.length.toLocaleString()} chars</span>
        ) : (
          `${value.slice(0, LONG_STRING_CHARS)}… (${value.length.toLocaleString()} chars)`
        )}
      </button>
      {open && <CodeBlock text={value} />}
    </div>
  );
}
