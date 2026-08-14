import { useEffect, useState } from "react";
import { getDashboard, type DashboardData, type DashboardRequest } from "../api";

export function useDashboard(request: DashboardRequest): {
  data: DashboardData | null;
  error: string | null;
} {
  const [data, setData] = useState<DashboardData | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    setError(null);
    getDashboard(request)
      .then((payload) => {
        if (!cancelled) {
          setData(payload);
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
  }, [request]);

  return { data, error };
}
