export interface DaemonUnreachableProps {
  error: string;
  showStartHint?: boolean;
}

export function DaemonUnreachable({ error, showStartHint = false }: DaemonUnreachableProps) {
  return (
    <div className="mx-auto my-auto max-w-md rounded-lg border border-red-500 bg-zinc-900 p-6 light:bg-white">
      <p className="mb-1.5 text-sm font-bold text-red-500">Cannot reach the traceknot daemon</p>
      <p className="mt-1 text-xs text-zinc-400 light:text-zinc-500">{error}</p>
      {showStartHint ? (
        <p className="mt-1 text-xs text-zinc-400 light:text-zinc-500">
          Start it with <code className="font-mono text-xs">go run ./cmd/traceknot</code>, then reload this page.
        </p>
      ) : null}
    </div>
  );
}
