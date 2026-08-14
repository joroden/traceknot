import { usePaginatedList, type PaginatedPage } from "../../../hooks/usePaginatedList";
import { getUnclaimedSessions, type UnclaimedFilter, type UnclaimedSession } from "../api";

export interface UseUnclaimedSessionsResult {
  rows: UnclaimedSession[];
  totalCount: number;
  loading: boolean;
  loadingMore: boolean;
  error: string | null;
  loadMore: () => void;
  removeRow: (sessionId: string) => void;
}

function fetchUnclaimedPage(
  filter: UnclaimedFilter,
  offset: number,
  limit: number,
): Promise<PaginatedPage<UnclaimedSession>> {
  return getUnclaimedSessions({ ...filter, offset, limit }).then((response) => ({
    items: response.sessions,
    totalCount: response.total_count,
  }));
}

export function useUnclaimedSessions(filter: UnclaimedFilter): UseUnclaimedSessionsResult {
  const { items, totalCount, loading, loadingMore, error, loadMore, removeItem } = usePaginatedList(
    filter,
    fetchUnclaimedPage,
  );
  return {
    rows: items,
    totalCount,
    loading,
    loadingMore,
    error,
    loadMore,
    removeRow: (sessionId: string) => removeItem((row) => row.session_id === sessionId),
  };
}
