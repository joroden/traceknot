import { ChevronDown, ChevronRight } from "lucide-react";
import { Badge } from "../../../../components/Badge";
import type { WorkItemGroup } from "../../api";

export interface GroupRowLabelProps {
  group: WorkItemGroup;
  expanded: boolean;
}

export function GroupRowLabel({ group, expanded }: GroupRowLabelProps) {
  const Chevron = expanded ? ChevronDown : ChevronRight;
  return (
    <div className="flex min-w-0 items-center gap-2">
      <Chevron className="h-3.5 w-3.5 shrink-0 text-zinc-500" />
      {group.is_unclaimed ? (
        <Badge tone="amber">Unclaimed</Badge>
      ) : (
        <span className="truncate font-mono text-xs font-semibold text-zinc-100 light:text-zinc-900">
          {group.work_item_key}
        </span>
      )}
    </div>
  );
}
