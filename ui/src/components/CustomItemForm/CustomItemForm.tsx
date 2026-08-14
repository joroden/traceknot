import { useState } from "react";
import type { WorkItemRow } from "../../types/workItem";

export interface CustomItemFormProps {
  onSubmit: (item: WorkItemRow) => void;
}

const inputClassName =
  "w-full rounded-lg border border-zinc-700 bg-zinc-950 px-3 py-2 text-sm text-zinc-100 outline-none transition-colors placeholder:text-zinc-500 focus:border-violet-500 light:border-zinc-300 light:bg-white light:text-zinc-900 light:placeholder:text-zinc-400";

export function CustomItemForm({ onSubmit }: CustomItemFormProps) {
  const [key, setKey] = useState("");
  const [title, setTitle] = useState("");

  const canSubmit = key.trim() !== "" && title.trim() !== "";

  const submit = () => {
    if (!canSubmit) {
      return;
    }
    onSubmit({
      key: key.trim(),
      title: title.trim(),
      project: "",
      provider: "custom",
      updatedAtUnixMs: null,
    });
  };

  return (
    <div className="flex flex-1 flex-col gap-3 rounded-lg border border-zinc-800 bg-zinc-900 p-5 light:border-zinc-200 light:bg-white">
      <div className="flex flex-col gap-1.5">
        <label className="text-xs font-semibold text-zinc-400 light:text-zinc-500" htmlFor="custom-item-key">
          ID
        </label>
        <input
          id="custom-item-key"
          className={inputClassName}
          type="text"
          autoFocus
          placeholder="e.g. TK-101"
          value={key}
          onChange={(event) => setKey(event.target.value)}
        />
      </div>
      <div className="flex flex-col gap-1.5">
        <label className="text-xs font-semibold text-zinc-400 light:text-zinc-500" htmlFor="custom-item-title">
          Title
        </label>
        <input
          id="custom-item-title"
          className={inputClassName}
          type="text"
          placeholder="What is this session about?"
          value={title}
          onChange={(event) => setTitle(event.target.value)}
          onKeyDown={(event) => {
            if (event.key === "Enter") {
              submit();
            }
          }}
        />
      </div>
      <button
        type="button"
        className="mt-1 cursor-pointer self-end rounded-lg border border-violet-600 bg-violet-600 px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-violet-500 disabled:cursor-default disabled:opacity-60"
        onClick={submit}
        disabled={!canSubmit}
      >
        Use this item
      </button>
    </div>
  );
}
