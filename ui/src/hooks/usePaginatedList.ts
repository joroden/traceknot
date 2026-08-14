import { useCallback, useEffect, useRef, useState } from "react";

export interface PaginatedPage<T> {
  items: T[];
  totalCount: number;
}

export interface UsePaginatedListOptions {
  pageSize?: number;
  refetchOnVisible?: boolean;
}

export interface UsePaginatedListResult<T> {
  items: T[];
  totalCount: number;
  loading: boolean;
  loadingMore: boolean;
  error: string | null;
  loadMore: () => void;
  removeItem: (predicate: (item: T) => boolean) => void;
}

export function usePaginatedList<T, F>(
  filter: F,
  fetchPage: (filter: F, offset: number, limit: number) => Promise<PaginatedPage<T>>,
  options: UsePaginatedListOptions = {},
): UsePaginatedListResult<T> {
  const pageSize = options.pageSize ?? 100;
  const [items, setItems] = useState<T[]>([]);
  const [totalCount, setTotalCount] = useState(0);
  const [loading, setLoading] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const filterKey = JSON.stringify(filter);
  const offsetRef = useRef(0);
  const totalCountRef = useRef(0);
  const loadingMoreRef = useRef(false);
  const filterKeyRef = useRef(filterKey);
  filterKeyRef.current = filterKey;

  const loadFirstPage = useCallback(() => {
    let cancelled = false;
    setLoading(true);
    setError(null);
    offsetRef.current = 0;

    fetchPage(filter, 0, pageSize)
      .then((page) => {
        if (cancelled) {
          return;
        }
        setItems(page.items);
        setTotalCount(page.totalCount);
        totalCountRef.current = page.totalCount;
        offsetRef.current = page.items.length;
      })
      .catch((reason: unknown) => {
        if (!cancelled) {
          setError(reason instanceof Error ? reason.message : String(reason));
        }
      })
      .finally(() => {
        if (!cancelled) {
          setLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [filterKey]);

  useEffect(() => loadFirstPage(), [loadFirstPage]);

  useEffect(() => {
    if (!options.refetchOnVisible) {
      return;
    }
    function onVisible() {
      if (document.visibilityState === "visible") {
        loadFirstPage();
      }
    }
    document.addEventListener("visibilitychange", onVisible);
    return () => document.removeEventListener("visibilitychange", onVisible);
  }, [loadFirstPage, options.refetchOnVisible]);

  const loadMore = useCallback(() => {
    if (loadingMoreRef.current || offsetRef.current >= totalCountRef.current) {
      return;
    }
    loadingMoreRef.current = true;
    setLoadingMore(true);
    const requestFilterKey = filterKey;

    fetchPage(filter, offsetRef.current, pageSize)
      .then((page) => {
        if (requestFilterKey !== filterKeyRef.current) {
          return;
        }
        setItems((previous) => [...previous, ...page.items]);
        offsetRef.current += page.items.length;
      })
      .catch(() => {})
      .finally(() => {
        loadingMoreRef.current = false;
        setLoadingMore(false);
      });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [filterKey]);

  const removeItem = useCallback((predicate: (item: T) => boolean) => {
    setItems((previous) => previous.filter((item) => !predicate(item)));
    setTotalCount((previous) => Math.max(0, previous - 1));
    totalCountRef.current = Math.max(0, totalCountRef.current - 1);
    offsetRef.current = Math.max(0, offsetRef.current - 1);
  }, []);

  return { items, totalCount, loading, loadingMore, error, loadMore, removeItem };
}
