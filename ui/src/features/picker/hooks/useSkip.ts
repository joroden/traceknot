import { useCallback } from "react";
import { postSkip } from "../../../types/workItem";

export function useSkip(sessionID: string | null): {
  submit: () => Promise<void>;
} {
  return {
    submit: useCallback(() => postSkip(sessionID), [sessionID]),
  };
}
