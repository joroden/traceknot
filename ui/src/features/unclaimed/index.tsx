import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useNavigate } from "react-router";
import { ClaimDialog } from "../../components/ClaimDialog";
import { DaemonUnreachable } from "../../components/DaemonUnreachable";
import { EmptyState } from "../../components/EmptyState";
import { SearchInput } from "../../components/SearchInput";
import { SegmentedControl, type SegmentedOption } from "../../components/SegmentedControl";
import { Skeleton } from "../../components/Skeleton";
import { postClaim, type WorkItemRow } from "../../types/workItem";
import type { UnclaimedFilter } from "./api";
import { QueueRow } from "./components/QueueRow";
import { SearchDialog } from "./components/SearchDialog";
import { useDebouncedValue } from "../../hooks/useDebouncedValue";
import { useProviders } from "./hooks/useProviders";
import { useUnclaimedSessions } from "./hooks/useUnclaimedSessions";

const PROVIDER_OPTIONS: SegmentedOption<string>[] = [
  { value: "", label: "All" },
  { value: "claude", label: "Claude" },
  { value: "codex", label: "Codex" },
  { value: "copilot", label: "Copilot" },
];

interface ConfirmTarget {
  sessionId: string;
  item: WorkItemRow;
}

export function UnclaimedPage() {
  const navigate = useNavigate();
  const [provider, setProvider] = useState("");
  const [promptQuery, setPromptQuery] = useState("");
  const debouncedQuery = useDebouncedValue(promptQuery, 250);

  const filter = useMemo<UnclaimedFilter>(
    () => ({ provider: provider || undefined, q: debouncedQuery || undefined, sort: "started:desc" }),
    [provider, debouncedQuery],
  );

  const queue = useUnclaimedSessions(filter);
  const providers = useProviders();

  const [searchSessionId, setSearchSessionId] = useState<string | null>(null);
  const [confirming, setConfirming] = useState<ConfirmTarget | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [claimError, setClaimError] = useState<string | null>(null);

  const openSearch = useCallback((sessionId: string) => setSearchSessionId(sessionId), []);
  const closeSearch = useCallback(() => setSearchSessionId(null), []);

  const selectFromSearch = useCallback(
    (item: WorkItemRow) => {
      if (!searchSessionId) {
        return;
      }
      setConfirming({ sessionId: searchSessionId, item });
      setSearchSessionId(null);
    },
    [searchSessionId],
  );

  const closeConfirm = useCallback(() => {
    if (submitting) {
      return;
    }
    setConfirming(null);
    setClaimError(null);
  }, [submitting]);

  const confirmClaim = useCallback(async () => {
    if (!confirming) {
      return;
    }
    setSubmitting(true);
    setClaimError(null);
    try {
      await postClaim(confirming.sessionId, confirming.item);
      queue.removeRow(confirming.sessionId);
      setConfirming(null);
    } catch (reason) {
      setClaimError(reason instanceof Error ? reason.message : String(reason));
    } finally {
      setSubmitting(false);
    }
  }, [confirming, queue]);

  const sentinelRef = useRef<HTMLDivElement>(null);
  useEffect(() => {
    const node = sentinelRef.current;
    if (!node) {
      return;
    }
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0]?.isIntersecting) {
          queue.loadMore();
        }
      },
      { rootMargin: "200px" },
    );
    observer.observe(node);
    return () => observer.disconnect();
  }, [queue]);

  if (queue.error) {
    return <DaemonUnreachable error={queue.error} />;
  }

  const filtersActive = provider !== "" || promptQuery !== "";
  const emptyTitle = filtersActive
    ? "No unclaimed sessions match these filters."
    : "Queue clear — every session is claimed.";

  return (
    <div className="flex h-full min-h-0 flex-col gap-4">
      <div className="flex flex-wrap items-center justify-between gap-2">
        <div className="flex flex-wrap items-center gap-4">
          <SegmentedControl options={PROVIDER_OPTIONS} value={provider} onChange={setProvider} />
          <div className="w-[240px]">
            <SearchInput value={promptQuery} onChange={setPromptQuery} placeholder="Search titles…" />
          </div>
        </div>
        <span className="text-xs text-zinc-500">
          {queue.totalCount} unclaimed session{queue.totalCount === 1 ? "" : "s"}
        </span>
      </div>

      {queue.loading && queue.rows.length === 0 ? (
        <div className="flex flex-col gap-2.5">
          {Array.from({ length: 6 }, (_, index) => (
            <Skeleton key={index} className="h-[84px] w-full" />
          ))}
        </div>
      ) : queue.rows.length === 0 ? (
        <EmptyState title={emptyTitle} />
      ) : (
        <div className="flex min-h-0 flex-1 flex-col gap-2.5 overflow-y-auto pr-1">
          {queue.rows.map((session) => (
            <QueueRow
              key={session.session_id}
              session={session}
              onOpenSession={(sessionId) => navigate(`/sessions/${encodeURIComponent(sessionId)}`)}
              onFindWorkItem={() => openSearch(session.session_id)}
            />
          ))}
          <div ref={sentinelRef} />
          {queue.loadingMore ? <Skeleton className="h-[84px] w-full" /> : null}
        </div>
      )}

      <SearchDialog
        open={searchSessionId !== null}
        providers={providers}
        onSelect={selectFromSearch}
        onClose={closeSearch}
      />
      <ClaimDialog
        item={confirming?.item ?? null}
        submitting={submitting}
        error={claimError}
        onConfirm={() => void confirmClaim()}
        onClose={closeConfirm}
      />
    </div>
  );
}
