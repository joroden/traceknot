import { useCallback, useEffect, useMemo, useState } from "react";
import { useSearchParams } from "react-router";
import {
  fetchSessionTree,
  type SessionMeta,
  type TreeNodeRow,
} from "../api";

export interface MergedRow extends TreeNodeRow {
  launchNodeId: string | null;
}

export interface TreeDerived {
  byId: Map<string, MergedRow>;
  childrenOf: Map<string, MergedRow[]>;
  roots: MergedRow[];
  agentOfLaunch: Map<string, string>;
}

export interface VisibleRow {
  row: MergedRow;
  depth: number;
}

function mergeLaunches(
  byId: Map<string, MergedRow>,
  childrenOf: Map<string, MergedRow[]>,
  roots: MergedRow[],
  agentOfLaunch: Map<string, string>,
): void {
  for (const launch of byId.values()) {
    if (launch.kind !== "tool_call") {
      continue;
    }
    const kids = childrenOf.get(launch.node_id) ?? [];
    if (kids.length !== 1 || kids[0].kind !== "agent") {
      continue;
    }
    const agent = kids[0];
    agent.launchNodeId = launch.node_id;
    agentOfLaunch.set(launch.node_id, agent.node_id);
    const parentID = launch.parent_node_id;
    if (parentID && byId.has(parentID)) {
      const siblings = childrenOf.get(parentID) ?? [];
      childrenOf.set(
        parentID,
        siblings.map((sibling) => (sibling.node_id === launch.node_id ? agent : sibling)),
      );
    } else {
      roots.splice(
        roots.findIndex((root) => root.node_id === launch.node_id),
        1,
        agent,
      );
    }
  }
}

function buildDerived(nodes: TreeNodeRow[]): TreeDerived {
  const byId = new Map<string, MergedRow>();
  const childrenOf = new Map<string, MergedRow[]>();
  const agentOfLaunch = new Map<string, string>();
  for (const node of nodes) {
    const row: MergedRow = { ...node, launchNodeId: null };
    byId.set(row.node_id, row);
  }
  const roots: MergedRow[] = [];
  for (const node of byId.values()) {
    const parent = node.parent_node_id ? byId.get(node.parent_node_id) : undefined;
    if (parent) {
      const siblings = childrenOf.get(parent.node_id) ?? [];
      siblings.push(node);
      childrenOf.set(parent.node_id, siblings);
    } else {
      roots.push(node);
    }
  }
  const order = (list: MergedRow[]) =>
    list.sort(
      (a, b) =>
        (a.started_at_unix_ms ?? 0) - (b.started_at_unix_ms ?? 0),
    );
  order(roots);
  for (const siblings of childrenOf.values()) {
    order(siblings);
  }
  mergeLaunches(byId, childrenOf, roots, agentOfLaunch);
  return { byId, childrenOf, roots, agentOfLaunch };
}

function matchesQuery(row: TreeNodeRow, query: string): boolean {
  const needle = query.trim().toLowerCase();
  if (!needle) {
    return true;
  }
  return [row.name, row.tool_name, row.agent_name].some(
    (value) => value !== null && value.toLowerCase().includes(needle),
  );
}

function emitTreeRows(
  derived: TreeDerived,
  scope: MergedRow[],
  expanded: Set<string>,
  query: string,
  out: VisibleRow[],
): void {
  const needle = query.trim().toLowerCase();
  const matched = new Set<string>();
  const ancestorSet = new Set<string>();
  if (needle) {
    for (const node of derived.byId.values()) {
      if (matchesQuery(node, needle)) {
        matched.add(node.node_id);
        const agentID = derived.agentOfLaunch.get(node.node_id);
        if (agentID) {
          matched.add(agentID);
        }
      }
    }
    for (const node of derived.byId.values()) {
      if (!matched.has(node.node_id)) {
        continue;
      }
      let current = node.parent_node_id;
      while (current !== null && current !== undefined) {
        if (ancestorSet.has(current)) {
          break;
        }
        ancestorSet.add(current);
        current = derived.byId.get(current)?.parent_node_id ?? null;
      }
    }
  }
  const relevant = (node: MergedRow) =>
    !needle || matched.has(node.node_id) || ancestorSet.has(node.node_id);
  const emit = (node: MergedRow, depth: number) => {
    if (relevant(node)) {
      out.push({ row: node, depth });
    }
    const kids = derived.childrenOf.get(node.node_id) ?? [];
    if (kids.length === 0 || (!needle && !expanded.has(node.node_id))) {
      return;
    }
    for (const child of kids) {
      if (relevant(child)) {
        emit(child, depth + 1);
      }
    }
  };
  for (const root of scope) {
    emit(root, 0);
  }
}

export function useSessionDetail(sessionId: string | undefined) {
  const [searchParams, setSearchParams] = useSearchParams();
  const [meta, setMeta] = useState<SessionMeta | null>(null);
  const [nodes, setNodes] = useState<TreeNodeRow[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [expanded, setExpanded] = useState<Set<string>>(new Set());
  const [query, setQuery] = useState("");

  const selectedId = searchParams.get("item");
  const focusId = searchParams.get("focus");

  useEffect(() => {
    if (!sessionId) {
      return;
    }
    let cancelled = false;
    setLoading(true);
    setError(null);
    setMeta(null);
    setNodes([]);
    setExpanded(new Set());
    setQuery("");
    fetchSessionTree(sessionId)
      .then((payload) => {
        if (!cancelled) {
          setMeta(payload.session);
          setNodes(payload.nodes);
          setLoading(false);
        }
      })
      .catch((reason: unknown) => {
        if (!cancelled) {
          setError(reason instanceof Error ? reason.message : String(reason));
          setLoading(false);
        }
      });
    return () => {
      cancelled = true;
    };
  }, [sessionId]);

  const derived = useMemo(() => buildDerived(nodes), [nodes]);

  const scope = useMemo<MergedRow[]>(() => {
    if (!focusId) {
      return derived.roots;
    }
    const focused = derived.byId.get(focusId);
    return focused ? [focused] : derived.roots;
  }, [derived, focusId]);

  const visible = useMemo<VisibleRow[]>(() => {
    const out: VisibleRow[] = [];
    emitTreeRows(derived, scope, expanded, query, out);
    return out;
  }, [derived, scope, expanded, query]);

  const toggleExpanded = useCallback((nodeId: string) => {
    setExpanded((prev) => {
      const next = new Set(prev);
      if (next.has(nodeId)) {
        next.delete(nodeId);
      } else {
        next.add(nodeId);
      }
      return next;
    });
  }, []);

  const expandAll = useCallback(() => {
    setExpanded(new Set(nodes.map((node) => node.node_id)));
  }, [nodes]);

  const collapseAll = useCallback(() => {
    setExpanded(new Set());
  }, []);

  const updateParams = useCallback(
    (patch: { item?: string | null; focus?: string | null }) => {
      setSearchParams(
        (prev) => {
          const next = new URLSearchParams(prev);
          for (const [key, value] of Object.entries(patch)) {
            if (value === null || value === undefined || value === "") {
              next.delete(key);
            } else {
              next.set(key, value);
            }
          }
          return next;
        },
        { replace: true },
      );
    },
    [setSearchParams],
  );

  const selectNode = useCallback(
    (nodeId: string | null) => {
      updateParams({ item: nodeId });
    },
    [updateParams],
  );

  return {
    loading,
    error,
    meta,
    derived,
    visible,
    expanded,
    toggleExpanded,
    expandAll,
    collapseAll,
    query,
    setQuery,
    selectedId,
    selectNode,
  };
}
