import { Dialog } from "../Dialog";
import type { WorkItemRow } from "../../types/workItem";

export interface ClaimDialogProps {
  item: WorkItemRow | null;
  submitting: boolean;
  error: string | null;
  onConfirm: () => void;
  onClose: () => void;
}

export function ClaimDialog({ item, submitting, error, onConfirm, onClose }: ClaimDialogProps) {
  return (
    <Dialog
      open={item !== null}
      onOpenChange={(open) => {
        if (!open && !submitting) {
          onClose();
        }
      }}
      title="Attach this item?"
      description={
        item ? `Attach ${item.key} — ${item.title} — to this session?` : undefined
      }
    >
      <div className="mt-4 flex flex-wrap justify-end gap-2">
        {error ? <p className="mb-1 basis-full text-xs text-red-500">{error}</p> : null}
        <button
          type="button"
          className="cursor-pointer rounded-lg border border-zinc-700 bg-transparent px-4 py-2 text-sm font-semibold text-zinc-400 transition-colors hover:text-zinc-100 disabled:cursor-default disabled:opacity-60 light:border-zinc-300 light:text-zinc-500 light:hover:text-zinc-900"
          onClick={onClose}
          disabled={submitting}
        >
          Go back
        </button>
        <button
          type="button"
          className="cursor-pointer rounded-lg border border-violet-600 bg-violet-600 px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-violet-500 disabled:cursor-default disabled:opacity-60"
          onClick={onConfirm}
          disabled={submitting}
        >
          {submitting ? "Attaching…" : "Confirm"}
        </button>
      </div>
    </Dialog>
  );
}
