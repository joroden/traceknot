import { usePaginatedList, type PaginatedPage } from "../../../hooks/usePaginatedList";
import { getSessions, type SessionRow, type SessionsFilter } from "../api";

export interface UseSessionsResult {
  rows: SessionRow[];
  totalCount: number;
  loading: boolean;
  loadingMore: boolean;
  error: string | null;
  loadMore: () => void;
}

function fetchSessionsPage(filter: SessionsFilter, offset: number, limit: number): Promise<PaginatedPage<SessionRow>> {
  return getSessions({ ...filter, offset, limit }).then((response) => ({
    items: response.sessions,
    totalCount: response.total_count,
  }));
}

export function useSessions(filter: SessionsFilter): UseSessionsResult {
  const { items, totalCount, loading, loadingMore, error, loadMore } = usePaginatedList(
    filter,
    fetchSessionsPage,
    { refetchOnVisible: true },
  );
  return { rows: items, totalCount, loading, loadingMore, error, loadMore };
}
