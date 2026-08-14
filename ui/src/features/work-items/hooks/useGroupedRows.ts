import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import {
  encodeSort,
  getSessions,
  groupIdentity,
  type GroupSortKey,
  type SessionsFilter,
  type SortDir,
  type WorkItemGroup,
  type WorkItemGroupsFilter,
} from "../api";
import type { SortEntry, WorkItemsRow } from "../columns";
import { useWorkItemGroups } from "./useWorkItemGroups";

const SESSIONS_PAGE_SIZE = 100;

interface GroupSessionsState {
  sessions: WorkItemsRow[];
  totalCount: number;
  loading: boolean;
}

export interface WorkItemScope {
  key: string;
  provider: string;
}

export interface GroupSort {
  key: GroupSortKey;
  dir: SortDir;
}

export interface UseGroupedRowsArgs {
  baseFilter: Pick<SessionsFilter, "provider" | "model" | "q" | "start_unix_ms" | "end_unix_ms">;
  sorts: SortEntry[];
  groupSort: GroupSort;
  scope: WorkItemScope | null;
}

export interface UseGroupedRowsResult {
  rows: WorkItemsRow[];
  loading: boolean;
  error: string | null;
  expandedGroups: Set<string>;
  toggleGroup: (group: WorkItemGroup) => void;
  onEndReached: (row: WorkItemsRow) => void;
}

export function useGroupedRows({ baseFilter, sorts, groupSort, scope }: UseGroupedRowsArgs): UseGroupedRowsResult {
  const groupsFilter = useMemo<WorkItemGroupsFilter>(
    () => ({
      ...baseFilter,
      work_item_key: scope?.key,
      work_item_provider: scope?.provider,
      sort: `${groupSort.key}:${groupSort.dir}`,
    }),
    [baseFilter, scope, groupSort],
  );

  const {
    groups,
    totalCount: groupsTotalCount,
    loading: groupsLoading,
    loadingMore: groupsLoadingMore,
    error,
    loadMore: loadMoreGroups,
  } = useWorkItemGroups(groupsFilter);

  const [expandedGroups, setExpandedGroups] = useState<Set<string>>(
    () => new Set(scope ? [`${scope.provider}:${scope.key}`] : []),
  );
  const [sessionsByGroup, setSessionsByGroup] = useState<Record<string, GroupSessionsState>>({});

  const cursorsRef = useRef<Map<string, number>>(new Map());
  const loadingRef = useRef<Set<string>>(new Set());
  const sortKey = sorts.map((entry) => `${entry.key}:${entry.dir}`).join(",");

  const sessionsFilterFor = useCallback(
    (group: WorkItemGroup): SessionsFilter => ({
      ...baseFilter,
      ...(group.is_unclaimed
        ? { unclaimed: "true" as const }
        : { work_item_key: group.work_item_key, work_item_provider: group.work_item_provider }),
      sort: encodeSort(sorts),
    }),
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [baseFilter, sortKey],
  );

  const loadGroupSessions = useCallback(
    (group: WorkItemGroup, append: boolean) => {
      const identity = groupIdentity(group);
      if (loadingRef.current.has(identity)) {
        return;
      }
      loadingRef.current.add(identity);
      const offset = append ? (cursorsRef.current.get(identity) ?? 0) : 0;
      if (!append) {
        cursorsRef.current.set(identity, 0);
      }
      setSessionsByGroup((current) => ({
        ...current,
        [identity]: {
          sessions: append ? (current[identity]?.sessions ?? []) : [],
          totalCount: current[identity]?.totalCount ?? 0,
          loading: true,
        },
      }));

      getSessions({ ...sessionsFilterFor(group), offset, limit: SESSIONS_PAGE_SIZE })
        .then((response) => {
          const newRows: WorkItemsRow[] = response.sessions.map((session) => ({
            kind: "session",
            session,
            groupIdentity: identity,
          }));
          cursorsRef.current.set(identity, offset + response.sessions.length);
          setSessionsByGroup((current) => ({
            ...current,
            [identity]: {
              sessions: append ? [...(current[identity]?.sessions ?? []), ...newRows] : newRows,
              totalCount: response.total_count,
              loading: false,
            },
          }));
        })
        .catch(() => {
          setSessionsByGroup((current) => ({
            ...current,
            [identity]: { ...(current[identity] ?? { sessions: [], totalCount: 0 }), loading: false },
          }));
        })
        .finally(() => {
          loadingRef.current.delete(identity);
        });
    },
    [sessionsFilterFor],
  );

  useEffect(() => {
    setSessionsByGroup({});
    cursorsRef.current.clear();
    for (const group of groups) {
      const identity = groupIdentity(group);
      if (expandedGroups.has(identity)) {
        loadGroupSessions(group, false);
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [groupsFilter.provider, groupsFilter.model, groupsFilter.q, groupsFilter.start_unix_ms, groupsFilter.end_unix_ms, sortKey]);

  useEffect(() => {
    for (const group of groups) {
      const identity = groupIdentity(group);
      if (expandedGroups.has(identity) && !sessionsByGroup[identity]) {
        loadGroupSessions(group, false);
      }
    }
  }, [groups, expandedGroups, sessionsByGroup, loadGroupSessions]);

  const toggleGroup = useCallback(
    (group: WorkItemGroup) => {
      const identity = groupIdentity(group);
      const willExpand = !expandedGroups.has(identity);
      setExpandedGroups((current) => {
        const next = new Set(current);
        if (next.has(identity)) {
          next.delete(identity);
        } else {
          next.add(identity);
        }
        return next;
      });
      if (willExpand && !sessionsByGroup[identity]) {
        loadGroupSessions(group, false);
      }
    },
    [expandedGroups, sessionsByGroup, loadGroupSessions],
  );

  const rows = useMemo<WorkItemsRow[]>(() => {
    const flat: WorkItemsRow[] = [];
    for (const group of groups) {
      flat.push({ kind: "group", group });
      const identity = groupIdentity(group);
      if (expandedGroups.has(identity)) {
        flat.push(...(sessionsByGroup[identity]?.sessions ?? []));
      }
    }
    return flat;
  }, [groups, expandedGroups, sessionsByGroup]);

  const onEndReached = useCallback(
    (row: WorkItemsRow) => {
      if (row.kind === "session" && row.groupIdentity) {
        const state = sessionsByGroup[row.groupIdentity];
        if (state && !state.loading && state.sessions.length < state.totalCount) {
          const group = groups.find((candidate) => groupIdentity(candidate) === row.groupIdentity);
          if (group) {
            loadGroupSessions(group, true);
          }
        }
      }
      if (groups.length < groupsTotalCount && !groupsLoadingMore) {
        loadMoreGroups();
      }
    },
    [groups, groupsTotalCount, groupsLoadingMore, sessionsByGroup, loadGroupSessions, loadMoreGroups],
  );

  return { rows, loading: groupsLoading, error, expandedGroups, toggleGroup, onEndReached };
}
