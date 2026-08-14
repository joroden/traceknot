import { useEffect, useState } from "react";
import { getContext } from "../api";
import { getProviders, type ProviderProbe } from "../../../types/workItem";

export interface PickerLoad {
  providers: ProviderProbe[];
  prompt: string;
}

export function usePickerLoad(sessionID: string | null): {
  load: PickerLoad | null;
  error: string | null;
} {
  const [load, setLoad] = useState<PickerLoad | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    Promise.all([getProviders(), getContext(sessionID)])
      .then(([providers, prompt]) => {
        if (!cancelled) {
          setLoad({ providers, prompt });
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
