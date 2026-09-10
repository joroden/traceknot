import { LogoMark } from "../../components/Logo";
import { ThemeToggle } from "../../components/ThemeToggle";
import { DaemonUnreachable } from "../../components/DaemonUnreachable";
import { SettingsBody } from "../settings/components/SettingsBody";
import { useSetupDraft } from "./hooks/useSetupDraft";
import { CompletionScreen } from "./components/CompletionScreen";

export function SetupPage() {
  const { draft, loadError, patch, finish, finishing, finishError, finished } = useSetupDraft();

  if (loadError) {
    return (
      <div className="mx-auto flex h-full max-w-2xl flex-col gap-3.5 px-7 py-5.5">
        <DaemonUnreachable error={loadError} />
      </div>
    );
  }

  if (!draft) {
    return (
      <div className="mx-auto flex h-full max-w-2xl flex-col gap-3.5 px-7 py-5.5">
        <div className="m-auto flex flex-col items-center gap-2.5 text-sm text-zinc-400">
          <span className="size-[26px] animate-spin rounded-full border-[3px] border-zinc-700 border-t-violet-500" />
          <p>Loading setup…</p>
        </div>
      </div>
    );
  }

  return (
    <div className="mx-auto flex h-full max-w-2xl flex-col gap-5 px-7 py-8">
      <header className="flex items-center justify-between gap-4 border-b border-zinc-800 pb-4 light:border-zinc-200">
        <div className="flex min-w-0 items-center gap-3">
          <span className="inline-flex size-[38px] flex-shrink-0 items-center justify-center rounded-lg bg-violet-600 text-white">
            <LogoMark size={18} />
          </span>
          <div>
            <h1 className="text-lg font-bold">Set up traceknot</h1>
            <p className="mt-0.5 text-sm text-zinc-400 light:text-zinc-500">
              Local telemetry for your coding agents — nothing leaves this machine.
            </p>
          </div>
        </div>
        <ThemeToggle />
      </header>

      <main className="flex-1 overflow-y-auto">
        <SettingsBody state={draft} onPatch={patch} restartNoticeActionable={false} />
      </main>

      <footer className="flex flex-col items-end gap-2 border-t border-zinc-800 pt-4 light:border-zinc-200">
        {finishError ? <p className="text-xs text-red-500">{finishError}</p> : null}
        <button
          type="button"
          className="cursor-pointer rounded-lg border border-violet-600 bg-violet-600 px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-violet-500 disabled:cursor-default disabled:opacity-60"
          onClick={finish}
          disabled={finishing}
        >
          {finishing ? "Finishing setup…" : "Finish setup"}
        </button>
      </footer>

      {finished ? (
        <CompletionScreen copilotEnabled={draft.hooks.some((hook) => hook.binary === "copilot" && hook.enabled)} />
      ) : null}
    </div>
  );
}
