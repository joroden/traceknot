import type { ReactNode } from "react";

export interface SettingsCardProps {
  title: string;
  description?: string;
  children: ReactNode;
}

export function SettingsCard({ title, description, children }: SettingsCardProps) {
  return (
    <section className="rounded-lg border border-zinc-800 bg-zinc-900 light:border-zinc-200 light:bg-white">
      <div className="border-b border-zinc-800 px-5 py-4 light:border-zinc-200">
        <h2 className="text-sm font-semibold">{title}</h2>
        {description ? (
          <p className="mt-0.5 text-xs text-zinc-400 light:text-zinc-500">{description}</p>
        ) : null}
      </div>
      <div className="flex flex-col divide-y divide-zinc-800 light:divide-zinc-200">{children}</div>
    </section>
  );
}
