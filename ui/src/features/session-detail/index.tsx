import { useParams } from "react-router";
import { EmptyState } from "../../components/EmptyState";
import { SessionHeader } from "./components/SessionHeader";
import { TreeToolbar } from "./components/TreeToolbar";
import { TreeTable } from "./components/TreeTable";
import { DetailDrawer } from "./components/DetailDrawer";
import { useSessionDetail } from "./hooks/useSessionDetail";

export function SessionDetailPage() {
  const { id } = useParams();
  const {
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
  } = useSessionDetail(id);

  if (loading) {
    return <EmptyState title="Loading session tree…" />;
  }
  if (error || !meta) {
    return (
      <EmptyState
        title={error ?? "Session not found"}
      />
    );
  }

  return (
    <div className="flex h-full min-h-0 flex-col gap-3">
      <div className="flex flex-wrap items-center gap-4 rounded-lg border border-zinc-800 bg-zinc-900 px-4 py-3 light:border-zinc-200 light:bg-white">
        <SessionHeader meta={meta} />
        <div className="hidden h-5 w-px shrink-0 bg-zinc-800 sm:block light:bg-zinc-200" />
        <TreeToolbar
          query={query}
          onQueryChange={setQuery}
          visibleCount={visible.length}
          onExpandAll={expandAll}
          onCollapseAll={collapseAll}
        />
      </div>
      <div className="flex min-h-0 flex-1 gap-3">
        <TreeTable
          rows={visible}
          derived={derived}
          expanded={expanded}
          selectedId={selectedId}
          onToggle={toggleExpanded}
          onSelect={selectNode}
        />
        {selectedId && (
          <DetailDrawer
            nodeId={selectedId}
            launchNodeId={derived.byId.get(selectedId)?.launchNodeId ?? null}
            treeRow={derived.byId.get(selectedId) ?? null}
            onClose={() => selectNode(null)}
          />
        )}
      </div>
    </div>
  );
}
