import { SettingsCard, SettingRow } from "../../../../components/SettingsCard";
import { Switch } from "../../../../components/Switch";
import type { ProviderState } from "../../types";
import { SKILL_COPY, providerCopy } from "../../providerCopy";

export interface SkillCardProps {
  skills: ProviderState[];
  onToggle: (binary: string, enabled: boolean) => void;
}

export function SkillCard({ skills, onToggle }: SkillCardProps) {
  return (
    <SettingsCard
      title="Session analysis skill"
      description="Lets an agent inspect its own session cost and telemetry when you ask it to."
    >
      {skills.map((skill) => {
        const copy = providerCopy(SKILL_COPY, skill.binary);
        return (
          <SettingRow
            key={skill.binary}
            title={copy.label}
            description={copy.description}
            control={
              <Switch
                checked={skill.enabled}
                onChange={(enabled) => onToggle(skill.binary, enabled)}
                label={copy.label}
              />
            }
          />
        );
      })}
    </SettingsCard>
  );
}
