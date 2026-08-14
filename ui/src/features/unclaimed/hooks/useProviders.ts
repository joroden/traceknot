import { useEffect, useState } from "react";
import { getProviders, type ProviderProbe } from "../../../types/workItem";

export function useProviders(): ProviderProbe[] {
  const [providers, setProviders] = useState<ProviderProbe[]>([]);

  useEffect(() => {
    let cancelled = false;
    getProviders()
      .then((items) => {
        if (!cancelled) {
          setProviders(items);
        }
      })
      .catch(() => {});
    return () => {
      cancelled = true;
    };
  }, []);

  return providers;
}
