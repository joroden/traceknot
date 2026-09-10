import { DaemonUnreachable } from "../../components/DaemonUnreachable";
import { SettingsBody } from "./components/SettingsBody";
import { useSettingsState } from "./hooks/useSettingsState";

export function SettingsPage() {
  const { state, error, patch } = useSettingsState();

  if (error) {
    return <DaemonUnreachable error={error} showStartHint />;
  }

  if (!state) {
    return (
      <div className="m-auto flex flex-col items-center gap-2.5 text-sm text-zinc-400">
        <span className="size-[26px] animate-spin rounded-full border-[3px] border-zinc-700 border-t-violet-500" />
        <p>Loading settings…</p>
      </div>
    );
  }

  return (
    <div className="mx-auto w-full max-w-2xl">
      <SettingsBody state={state} onPatch={patch} />
    </div>
  );
}
