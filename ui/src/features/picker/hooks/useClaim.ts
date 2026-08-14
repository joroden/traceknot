import { useCallback } from "react";
import { postClaim, type WorkItemRow } from "../../../types/workItem";

export function useClaim(sessionID: string | null): {
  submit: (item: WorkItemRow) => Promise<void>;
} {
  return {
    submit: useCallback(
      (item: WorkItemRow) => postClaim(sessionID, item),
      [sessionID],
    ),
  };
}
