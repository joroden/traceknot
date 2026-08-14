import { useEffect, useState } from "react";
import { getRecent, recentToRow, type WorkItemRow } from "../types/workItem";

export interface UseRecentItemsResult {
  rows: WorkItemRow[];
  error: string | null;
}

export function useRecentItems(enabled: boolean, query: string): UseRecentItemsResult {
  const [allRows, setAllRows] = useState<WorkItemRow[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!enabled) {
      return;
    }
    let cancelled = false;
    setError(null);
    getRecent(null)
      .then((items) => {
        if (!cancelled) {
          setAllRows(items.map(recentToRow));
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
  }, [enabled]);

  const lowerQuery = query.trim().toLowerCase();
  const rows = lowerQuery
    ? allRows.filter(
        (item) => item.key.toLowerCase().includes(lowerQuery) || item.title.toLowerCase().includes(lowerQuery),
      )
    : allRows;

  return { rows, error };
}
