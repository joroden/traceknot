import { useEffect, useState } from "react";
import { CustomItemForm } from "../../../../components/CustomItemForm";
import { Dialog } from "../../../../components/Dialog";
import { ProviderPanel } from "../../../../components/ProviderPanel";
import { ResultList } from "../../../../components/ResultList";
import { TabBar } from "../../../../components/TabBar";
import type { ProviderProbe, WorkItemRow } from "../../../../types/workItem";
import { useDebouncedValue } from "../../../../hooks/useDebouncedValue";
import { useRecentItems } from "../../../../hooks/useRecentItems";
import { useWorkItemSearch } from "../../../../hooks/useWorkItemSearch";

const RECENT_TAB = "recent";
const CUSTOM_TAB = "custom";

export interface SearchDialogProps {
  open: boolean;
  providers: ProviderProbe[];
  onSelect: (item: WorkItemRow) => void;
  onClose: () => void;
}

export function SearchDialog({ open, providers, onSelect, onClose }: SearchDialogProps) {
  const [activeTab, setActiveTab] = useState(RECENT_TAB);
  const [query, setQuery] = useState("");
  const debouncedQuery = useDebouncedValue(query, 250);

  useEffect(() => {
    if (open) {
      setActiveTab(RECENT_TAB);
      setQuery("");
    }
  }, [open]);

  const providerProbe =
    activeTab === RECENT_TAB ? null : providers.find((item) => item.provider === activeTab) ?? null;
  const providerAvailable = providerProbe === null || providerProbe.status === "available";

  const recent = useRecentItems(open, debouncedQuery);
  const search = useWorkItemSearch(
    activeTab === RECENT_TAB || activeTab === CUSTOM_TAB
      ? null
      : providerAvailable
        ? activeTab
        : null,
    debouncedQuery,
  );

  const rows =
    activeTab === RECENT_TAB
      ? recent.rows
      : activeTab === CUSTOM_TAB
        ? []
        : providerAvailable
          ? search.rows
          : [];
  const emptyTitle =
    activeTab === RECENT_TAB
      ? "No recent work items yet — pick a platform tab to browse."
      : query.trim()
        ? `No matches for "${query.trim()}" in ${activeTab}`
        : `No items found in ${activeTab}.`;

  return (
    <Dialog
      open={open}
      onOpenChange={(next) => {
        if (!next) {
          onClose();
        }
      }}
      title="Find a work item"
      widthClassName="w-[min(560px,calc(100vw-48px))]"
    >
      <div className="mt-3 flex max-h-[60vh] flex-col gap-3">
        <TabBar providers={providers} activeTab={activeTab} onSelect={setActiveTab} />
        {activeTab === CUSTOM_TAB ? null : providerAvailable ? (
          <input
            className="w-full rounded-lg border border-zinc-700 bg-zinc-900 px-3 py-2 text-sm text-zinc-100 outline-none transition-colors placeholder:text-zinc-500 focus:border-violet-500 light:border-zinc-300 light:bg-white light:text-zinc-900 light:placeholder:text-zinc-400"
            type="text"
            autoFocus
            placeholder={activeTab === RECENT_TAB ? "Filter recent items" : `Search ${activeTab}`}
            value={query}
            onChange={(event) => setQuery(event.target.value)}
          />
        ) : null}
        <div className="min-h-0 flex-1 overflow-y-auto">
          {activeTab === CUSTOM_TAB ? (
            <CustomItemForm onSubmit={onSelect} />
          ) : providerAvailable ? (
            <ResultList
              rows={rows}
              loading={activeTab !== RECENT_TAB && providerAvailable && search.searching}
              error={activeTab === RECENT_TAB ? recent.error : search.error}
              highlighted={-1}
              onHighlight={() => {}}
              onSelect={onSelect}
              emptyTitle={emptyTitle}
            />
          ) : providerProbe ? (
            <ProviderPanel probe={providerProbe} />
          ) : null}
        </div>
      </div>
    </Dialog>
  );
}
