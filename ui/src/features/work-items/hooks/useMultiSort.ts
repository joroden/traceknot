import { useCallback, useState } from "react";
import type { SortEntry } from "../columns";
import type { SortKey } from "../api";

export function useMultiSort(initial: SortEntry[] = []): [SortEntry[], (key: SortKey, additive: boolean) => void] {
  const [sorts, setSorts] = useState<SortEntry[]>(initial);

  const onSortToggle = useCallback((key: SortKey, additive: boolean) => {
    setSorts((current) => {
      const existingIndex = current.findIndex((entry) => entry.key === key);
      if (!additive) {
        if (existingIndex === -1) {
          return [{ key, dir: "desc" }];
        }
        if (current.length > 1) {
          return [{ key, dir: "desc" }];
        }
        return current[0].dir === "desc" ? [{ key, dir: "asc" }] : [];
      }
      if (existingIndex === -1) {
        return [...current, { key, dir: "desc" }];
      }
      const next = [...current];
      if (next[existingIndex].dir === "desc") {
        next[existingIndex] = { key, dir: "asc" };
        return next;
      }
      next.splice(existingIndex, 1);
      return next;
    });
  }, []);

  return [sorts, onSortToggle];
}
