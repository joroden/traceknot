export function StatusDot({ status }: { status: string | null }) {
  const color =
    status === "ok"
      ? "bg-emerald-500"
      : status === "error"
        ? "bg-red-500"
        : status === "running" || status === "in_progress"
          ? "bg-amber-500"
          : "bg-zinc-600";
  return (
    <span className="flex items-center justify-center">
      <span className={`h-2 w-2 rounded-full ${color}`} title={status ?? "no status"} />
    </span>
  );
}
