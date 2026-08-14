import { useEffect, useState } from "react";
import { SegmentedControl, type SegmentedOption } from "../../../../components/SegmentedControl";

export type RangeKey = "today" | "week" | "month" | "all" | "custom";

export interface DashboardRangeRequest {
  key: RangeKey;
  startMs: number | null;
  endMs: number | null;
}

export interface RangeSelectorProps {
  onRequest: (request: DashboardRangeRequest) => void;
}

const RANGE_OPTIONS: SegmentedOption<RangeKey>[] = [
  { value: "today", label: "Today" },
  { value: "week", label: "This week" },
  { value: "month", label: "This month" },
  { value: "all", label: "All time" },
  { value: "custom", label: "Custom" },
];

function dateToDayStart(value: string): number | null {
  if (!value) {
    return null;
  }
  const [year, month, day] = value.split("-").map(Number);
  if (!year || !month || !day) {
    return null;
  }
  return new Date(year, month - 1, day, 0, 0, 0, 0).getTime();
}

export function RangeSelector({ onRequest }: RangeSelectorProps) {
  const [key, setKey] = useState<RangeKey>("week");
  const [startDate, setStartDate] = useState("");
  const [endDate, setEndDate] = useState("");

  useEffect(() => {
    if (key === "custom") {
      const startMs = dateToDayStart(startDate);
      const endMsRaw = dateToDayStart(endDate);
      if (startMs === null || endMsRaw === null) {
        return;
      }
      const endMs = endMsRaw + 24 * 60 * 60 * 1000;
      if (endMs <= startMs) {
        return;
      }
      onRequest({ key, startMs, endMs });
      return;
    }
    onRequest({ key, startMs: null, endMs: null });
  }, [key, startDate, endDate, onRequest]);

  return (
    <div className="flex flex-wrap items-center gap-2">
      <SegmentedControl
        options={RANGE_OPTIONS}
        value={key}
        onChange={setKey}
      />
      {key === "custom" ? (
        <div className="flex items-center gap-2">
          <input
            type="date"
            value={startDate}
            onChange={(event) => setStartDate(event.target.value)}
            className="rounded-md border border-zinc-700 bg-zinc-900 px-2 py-1.5 text-xs text-zinc-100 outline-none transition-colors focus:border-violet-500 [color-scheme:dark] light:border-zinc-300 light:bg-white light:text-zinc-900 light:[color-scheme:light]"
          />
          <span className="text-xs text-zinc-500">to</span>
          <input
            type="date"
            value={endDate}
            onChange={(event) => setEndDate(event.target.value)}
            className="rounded-md border border-zinc-700 bg-zinc-900 px-2 py-1.5 text-xs text-zinc-100 outline-none transition-colors focus:border-violet-500 [color-scheme:dark] light:border-zinc-300 light:bg-white light:text-zinc-900 light:[color-scheme:light]"
          />
        </div>
      ) : null}
    </div>
  );
}
