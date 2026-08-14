import { LogoMark } from "../../../../components/Logo";

const OTLP_COMMANDS = [
  {
    label: "Claude Code",
    env: "export ANTHROPIC_OTEL_EXPORTER_HTTP_HEADERS=\"x-otel-session-id=$(uuidgen)\"\nexport ANTHROPIC_OTEL_EXPORTER_BASE_URL=http://127.0.0.1:4318/v1/traces",
  },
  {
    label: "Any OTLP emitter",
    env: "export OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=http://127.0.0.1:4318/v1/traces",
  },
];

export function FirstRunHero() {
  return (
    <div className="mx-auto my-auto flex max-w-[560px] flex-col items-center gap-4 text-center">
      <span className="grid size-12 place-items-center rounded-lg bg-violet-600 text-white">
        <LogoMark size={24} />
      </span>
      <h2 className="text-lg font-bold">Point an agent at traceknot</h2>
      <p className="text-sm text-zinc-400 light:text-zinc-500">
        traceknot collects telemetry over OTLP on{" "}
        <code className="font-mono text-xs">127.0.0.1:4318</code>. Once a
        session arrives, this dashboard shows cost, tokens, coverage, and
        what the money bought.
      </p>
      <div className="w-full divide-y divide-zinc-800 overflow-hidden rounded-lg border border-zinc-800 bg-zinc-900 text-left light:divide-zinc-200 light:border-zinc-200 light:bg-white">
        {OTLP_COMMANDS.map((entry) => (
          <div key={entry.label} className="p-4">
            <p className="mb-2 text-xs font-semibold text-zinc-100 light:text-zinc-900">
              {entry.label}
            </p>
            <pre className="m-0 overflow-x-auto font-mono text-xs leading-relaxed text-zinc-400 light:text-zinc-500">
              {entry.env}
            </pre>
          </div>
        ))}
      </div>
    </div>
  );
}
