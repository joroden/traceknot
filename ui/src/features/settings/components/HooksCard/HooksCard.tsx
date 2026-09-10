import { SettingsCard, SettingRow } from "../../../../components/SettingsCard";
import { Switch } from "../../../../components/Switch";
import type { ProviderState } from "../../types";
import { HOOK_COPY, providerCopy } from "../../providerCopy";
import { CopilotRestartNotice } from "../CopilotRestartNotice";

export interface HooksCardProps {
  hooks: ProviderState[];
  onToggle: (binary: string, enabled: boolean) => void;
  restartNoticeActionable?: boolean;
}

export function HooksCard({ hooks, onToggle, restartNoticeActionable = true }: HooksCardProps) {
  return (
    <SettingsCard title="Agent hooks" description="Which agents traceknot should collect telemetry from.">
      {hooks.map((hook) => {
        const copy = providerCopy(HOOK_COPY, hook.binary);
        return (
          <div key={hook.binary}>
            <SettingRow
              title={copy.label}
              description={copy.description}
              control={
                <Switch
                  checked={hook.enabled}
                  onChange={(enabled) => onToggle(hook.binary, enabled)}
                  label={copy.label}
                />
              }
            />
            {hook.binary === "copilot" && hook.enabled ? (
              <div className="px-5 pb-4">
                <CopilotRestartNotice actionable={restartNoticeActionable} />
              </div>
            ) : null}
          </div>
        );
      })}
    </SettingsCard>
  );
}
