import {
  Area,
  AreaChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { formatBucketLabel, formatUSD } from "../../../../utils/format";
import type { TrendBucket } from "../../api";

const MODEL_COLORS = [
  "#8b5cf6",
  "#0ea5e9",
  "#f59e0b",
  "#10b981",
  "#f43f5e",
  "#14b8a6",
  "#d946ef",
  "#71717a",
];

export interface TrendChartProps {
  buckets: TrendBucket[];
  granularityMs: number;
}

export function TrendChart({ buckets, granularityMs }: TrendChartProps) {
  const modelTotals = new Map<string, number>();
  for (const bucket of buckets) {
    for (const model of bucket.models) {
      modelTotals.set(model.name, (modelTotals.get(model.name) ?? 0) + model.cost);
    }
  }
  const series = [...modelTotals.entries()]
    .sort((left, right) => right[1] - left[1])
    .map(([name]) => name);

  const data = buckets.map((bucket) => {
    const row: Record<string, number | string> = {
      label: formatBucketLabel(bucket.start_unix_ms, granularityMs),
    };
    for (const name of series) {
      row[name] = 0;
    }
    for (const model of bucket.models) {
      row[model.name] = model.cost;
    }
    return row;
  });

  return (
    <div className="h-56">
      <ResponsiveContainer width="100%" height="100%">
        <AreaChart data={data} margin={{ top: 4, right: 4, left: 4, bottom: 0 }}>
          <defs>
            {series.map((name, index) => (
              <linearGradient key={name} id={`trend-fill-${index}`} x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stopColor={MODEL_COLORS[index % MODEL_COLORS.length]} stopOpacity={0.45} />
                <stop offset="100%" stopColor={MODEL_COLORS[index % MODEL_COLORS.length]} stopOpacity={0.05} />
              </linearGradient>
            ))}
          </defs>
          <CartesianGrid stroke="#26262b" strokeDasharray="3 3" vertical={false} />
          <XAxis
            dataKey="label"
            tick={{ fontSize: 10, fill: "#71717a" }}
            tickLine={false}
            axisLine={{ stroke: "#3f3f46" }}
            interval="preserveStartEnd"
            minTickGap={24}
          />
          <YAxis
            tick={{ fontSize: 10, fill: "#71717a" }}
            tickLine={false}
            axisLine={false}
            width={44}
            tickFormatter={(value: number) => formatUSD(value)}
          />
          <Tooltip
            cursor={{ fill: "#1c1c20", opacity: 0.6 }}
            contentStyle={{
              background: "#141416",
              border: "1px solid #3a3a42",
              borderRadius: 8,
              fontSize: 12,
              fontFamily: "JetBrains Mono, monospace",
            }}
            labelStyle={{ color: "#ededef", marginBottom: 4 }}
            itemStyle={{ color: "#9e9ea7" }}
            formatter={(value) =>
              typeof value === "number" ? formatUSD(value) : ""
            }
          />
          {series.map((name, index) => (
            <Area
              key={name}
              type="monotone"
              dataKey={name}
              stackId="cost"
              stroke={MODEL_COLORS[index % MODEL_COLORS.length]}
              strokeWidth={1.5}
              fill={`url(#trend-fill-${index})`}
            />
          ))}
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
}
