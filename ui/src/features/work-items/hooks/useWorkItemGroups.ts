import { usePaginatedList, type PaginatedPage } from "../../../hooks/usePaginatedList";
import { getWorkItemGroups, type WorkItemGroup, type WorkItemGroupsFilter } from "../api";

export interface UseWorkItemGroupsResult {
  groups: WorkItemGroup[];
  totalCount: number;
  loading: boolean;
  loadingMore: boolean;
  error: string | null;
  loadMore: () => void;
}

function fetchGroupsPage(
  filter: WorkItemGroupsFilter,
  offset: number,
  limit: number,
): Promise<PaginatedPage<WorkItemGroup>> {
  return getWorkItemGroups({ ...filter, offset, limit }).then((response) => ({
    items: response.groups,
    totalCount: response.total_count,
  }));
}

export function useWorkItemGroups(filter: WorkItemGroupsFilter): UseWorkItemGroupsResult {
  const { items, totalCount, loading, loadingMore, error, loadMore } = usePaginatedList(filter, fetchGroupsPage);
  return { groups: items, totalCount, loading, loadingMore, error, loadMore };
}
