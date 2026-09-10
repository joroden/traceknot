import { useCallback, useEffect, useState } from "react";
import { getSettings, patchSettings } from "../../settings/api";
import { mergePatch, stateToPatch } from "../../settings/patch";
import type { SettingsPatch, SettingsState } from "../../settings/types";

export interface UseSetupDraft {
  draft: SettingsState | null;
  loadError: string | null;
  patch: (patch: SettingsPatch) => void;
  finish: () => void;
  finishing: boolean;
  finishError: string | null;
  finished: boolean;
}

function isUnconfigured(state: SettingsState): boolean {
  return !state.autostart_on && state.hooks.every((hook) => !hook.enabled);
}

function recommendedDefaults(state: SettingsState): SettingsState {
  if (!isUnconfigured(state)) {
    return state;
  }
  return {
    ...state,
    autostart_on: true,
    hooks: state.hooks.map((hook) => ({ ...hook, enabled: true })),
  };
}

export function useSetupDraft(): UseSetupDraft {
  const [draft, setDraft] = useState<SettingsState | null>(null);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [finishing, setFinishing] = useState(false);
  const [finishError, setFinishError] = useState<string | null>(null);
  const [finished, setFinished] = useState(false);

  useEffect(() => {
    let cancelled = false;
    getSettings()
      .then((data) => {
        if (!cancelled) {
          setDraft(recommendedDefaults(data));
        }
      })
      .catch((reason: unknown) => {
        if (!cancelled) {
          setLoadError(reason instanceof Error ? reason.message : String(reason));
        }
      });
    return () => {
      cancelled = true;
    };
  }, []);

  const patch = useCallback((next: SettingsPatch) => {
    setDraft((current) => (current ? mergePatch(current, next) : current));
  }, []);

  const finish = useCallback(() => {
    if (!draft) {
      return;
    }
    setFinishing(true);
    setFinishError(null);
    patchSettings(stateToPatch(draft))
      .then(() => {
        setFinished(true);
      })
      .catch((reason: unknown) => {
        setFinishError(reason instanceof Error ? reason.message : String(reason));
      })
      .finally(() => {
        setFinishing(false);
      });
  }, [draft]);

  return { draft, loadError, patch, finish, finishing, finishError, finished };
}
