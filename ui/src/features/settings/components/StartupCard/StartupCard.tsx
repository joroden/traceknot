import { SettingsCard, SettingRow } from "../../../../components/SettingsCard";
import { Switch } from "../../../../components/Switch";

export interface StartupCardProps {
  autostartOn: boolean;
  onChange: (on: boolean) => void;
}

export function StartupCard({ autostartOn, onChange }: StartupCardProps) {
  return (
    <SettingsCard title="Startup">
      <SettingRow
        title="Start on login"
        description="Start traceknot automatically when your computer boots up."
        control={<Switch checked={autostartOn} onChange={onChange} label="Start on login" />}
      />
    </SettingsCard>
  );
}
