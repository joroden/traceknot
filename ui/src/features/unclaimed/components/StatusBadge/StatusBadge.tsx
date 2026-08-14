import { Badge } from "../../../../components/Badge";
import type { UnclaimedStatus } from "../../api";

const LABEL: Record<UnclaimedStatus, string> = {
  pending: "Pending",
  skipped: "Skipped",
  unclaimed: "Unclaimed",
};

export function StatusBadge({ status }: { status: UnclaimedStatus }) {
  return <Badge tone={status === "pending" ? "amber" : "neutral"}>{LABEL[status]}</Badge>;
}
