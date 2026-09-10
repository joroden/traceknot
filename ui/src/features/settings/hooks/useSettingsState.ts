import { useCallback, useEffect, useState } from "react";
import { getSettings, patchSettings } from "../api";
import type { SettingsPatch, SettingsState } from "../types";

export interface UseSettingsState {
  state: SettingsState | null;
  error: string | null;
  patch: (patch: SettingsPatch) => void;
}

export function useSettingsState(): UseSettingsState {
  const [state, setState] = useState<SettingsState | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    getSettings()
      .then((data) => {
        if (!cancelled) {
          setState(data);
        }
      })
      .catch((reason: unknown) => {
        if (!cancelled) {
          setError(reason instanceof Error ? reason.message : String(reason));
        }
      });
    return () => {
      cancelled = true;
    };
  }, []);

  const patch = useCallback((next: SettingsPatch) => {
    patchSettings(next)
      .then((data) => {
        setState(data);
      })
      .catch((reason: unknown) => {
        setError(reason instanceof Error ? reason.message : String(reason));
      });
  }, []);

  return { state, error, patch };
}
