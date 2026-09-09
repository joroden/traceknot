import { useEffect, useState } from "react";
import { getContext } from "../api";
import { getSettings } from "../../settings/api";
import { getProviders, type ProviderProbe } from "../../../types/workItem";

export interface PickerLoad {
  providers: ProviderProbe[];
  prompt: string;
  requireWorkItem: boolean;
}

export function usePickerLoad(sessionID: string | null): {
  load: PickerLoad | null;
  error: string | null;
} {
  const [load, setLoad] = useState<PickerLoad | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    Promise.all([getProviders(), getContext(sessionID), getSettings()])
      .then(([providers, prompt, settings]) => {
        if (!cancelled) {
          setLoad({ providers, prompt, requireWorkItem: settings.require_work_item });
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
  }, [sessionID]);

  return { load, error };
}
