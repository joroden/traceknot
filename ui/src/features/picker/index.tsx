import { useCallback, useMemo, useState } from "react";
import { Forward } from "lucide-react";
import { LogoMark } from "../../components/Logo";
import { ThemeToggle } from "../../components/ThemeToggle";
import { ClaimDialog } from "../../components/ClaimDialog";
import { CustomItemForm } from "../../components/CustomItemForm";
import { DaemonUnreachable } from "../../components/DaemonUnreachable";
import { ProviderPanel } from "../../components/ProviderPanel";
import { ResultList } from "../../components/ResultList";
import { TabBar } from "../../components/TabBar";
import type { ProviderProbe, WorkItemRow } from "../../types/workItem";
import { useRecentItems } from "../../hooks/useRecentItems";
import { useWorkItemSearch } from "../../hooks/useWorkItemSearch";
import { useClaim } from "./hooks/useClaim";
import { usePickerLoad } from "./hooks/usePickerLoad";
import { usePickerKeyboard } from "./hooks/usePickerKeyboard";
import { useSkip } from "./hooks/useSkip";
import { OutcomeOverlay, type Outcome } from "./components/OutcomeOverlay";
import { PromptBanner } from "./components/PromptBanner";

const RECENT_TAB = "recent";
const CUSTOM_TAB = "custom";

export interface PickerPageProps {
  sessionID: string | null;
}

export function PickerPage({ sessionID }: PickerPageProps) {
  const { load, error: loadError } = usePickerLoad(sessionID);
  const [activeTab, setActiveTab] = useState(RECENT_TAB);
  const [query, setQuery] = useState("");
  const [confirming, setConfirming] = useState<WorkItemRow | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [claimError, setClaimError] = useState<string | null>(null);
  const [outcome, setOutcome] = useState<Outcome>(null);

  const providers: ProviderProbe[] = load?.providers ?? [];
  const prompt = load?.prompt ?? "";

  const providerProbe = useMemo(
    () =>
      activeTab === RECENT_TAB
        ? null
        : providers.find((item) => item.provider === activeTab) ?? null,
    [activeTab, providers],
  );
  const providerAvailable =
    providerProbe === null || providerProbe.status === "available";

  const recent = useRecentItems(true, query);
  const search = useWorkItemSearch(
    activeTab === RECENT_TAB || activeTab === CUSTOM_TAB
      ? null
      : providerAvailable
        ? activeTab
        : null,
    query,
    { debounceMs: 250 },
  );

  const rows = useMemo(() => {
    if (activeTab === RECENT_TAB) {
      return recent.rows;
    }
    if (activeTab === CUSTOM_TAB || !providerAvailable) {
      return [];
    }
    return search.rows;
  }, [activeTab, providerAvailable, recent.rows, search.rows]);

  const { submit: submitClaim } = useClaim(sessionID);
  const { submit: submitSkip } = useSkip(sessionID);

  const skip = useCallback(() => {
    setOutcome({ kind: "skipped", item: null });
    void submitSkip();
  }, [submitSkip]);

  const handleEscape = useCallback(() => {
    if (query) {
      setQuery("");
    } else {
      skip();
    }
  }, [query, skip]);

  const handleConfirm = useCallback(
    (index: number) => setConfirming(rows[index]),
    [rows],
  );

  const { highlighted, onHighlight } = usePickerKeyboard({
    enabled: outcome === null && confirming === null,
    rows,
    onEscape: handleEscape,
    onConfirm: handleConfirm,
  });

  const confirmClaim = useCallback(async () => {
    if (!confirming) {
      return;
    }
    setSubmitting(true);
    setClaimError(null);
    try {
      await submitClaim(confirming);
      setOutcome({ kind: "claimed", item: confirming });
      setConfirming(null);
    } catch (reason) {
      setClaimError(reason instanceof Error ? reason.message : String(reason));
    } finally {
      setSubmitting(false);
    }
  }, [confirming, submitClaim]);

  const closeDialog = useCallback(() => {
    setConfirming(null);
    setClaimError(null);
  }, []);

  if (loadError) {
    return (
      <div className="mx-auto flex h-full max-w-6xl flex-col gap-3.5 px-7 py-5.5">
        <DaemonUnreachable error={loadError} showStartHint />
      </div>
    );
  }

  if (!load) {
    return (
      <div className="mx-auto flex h-full max-w-6xl flex-col gap-3.5 px-7 py-5.5">
        <div className="m-auto flex flex-col items-center gap-2.5 text-sm text-zinc-400">
          <span className="size-[26px] animate-spin rounded-full border-[3px] border-zinc-700 border-t-violet-500" />
          <p>Loading picker…</p>
        </div>
      </div>
    );
  }

  const emptyTitle =
    activeTab === RECENT_TAB
      ? "No recent items yet — pick a platform tab to browse, or skip."
      : query.trim()
        ? `No matches for "${query.trim()}" in ${activeTab}`
        : `No items found in ${activeTab}.`;

  const listLoading =
    activeTab !== RECENT_TAB && providerAvailable && search.searching;

  return (
    <div className="mx-auto flex h-full max-w-6xl flex-col gap-3.5 px-7 py-5.5">
      <header className="flex items-center justify-between gap-4 border-b border-zinc-800 pb-4 light:border-zinc-200">
        <div className="flex min-w-0 items-center gap-3">
          <span className="inline-flex size-[38px] flex-shrink-0 items-center justify-center rounded-lg bg-violet-600 text-white">
            <LogoMark size={18} />
          </span>
          <div>
            <h1 className="text-lg font-bold">Attach Context</h1>
            <p className="mt-0.5 text-sm text-zinc-400 light:text-zinc-500">
              Match this session to a work item, or skip
            </p>
          </div>
        </div>
        <div className="flex flex-shrink-0 items-center gap-2">
          <ThemeToggle />
          <button
            type="button"
            className="inline-flex cursor-pointer items-center gap-[7px] rounded-lg border border-zinc-700 bg-zinc-800 px-3.5 py-2 text-sm font-bold text-zinc-100 transition-colors hover:border-amber-500 light:border-zinc-300 light:bg-zinc-100 light:text-zinc-900 light:hover:border-amber-600"
            onClick={skip}
          >
            <Forward size={13} />
            Skip Context
          </button>
        </div>
      </header>

      {prompt ? <PromptBanner prompt={prompt} /> : null}

      <main className="flex min-h-0 flex-1 flex-col gap-3">
        <TabBar
          providers={providers}
          activeTab={activeTab}
          onSelect={setActiveTab}
        />

        {activeTab === CUSTOM_TAB ? null : providerAvailable ? (
          <div className="relative">
            <input
              className="w-full rounded-lg border border-zinc-700 bg-zinc-900 px-3.5 py-2.5 text-sm text-zinc-100 outline-none transition-colors placeholder:text-zinc-500 focus:border-violet-500 light:border-zinc-300 light:bg-white light:text-zinc-900 light:placeholder:text-zinc-400"
              type="text"
              autoFocus
              placeholder={
                activeTab === RECENT_TAB
                  ? "Filter recent items by key or title"
                  : `Search ${activeTab} by key or title`
              }
              value={query}
              onChange={(event) => setQuery(event.target.value)}
            />
            {search.searching ? (
              <span className="absolute right-3.5 top-1/2 -translate-y-1/2 font-mono text-xs text-zinc-500">
                searching…
              </span>
            ) : null}
          </div>
        ) : null}

        {activeTab === CUSTOM_TAB ? (
          <CustomItemForm onSubmit={setConfirming} />
        ) : providerAvailable ? (
          <ResultList
            rows={rows}
            loading={listLoading}
            error={search.error}
            highlighted={highlighted}
            onHighlight={onHighlight}
            onSelect={setConfirming}
            emptyTitle={emptyTitle}
          />
        ) : providerProbe ? (
          <ProviderPanel probe={providerProbe} />
        ) : null}
      </main>

      <ClaimDialog
        item={confirming}
        submitting={submitting}
        error={claimError}
        onConfirm={() => void confirmClaim()}
        onClose={closeDialog}
      />
      <OutcomeOverlay outcome={outcome} />
    </div>
  );
}
