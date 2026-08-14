import { History, PenLine } from "lucide-react";
import type { ProviderProbe } from "../../types/workItem";
import { BrandIcon } from "../BrandIcon";

export interface TabBarProps {
  providers: ProviderProbe[];
  activeTab: string;
  onSelect: (tab: string) => void;
}

const RECENT_TAB = "recent";
const CUSTOM_TAB = "custom";

function statusChip(status: ProviderProbe["status"]) {
  if (status === "not_authenticated") {
    return (
      <span className="rounded-full border border-amber-500/35 bg-amber-500/15 px-1.5 py-0.5 font-mono text-xs text-amber-500">
        Auth Req
      </span>
    );
  }
  if (status === "cli_missing") {
    return (
      <span className="rounded-full border border-zinc-700 bg-zinc-800 px-1.5 py-0.5 font-mono text-xs text-zinc-500 light:border-zinc-300 light:bg-zinc-100 light:text-zinc-400">
        No CLI
      </span>
    );
  }
  return null;
}

export function TabBar({ providers, activeTab, onSelect }: TabBarProps) {
  return (
    <div className="flex gap-2 overflow-x-auto pb-1">
      <button
        type="button"
        className={
          activeTab === RECENT_TAB
            ? "flex cursor-pointer items-center gap-[7px] whitespace-nowrap rounded-lg border border-violet-600 bg-violet-600 px-3.5 py-2 text-sm font-semibold text-white transition-colors"
            : "flex cursor-pointer items-center gap-[7px] whitespace-nowrap rounded-lg border border-zinc-800 bg-zinc-900 px-3.5 py-2 text-sm font-semibold text-zinc-400 transition-colors hover:border-zinc-700 hover:bg-zinc-800 hover:text-zinc-100 light:border-zinc-200 light:bg-white light:text-zinc-500 light:hover:border-zinc-300 light:hover:bg-zinc-100 light:hover:text-zinc-900"
        }
        onClick={() => onSelect(RECENT_TAB)}
      >
        <span className="inline-flex">
          <History size={14} />
        </span>
        Recent
      </button>
      {providers.map((item) => {
        const active = activeTab === item.provider;
        return (
          <button
            key={item.provider}
            type="button"
            className={
              active
                ? "flex cursor-pointer items-center gap-[7px] whitespace-nowrap rounded-lg border border-violet-600 bg-violet-600 px-3.5 py-2 text-sm font-semibold text-white transition-colors"
                : "flex cursor-pointer items-center gap-[7px] whitespace-nowrap rounded-lg border border-zinc-800 bg-zinc-900 px-3.5 py-2 text-sm font-semibold text-zinc-400 transition-colors hover:border-zinc-700 hover:bg-zinc-800 hover:text-zinc-100 light:border-zinc-200 light:bg-white light:text-zinc-500 light:hover:border-zinc-300 light:hover:bg-zinc-100 light:hover:text-zinc-900"
            }
            onClick={() => onSelect(item.provider)}
          >
            <span className="inline-flex">
              <BrandIcon provider={item.provider} />
            </span>
            {item.provider}
            {item.status === "available" ? (
              <span className="size-[7px] rounded-full bg-emerald-500 shadow-[0_0_6px_rgba(16,185,129,0.5)]" />
            ) : (
              statusChip(item.status)
            )}
          </button>
        );
      })}
      <button
        type="button"
        className={
          activeTab === CUSTOM_TAB
            ? "flex cursor-pointer items-center gap-[7px] whitespace-nowrap rounded-lg border border-violet-600 bg-violet-600 px-3.5 py-2 text-sm font-semibold text-white transition-colors"
            : "flex cursor-pointer items-center gap-[7px] whitespace-nowrap rounded-lg border border-zinc-800 bg-zinc-900 px-3.5 py-2 text-sm font-semibold text-zinc-400 transition-colors hover:border-zinc-700 hover:bg-zinc-800 hover:text-zinc-100 light:border-zinc-200 light:bg-white light:text-zinc-500 light:hover:border-zinc-300 light:hover:bg-zinc-100 light:hover:text-zinc-900"
        }
        onClick={() => onSelect(CUSTOM_TAB)}
      >
        <span className="inline-flex">
          <PenLine size={14} />
        </span>
        Custom
      </button>
    </div>
  );
}
