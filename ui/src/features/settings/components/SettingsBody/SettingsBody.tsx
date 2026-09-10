import type { SettingsPatch, SettingsState } from "../../types";
import { StartupCard } from "../StartupCard";
import { HooksCard } from "../HooksCard";
import { WorkItemCard } from "../WorkItemCard";
import { SkillCard } from "../SkillCard";

export interface SettingsBodyProps {
  state: SettingsState;
  onPatch: (patch: SettingsPatch) => void;
  restartNoticeActionable?: boolean;
}

export function SettingsBody({ state, onPatch, restartNoticeActionable = true }: SettingsBodyProps) {
  return (
    <div className="flex flex-col gap-5">
      <StartupCard autostartOn={state.autostart_on} onChange={(on) => onPatch({ autostart_on: on })} />
      <HooksCard
        hooks={state.hooks}
        onToggle={(binary, enabled) => onPatch({ hooks: { [binary]: enabled } })}
        restartNoticeActionable={restartNoticeActionable}
      />
      <WorkItemCard
        requireWorkItem={state.require_work_item}
        onChange={(required) => onPatch({ require_work_item: required })}
      />
      <SkillCard skills={state.skills} onToggle={(binary, enabled) => onPatch({ skills: { [binary]: enabled } })} />
    </div>
  );
}
