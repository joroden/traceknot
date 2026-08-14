import { BookOpen, KeyRound, PackageX } from "lucide-react";
import type { ProviderProbe } from "../../types/workItem";

export interface ProviderPanelProps {
  probe: ProviderProbe;
}

function iconFor(status: ProviderProbe["status"]) {
  if (status === "not_authenticated") {
    return <KeyRound size={18} />;
  }
  return <PackageX size={18} />;
}

export function ProviderPanel({ probe }: ProviderPanelProps) {
  const missing = probe.status === "cli_missing";
  const heading = missing
    ? `${probe.provider} CLI is not installed`
    : `${probe.provider} is not authenticated`;
  const docsHref = missing
    ? probe.install_docs_url
    : probe.auth_docs_url;
  const docsLabel = missing ? "View install docs" : "View auth docs";

  return (
    <div className="flex flex-1 flex-col items-center justify-center gap-2.5 rounded-lg border border-zinc-800 bg-zinc-900 p-8 text-center light:border-zinc-200 light:bg-white">
      <span className="inline-flex rounded-lg border border-zinc-700 p-2.5 text-amber-500 light:border-zinc-300 light:text-amber-600">
        {iconFor(probe.status)}
      </span>
      <h2 className="text-sm font-semibold">{heading}</h2>
      {probe.hint ? <p className="max-w-md font-mono text-xs text-zinc-400 light:text-zinc-500">{probe.hint}</p> : null}
      {docsHref ? (
        <a
          className="mt-1 inline-flex items-center gap-1.5 rounded-lg border border-zinc-700 bg-zinc-800 px-3 py-1.5 text-xs font-semibold text-zinc-100 no-underline transition-colors hover:border-violet-500 light:border-zinc-300 light:bg-zinc-100 light:text-zinc-900"
          href={docsHref}
          target="_blank"
          rel="noreferrer"
        >
          <BookOpen size={12} />
          {docsLabel}
        </a>
      ) : null}
    </div>
  );
}
