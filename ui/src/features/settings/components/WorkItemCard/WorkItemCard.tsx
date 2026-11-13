import { SettingsCard, SettingRow } from "../../../../components/SettingsCard";
import { SegmentedControl } from "../../../../components/SegmentedControl";
import { Switch } from "../../../../components/Switch";

export interface WorkItemCardProps {
  requireWorkItem: boolean;
  skipNonInteractivePicker: boolean;
  onChange: (required: boolean) => void;
  onSkipNonInteractivePickerChange: (skip: boolean) => void;
}

export function WorkItemCard({
  requireWorkItem,
  skipNonInteractivePicker,
  onChange,
  onSkipNonInteractivePickerChange,
}: WorkItemCardProps) {
  return (
    <SettingsCard title="Work items">
      <SettingRow
        title="Work item assignment"
        description="When required, your agent won't be able to start working on a session until a work item is selected for it."
        control={
          <SegmentedControl
            options={[
              { value: "optional", label: "Optional" },
              { value: "required", label: "Required" },
            ]}
            value={requireWorkItem ? "required" : "optional"}
            onChange={(value) => onChange(value === "required")}
          />
        }
      />
      <SettingRow
        title="Skip picker for non-interactive runs"
        description="Let non-interactive CLI runs continue without selecting a work item, even when assignment is required. Telemetry is still collected."
        control={
          <Switch
            label="Skip picker for non-interactive runs"
            checked={skipNonInteractivePicker}
            onChange={onSkipNonInteractivePickerChange}
          />
        }
      />
    </SettingsCard>
  );
}
