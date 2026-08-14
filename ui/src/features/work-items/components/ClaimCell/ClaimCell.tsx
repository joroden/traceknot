import { Badge } from "../../../../components/Badge";
import type { ClaimState } from "../../api";

export interface ClaimCellProps {
  claim: ClaimState;
}

const STATUS_LABEL: Record<Exclude<ClaimState["status"], "claimed">, string> = {
  pending: "Pending",
  skipped: "Skipped",
  unclaimed: "Unclaimed",
};

export function ClaimCell({ claim }: ClaimCellProps) {
  if (claim.status === "claimed") {
    return (
      <div className="flex min-w-0 flex-col">
        <span className="truncate font-mono text-xs font-semibold text-zinc-100 light:text-zinc-900">
          {claim.work_item_key}
        </span>
        <span className="truncate text-xs text-zinc-500">{claim.work_item_title}</span>
      </div>
    );
  }

  return <Badge tone={claim.status === "pending" ? "amber" : "neutral"}>{STATUS_LABEL[claim.status]}</Badge>;
}
