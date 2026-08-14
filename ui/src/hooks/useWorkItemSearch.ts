import { useEffect, useRef, useState } from "react";
import { searchItems, toRow, type WorkItem, type WorkItemRow } from "../types/workItem";

export interface UseWorkItemSearchOptions {
  debounceMs?: number;
}

export interface UseWorkItemSearchResult {
  rows: WorkItemRow[];
  searching: boolean;
  error: string | null;
}

export function useWorkItemSearch(
  provider: string | null,
  query: string,
  options: UseWorkItemSearchOptions = {},
): UseWorkItemSearchResult {
  const debounceMs = options.debounceMs ?? 0;
  const [rows, setRows] = useState<WorkItemRow[]>([]);
  const [searching, setSearching] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const previousProvider = useRef<string | null>(null);

  useEffect(() => {
    if (previousProvider.current !== provider) {
      setRows([]);
      previousProvider.current = provider;
    }
    if (!provider) {
      return;
    }
    let cancelled = false;
    const timer = setTimeout(
      () => {
        setSearching(true);
        setError(null);
        searchItems(provider, query.trim())
          .then((items: WorkItem[]) => {
            if (!cancelled) {
              setRows(items.map((item) => toRow(item, provider)));
            }
          })
          .catch((reason: unknown) => {
            if (!cancelled) {
              setError(reason instanceof Error ? reason.message : String(reason));
            }
          })
          .finally(() => {
            if (!cancelled) {
              setSearching(false);
            }
          });
      },
      query.trim() ? debounceMs : 0,
    );
    return () => {
      cancelled = true;
      clearTimeout(timer);
    };
  }, [provider, query, debounceMs]);

  return { rows, searching, error };
}
