function toDateInputValue(unixMs: number | undefined): string {
  if (!unixMs) {
    return "";
  }
  const date = new Date(unixMs);
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, "0")}-${String(date.getDate()).padStart(2, "0")}`;
}

function fromDateInputValue(value: string): number | null {
  if (!value) {
    return null;
  }
  const [year, month, day] = value.split("-").map(Number);
  if (!year || !month || !day) {
    return null;
  }
  return new Date(year, month - 1, day, 0, 0, 0, 0).getTime();
}

export interface DateRangeFilterProps {
  startUnixMs: number | undefined;
  endUnixMs: number | undefined;
  onChange: (startUnixMs: number | undefined, endUnixMs: number | undefined) => void;
}

const inputClass =
  "w-full rounded-md border border-zinc-700 bg-zinc-950 px-2 py-1.5 text-xs text-zinc-100 outline-none focus:border-violet-500 [color-scheme:dark] light:border-zinc-300 light:bg-white light:text-zinc-900 light:[color-scheme:light]";

export function DateRangeFilter({ startUnixMs, endUnixMs, onChange }: DateRangeFilterProps) {
  return (
    <div className="flex flex-col gap-2">
      <label className="flex flex-col gap-1 text-xs text-zinc-500">
        From
        <input
          type="date"
          className={inputClass}
          value={toDateInputValue(startUnixMs)}
          onChange={(event) => {
            const startMs = fromDateInputValue(event.target.value) ?? undefined;
            onChange(startMs, endUnixMs);
          }}
        />
      </label>
      <label className="flex flex-col gap-1 text-xs text-zinc-500">
        To
        <input
          type="date"
          className={inputClass}
          value={toDateInputValue(endUnixMs)}
          onChange={(event) => {
            const dayMs = fromDateInputValue(event.target.value);
            const endMs = dayMs === null ? undefined : dayMs + 24 * 60 * 60 * 1000;
            onChange(startUnixMs, endMs);
          }}
        />
      </label>
      {startUnixMs || endUnixMs ? (
        <button
          type="button"
          onClick={() => onChange(undefined, undefined)}
          className="cursor-pointer self-start text-xs text-violet-400 hover:text-violet-300"
        >
          Clear
        </button>
      ) : null}
    </div>
  );
}
