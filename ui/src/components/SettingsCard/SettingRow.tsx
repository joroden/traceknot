import type { ReactNode } from "react";

export interface SettingRowProps {
  title: string;
  description?: string;
  control: ReactNode;
}

export function SettingRow({ title, description, control }: SettingRowProps) {
  return (
    <div className="flex items-start justify-between gap-4 px-5 py-4">
      <div className="min-w-0">
        <p className="text-sm font-medium">{title}</p>
        {description ? (
          <p className="mt-0.5 text-xs text-zinc-400 light:text-zinc-500">{description}</p>
        ) : null}
      </div>
      <div className="flex-shrink-0 pt-0.5">{control}</div>
    </div>
  );
}
