import { useEffect, useState } from "react";
import type { WorkItemRow } from "../../../types/workItem";

export interface PickerKeyboardOptions {
  enabled: boolean;
  rows: WorkItemRow[];
  onEscape: () => void;
  onConfirm: (index: number) => void;
}

export function usePickerKeyboard({
  enabled,
  rows,
  onEscape,
  onConfirm,
}: PickerKeyboardOptions): {
  highlighted: number;
  onHighlight: (index: number) => void;
} {
  const [highlighted, setHighlighted] = useState(rows.length > 0 ? 0 : -1);

  useEffect(() => {
    setHighlighted(rows.length > 0 ? 0 : -1);
  }, [rows.length]);

  useEffect(() => {
    if (!enabled) {
      return;
    }
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        onEscape();
        return;
      }
      if (rows.length === 0) {
        return;
      }
      if (event.key === "ArrowDown") {
        event.preventDefault();
        setHighlighted((index) => (index + 1) % rows.length);
      } else if (event.key === "ArrowUp") {
        event.preventDefault();
        setHighlighted((index) => (index - 1 + rows.length) % rows.length);
      } else if (event.key === "Enter") {
        event.preventDefault();
        const index = highlighted >= 0 ? highlighted : 0;
        if (index < rows.length) {
          onConfirm(index);
        }
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [enabled, rows.length, highlighted, onEscape, onConfirm]);

  return { highlighted, onHighlight: setHighlighted };
}
