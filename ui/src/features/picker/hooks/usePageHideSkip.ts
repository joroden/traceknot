import { useEffect, useRef } from "react";
import { postSkipBeacon } from "../../../types/workItem";

export function usePageHideSkip(sessionID: string | null, resolved: boolean): void {
  const resolvedRef = useRef(resolved);
  resolvedRef.current = resolved;

  useEffect(() => {
    if (!sessionID) {
      return;
    }
    const handlePageHide = (event: PageTransitionEvent) => {
      if (event.persisted || resolvedRef.current) {
        return;
      }
      postSkipBeacon(sessionID);
    };
    window.addEventListener("pagehide", handlePageHide);
    return () => window.removeEventListener("pagehide", handlePageHide);
  }, [sessionID]);
}
